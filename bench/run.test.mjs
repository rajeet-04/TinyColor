import assert from "node:assert/strict";
import { mkdtempSync, readFileSync } from "node:fs";
import { tmpdir } from "node:os";
import { join } from "node:path";
import test from "node:test";
import { execFileSync } from "node:child_process";
import { parseWorkload, percentile99, validateResults } from "./run.mjs";

test("percentile99 rounds up", () => assert.equal(percentile99([1, 2, 3, 4]), 4));
test("workload parser accepts representative requests", () => assert.equal(parseWorkload('{"id":"x","operation":"inspect","input":"red"}\n').length, 1));
test("results schema requires both implementations", () => assert.throws(() => validateResults({ samples: {} })));
test("quick benchmark writes a temporary measurement", () => {
  const output = join(mkdtempSync(join(tmpdir(), "tinycolor-bench-")), "results.json");
  execFileSync(process.execPath, ["bench/run.mjs", "--quick", "--output", output], { stdio: "inherit" });
  const results = JSON.parse(readFileSync(output, "utf8"));
  for (const implementation of Object.values(results.implementations)) {
    assert.ok(implementation.startupP99Ms >= 0);
    assert.ok(implementation.latencyP99Ms >= 0);
    assert.ok(implementation.throughputOpsPerSecond > 0);
  }
});
