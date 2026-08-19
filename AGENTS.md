# Repository Guidelines

## Project Structure & Module Organization

`pkg/interpreter/` contains the product code: token definitions, lexer, AST, parser, and evaluator. Keep language behavior there rather than in command or test helpers. `cmd/go-basic/` is the thin CLI entry point; `cmd/corpus-*` and `internal/corpus/` implement the pinned external-corpus test tier. Unit tests live beside their implementation as `*_test.go`; `test/` contains black-box CLI tests and BASIC fixtures under `test/scripts/`. `README.md` documents user-facing behavior and the Microsoft 8K/6502 BASIC compatibility goal.

## Build, Test, and Development Commands

- `make fmt`: format all Go packages before committing.
- `go test ./pkg/interpreter`: run the fast lexer, parser, and evaluator unit tests.
- `make test`: run all tests with race detection and write `coverage.out`.
- `make check`: enforce formatting, vet, 80% coverage, and golangci-lint.
- `make fuzz`: fuzz the lexer and parser for ten seconds each.
- `make build`: build the CLI as `bin/go-basic`.
- `make corpus-smoke`: fetch and run all 112 byte-distinct pinned corpus variants.
- `make corpus-playable`: run the complete deterministic CLI gameplay suite.

Pushing a `v*` tag triggers `.github/workflows/release.yml`, which cross-compiles Linux, Windows, and macOS release binaries.

## Coding Style & Naming Conventions

Follow standard Go conventions and let `gofmt` control indentation. Use short, lowercase package names; exported identifiers use `PascalCase`, internal identifiers use `camelCase`, and files use descriptive lowercase names such as `evaluator_test.go`. Prefer small parser/evaluator methods with explicit error returns. Unsupported or malformed BASIC must produce useful diagnostics—never silent omissions or panics.

## Testing Guidelines

Use Go's `testing` package. Start each change with a failing behavior test, make the smallest product change that passes, and refactor only while the suite remains green. Name tests `TestBehavior` and favor table-driven cases for tokens, precedence, and statement variants. Tests should exercise public product behavior or the CLI, not parse documentation or workflow files. Inject input/output and other nondeterministic dependencies so interactive game transcripts remain repeatable. Add a regression test with every bug fix.

## Commit & Pull Request Guidelines

Recent commits use concise, imperative, sentence-case subjects, for example `Add cross-platform build and release workflow`. Keep each commit focused. Pull requests should explain behavior changes, link relevant issues, list commands run, and include before/after terminal transcripts when output or interaction changes. Screenshots are unnecessary for this command-line project.

## Gotchas & Troubleshooting

- **`make check` is not the CI gate.** `check` is `fmt-check vet coverage-check lint`, which is what `release.yml` runs. `ci.yml` runs `corpus-smoke`, `fmt-check`, `vet`, `coverage-check`, `fuzz`, `build`, and `release-check VERSION=ci` in one job, plus `golangci-lint` and `make vuln` in separate jobs. A green `make check` does not imply a green pull request — see the **Verification Map** for the full list.
- **First tool and corpus commands need network.** `make lint`, `make vuln`, and `make tools` `go install` pinned binaries into `.tools/bin`; `make corpus-fetch` (a prerequisite of `corpus-smoke`) downloads the pinned tarball into `.cache/`. Both are cached afterward. `make clean` deletes `.tools`, so the next lint or vuln run needs network again.
- **The corpus cache is keyed by commit.** `internal/corpus.Fetch` refuses a cache directory holding a different SHA (`cache <dir> contains commit X, want Y`) and refuses a target directory that exists without a `.corpus-commit` marker. When `CORPUS_COMMIT` changes or a fetch is interrupted, delete the stale directory under `.cache/basic-computer-games/` rather than re-running the fetch.
- **`make coverage-check` is not a cheap check.** It runs `make test` first — the whole suite under `-race` — then fails if total statement coverage is below `COVERAGE_MIN` (80%). Budget for it accordingly.
- **`test/` is not a unit tier.** Its tests build the real `go-basic` binary and execute `test/scripts/*.bas` fixtures through it, so a change to `cmd/go-basic/` or `pkg/interpreter/` can break them without any test file changing.
- **Fuzz seeds are committed.** Corpus entries under `pkg/interpreter/testdata/fuzz/` are regression seeds. Preserve them when a fuzzer finds a defect; do not delete them to make `make fuzz` quiet.

## Verification Map

The complete gate is what CI runs, in this order:

```bash
make corpus-smoke && make fmt-check && make vet && make coverage-check \
  && make fuzz && make build && make release-check VERSION=ci \
  && make lint && make vuln
```

Use this table in stage 7 to rerun only what a fix could have invalidated.

| A fix touches | Rerun |
|---|---|
| `pkg/interpreter/**` | `make fmt-check`, `make vet`, `make coverage-check`, `make lint`, `make fuzz`, `make corpus-smoke`, `make build` |
| `cmd/go-basic/**` | `make fmt-check`, `make vet`, `make coverage-check`, `make lint`, `make build`, `make release-check VERSION=ci` |
| `internal/corpus/**` or `cmd/corpus-*/**` | `make fmt-check`, `make vet`, `make coverage-check`, `make lint`, `make corpus-smoke` |
| `test/**` (`*_test.go` or `scripts/*.bas`) | `make fmt-check`, `make vet`, `make coverage-check`, `make lint` |
| `.golangci.yml` | `make lint` |
| `go.mod` | the complete gate, and `make vuln` in particular |
| `Makefile` | the complete gate — it defines every gate command |
| `.github/workflows/**` | nothing runs, but re-read the workflow: it defines the gate above |
| `*.md`, `docs/**`, `LICENSE` | nothing. `go test ./...`, `go vet ./...`, and `golangci-lint run ./...` collect Go packages only; `make fmt-check` globs `-name '*.go'`; and no Go file in the repository opens a Markdown, YAML, or Makefile path |
| anything else | the complete gate |

