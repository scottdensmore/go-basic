# Security Policy

## Supported Versions

Security fixes target the `main` branch and, when practical, the latest published
release. Older releases are not supported. Users should reproduce a suspected
issue against the latest release or `main` before reporting it.

## Reporting a Vulnerability

Report suspected vulnerabilities through
[GitHub private vulnerability reporting](https://github.com/scottdensmore/go-basic/security/advisories/new).
Do not disclose security-sensitive details in a public issue.

Include the affected version or commit, operating system and architecture,
realistic impact, reproduction steps or a minimal BASIC program, and any known
mitigations. Reports should avoid unrelated personal data or secrets.

The maintainer will aim to acknowledge a report within seven days, confirm
whether it is accepted, and coordinate remediation and disclosure when
appropriate.

## System and Scope

This policy covers:

- the `go-basic` command-line application;
- the lexer, parser, structured-source lowering, and evaluator under
  `pkg/interpreter`;
- pinned corpus acquisition and execution tooling;
- GitHub Actions workflows and published release artifacts.

go-basic is a local command-line program, not a network service. It reads a
source file selected by the user, may read program input from standard input,
and writes program output and diagnostics to the provided streams.

## Threat Model and Trust Boundaries

BASIC source, program input, corpus archives, archive paths, and command-line
arguments are potentially attacker-controlled. Important assets include the
host filesystem and processes, developer and CI credentials, build integrity,
and published release artifacts.

Running a BASIC program intentionally grants it CPU time, memory, and output
within the go-basic process. It must not grant access to unrelated host
capabilities.

## Security Invariants

The following properties must hold:

- BASIC programs cannot read or modify arbitrary host files, start processes,
  or initiate network connections.
- Malformed source and input return actionable errors rather than panicking or
  corrupting interpreter state.
- Array dimensions and other attacker-controlled allocations remain validated
  and bounded.
- Corpus extraction cannot escape its target directory, follow archive
  symlinks, or overwrite existing files.
- Corpus and release workflows use pinned inputs, least-privilege permissions,
  and immutable GitHub Action revisions.
- Release artifacts are produced only after the repository verification gates
  pass.

## Reportable Findings and Severity Context

Reportable findings include:

- arbitrary file access, process execution, network access, or directory
  traversal caused by BASIC input or a corpus archive;
- a panic or disproportionate resource exhaustion caused by a small,
  well-formed input that bypasses documented limits;
- a CI or release-integrity weakness that creates a realistic path to modifying
  published artifacts or exposing credentials; and
- validation failures that cross the documented interpreter or tooling trust
  boundaries.

Severity depends on realistic reachability and impact. An issue requiring a user
to explicitly execute a malicious local program is generally less severe than a
repository or release compromise affecting downstream users.

## Out of Scope

The following are not security findings without additional security impact:

- BASIC dialect incompatibilities or incorrect game behavior;
- upstream corpus defects;
- intentionally non-terminating programs run without `-max-statements`;
- large or continuous output explicitly produced by the BASIC program; and
- findings that require an already-compromised host, GitHub account, or build
  runner without crossing another boundary.

## Known Limitations and Compensating Controls

go-basic is an interpreter, not a security sandbox. Run untrusted programs
inside an operating-system sandbox when stronger isolation is required.

Execution is unlimited unless the caller supplies `-max-statements`, and program
output is not capped. The corpus acceptance tools apply deterministic statement
and wall-clock bounds.

The corpus downloader selects an immutable upstream commit over HTTPS and
restricts extraction to expected regular BASIC files. The commit pin stabilizes
source selection but is not an independent cryptographic signature over the
downloaded archive.
