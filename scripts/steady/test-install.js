"use strict";

const assert = require("assert");
const childProcess = require("child_process");
const crypto = require("crypto");
const fs = require("fs");
const os = require("os");
const path = require("path");

const toolDirectory = __dirname;
const manifest = JSON.parse(fs.readFileSync(path.join(toolDirectory, "package.json")));
const lockfile = JSON.parse(fs.readFileSync(path.join(toolDirectory, "package-lock.json")));
const expectedVersion = manifest.dependencies["@stdy/cli"];
const platforms = ["darwin-arm64", "darwin-x64", "linux-arm64", "linux-x64", "win32-x64"];

const npmStub = `#!/usr/bin/env node
const fs = require("fs");
const path = require("path");
const args = process.argv.slice(2);

fs.appendFileSync(process.env.STEADY_NPM_LOG, JSON.stringify(args) + "\\n");
if (args[0] !== "ci" || !args.includes("--ignore-scripts")) {
  fs.writeFileSync(process.env.STEADY_LIFECYCLE_MARKER, "executed");
  process.exit(91);
}
if (process.env.STEADY_INSTALL_DELAY) {
  Atomics.wait(new Int32Array(new SharedArrayBuffer(4)), 0, 0,
    Number(process.env.STEADY_INSTALL_DELAY));
}

const root = args[args.indexOf("--prefix") + 1];
const modules = path.join(root, "node_modules");
const cliDirectory = path.join(modules, "@stdy", "cli");
fs.mkdirSync(path.join(modules, ".bin"), { recursive: true });
fs.mkdirSync(cliDirectory, { recursive: true });
fs.writeFileSync(path.join(cliDirectory, "steady.js"), "#!/bin/sh\\nexit 0\\n", { mode: 0o755 });
fs.symlinkSync("../@stdy/cli/steady.js", path.join(modules, ".bin", "steady"));
if (process.env.STEADY_FAIL_AFTER_PARTIAL_INSTALL) process.exit(92);

const expected = require(path.join(root, "package.json")).dependencies["@stdy/cli"];
const nativeName = "@stdy/cli-" + process.platform + "-" + process.arch;
const nativeDirectory = path.join(modules, nativeName);
fs.mkdirSync(path.join(nativeDirectory, "bin"), { recursive: true });
fs.writeFileSync(path.join(cliDirectory, "package.json"), JSON.stringify({
  name: "@stdy/cli", version: expected, optionalDependencies: { [nativeName]: expected },
}));
fs.writeFileSync(path.join(nativeDirectory, "package.json"), JSON.stringify({
  name: nativeName, version: process.env.STEADY_WRONG_NATIVE_VERSION ? "0.0.0" : expected,
}));
fs.writeFileSync(path.join(nativeDirectory, "bin", "steady"), "#!/bin/sh\\nexit 0\\n", {
  mode: 0o755,
});
`;

function fixture() {
  const directory = fs.realpathSync(fs.mkdtempSync(
    path.join(os.tmpdir(), "openai-go-steady-test-")));
  const binaryDirectory = path.join(directory, "fake-bin");
  fs.mkdirSync(binaryDirectory);

  for (const file of ["install", "update", "package.json", "package-lock.json", ".npmrc"]) {
    fs.copyFileSync(path.join(toolDirectory, file), path.join(directory, file));
  }

  const fixtureManifestPath = path.join(directory, "package.json");
  const fixtureManifest = JSON.parse(fs.readFileSync(fixtureManifestPath));
  const digest = crypto.createHash("sha256").update("#!/bin/sh\nexit 0\n").digest("hex");
  fixtureManifest.steadyIntegrity.wrapperSha256 = digest;
  fixtureManifest.steadyIntegrity.nativeBinarySha256[`${process.platform}-${process.arch}`] = digest;
  fs.writeFileSync(fixtureManifestPath, JSON.stringify(fixtureManifest));
  fs.chmodSync(path.join(directory, "install"), 0o755);
  fs.chmodSync(path.join(directory, "update"), 0o755);
  fs.writeFileSync(path.join(binaryDirectory, "npm"), npmStub, { mode: 0o755 });

  return {
    directory,
    executable: path.join(directory, "install"),
    log: path.join(directory, "npm.log"),
    lifecycleMarker: path.join(directory, "lifecycle-ran"),
    environment: {
      ...process.env,
      PATH: `${binaryDirectory}${path.delimiter}${process.env.PATH}`,
      STEADY_NPM_LOG: path.join(directory, "npm.log"),
      STEADY_LIFECYCLE_MARKER: path.join(directory, "lifecycle-ran"),
    },
  };
}

