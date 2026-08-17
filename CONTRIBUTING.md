# Contributing to go-basic

Thank you for improving go-basic. Changes should preserve useful diagnostics,
classic BASIC compatibility, and deterministic behavior across supported
platforms.

## Set up the repository

Install Go 1.26.6 or newer, then clone and verify the checkout:

```bash
git clone https://github.com/scottdensmore/go-basic.git
cd go-basic
make tools
make check
```

The pinned tools are installed under `.tools/`; no global installation is
required. The first tool or corpus command needs network access. Subsequent
corpus runs reuse the verified `.cache/` checkout.

## Make a change

1. Start with a failing behavior test that uses `pkg/interpreter` or the real
   `go-basic` CLI.
2. Make the smallest product-code change that passes the test.
3. Refactor only while the tests remain green.
4. Update the relevant user guide, language reference, and compatibility
   documentation when behavior or supported syntax changes. Keep `README.md`
   focused on project orientation and links into those documents.
5. Run the appropriate validation gates before opening a pull request.

Keep interpreter behavior in `pkg/interpreter/`. The command under
`cmd/go-basic/` should remain a thin adapter for files, flags, standard streams,
and exit codes. Tests should verify product behavior rather than inspect source,
documentation, or workflow text.

## Validate the change

For most changes, run:

```bash
make fmt
make check
make fuzz
make build
make vuln
```

Changes that affect lexing, parsing, evaluation, input/output, randomness, or
control flow must also run:

```bash
make corpus-smoke
make corpus-playable
```

The full test suite uses the race detector and enforces at least 80% total
statement coverage. Fuzzing is deliberately time-bounded; preserve useful
regression seeds when a fuzzer finds a defect.

## Releases

Tags matching `v*` trigger `.github/workflows/release.yml`. The workflow repeats
the quality, fuzz, compatibility, and vulnerability gates, then publishes
checksummed Linux (`amd64`, `arm64`), macOS (`amd64`, `arm64`), and Windows
(`amd64`) archives.

Build and verify the same artifact set locally before tagging:

```bash
make release-check VERSION=v0.1.0
```

## Pull requests

Use a concise, imperative, sentence-case title. Keep each change focused and
include:

- the behavior and motivation;
- linked issues, when applicable;
- the exact validation commands run; and
- before/after terminal output when CLI interaction changes.

Unsupported or malformed BASIC must return an actionable diagnostic. Do not
silently skip syntax or introduce panic paths.
