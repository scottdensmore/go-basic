# BASIC Computer Games compatibility

go-basic uses the original public-domain programs from
[`coding-horror/basic-computer-games`](https://github.com/coding-horror/basic-computer-games)
as an external acceptance suite. The input is pinned at commit
`5301155192d91d74d337899cecc59dbda59c4c17`, so compatibility results are
repeatable and cannot drift with the upstream default branch.

## Audited inventory

The pinned repository contains 96 games, 105 main-tree BASIC sources, and 210
paths when mirrored alternate-language copies are included. Byte comparison
reduces those paths to 112 distinct sources: all 105 main-tree variants plus 7
byte-different alternates. Every distinct source remains independently tested.

The repository keeps corresponding fixtures in `test/scripts/`. Besides the
112 corpus fixtures, that directory contains four small local programs used for
general CLI and numeric-output tests.

## Acceptance tiers

### Pinned smoke tier

```bash
make corpus-smoke
```

This command:

1. downloads only original `.bas` sources from the pinned archive;
2. verifies a commit marker before reusing the ignored local cache;
3. discovers main sources and only byte-different alternates;
4. parses and executes every variant with deterministic randomness, no sleeping,
   a statement limit, and a wall-clock limit; and
5. reports status, last BASIC line, duration, and any error per variant.

An accepted smoke result reaches normal completion, an expected input boundary,
or the configured statement bound. Parse failures, unexpected runtime failures,
panics, and timeouts fail the tier.

### Deterministic gameplay tier

```bash
make corpus-playable
```

This black-box suite builds the real CLI and runs all corpus fixtures one by
one. Inputs, random seeds, statement limits, and expected termination are
controlled for repeatability. Assertions cover exact transcripts where stable
and complete gameplay milestones where a full transcript would obscure the
behavior under test. This tier exercises the product executable rather than a
test-only interpreter path.

## Interpreting bounded runs

Some original programs intentionally repeat forever or rely on Control-C.
Their deterministic scenario verifies meaningful behavior before stopping at a
statement limit. Other programs have no quit command and intentionally end at a
well-defined next-input boundary. These are accepted only when their expected
transcript or gameplay milestones have already been observed.

The byte-different alternate Life for Two source jumps to a missing line 800
after completing and declaring a draw. The main source exits normally; the
alternate is accepted only after the complete draw behavior, then bounded
before that upstream defect. It is tracked in
[`go-basic#67`](https://github.com/scottdensmore/go-basic/issues/67) and
[`basic-computer-games#957`](https://github.com/coding-horror/basic-computer-games/issues/957).

## CI contract

Pull requests and pushes to `main` run the race-enabled unit and CLI suite,
coverage threshold, formatting, vet, lint, fuzzing, release builds, vulnerability
scan, and pinned smoke tier. A release tag repeats the quality, fuzz, corpus, and
vulnerability gates before publishing artifacts.

Passing the corpus is strong evidence for the supported language surface, but
it is not a claim that every Microsoft BASIC dialect or hardware-specific
feature is implemented. Supported behavior is documented in the
[language reference](language-reference.md) and enforced by the product and
compatibility tests.
