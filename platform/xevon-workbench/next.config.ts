import type { NextConfig } from 'next';

const nextConfig: NextConfig = {
  allowedDevOrigins: ['host.docker.internal', 'local-testing.xevon.live'],
  output: 'export',
  distDir: 'dist',
  trailingSlash: true,
  skipTrailingSlashRedirect: false,
  images: {
    unoptimized: true,
  },
};

export default nextConfig;
