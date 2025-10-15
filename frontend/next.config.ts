import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  async rewrites() {
    return [
      {
        // This rule will proxy any request starting with /api/
        // to the Go backend server running on localhost:8080
        source: "/api/:path*",
        destination: `${process.env.API_BASE_URL}/api/:path*`,
      },
    ];
  },
};

export default nextConfig;
