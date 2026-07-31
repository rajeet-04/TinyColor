# Decisions

| ID | Decision | Why | Status |
|---|---|---|---|
| D-001 | Port the checked-out `mod.js` behavior to Go. | It is the exact source and test target selected for the project. | Accepted |
| D-002 | Keep JavaScript oracle files immutable. | Differential evidence is credible only when the reference does not move with the port. | Accepted |
| D-003 | Put the Go module in `src/`. | It keeps the port idiomatic while following the Track H submission layout; the JavaScript oracle remains at the root. | Accepted |
| D-004 | Use JSON Lines adapters over stdin/stdout. | Node and Go can compare dynamic inputs and operations without modifying source tests. | Accepted |
| D-005 | Use Go standard library first. | The port's algorithms are small and self-contained; dependencies add compatibility and supply-chain surface. | Accepted |
| D-006 | Preserve source quirks at the adapter boundary; expose an idiomatic Go API separately. | Compatibility and Go usability are both required, but they should not leak dynamic JavaScript semantics into every caller. | Accepted |
| D-007 | Deno's source suite is a required later check, not a current pass. | `deno` is not installed in the checked environment as of 2026-08-01. | Accepted |
