package interpreter

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type structuredBlock struct {
	kind  string
	start string
	end   string
	line  int
}

type loweredStructuredLine struct {
	labels []string
	text   string
}

var (
	structuredJumpTarget     = regexp.MustCompile(`(?i)\b(GOTO|GOSUB|THEN)\s+([0-9]+)\b`)
	structuredInternalTarget = regexp.MustCompile(`@[A-Za-z0-9_]+`)
)

// PrepareSource lowers the explicitly marked structured dialect while leaving
// ordinary numbered source and its diagnostics unchanged.
func PrepareSource(source string) (string, error) {
	structured := false
	for _, raw := range strings.Split(strings.ReplaceAll(source, "\r\n", "\n"), "\n") {
		text := strings.TrimSpace(stripStructuredComment(raw))
		_, statement := splitStructuredLabel(text)
		if strings.EqualFold(strings.TrimSpace(statement), "sub_start") {
			structured = true
			break
		}
	}
	if !structured {
		return source, nil
	}
	return LowerStructuredBASIC(source)
}

// LowerStructuredBASIC converts the annotated corpus dialect into strict
// numbered BASIC. Numbered Microsoft BASIC should bypass this compatibility
// lowering and continue to be parsed directly.
func LowerStructuredBASIC(source string) (string, error) {
	var (
		lines   []loweredStructuredLine
		blocks  []structuredBlock
		counter int
	)
	appendLine := func(text string, labels ...string) {
		lines = append(lines, loweredStructuredLine{labels: labels, text: text})
	}
	newLabel := func(prefix string) string {
		counter++
		return fmt.Sprintf("__%s_%d", prefix, counter)
	}
	nearestLoopEnd := func() string {
		for blockIndex := len(blocks) - 1; blockIndex >= 0; blockIndex-- {
			if blocks[blockIndex].kind == "LOOP" {
				return blocks[blockIndex].end
			}
		}
		return ""
	}

	for index, raw := range strings.Split(strings.ReplaceAll(source, "\r\n", "\n"), "\n") {
		sourceLine := index + 1
		text := strings.TrimSpace(stripStructuredComment(raw))
		if text == "" {
			continue
		}
		label, statement := splitStructuredLabel(text)
		text = strings.TrimSpace(statement)
		sourceLabels := []string{}
		if label != "" {
			sourceLabels = append(sourceLabels, label)
		}
		if text == "" {
			appendLine("REM LABEL", sourceLabels...)
			continue
		}
		text = replaceStructuredEquality(text)
		lower := strings.ToLower(strings.TrimSpace(text))

		switch lower {
		case "sub_start":
			appendLine("REM SUB_START", sourceLabels...)
			continue
		case "loop":
			block := structuredBlock{kind: "LOOP", start: newLabel("loop"), end: newLabel("endloop"), line: sourceLine}
			appendLine("REM LOOP", append(sourceLabels, block.start)...)
			blocks = append(blocks, block)
			continue
		case "endloop":
			if len(blocks) == 0 || blocks[len(blocks)-1].kind != "LOOP" {
				return "", fmt.Errorf("source line %d: ENDLOOP without LOOP", sourceLine)
			}
			block := blocks[len(blocks)-1]
			blocks = blocks[:len(blocks)-1]
			appendLine("GOTO @"+block.start, sourceLabels...)
			appendLine("REM ENDLOOP", block.end)
			continue
		case "break":
			end := nearestLoopEnd()
			if end == "" {
				return "", fmt.Errorf("source line %d: BREAK without LOOP", sourceLine)
			}
			appendLine("GOTO @"+end, sourceLabels...)
			continue
		case "endif":
			if len(blocks) == 0 || blocks[len(blocks)-1].kind != "IF" {
				return "", fmt.Errorf("source line %d: ENDIF without IF", sourceLine)
			}
			block := blocks[len(blocks)-1]
			blocks = blocks[:len(blocks)-1]
			appendLine("REM ENDIF", append(sourceLabels, block.end)...)
			continue
		}

		if condition, ok := structuredBlockCondition(text); ok {
			block := structuredBlock{kind: "IF", start: newLabel("if"), end: newLabel("endif"), line: sourceLine}
			appendLine("IF "+condition+" THEN GOTO @"+block.start, sourceLabels...)
			appendLine("GOTO @" + block.end)
			appendLine("REM IF", block.start)
			blocks = append(blocks, block)
			continue
		}

		if strings.Contains(strings.ToLower(text), "then break") {
			end := nearestLoopEnd()
			if end == "" {
				return "", fmt.Errorf("source line %d: BREAK without LOOP", sourceLine)
			}
			text = replaceThenBreak(text, "THEN GOTO @"+end)
		}
		appendLine(text, sourceLabels...)
	}

	if len(blocks) != 0 {
		block := blocks[len(blocks)-1]
		return "", fmt.Errorf("source line %d: unclosed %s", block.line, block.kind)
	}

	labelLines := make(map[string]int)
	for index, line := range lines {
		lineNumber := (index + 1) * 10
		for _, label := range line.labels {
			if _, duplicate := labelLines[label]; duplicate {
				return "", fmt.Errorf("duplicate structured label %s", label)
			}
			labelLines[label] = lineNumber
		}
	}

	var output strings.Builder
	for index, line := range lines {
		text, err := rewriteStructuredTargets(line.text, labelLines)
		if err != nil {
			return "", err
		}
		fmt.Fprintf(&output, "%d %s\n", (index+1)*10, text)
	}
	return output.String(), nil
}

