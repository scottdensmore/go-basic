package main

import (
	"fmt"
	"math"
	"strings"
	"time"
)

type Environment struct {
	Vars map[string]float64 // Basic variables are usually numbers. Strings supported too? 
    // "program.bas" only uses 'a' as number.
    // But print has string literal "*".
    // Let's use interface{}
    Store map[string]interface{}
}

func NewEnvironment() *Environment {
	return &Environment{Store: make(map[string]interface{})}
}

type LoopContext struct {
	VarName         string
	End             float64
	Step            float64
	BodyLineIndex   int // Index in LineNumbers where the loop body starts (the line AFTER ForStmt)
}

type Evaluator struct {
	Env              *Environment
	Program          *Program
	CurrentLineIndex int
	LoopStack        []*LoopContext
    OutputColumn     int // To track for TAB
}

func NewEvaluator(p *Program) *Evaluator {
	return &Evaluator{
		Env:              NewEnvironment(),
		Program:          p,
		CurrentLineIndex: 0,
		LoopStack:        []*LoopContext{},
        OutputColumn:     0,
	}
}

func (e *Evaluator) Run() {
	for e.CurrentLineIndex < len(e.Program.LineNumbers) {
		lineNum := e.Program.LineNumbers[e.CurrentLineIndex]
		stmt := e.Program.Lines[lineNum]
        
        startIdx := e.CurrentLineIndex
		e.evalStatement(stmt)
        
        if e.CurrentLineIndex == startIdx {
            e.CurrentLineIndex++
        }
	}
}

func (e *Evaluator) evalStatement(stmt Statement) {
	switch s := stmt.(type) {
	case *LetStatement:
		val := e.evalExpression(s.Value)
		e.Env.Store[strings.ToLower(s.Name.Value)] = val
	case *PrintStmt:
		e.evalPrintStatement(s)
	case *ForStatement:
		e.evalForStatement(s)
	case *NextStatement:
		e.evalNextStatement(s)
	case *SleepStatement:
		e.evalSleepStatement(s)
	}
}

func (e *Evaluator) evalPrintStatement(stmt *PrintStmt) {
    for _, item := range stmt.Items {
        if item.IsSeparator {
            // Semicolon: do nothing (suppress newline, which is default at end)
            // But if it's in middle? `PRINT A; B` -> A then B immediately.
            // TAB relies on OutputColumn.
            continue
        }
        
        val := e.evalExpression(item.Expr)
        var str string
        
        if tab, ok := val.(TabValue); ok {
            // TAB(n)
            target := tab.Pos
            if target > e.OutputColumn {
                str = strings.Repeat(" ", target - e.OutputColumn)
            }
        } else {
            str = fmt.Sprintf("%v", val)
            // Basic numbers often print with a leading space if positive? 6502 style.
            // MS Basic: Numbers are printed with a preceding space if positive, or minus if negative, and a trailing space.
            // Let's stick to simple fmt for now unless requested.
            // If it's a number
            if f, ok := val.(float64); ok {
                // If integer, print as integer
                if f == math.Trunc(f) {
                     str = fmt.Sprintf("%d", int64(f))
                } else {
                     str = fmt.Sprintf("%g", f)
                }
                
                // Add explicit spacing if this is what strict basic does, but `program.bas` uses TAB mostly.
                // However, `print tab(...); "*"` uses semi-colon.
            }
        }
        
        fmt.Print(str)
        e.OutputColumn += len(str)
    }
    
    // Check if last item is separator
    lastIsSeparator := false
    if len(stmt.Items) > 0 {
        if stmt.Items[len(stmt.Items)-1].IsSeparator {
            lastIsSeparator = true
        }
    }
    
    if !lastIsSeparator {
        fmt.Println()
        e.OutputColumn = 0
    }
}

func (e *Evaluator) evalForStatement(stmt *ForStatement) {
	startVal := e.evalExpression(stmt.Start)
	endVal := e.evalExpression(stmt.End)
	stepVal := e.evalExpression(stmt.Step)

	var start, end, step float64
	start = toFloat(startVal)
	end = toFloat(endVal)
	step = toFloat(stepVal)

	varName := strings.ToLower(stmt.Var.Value)
	e.Env.Store[varName] = start
    
    // Push context
    // The body starts at the NEXT line index
    e.LoopStack = append(e.LoopStack, &LoopContext{
        VarName: varName,
        End: end,
        Step: step,
        BodyLineIndex: e.CurrentLineIndex + 1,
    })
}

func (e *Evaluator) evalNextStatement(stmt *NextStatement) {
	if len(e.LoopStack) == 0 {
		return // Error: NEXT without FOR
	}
	// Peek stack
	ctx := e.LoopStack[len(e.LoopStack)-1]
    
    // Check if var matches (if provided)
    if stmt.Var != nil {
        if strings.ToLower(stmt.Var.Value) != ctx.VarName {
            // Mismatch variable. In strict BASIC this is error or it pops until match.
            // We'll ignore for now or pop?
            // Let's assume correct nesting.
        }
    }

	val := toFloat(e.Env.Store[ctx.VarName])
	val += ctx.Step
	e.Env.Store[ctx.VarName] = val

	// Check condition
    // If Step > 0: continue if val <= End
    // If Step < 0: continue if val >= End
    loop := false
    if ctx.Step > 0 {
        if val <= ctx.End + 1e-9 { // epsilon for float comparison
            loop = true
        }
    } else {
        if val >= ctx.End - 1e-9 {
            loop = true
        }
    }

	if loop {
		// Jump back to body
		e.CurrentLineIndex = ctx.BodyLineIndex
        // We modified CurrentLineIndex, so the main loop won't increment it further (due to our check logic)
        // Wait, my logic in Run() was: if changed, don't increment.
        // So we set it to body index.
        // The main loop will then execute the body.
	} else {
		// Loop finished, pop stack
		e.LoopStack = e.LoopStack[:len(e.LoopStack)-1]
        // Execution falls through to next statement (CurrentLineIndex increments naturally in Run)
	}
}

func (e *Evaluator) evalSleepStatement(stmt *SleepStatement) {
	dur := toFloat(e.evalExpression(stmt.Duration))
	time.Sleep(time.Duration(dur * float64(time.Second)))
}

type TabValue struct {
    Pos int
}

func (e *Evaluator) evalExpression(expr Expression) interface{} {
	switch n := expr.(type) {
	case *IntegerLiteral:
		return float64(n.Value)
	case *FloatLiteral:
		return n.Value
	case *StringLiteral:
		return n.Value
	case *Identifier:
		val, ok := e.Env.Store[strings.ToLower(n.Value)]
		if !ok {
			return float64(0) // Default 0
		}
		return val
    case *PrefixExpression:
        right := e.evalExpression(n.Right)
        if n.Operator == "-" {
            return -toFloat(right)
        }
        return right
	case *InfixExpression:
		left := toFloat(e.evalExpression(n.Left))
		right := toFloat(e.evalExpression(n.Right))
		switch n.Operator {
		case "+":
			return left + right
		case "-":
			return left - right
		case "*":
			return left * right
		case "/":
			return left / right
		}
    case *CallExpression:
        fn := strings.ToUpper(n.Function)
        if fn == "TAB" {
             arg := toFloat(e.evalExpression(n.Arguments[0]))
             return TabValue{Pos: int(arg)}
        }
        if fn == "SIN" {
             arg := toFloat(e.evalExpression(n.Arguments[0]))
             return math.Sin(arg)
        }
	}
	return nil
}

func toFloat(val interface{}) float64 {
	switch v := val.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	case int64:
		return float64(v)
	}
	return 0
}