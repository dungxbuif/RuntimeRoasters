import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  env: {
    NEXT_PUBLIC_KRATOS_PUBLIC_URL: process.env.NEXT_PUBLIC_KRATOS_PUBLIC_URL ?? 'http://localhost:4433',
    NEXT_PUBLIC_HYDRA_PUBLIC_URL: process.env.NEXT_PUBLIC_HYDRA_PUBLIC_URL ?? 'http://localhost:4444',
    NEXT_PUBLIC_GATEWAY_URL: process.env.NEXT_PUBLIC_GATEWAY_URL ?? 'http://localhost:8081',
    NEXT_PUBLIC_APP_URL: process.env.NEXT_PUBLIC_APP_URL ?? 'http://localhost:3000',
    NEXT_PUBLIC_SIGNOZ_URL: process.env.NEXT_PUBLIC_SIGNOZ_URL ?? 'http://localhost:3301',
    NEXT_PUBLIC_JAEGER_URL: process.env.NEXT_PUBLIC_JAEGER_URL ?? 'http://localhost:3301',
    NEXT_PUBLIC_SWAGGER_JSON_URL: process.env.NEXT_PUBLIC_SWAGGER_JSON_URL ?? 'http://localhost:8081/swagger/demo.swagger.json',
  },
  images: {
    remotePatterns: [
      {
        protocol: 'https',
        hostname: 'lh3.googleusercontent.com',
      },
    ],
  },
};

export default nextConfig;
