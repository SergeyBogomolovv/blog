// Проверяем реальную конфигурацию после предварительной загрузки env самим Next.js.
import assert from "node:assert/strict";
import { spawnSync } from "node:child_process";
import { mkdirSync, mkdtempSync, rmSync, writeFileSync } from "node:fs";
import { tmpdir } from "node:os";
import path from "node:path";
import { createRequire } from "node:module";
import { test } from "node:test";
import { fileURLToPath, pathToFileURL } from "node:url";

const require = createRequire(import.meta.url);
const envURL = pathToFileURL(require.resolve("@next/env")).href;
const configPath = fileURLToPath(new URL("../next.config.ts", import.meta.url));
const loaderURL = pathToFileURL(
  require.resolve("next/dist/build/next-config-ts/transpile-config.js"),
).href;

for (const fromProcess of [false, true]) {
  test(`общий .env загружается после кеша Next.js; ключ из окружения: ${fromProcess}`, (t) => {
    const root = mkdtempSync(path.join(tmpdir(), "studyblog-env-"));
    t.after(() => rmSync(root, { recursive: true, force: true }));
    mkdirSync(path.join(root, "frontend"));
    writeFileSync(
      path.join(root, ".env"),
      `JWT_SECRET=${"01".repeat(32)}\nCOOKIE_SECURE=true\nAPP_ORIGIN=https://blog.example.test\n`,
    );
    const childEnv = {
      ...process.env,
      NODE_ENV: "development",
      __NEXT_NODE_NATIVE_TS_LOADER_ENABLED: "false",
    };

    for (const key of [
      "JWT_SECRET",
      "COOKIE_SECURE",
      "APP_ORIGIN",
      "__NEXT_PROCESSED_ENV",
    ]) {
      delete childEnv[key];
    }

    if (fromProcess) {
      childEnv.JWT_SECRET = "02".repeat(32);
    }

    const script = `
      import assert from "node:assert/strict";
      import env from ${JSON.stringify(envURL)};
      import loader from ${JSON.stringify(loaderURL)};
      env.loadEnvConfig(process.cwd(), true);
      assert.equal(process.env.COOKIE_SECURE, undefined);
      await loader.transpileConfig({
        nextConfigPath: ${JSON.stringify(configPath)},
        dir: ${JSON.stringify(path.dirname(configPath))},
      });
      const expectedKey = ${JSON.stringify((fromProcess ? "02" : "01").repeat(32))};
      assert.equal(process.env.JWT_SECRET, expectedKey);
      assert.equal(process.env.COOKIE_SECURE, "true");
      assert.equal(process.env.APP_ORIGIN, "https://blog.example.test");
      env.resetEnv();
      assert.equal(process.env.JWT_SECRET, expectedKey);
      assert.equal(process.env.COOKIE_SECURE, "true");
      assert.equal(process.env.APP_ORIGIN, "https://blog.example.test");
    `;
    const result = spawnSync(
      process.execPath,
      ["--experimental-strip-types", "--input-type=module", "-e", script],
      { cwd: path.join(root, "frontend"), env: childEnv, encoding: "utf8" },
    );

    assert.equal(result.status, 0, result.stderr);
  });
}