<!-- agent-skills:begin workflow 8881bee7 — managed block, edits here are overwritten -->
## Development Workflow

Follow these stages in order (governed by the global `agent-workflow-skills`). Scale the pipeline to the
size of the change using the triage table — skipping a stage is a decision to
state out loud, never a shortcut taken silently.

| Track | When | Stages |
|---|---|---|
| **Trivial** | Docs, comments, typos, config with no logic change | 1 → 7 → 9 |
| **Single fix** | One bug or small change with a clear, contained cause | 1 → 2 → 5 → 7 → 8 → 9 |
| **Feature** | New behavior, several files, or an architectural choice | All stages; repeat 5–8 per slice |

**Division of labor.** The main agent runs only focused checks — the single test
it just wrote, a formatter over the files it just touched. Whole suites, builds,
dependency audits, and repository-wide lint belong to `verifier`, and reviews
belong to `ui-review` and `code-review`. This is not ceremony: it keeps routine
command output out of the implementation context, and it means each gate is read
by something that has not already convinced itself the change is correct.
Sub-agents report successes in one line and include only the evidence needed to
diagnose a failure.

**Preserve what you did not change.** A worktree may hold work that is not yours.
Never stage, revert, or "clean up" a change you did not make; when something
unrelated is in the way, name it and leave it alone.

1. **Inspect & Branch**: Inspect `git status`, the current branch, and every
   applicable instruction file before touching anything. Note unrelated staged,
   unstaged, and untracked work so you can preserve it. Fetch the base branch
   (`git fetch origin main`) and create a dedicated branch:
   `git checkout -b <owner>/<type>/<short-description> origin/main`.
   `<owner>` is your GitHub login (`gh api user --jq .login`); `<type>` is one of
   `feat`, `fix`, `refactor`, `chore`, `test`, `docs`. Never commit to `main`.
2. **Plan & Slice (`plan-and-prototype`)**: Formulate a clear step-by-step plan
   before writing code. Define the smallest end-to-end slice that can be reviewed,
   tested, and shipped independently; if the work is too large for one pull
   request, order the slices and complete only the current one.
3. **Prototype Options (if needed)**: When facing architectural choices, unfamiliar
   APIs, or UX alternatives, spike lightweight prototypes and compare trade-offs
   before committing to an approach.
4. **Track Bugs & Follow-ups**: When bugs, edge cases, technical debt, or follow-up
   tasks surface mid-change, file them immediately (`gh issue create`, the project's
   tracker, or `ISSUES.md` when none is configured) instead of expanding the current
   slice.
5. **Test-Driven Development (`tdd-workflow`)**:
   - Write/update a focused test first → confirm it fails for the expected reason →
     minimal implementation → iterate until passing → refactor. A test that passes
     before the code exists is testing the wrong thing.
   - **When the change replaces an existing contract, find the tests pinning the old
     one first.** A new failing test proves the new behavior is missing; it says
     nothing about tests still asserting the behavior being removed. Search for
     assertions on the symbol, attribute, label, or role being changed and update
     them inside the same red/green loop. Skipping this is silently safe — the new
     test goes green, the loop looks complete, and the contradiction only surfaces a
     full gate cycle later.
   - Run only the test you authored or changed, filtered by file and name. Whole
     suites are stage 7's job.
   - Pure logic (calculations, state machines, business rules) must be unit-tested.
     Non-testable areas (rendering, audio) must be visually/interactively verified.
6. **UI Review (`ui-review`)**:
   - Audit layout, visual hierarchy, contrast (WCAG AA), interaction states, and
     accessibility according to the project's UI domain.
   - For a change with no user-visible surface, say so and return. Do not invent
     findings to justify the stage.
7. **Verification (`verifier`)**:
   - Run the project's full gate: lint, type-check, test suites, build. Focused runs
     from stage 5 do not substitute for it.
   - Fix or explicitly resolve every actionable finding before code review. When a
     fix changes code, rerun the affected focused tests, then rerun the gate commands
     whose inputs the fix touched — see **Verification Map** below if this project
     defines one. The complete gate must run in full at least once on the state that
     enters code review.
   - Some findings are environmental and no code change resolves them (browsers that
     will not install, no network, a missing credential). Resolving those means
     naming them precisely — what ran, what did not, and why — not retrying them.
8. **Code Review (`code-review`)**:
   - Read the complete change: `git diff origin/main...HEAD`, plus staged
     and unstaged edits (`git diff HEAD`) and untracked files (`git status
     --porcelain`). Remove accidental or unrelated edits of your own; preserve
     anything that belongs to the user.
   - Enforce architectural boundaries, language idioms, defensive error handling,
     and zero committed secrets.
   - Do not repeat this review on an unchanged state. Rerun it only when the
     reviewed content actually changed.
9. **Commit & PR Lifecycle (`slice-and-pr`)**:
   - Commit using Conventional Commits (`<type>(<scope>): <summary>`). Stage files
     explicitly; never `git add -A` when unrelated work is present.
   - Open the PR with `gh pr create` and watch CI with `gh pr checks --watch`.
   - **Stop there and report.** Merging (`gh pr merge`) and force-pushing require
     explicit approval from the user in the current conversation.
<!-- agent-skills:end workflow -->
