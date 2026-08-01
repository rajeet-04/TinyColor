# Decisions

| ID | Decision | Why | Status |
|---|---|---|---|
| D-001 | Port the checked-out `mod.js` behavior to Go. | It is the exact source and test target selected for the project. | Accepted |
| D-002 | Keep JavaScript oracle files immutable. | Differential evidence is credible only when the reference does not move with the port. | Accepted |
| D-003 | Put the Go module in `src/`. | It keeps the port idiomatic while following the Track H submission layout; the JavaScript oracle remains at the root. | Accepted |
| D-004 | Use JSON Lines adapters over stdin/stdout. | Node and Go can compare dynamic inputs and operations without modifying source tests. | Accepted |
| D-005 | Use Go standard library first. | The port's algorithms are small and self-contained; dependencies add compatibility and supply-chain surface. | Accepted |
| D-006 | Preserve source quirks at the adapter boundary; expose an idiomatic Go API separately. | Compatibility and Go usability are both required, but they should not leak dynamic JavaScript semantics into every caller. | Accepted |
| D-007 | Run the immutable Deno source suite explicitly as `deno test test.js`. | Deno 2 broad discovery includes generated npm and Node port tests outside the pinned source suite. | Accepted |
| D-008 | Use a finite V8-derived luminance transfer table. | TinyColor rounds RGB to 8-bit channels before luminance, while Go `math.Pow` differs from V8 for 93 nonlinear values; the table preserves exact parity without tolerance. | Accepted |
