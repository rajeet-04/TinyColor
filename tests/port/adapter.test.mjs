import assert from "node:assert/strict";
import { spawnSync } from "node:child_process";

const run = (input) => {
  const result = spawnSync(process.execPath, ["compat/js-runner.mjs"], {
    cwd: process.cwd(),
    input,
    encoding: "utf8",
  });
  return result.stdout.trim().split("\n").filter(Boolean).map(JSON.parse);
};

const [success] = run('{"id":"red","operation":"inspect","input":"red"}\n');
assert.deepEqual(success, {
  id: "red",
  result: {
    valid: true,
    format: "name",
    alpha: 1,
    rgb: { r: 255, g: 0, b: 0, a: 1 },
    value: "red",
    original: "red",
  },
});

const [malformed, unknown] = run('{bad json\n{"id":"unknown","operation":"unknown"}\n');
for (const response of [malformed, unknown]) {
  assert.notEqual(Object.hasOwn(response, "result"), Object.hasOwn(response, "error"));
  assert.equal(typeof response.error, "string");
}
