# Benchmark methodology

Measure the compiled compatibility CLI on a quiet host with the same corpus for
each run. Record median and p99 latency, RSS, startup time, and throughput.
Do not compare results across different hosts or Go versions.

These are same-host observations and not universal speedup claims.
Cold starts measure process creation, engine initialization, and one request/response.
Persistent requests measure latency and throughput over a single long-lived process.
p99 is calculated by sorting the samples and selecting `ceil(0.99*n)-1`.
Peak RSS is sampled using `/proc/<pid>/status` on Linux and `Get-Process -Id <pid>` on Windows. If unsupported, it emits null along with a limitation string.
