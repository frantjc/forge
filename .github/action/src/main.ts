import * as core from "@actions/core";
import * as tc from "@actions/tool-cache";
import * as cp from "@actions/exec";

import fs from "fs";
import path from "path";
import yaml from "yaml";

import packageJSON from "../package.json";

const goreleaserYML = yaml.parse(
  fs.readFileSync(path.join(__dirname, "../../../.goreleaser.yaml")).toString()
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

    const action = get ? "get" : "put";
    const resource = get || put;

    const params = core.getMultilineInput("params");
    const config = core.getInput("config");

    // Turn RUNNER_ARCH into GOARCH.
    let arch;
    switch (process.env.RUNNER_ARCH) {
      case "X86":
      case "X64":
        arch = "amd64";
        break;
    }

    // Before we even attempt the download, check if goreleaser was configured
    // to build the GOARCH that we are trying to download.
    //
    // Note that this would become non-backwards-compatible if we remove support for
    // a GOARCH and it acts funny if we add support for one and someone uses it like so:
    //
    //  - uses: frantjc/forge@v0.6.0
    //    with:
    //      version: 0.7.0
    if (!goreleaserYML.builds[0].goarch.includes(arch)) {
      throw new Error(`unsupported architecture ${process.env.RUNNER_ARCH}`);
    }

    // Turn RUNNER_OS into GOOS.
    let os;
    switch (process.env.RUNNER_OS) {
      case "Linux":
        os = "linux";
        break;
      case "Windows":
        os = "windows";
        break;
      case "macOS":
        os = "darwin";
        break;
    }

    const versionOs = `${version}_${os}`;

    // Before we even attempt the download, check if goreleaser was configured
    // to build the GOOS that we are trying to download
    //
    // Note that this would become non-backwards-compatible if we remove support for
    // a GOOS and it acts funny if we add support for one and someone uses it like so:
    //
    //  - uses: frantjc/forge@v0.6.0
    //    with:
    //      version: 0.7.0
    if (!goreleaserYML.builds[0].goos.includes(os)) {
      throw new Error(`unsupported OS ${process.env.RUNNER_OS}`);
    }

    // Default to looking it up on PATH if install is explicitly set to false.
    let forge = "forge";
    if (core.getBooleanInput("install")) {
      core.startGroup("install");

      // Look for forge in the cache.
      forge = tc.find(tool, versionOs);

      // If we don't find forge in the cache, download, extract and cache it
      // from its GitHub release.
      if (!forge) {
        forge = await tc.cacheFile(
          path.join(
            await tc.extractTar(
              await tc.downloadTool(
                `https://github.com/frantjc/forge/releases/download/v${version}/forge_${version}_${os}_${arch}.tar.gz`
              )
            ),
            tool
          ),
          tool,
          tool,
          versionOs
        );
      }

      core.endGroup();
    }

    // Sanity check that forge was installed correctly.
    await cp.exec(forge, ["-v"]);

    // Inputs for `get` a `put` are not required so that this action can be used to
    // only install forge. Note that we checked above if both were set, so at most
    // one of these conditions could evaluate to true.
    if (resource) {
      await cp.exec(
        forge,
        [
          action,
          resource,
          `-c=${config}`,
          ...params.map((param) => `-p=${param}`),
        ],
        { cwd }
      );
    }
  } catch (err) {
    if (typeof err === "string" || err instanceof Error) core.setFailed(err);
  }
}

run();
