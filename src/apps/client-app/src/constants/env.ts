export const ENV = {
  KRATOS_PUBLIC_URL: process.env.NEXT_PUBLIC_KRATOS_PUBLIC_URL || 'http://localhost:4433',
  HYDRA_PUBLIC_URL: process.env.NEXT_PUBLIC_HYDRA_PUBLIC_URL || 'http://localhost:4444',
  GATEWAY_URL: process.env.NEXT_PUBLIC_GATEWAY_URL || 'http://localhost:8081',
  APP_URL: process.env.NEXT_PUBLIC_APP_URL || 'http://localhost:3000',
  JAEGER_URL: process.env.NEXT_PUBLIC_JAEGER_URL || 'http://localhost:16686',
  SWAGGER_JSON_URL: process.env.NEXT_PUBLIC_SWAGGER_JSON_URL || 'http://localhost:8081/swagger/demo.swagger.json',
  
  // Server-side only
  HYDRA_ADMIN_URL: process.env.HYDRA_ADMIN_URL || 'http://localhost:4445',
} as const;
