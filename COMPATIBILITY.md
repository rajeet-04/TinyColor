# Compatibility evidence

Target: `bgrins/TinyColor` commit
`b49018c9f2dbca313d80d7a4dad25e26143cfe01`, pinned by
`tests/original/manifest.sha256`.

## Fixed differential corpora

Fresh command: `node compat/run.mjs <corpus>`.

| Corpus | Passed/total | Mismatches |
|---|---:|---:|
| `compat/cases/smoke.jsonl` | 9/9 | 0 |
| `compat/cases/parser-hex-rgb.jsonl` | 26/26 | 0 |
| `compat/cases/parser.jsonl` | 23/23 | 0 |
| `compat/cases/conversion.jsonl` | 35/35 | 0 |
| `compat/cases/operations.jsonl` | 71/71 | 0 |
| **Total** | **164/164** | **0** |

Random-color behavior is checked through validity, alpha, and channel-range
invariants rather than exact equality between independent random generators.

## Source suite and oracle integrity

- `node tests/original/verify.mjs`: 3/3 kickoff hashes verified.
- `deno test test.js`: 45 passed, 0 failed, 1 ignored.
- The ignored `polyad` test is also ignored by the pinned upstream suite.

## Go coverage and safety

Fresh command: `go -C src test -cover ./...`.

| Package | Statement coverage |
|---|---:|
| `cmd/tinycolor-compat` | 33.6% |
| `internal/color` | 90.5% |
| `internal/compat` | 57.1% |
| `internal/parser` | 85.8% |
| `tinycolor` | 92.8% |

These are Go package coverage figures, not a JavaScript coverage comparison.
`rg -n '\bunsafe\b' src -g '*.go'` reports **0 Go source occurrences**.

## Differential fuzzing

[`fuzz/log.txt`](fuzz/log.txt) records 60.012 seconds, seed 20260801,
1,091,630 cases, and zero divergences. Validate it with
`node fuzz/validate-log.mjs fuzz/log.txt`. This supports the Differential Fuzz
Survivor claim; exact comparison remains enabled.

## Same-host benchmark

The committed [results](bench/results.json) were measured on Windows x64 with
20 cold starts and 1,000 persistent requests per implementation:

| Implementation | Startup p99 | Latency p99 | Throughput | Peak RSS |
|---|---:|---:|---:|---:|
| JavaScript | 41.4131 ms | 0.2795 ms | 7,834.77 ops/s | 42,868,736 bytes |
| Go | 11.0109 ms | 0.1808 ms | 16,756.79 ops/s | 12,668,928 bytes |

See [`bench/methodology.md`](bench/methodology.md). These are same-host
observations, not universal speedup claims.

## Bonus and external evidence status

- Differential Fuzz Survivor: eligible from the validated 60-second log.
- Zero Unsafe: eligible from zero Go source occurrences.
- Decision Log: eligible from 15 substantive decisions.
- GitHub Actions: [run 30686980719](https://github.com/rajeet-04/TinyColor/actions/runs/30686980719) passed the full gate, Deno, build, and artifact upload for commit `1c218b6`.
- Public repository visibility: unauthenticated `git ls-remote https://github.com/rajeet-04/TinyColor.git refs/heads/rajeet` returned exact commit `1c218b6`.
- Five-minute demo video: not supplied.

No known mismatch remains in the fixed corpus or recorded fuzz session.
