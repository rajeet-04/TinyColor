import readline from "node:readline";
import tinycolor from "../mod.js";

const inspect = (color) => ({
  valid: color.isValid(),
  format: color.getFormat(),
  alpha: color.getAlpha(),
  rgb: color.toRgb(),
  value: color.toString(),
  original: color.getOriginalInput(),
});

const respond = (id, result, error) => {
  if ((result === undefined) === (error === undefined)) {
    throw new Error("response must contain exactly one of result or error");
  }
  return error === undefined ? { id, result } : { id, error };
};

const handle = (request) => {
  switch (request.operation) {
    case "inspect":
      return respond(request.id, inspect(tinycolor(request.input)));
    case "string":
      return respond(request.id, tinycolor(request.input).toString(request.args?.format));
    case "fromRatio":
      return respond(request.id, inspect(tinycolor.fromRatio(request.input, request.args)));
    default:
      return respond(request.id, undefined, "unsupported operation");
  }
};

const lines = readline.createInterface({ input: process.stdin, crlfDelay: Infinity });
for await (const line of lines) {
  if (!line) continue;
  try {
    const request = JSON.parse(line);
    if (!request.id || !request.operation) throw new Error("id and operation are required");
    console.log(JSON.stringify(handle(request)));
  } catch (error) {
    console.log(JSON.stringify(respond("", undefined, error instanceof SyntaxError ? "malformed JSON" : error.message)));
  }
}
