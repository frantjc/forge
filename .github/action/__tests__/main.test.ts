import { afterEach, beforeAll, expect, test } from "@jest/globals";

import cp from "child_process";
import fs from "fs";
import os from "os";
import process from "process";
import path from "path";

let RUNNER_OS = "Linux";
switch (os.platform()) {
  case "darwin":
    RUNNER_OS = "macOS";
    break;
  case "win32":
    RUNNER_OS = "Windows";
    break;
}

let RUNNER_ARCH = "X64";
switch (os.arch()) {
  case "arm":
    RUNNER_ARCH = "ARM";
    break;
  case "arm64":
    RUNNER_ARCH = "ARM64";
    break;
}

const env = {
  RUNNER_TEMP: `${os.tmpdir()}/tmp`,
  RUNNER_TOOL_CACHE: `${os.tmpdir()}/tc`,
  RUNNER_OS,
  RUNNER_ARCH,
};

beforeAll(() => {
  fs.mkdirSync(env.RUNNER_TEMP, { recursive: true });
  fs.mkdirSync(env.RUNNER_TOOL_CACHE, { recursive: true });
});

afterEach(() => {
  fs.rmSync(env.RUNNER_TEMP, { recursive: true });
  fs.rmSync(env.RUNNER_TOOL_CACHE, { recursive: true });
});

test("run install", () => {
  expect(
    cp
      .execFileSync(
        process.execPath,
        [path.join(__dirname, "../dist/index.js")],
        {
          env: {
            ...process.env,
            ...env,
            INPUT_INSTALL: "true",
          },
        }
      )
      .toString()
  ).toContain("forge -v");
});