func stripStructuredComment(line string) string {
	inString := false
	for index, char := range line {
		if char == '"' {
			inString = !inString
		}
		if char == '#' && !inString {
			return line[:index]
		}
	}
	return line
}

func replaceStructuredEquality(line string) string {
	var output strings.Builder
	inString := false
	for index := 0; index < len(line); index++ {
		if line[index] == '"' {
			inString = !inString
		}
		if !inString && line[index] == '=' && index+1 < len(line) && line[index+1] == '=' {
			output.WriteByte('=')
			index++
			continue
		}
		output.WriteByte(line[index])
	}
	return output.String()
}

func splitStructuredLabel(line string) (string, string) {
	index := 0
	for index < len(line) && line[index] >= '0' && line[index] <= '9' {
		index++
	}
	if index == 0 || index < len(line) && line[index] != ' ' && line[index] != '\t' {
		return "", line
	}
	return line[:index], strings.TrimSpace(line[index:])
}

func structuredBlockCondition(statement string) (string, bool) {
	trimmed := strings.TrimSpace(statement)
	lower := strings.ToLower(trimmed)
	if !strings.HasPrefix(lower, "if ") {
		return "", false
	}
	thenIndex := strings.LastIndex(lower, "then")
	if thenIndex >= 0 {
		if strings.TrimSpace(trimmed[thenIndex+len("then"):]) != "" {
			return "", false
		}
		return strings.TrimSpace(trimmed[len("if "):thenIndex]), true
	}
	return strings.TrimSpace(trimmed[len("if "):]), true
}

func replaceThenBreak(statement, replacement string) string {
	lower := strings.ToLower(statement)
	index := strings.Index(lower, "then break")
	if index < 0 {
		return statement
	}
	return statement[:index] + replacement + statement[index+len("then break"):]
}

func rewriteStructuredTargets(statement string, labels map[string]int) (string, error) {
	statement = structuredJumpTarget.ReplaceAllStringFunc(statement, func(target string) string {
		parts := strings.Fields(target)
		line, ok := labels[parts[len(parts)-1]]
		if !ok {
			return target
		}
		return target[:len(target)-len(parts[len(parts)-1])] + strconv.Itoa(line)
	})
	var unresolved error
	statement = structuredInternalTarget.ReplaceAllStringFunc(statement, func(target string) string {
		line, ok := labels[strings.TrimPrefix(target, "@")]
		if !ok {
			unresolved = fmt.Errorf("undefined structured label %s", target)
			return target
		}
		return strconv.Itoa(line)
	})
	if unresolved != nil {
		return "", unresolved
	}
	return statement, nil
}
