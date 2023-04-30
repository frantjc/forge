import * as core from "@actions/core";
import * as tc from "@actions/tool-cache";
import * as cp from "@actions/exec";

import fs from "fs";
import path from "path";
import yaml from "yaml";

import packageJSON from "../package.json";

const goreleaserYML = yaml.parse(
  fs.readFileSync(path.join(__dirname, "../../.goreleaser.yaml")).toString()
);

async function run(): Promise<void> {
  try {
    const tool = "forge";
    const version = core.getInput("version") || packageJSON.version;

    const get = core.getInput("get");
    const put = core.getInput("put");

    const cwd = process.env.GITHUB_WORKSPACE;

    if (get && put) {
      throw new Error("cannot use with `get` and `put`");
    }

    let runnerArch;
    switch (process.env.RUNNER_ARCH) {
      case "X64":
        runnerArch = "amd64";
        break;
    }

    if (!goreleaserYML.builds[0].goarch.includes(runnerArch)) {
      throw new Error(`unsupported architecture ${process.env.RUNNER_ARCH}`);
    }

    let runnerOs;
    switch (process.env.RUNNER_OS) {
      case "Linux":
        runnerOs = "linux";
        break;
      case "Windows":
        runnerOs = "windows";
        break;
      case "macOS":
        runnerOs = "darwin";
        break;
    }

    if (!goreleaserYML.builds[0].goos.includes(runnerOs)) {
      throw new Error(`unsupported OS ${process.env.RUNNER_OS}`);
    }

    let forge = "forge";
    if (core.getBooleanInput("install")) {
      forge = tc.find(tool, version);
      if (!forge) {
        forge = await tc.cacheFile(
          path.join(
            await tc.extractTar(
              await tc.downloadTool(
                `https://github.com/frantjc/forge/releases/download/v${version}/forge_${version}_${runnerOs}_${runnerArch}.tar.gz`
              )
            ),
            tool
          ),
          tool,
          tool,
          version
        );
      }

      forge = path.join(forge, "forge");
    }

    await cp.exec(forge, ["-v"]);

    if (get) {
      await cp.exec(forge, ["get", get], { cwd });
    } else if (put) {
      await cp.exec(forge, ["put", put], { cwd });
    }
  } catch (err) {
    if (typeof err === "string" || err instanceof Error) core.setFailed(err);
  }
}

run();
