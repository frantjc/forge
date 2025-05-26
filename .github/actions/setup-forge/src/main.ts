import * as core from "@actions/core";
import * as tc from "@actions/tool-cache";
import * as cp from "@actions/exec";

import fs from "fs";
import path from "path";
import yaml from "yaml";

const packageJSON = JSON.parse(
  fs.readFileSync(path.join(__dirname, "../package.json")).toString(),
);

const goreleaserYML = yaml.parse(
  fs
    .readFileSync(path.join(__dirname, "../../../../.goreleaser.yaml"))
    .toString(),
);

async function run(): Promise<void> {
  try {
    const tool = "forge";
    const version = packageJSON.version;

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
    // to build the GOOS that we are trying to download.
    if (!goreleaserYML.builds[0].goos.includes(os)) {
      throw new Error(`unsupported OS ${process.env.RUNNER_OS}`);
    }

    core.startGroup("install");

    // Look for forge in the cache.
    let dir = tc.find(tool, versionOs);

    // If we don't find forge in the cache, download, extract and cache it
    // from its GitHub release.
    if (!dir) {
      dir = await tc.cacheFile(
        path.join(
          await tc.extractTar(
            await tc.downloadTool(
              `https://github.com/frantjc/${tool}/releases/download/v${version}/${tool}_${version}_${os}_${arch}.tar.gz`,
            ),
          ),
          tool,
        ),
        tool,
        tool,
        versionOs,
      );
    }

    const bin = path.join(dir, tool);

    core.addPath(dir);

    core.endGroup();

    // Sanity check that forge was installed correctly.
    await cp.exec(bin, ["version"]);
  } catch (err) {
    if (typeof err === "string" || err instanceof Error) {
      core.setFailed(err);
    } else {
      core.setFailed(`caught unknown error ${err}`);
    }
  }
}

run();
