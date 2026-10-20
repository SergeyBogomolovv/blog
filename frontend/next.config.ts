import path from "node:path";
import { loadEnvConfig, updateInitialEnv } from "@next/env";
import type { NextConfig } from "next";

// Next.js и Go читают один корневой .env, включая dev-воркеры Next.js.
const { parsedEnv } = loadEnvConfig(
  path.resolve(process.cwd(), ".."),
  process.env.NODE_ENV === "development",
  undefined,
  true,
);
updateInitialEnv(parsedEnv || {});

// Для HTML proxy.ts заменяет эту политику вариантом с уникальным nonce.
const defaultContentSecurityPolicy = [
  "default-src 'self'",
  "script-src 'self'",
  "object-src 'none'",
  "base-uri 'self'",
  "form-action 'self'",
  "frame-ancestors 'none'",
].join("; ");

const nextConfig: NextConfig = {
  poweredByHeader: false,
  devIndicators: false,
  turbopack: {
    rules: {
      "*.css": {
        loaders: ["@tailwindcss/turbopack"],
        as: "*.css",
      },
    },
  },
  async headers() {
    return [
      {
        source: "/:path*",
        headers: [
          { key: "X-Content-Type-Options", value: "nosniff" },
          { key: "X-Frame-Options", value: "DENY" },
          { key: "Referrer-Policy", value: "no-referrer" },
          {
            key: "Permissions-Policy",
            value: "camera=(), microphone=(), geolocation=()",
          },
          {
            key: "Content-Security-Policy",
            value: defaultContentSecurityPolicy,
          },
        ],
      },
    ];
  },
};

export default nextConfig;
