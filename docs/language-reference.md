# Language reference

go-basic targets the line-numbered Microsoft 8K BASIC family exercised by the
BASIC Computer Games corpus. Keywords and variable names are
case-insensitive. This document describes the implemented language surface; it
is not a claim of compatibility with every BASIC dialect.

## Program form

Ordinary source uses an integer line number followed by one or more statements:

```basic
10 LET MESSAGE$="HELLO"
20 FOR I=1 TO 5
30 PRINT MESSAGE$;" ";I
40 NEXT I
50 END
```

Lines execute in numeric order regardless of their order in the source file.
Duplicate line numbers are errors. Separate multiple statements on one line
with `:`. `REM` ignores the remainder of its BASIC line. `LET` is optional for
assignments.

## Values, variables, and arrays

| Form | Behavior |
| --- | --- |
| `A`, `TOTAL` | Numeric scalar; unset values read as `0` |
| `A$`, `NAME$` | String scalar; unset values read as an empty string |
| `A(I)`, `A(I,J)` | Numeric array reference |
| `A$(I)` | String array reference |
| `DIM A(20),B$(4,4)` | Explicit inclusive upper bounds, starting at zero |

An array first referenced without `DIM` is created with an inclusive upper
bound of 10 in each supplied dimension. Array subscripts are truncated toward
zero after negative and non-finite values are rejected. The interpreter limits
each array to 1,000,000 elements.

Numeric and string values do not coerce silently during assignment or
comparison. `+` concatenates two strings and adds two numbers.

## Operators

| Category | Operators and behavior |
| --- | --- |
| Arithmetic | `+`, `-`, `*`, `/`, right-associative `^` |
| Comparison | `=`, `<>`, `<`, `<=`, `>`, `>=` |
| Logic | `NOT`, `AND`, `OR` over signed 16-bit integer operands |
| Grouping | Parentheses |

Comparisons return Microsoft-style truth values: `-1` for true and `0` for
false. Division by zero, non-real exponentiation, overflow, type mismatches,
and invalid logical operands are runtime errors.

## Statements

| Statement | Supported forms |
| --- | --- |
| Assignment | `LET A=1`, `A=1`, `A$(I)="X"` |
| Output | `PRINT`, comma print zones, semicolon suppression, adjacent items, `TAB`, `SPC` |
| Input | `INPUT A`, `INPUT "PROMPT";A$`, multiple scalar or array targets |
| Branching | `IF expression THEN line`, `IF expression THEN statement`, `GOTO` |
| Subroutines | `GOSUB`, `RETURN`, computed `ON expression GOSUB` |
| Computed branch | `ON expression GOTO line,...` |
| Loops | `FOR variable=start TO end [STEP step]`, `NEXT`, `NEXT I,J` |
| Arrays | `DIM` with one or more numeric or string arrays |
| Data | `DATA`, `READ`, `RESTORE` |
| Functions | `DEF FNname(parameter)=expression` |
| Termination | `END`, `STOP` |
| Comments | `REM` |
| Extension | `SLEEP seconds` |

`RESTORE` resets reading to the first `DATA` item. Line-targeted `RESTORE` is
not implemented. An `ON` selector outside the available positive target range
falls through without jumping; fractional or negative selectors are errors.

## Numeric functions

| Function | Purpose |
| --- | --- |
| `ABS(x)` | Absolute value |
| `SGN(x)` | Sign as `-1`, `0`, or `1` |
| `INT(x)` | Floor |
| `SIN(x)`, `COS(x)`, `TAN(x)` | Trigonometric functions in radians |
| `ATN(x)` | Arctangent in radians |
| `SQR(x)` | Square root |
| `LOG(x)` | Natural logarithm |
| `EXP(x)` | Natural exponential |
| `RND(x)` | Random value in `[0,1)` with classic repeat/reseed behavior |
| `POS(x)` | Current output column, counting from zero; `x` is ignored |

Supplying the CLI `-seed` option replaces the initial random source so a run
can be reproduced.

## String functions

| Function | Purpose |
| --- | --- |
| `LEFT$(text,n)` | Leftmost `n` bytes |
| `RIGHT$(text,n)` | Rightmost `n` bytes |
| `MID$(text,start[,length])` | One-based substring |
| `LEN(text)` | String length in bytes |
| `STR$(number)` | BASIC-formatted number converted to a string |
| `VAL(text)` | Numeric prefix converted to a number |
| `CHR$(number)` | Byte value `0..255` converted to a one-byte string |
| `ASC(text)` | Numeric value of the first byte |

## Input and printing

`INPUT` reads a line of comma-separated fields. It repeats the prompt after a
field-count or numeric-conversion error. A quoted prompt is followed by `? `;
without one, the prompt is simply `? `.

`PRINT` starts a new line unless its final item is followed by `;` or `,`.
Commas advance to 14-column print zones. `TAB(n)` advances when `n` is to the
right of the current output column and otherwise emits no spacing. `SPC(n)`
emits `n` spaces regardless of the current column, and `POS(x)` reports that
column as a number. Both truncate their argument toward zero; a negative `SPC`
count is a runtime error. `TAB` and `SPC` are print-list directives rather than
values, so using either outside `PRINT` is a type error.

## Annotated structured-source extension

The interpreter includes a narrow source-preparation extension for an
annotated Checkers variant in the pinned corpus. A standalone `Sub_Start`
marker opts a file into lowering for `LOOP`/`ENDLOOP`, block `IF`/`ENDIF`,
`BREAK`, `THEN BREAK`, numeric labels, `#` comments, and `==` equality.

Regular numbered BASIC does not use this lowering and remains strict. See
[how it works](how-it-works.md#2-source-preparation) for the pipeline details.

## Compatibility boundary

Unsupported or malformed syntax produces a diagnostic rather than being
silently omitted. Hardware-specific graphics, sound, memory access, and file
I/O are outside the currently implemented corpus-driven language surface.

For the exact external inventory and acceptance rules, see
[compatibility](compatibility.md).
