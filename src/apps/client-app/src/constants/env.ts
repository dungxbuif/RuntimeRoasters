export const ENV = {
  KRATOS_PUBLIC_URL: process.env.NEXT_PUBLIC_KRATOS_PUBLIC_URL || 'http://localhost:4433',
  HYDRA_PUBLIC_URL: process.env.NEXT_PUBLIC_HYDRA_PUBLIC_URL || 'http://localhost:4444',
  GATEWAY_URL: process.env.NEXT_PUBLIC_GATEWAY_URL || 'http://localhost:8081',
  APP_URL: process.env.NEXT_PUBLIC_APP_URL || 'http://localhost:3000',
  
  // Server-side only
  HYDRA_ADMIN_URL: process.env.HYDRA_ADMIN_URL || 'http://localhost:4445',
} as const;