async function withFixture(run) {
  const project = fixture();
  try {
    return await run(project);
  } finally {
    fs.rmSync(project.directory, { recursive: true, force: true });
  }
}

function install(project, overrides = {}) {
  return childProcess.spawnSync(project.executable, [], {
    encoding: "utf8", env: { ...project.environment, ...overrides },
  });
}

function npmCalls(project) {
  return fs.existsSync(project.log)
    ? fs.readFileSync(project.log, "utf8").trim().split("\n").map(JSON.parse)
    : [];
}

function expectSuccess(result) {
  assert.strictEqual(result.status, 0, result.stderr || result.stdout);
}

async function runWorkers(project, count, shouldSucceed = true) {
  const workers = Array.from({ length: count }, () => new Promise((resolve, reject) => {
    const worker = childProcess.spawn(project.executable, [], {
      env: { ...project.environment, STEADY_INSTALL_DELAY: "200" },
      stdio: ["ignore", "ignore", "pipe"],
    });
    let stderr = "";
    worker.stderr.on("data", (chunk) => { stderr += chunk; });
    worker.on("error", reject);
    worker.on("exit", (code) => {
      if ((code === 0) === shouldSucceed) resolve();
      else reject(new Error(`installer exited ${code}: ${stderr}`));
    });
  }));
  await Promise.all(workers);
}

async function test(name, run) {
  await run();
  process.stdout.write(`ok - ${name}\n`);
}

