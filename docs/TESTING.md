# Testing Strategy

## Test layers

| Layer | Purpose | Command when implemented |
|---|---|---|
| Go unit tests | Local conversion and parser invariants | `go test ./...` from `src/` |
| Differential corpus | Exact observable parity with local `mod.js` | `node compat/run.mjs <jsonl>` |
| Regression corpus | Permanently reproduce every found mismatch | `node compat/run.mjs compat/cases/regression.jsonl` |
| Seeded fuzz cases | Find coercion, bounds, and rounding gaps | documented seed command |
| Source suite | Validate original source remains runnable | `deno task test` when Deno exists |
| Static checks | Formatting and common Go defects | `gofmt`, `go vet ./...` |

## Comparison rules

- Compare strings byte-for-byte.
- Compare booleans, arrays, formats, and errors exactly.
- Serialize non-finite values explicitly if the source produces them; JSON's
  default inability to represent them must not silently erase a difference.
- Use a narrowly documented normalization only for JavaScript/Go JSON number
  serialization—not for color math. Never use a global epsilon.
- Seed random generation and print the seed on failure.

## Delivery evidence

Run the complete local gate from the repository root with `make verify`.
The differential harness is `node compat/run.mjs compat/cases/operations.jsonl`;
there is no `npm run test:compat` script in this checkout. Record a 60-second
fuzz session with `node fuzz/harness.mjs --duration 60 --seed 20260801` and
validate its saved output with `node fuzz/validate-log.mjs fuzz/log.txt`.

Create the same-host benchmark evidence with
`node bench/run.mjs --output bench/results.json`. First smoke-test it without
changing committed evidence using
`node bench/run.mjs --quick --output $env:TEMP\\tinycolor-bench-quick.json`.

## Corpus priorities

1. Inputs copied from `test.js` in parser and output groups.
2. Boundaries: `-1`, `0`, `1`, `1.0`, `100%`, `255`, overflow, blank, malformed,
   alpha `0`, and hue values below/above 360.
3. State transitions: clone, set alpha, all instance modifiers, default versus
   explicit zero amount.
4. Structured random RGB/RGBA/HSL/HSV inputs generated from a fixed seed.
