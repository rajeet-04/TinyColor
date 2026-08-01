import assert from "node:assert/strict";
import { createHash } from "node:crypto";
import { execFileSync, spawn } from "node:child_process";
import { mkdirSync, readFileSync, realpathSync, writeFileSync } from "node:fs";
import { dirname, resolve } from "node:path";
import { cpus } from "node:os";
import readline from "node:readline";

const root = resolve(import.meta.dirname, "..");
const committedOutput = resolve(root, "bench/results.json");
const runnerEnv = { ...process.env, GOCACHE: process.env.GOCACHE ?? resolve(root, ".cache", "go-build") };

export const percentile99 = (values) => {
  assert.ok(values.length > 0, "percentile requires at least one sample");
  const sorted = [...values].sort((a, b) => a - b);
  return sorted[Math.ceil(sorted.length * 0.99) - 1];
};

export const parseWorkload = (text) => text.split(/\r?\n/).filter(Boolean).map((line, index) => {
  const request = JSON.parse(line);
  if (!request.id || !request.operation) throw new Error(`workload line ${index + 1} requires id and operation`);
  return request;
});

export const validateResults = (results) => {
  assert.equal(typeof results.generatedAt, "string");
  assert.ok(results.samples.coldStarts >= 20);
  assert.ok(results.samples.requests >= 1000);
  for (const name of ["javascript", "go"]) {
    const metrics = results.implementations?.[name];
    assert.ok(metrics, `missing ${name} metrics`);
    for (const key of ["startupP99Ms", "latencyP99Ms", "throughputOpsPerSecond"]) {
      assert.ok(Number.isFinite(metrics[key]) && metrics[key] >= 0, `${name}.${key} must be nonnegative`);
    }
    assert.ok(metrics.peakRssBytes === null || (Number.isFinite(metrics.peakRssBytes) && metrics.peakRssBytes >= 0));
  }
  return true;
};

const nsToMs = (nanoseconds) => Number(nanoseconds) / 1e6;
const commandVersion = (command, args) => execFileSync(command, args, { cwd: root, encoding: "utf8" }).trim();
const runnerCommands = () => ({
  javascript: [process.execPath, ["compat/js-runner.mjs"]],
  go: [resolve(root, "bin", process.platform === "win32" ? "tinycolor.exe" : "tinycolor"), []],
});

class Runner {
  constructor(command, args) {
    this.child = spawn(command, args, { cwd: root, env: runnerEnv, stdio: ["pipe", "pipe", "pipe"] });
    this.pending = [];
    this.stderr = "";
    readline.createInterface({ input: this.child.stdout, crlfDelay: Infinity }).on("line", (line) => {
      const next = this.pending.shift();
      if (next) next.resolve(JSON.parse(line));
    });
    this.child.stderr.on("data", (chunk) => { this.stderr += chunk; });
  }
  request(value) {
    return new Promise((resolveRequest, reject) => {
      if (this.child.exitCode !== null) return reject(new Error(`runner exited: ${this.stderr}`));
      this.pending.push({ resolve: resolveRequest, reject });
      this.child.stdin.write(`${JSON.stringify(value)}\n`, (error) => error && reject(error));
    });
  }
  async close() {
    if (this.child.exitCode === null) this.child.kill();
  }
}

const readRss = (pid) => {
  try {
    if (process.platform === "linux") {
      const status = readFileSync(`/proc/${pid}/status`, "utf8");
      const match = /^VmRSS:\s+(\d+)\s+kB$/m.exec(status);
      return match ? Number(match[1]) * 1024 : null;
    }
    if (process.platform === "win32") {
      return Number(execFileSync("powershell.exe", ["-NoProfile", "-NonInteractive", "-Command", `(Get-Process -Id ${pid}).WorkingSet64`], { encoding: "utf8" }).trim());
    }
  } catch { /* process may have exited before the sample */ }
  return null;
};

const measure = async (name, command, args, workload, coldStarts, requests) => {
  const startup = [];
  let peakRss = 0;
  let rssSupported = false;
  for (let i = 0; i < coldStarts; i++) {
    const started = process.hrtime.bigint();
    const runner = new Runner(command, args);
    await runner.request(workload[i % workload.length]);
    startup.push(nsToMs(process.hrtime.bigint() - started));
    const rss = readRss(runner.child.pid);
    if (rss !== null) { rssSupported = true; peakRss = Math.max(peakRss, rss); }
    await runner.close();
  }
  const runner = new Runner(command, args);
  const latencies = [];
  const throughputStart = process.hrtime.bigint();
  for (let i = 0; i < requests; i++) {
    const started = process.hrtime.bigint();
    await runner.request(workload[i % workload.length]);
    latencies.push(nsToMs(process.hrtime.bigint() - started));
  }
  const totalSeconds = Number(process.hrtime.bigint() - throughputStart) / 1e9;
  // RSS sampling is intentionally outside the timed request loop. On Windows,
  // starting PowerShell for each sample would measure the sampler, not the runner.
  for (let i = 0; i < 10; i++) {
    const rss = readRss(runner.child.pid);
    if (rss !== null) { rssSupported = true; peakRss = Math.max(peakRss, rss); }
    await new Promise((resolveTick) => setImmediate(resolveTick));
  }
  await runner.close();
  return {
    startupP99Ms: percentile99(startup),
    latencyP99Ms: percentile99(latencies),
    throughputOpsPerSecond: requests / totalSeconds,
    peakRssBytes: rssSupported ? peakRss : null,
    ...(rssSupported ? {} : { rssLimitation: `RSS sampling is unsupported on ${process.platform}` }),
  };
};

const main = async () => {
  const options = process.argv.slice(2);
  const quick = options.includes("--quick");
  const outputIndex = options.indexOf("--output");
  if (outputIndex !== -1 && !options[outputIndex + 1]) throw new Error("usage: node bench/run.mjs [--quick] [--output <path>]");
  const output = outputIndex === -1 ? committedOutput : resolve(root, options[outputIndex + 1]);
  if (quick && output === committedOutput) throw new Error("--quick must not write bench/results.json");
  const coldStarts = quick ? 3 : 20;
  const requests = quick ? 30 : 1000;
  const workloadText = readFileSync(resolve(root, "bench/workload.jsonl"), "utf8");
  const workload = parseWorkload(workloadText);
  execFileSync("go", ["-C", "src", "build", "-o", `../bin/${process.platform === "win32" ? "tinycolor.exe" : "tinycolor"}`, "./cmd/tinycolor-compat"], { cwd: root, env: runnerEnv, stdio: "inherit" });
  const implementations = {};
  for (const [name, [command, args]] of Object.entries(runnerCommands())) implementations[name] = await measure(name, command, args, workload, coldStarts, requests);
  const results = {
    generatedAt: new Date().toISOString(), os: process.platform, architecture: process.arch,
    cpu: cpus()[0]?.model ?? "unknown", nodeVersion: process.version,
    goVersion: commandVersion("go", ["version"]), workloadSha256: createHash("sha256").update(workloadText).digest("hex"),
    samples: { coldStarts, requests }, implementations,
  };
  if (!quick) validateResults(results);
  mkdirSync(dirname(output), { recursive: true });
  writeFileSync(output, `${JSON.stringify(results, null, 2)}\n`);
  console.log(`wrote ${output}`);
};

if (process.argv[1] && realpathSync(process.argv[1]) === realpathSync(new URL(import.meta.url))) main().catch((error) => { console.error(error.message); process.exitCode = 1; });