async function main() {
  await test("locks reviewed licenses, SHA-512 archives, and trusted executable SHA-256s", () => {
    assert.match(expectedVersion, /^\d+\.\d+\.\d+$/);
    assert.strictEqual(manifest.license, "Apache-2.0");
    assert.strictEqual(lockfile.packages[""].license, "Apache-2.0");
    assert.deepStrictEqual(Object.keys(manifest.dependencies), ["@stdy/cli"]);
    assert.strictEqual(lockfile.packages[""].dependencies["@stdy/cli"], expectedVersion);
    assert.match(manifest.steadyIntegrity.wrapperSha256, /^[a-f0-9]{64}$/);
    assert.deepStrictEqual(Object.keys(manifest.steadyIntegrity.nativeBinarySha256).sort(), platforms);

    const cli = lockfile.packages["node_modules/@stdy/cli"];
    assert.deepStrictEqual(Object.keys(cli.optionalDependencies).sort(),
      platforms.map((platform) => `@stdy/cli-${platform}`));
    assert.strictEqual(Object.keys(lockfile.packages).length, platforms.length + 2);

    for (const name of ["@stdy/cli", ...platforms.map((platform) => `@stdy/cli-${platform}`)]) {
      const entry = lockfile.packages[`node_modules/${name}`];
      assert.ok(entry, `missing locked package ${name}`);
      assert.strictEqual(entry.version, expectedVersion);
      assert.strictEqual(entry.license, "MIT");
      assert.match(entry.integrity, /^sha512-[A-Za-z0-9+/]+={0,2}$/);
      assert.strictEqual(new URL(entry.resolved).hostname, "registry.npmjs.org");
      assert.notStrictEqual(entry.hasInstallScript, true);
      if (name !== "@stdy/cli") {
        assert.strictEqual(entry.optional, true);
        assert.strictEqual(cli.optionalDependencies[name], expectedVersion);
        assert.strictEqual(entry.os.length, 1);
        assert.strictEqual(entry.cpu.length, 1);
        assert.match(manifest.steadyIntegrity.nativeBinarySha256[
          name.replace("@stdy/cli-", "")
        ], /^[a-f0-9]{64}$/);
      }
    }

    assert.match(fs.readFileSync(path.join(toolDirectory, ".npmrc"), "utf8"),
      /^ignore-scripts=true$/m);
  });

  await test("installs once without lifecycle scripts and reuses verified bytes offline", () =>
    withFixture((project) => {
      expectSuccess(install(project));
      const [args] = npmCalls(project);
      assert.deepStrictEqual(args.slice(0, -1), [
        "ci", "--ignore-scripts", "--include=optional", "--no-audit", "--no-fund", "--prefix",
      ]);
      assert.ok(args[args.length - 1].startsWith(`${project.directory}/.install.`));
      assert.ok(!fs.existsSync(project.lifecycleMarker));
      expectSuccess(install(project, { STEADY_FAIL_AFTER_PARTIAL_INSTALL: "1" }));
      assert.strictEqual(npmCalls(project).length, 1);
    }));

  await test("rejects mismatched package metadata until explicitly removed", () =>
    withFixture((project) => {
      expectSuccess(install(project));
      const file = path.join(project.directory, "node_modules", "@stdy", "cli", "package.json");
      const installed = JSON.parse(fs.readFileSync(file));
      installed.version = "0.0.0";
      fs.writeFileSync(file, JSON.stringify(installed));
      assert.notStrictEqual(install(project).status, 0);
      assert.strictEqual(npmCalls(project).length, 1);
      fs.rmSync(path.join(project.directory, "node_modules"), { recursive: true, force: true });
      expectSuccess(install(project));
      assert.strictEqual(npmCalls(project).length, 2);
    }));

  await test("never publishes a partial or incorrectly versioned installation", () =>
    withFixture((project) => {
      assert.notStrictEqual(install(project, { STEADY_FAIL_AFTER_PARTIAL_INSTALL: "1" }).status, 0);
      assert.ok(!fs.existsSync(path.join(project.directory, "node_modules")));
      assert.notStrictEqual(install(project, { STEADY_WRONG_NATIVE_VERSION: "1" }).status, 0);
      assert.ok(!fs.existsSync(path.join(project.directory, "node_modules")));
      expectSuccess(install(project));
      assert.strictEqual(npmCalls(project).length, 3);
    }));

  await test("rejects native, wrapper, launcher, and forged-cache substitutions", () =>
    withFixture((project) => {
      expectSuccess(install(project));
      const modules = path.join(project.directory, "node_modules");
      const native = path.join(modules, "@stdy", `cli-${process.platform}-${process.arch}`,
        "bin", "steady");
      fs.writeFileSync(native, "#!/bin/sh\necho substituted-native\n", { mode: 0o755 });
      fs.writeFileSync(path.join(modules, ".forged-install-marker"), "forged\n");
      assert.notStrictEqual(install(project).status, 0);
      fs.rmSync(modules, { recursive: true, force: true });
      expectSuccess(install(project));

      const wrapper = path.join(modules, "@stdy", "cli", "steady.js");
      fs.writeFileSync(wrapper, "#!/bin/sh\necho substituted-wrapper\n", { mode: 0o755 });
      assert.notStrictEqual(install(project).status, 0);
      fs.rmSync(modules, { recursive: true, force: true });
      expectSuccess(install(project));

      const external = path.join(project.directory, "external-steady");
      const launcher = path.join(modules, ".bin", "steady");
      fs.writeFileSync(external, "#!/bin/sh\necho substituted-launcher\n", { mode: 0o755 });
      fs.unlinkSync(launcher);
      fs.symlinkSync(external, launcher);
      assert.notStrictEqual(install(project).status, 0);
      fs.rmSync(modules, { recursive: true, force: true });
      expectSuccess(install(project));
      assert.strictEqual(npmCalls(project).length, 4);
    }));

  await test("an abandoned partial staging directory never blocks a later install", () =>
    withFixture((project) => {
      fs.mkdirSync(path.join(project.directory, ".install.abandoned", "node_modules"),
        { recursive: true });
      expectSuccess(install(project));
      assert.strictEqual(npmCalls(project).length, 1);
    }));

  await test("rejects a node_modules symlink escaping the trusted project", () =>
    withFixture((project) => {
      const external = fs.mkdtempSync(path.join(os.tmpdir(), "openai-go-steady-external-"));
      try {
        expectSuccess(install(project));
        const modules = path.join(project.directory, "node_modules");
        const externalModules = path.join(external, "node_modules");
        fs.renameSync(modules, externalModules);
        fs.symlinkSync(externalModules, modules);
        const result = install(project);
        assert.notStrictEqual(result.status, 0);
        assert.match(result.stderr, /failed integrity verification/);
        assert.strictEqual(npmCalls(project).length, 1);
      } finally {
        fs.rmSync(external, { recursive: true, force: true });
      }
    }));

  await test("concurrent installers publish only complete verified directories", () =>
    withFixture(async (project) => {
      await runWorkers(project, 6);
      assert.ok(npmCalls(project).length >= 1);
      assert.ok(npmCalls(project).length <= 6);
      expectSuccess(install(project));
      assert.ok(!fs.readdirSync(project.directory).some((file) => file.startsWith(".install.")));
    }));

  await test("concurrent callers fail closed after executable substitution", () =>
    withFixture(async (project) => {
      expectSuccess(install(project));
      const native = path.join(project.directory, "node_modules", "@stdy",
        `cli-${process.platform}-${process.arch}`, "bin", "steady");
      fs.writeFileSync(native, "#!/bin/sh\necho substituted-native\n", { mode: 0o755 });
      await runWorkers(project, 6, false);
      assert.strictEqual(npmCalls(project).length, 1);
      fs.rmSync(path.join(project.directory, "node_modules"), { recursive: true, force: true });
      await runWorkers(project, 6);
      expectSuccess(install(project));
    }));

  await test("the updater rejects an untrusted archive without changing committed inputs", () =>
    withFixture((project) => {
      const manifestPath = path.join(project.directory, "package.json");
      const lockfilePath = path.join(project.directory, "package-lock.json");
      const originalManifest = fs.readFileSync(manifestPath, "utf8");
      const originalLockfile = fs.readFileSync(lockfilePath, "utf8");
      const updaterNpm = `#!/usr/bin/env node
const fs = require("fs");
const path = require("path");
const args = process.argv.slice(2);
fs.appendFileSync(process.env.STEADY_NPM_LOG, JSON.stringify(args) + "\\n");
if (!args.includes("--ignore-scripts")) process.exit(91);
if (args[0] === "install") process.exit(0);
if (args[0] !== "pack") process.exit(92);
const root = args[args.indexOf("--pack-destination") + 1];
const filename = "substituted.tgz";
fs.writeFileSync(path.join(root, filename), "untrusted package contents");
process.stdout.write(filename + "\\n");
`;
      fs.writeFileSync(path.join(project.directory, "fake-bin", "npm"), updaterNpm,
        { mode: 0o755 });

      const result = childProcess.spawnSync(path.join(project.directory, "update"),
        [expectedVersion], { encoding: "utf8", env: project.environment });
      assert.notStrictEqual(result.status, 0);
      assert.match(result.stderr, /does not match its locked SHA-512 integrity/);
      assert.strictEqual(fs.readFileSync(manifestPath, "utf8"), originalManifest);
      assert.strictEqual(fs.readFileSync(lockfilePath, "utf8"), originalLockfile);
      assert.ok(!fs.readdirSync(project.directory).some((file) => file.startsWith(".update.")));
    }));

  await test("npm rejects a substituted tarball against its original integrity hash", () => {
    const directory = fs.mkdtempSync(path.join(fs.realpathSync(os.tmpdir()),
      "openai-go-steady-integrity-"));
    try {
      const packageDirectory = path.join(directory, "fixture");
      const cacheDirectory = path.join(directory, "cache");
      fs.mkdirSync(packageDirectory);
      fs.writeFileSync(path.join(packageDirectory, "package.json"), JSON.stringify({
        name: "integrity-fixture", version: "1.0.0",
      }));
      fs.writeFileSync(path.join(packageDirectory, "index.js"), "module.exports = true;\n");
      fs.writeFileSync(path.join(directory, "package.json"), JSON.stringify({
        name: "integrity-root", private: true,
        dependencies: { "integrity-fixture": "file:./integrity-fixture-1.0.0.tgz" },
      }));

      const npm = (args) => childProcess.spawnSync("npm", args, { encoding: "utf8" });
      const pack = () => npm([
        "pack", "--silent", "--ignore-scripts", "--pack-destination", directory, packageDirectory,
      ]);
      expectSuccess(pack());
      expectSuccess(npm([
        "install", "--package-lock-only", "--ignore-scripts", "--offline", "--no-audit",
        "--no-fund", "--cache", cacheDirectory, "--prefix", directory,
      ]));
      fs.rmSync(cacheDirectory, { recursive: true, force: true });
      fs.writeFileSync(path.join(packageDirectory, "index.js"), "module.exports = false;\n");
      expectSuccess(pack());

      const result = npm([
        "ci", "--offline", "--ignore-scripts", "--no-audit", "--no-fund", "--cache",
        cacheDirectory, "--prefix", directory,
      ]);
      assert.notStrictEqual(result.status, 0, "npm accepted substituted package contents");
      assert.match(result.stderr, /EINTEGRITY/);
    } finally {
      fs.rmSync(directory, { recursive: true, force: true });
    }
  });
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
