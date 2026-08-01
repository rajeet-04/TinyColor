# 05-02 Summary

- Created `bench/run.test.mjs` for Node built-in tests for percentile calculation, workload parsing, results schema validation, and `--quick` mode.
- Created `bench/workload.jsonl` from representative parse, conversion, modifier, palette, and readability requests.
- Implemented `bench/run.mjs` to measure cold starts, persistent requests latency and throughput, and peak RSS.
- Updated `bench/methodology.md` to reflect same-host observations and the exact calculation methods used.
- Generated `bench/results.json` containing actual benchmark numbers.
- Updated `Makefile` to include the correct command for the `bench:` target.
