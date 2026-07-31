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

const [lightenDefault, lightenZero, spinOmitted, spinNull, mixed, badModifier] = run([
  '{"id":"lighten-default","operation":"modify","input":"red","args":{"method":"lighten"}}',
  '{"id":"lighten-zero","operation":"modify","input":"red","args":{"method":"lighten","amount":0}}',
  '{"id":"spin-omitted","operation":"modify","input":"red","args":{"method":"spin"}}',
  '{"id":"spin-null","operation":"modify","input":"red","args":{"method":"spin","amount":null}}',
  '{"id":"mix-default","operation":"mix","input":"red","args":{"other":"#000"}}',
  '{"id":"bad-modifier","operation":"modify","input":"red","args":{"method":"unknown"}}',
].join("\n"));

assert.deepEqual(lightenDefault, {
  id: "lighten-default",
  result: {
    before: success.result,
    after: {
      ...success.result,
      rgb: { r: 255, g: 51, b: 51, a: 1 },
      value: "#ff3333",
    },
    sameReceiver: true,
  },
});
assert.deepEqual(lightenZero.result, {
  before: success.result,
  after: success.result,
  sameReceiver: true,
});
assert.equal(spinOmitted.result.after.value, "black");
assert.equal(spinOmitted.result.after.original, "red");
assert.deepEqual(spinNull.result, lightenZero.result);
assert.deepEqual(mixed, {
  id: "mix-default",
  result: {
    valid: true,
    format: "rgb",
    alpha: 1,
    rgb: { r: 128, g: 0, b: 0, a: 1 },
    value: "rgb(128, 0, 0)",
    original: { r: 127.5, g: 0, b: 0, a: 1 },
  },
});
assert.deepEqual(badModifier, { id: "bad-modifier", error: "unsupported method" });

const [readability, readableDefault, readableMixedCase, readableInvalid, fallback, fallbackDisabled, emptyCandidates] = run([
  '{"id":"readability","operation":"readability","input":"#000","args":{"other":"#fff"}}',
  '{"id":"readable-default","operation":"isReadable","input":"#777","args":{"other":"#000","options":{}}}',
  '{"id":"readable-mixed","operation":"isReadable","input":"#000","args":{"other":"#fff","options":{"level":"aaa","size":"LARGE"}}}',
  '{"id":"readable-invalid","operation":"isReadable","input":"#777","args":{"other":"#000","options":{"level":false,"size":0}}}',
  '{"id":"fallback","operation":"mostReadable","input":"#777","args":{"candidates":["#777"],"options":{"includeFallbackColors":true}}}',
  '{"id":"fallback-disabled","operation":"mostReadable","input":"#777","args":{"candidates":["#777"],"options":{"includeFallbackColors":false}}}',
  '{"id":"empty-candidates","operation":"mostReadable","input":"#fff","args":{"candidates":[],"options":{"includeFallbackColors":true}}}',
].join("\n"));

assert.equal(readability.result, 21);
assert.equal(readableDefault.result, true);
assert.equal(readableMixedCase.result, true);
assert.equal(readableInvalid.result, true);
assert.equal(fallback.result.value, "#000000");
assert.equal(fallbackDisabled.result.value, "#777777");
assert.equal(emptyCandidates.result, null);

const [complement, splitComplement, triad, tetrad, analogous, analogousZero, monochromatic, monochromaticZero, badPalette] = run([
  '{"id":"complement","operation":"palette","input":"red","args":{"method":"complement"}}',
  '{"id":"split-complement","operation":"palette","input":"red","args":{"method":"splitcomplement"}}',
  '{"id":"triad","operation":"palette","input":"red","args":{"method":"triad"}}',
  '{"id":"tetrad","operation":"palette","input":"red","args":{"method":"tetrad"}}',
  '{"id":"analogous","operation":"palette","input":"red","args":{"method":"analogous"}}',
  '{"id":"analogous-zero","operation":"palette","input":"red","args":{"method":"analogous","results":0,"slices":0}}',
  '{"id":"monochromatic","operation":"palette","input":"red","args":{"method":"monochromatic"}}',
  '{"id":"monochromatic-zero","operation":"palette","input":"red","args":{"method":"monochromatic","results":0}}',
  '{"id":"bad-palette","operation":"palette","input":"red","args":{"method":"unknown"}}',
].join("\n"));

const paletteValues = (response) => response.result.map((color) => color.value);
assert.deepEqual(paletteValues(complement), ["#00ffff"]);
assert.deepEqual(paletteValues(splitComplement), ["red", "#ccff00", "#0066ff"]);
assert.deepEqual(paletteValues(triad), ["red", "#00ff00", "#0000ff"]);
assert.deepEqual(paletteValues(tetrad), ["red", "#80ff00", "#00ffff", "#7f00ff"]);
assert.deepEqual(paletteValues(analogous), ["red", "#ff0066", "#ff0033", "#ff0000", "#ff3300", "#ff6600"]);
assert.deepEqual(paletteValues(analogousZero), paletteValues(analogous));
assert.deepEqual(paletteValues(monochromatic), ["#ff0000", "#2a0000", "#550000", "#800000", "#aa0000", "#d40000"]);
assert.deepEqual(paletteValues(monochromaticZero), paletteValues(monochromatic));
assert.deepEqual(badPalette, { id: "bad-palette", error: "unsupported method" });
