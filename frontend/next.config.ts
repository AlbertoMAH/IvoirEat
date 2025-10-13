import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  async rewrites() {
    return [
      {
        source: "/ping",
        destination: "http://localhost:8080/ping",
      },
    ];
  },
};

export default nextConfig;
