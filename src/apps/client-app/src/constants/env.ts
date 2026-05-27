const getEnv = (key: string, required = true): string => {
  const demoDefaults: Record<string, string> = {
    NEXT_PUBLIC_KRATOS_PUBLIC_URL: 'http://localhost:4433',
    NEXT_PUBLIC_HYDRA_PUBLIC_URL: 'http://localhost:4444',
    NEXT_PUBLIC_GATEWAY_URL: 'http://localhost:8081',
    NEXT_PUBLIC_APP_URL: 'http://localhost:3000',
    NEXT_PUBLIC_SIGNOZ_URL: 'http://localhost:3301',
    NEXT_PUBLIC_JAEGER_URL: 'http://localhost:3301',
    NEXT_PUBLIC_SWAGGER_JSON_URL: 'http://localhost:8081/swagger/demo.swagger.json',
    NEXT_PUBLIC_SOCKET_URL: 'ws://localhost:8091',
  };
  const value = process.env[key] ?? demoDefaults[key];
  if (required && !value) {
    if (typeof window !== 'undefined') {
       // In browser, we might want to alert or just throw
       const error = `FATAL: Missing required environment variable: ${key}`;
       console.error(error);
       throw new Error(error);
    }
    throw new Error(`FATAL: Missing required environment variable: ${key}`);
  }
  return value || '';
};

export const ENV = {
  KRATOS_PUBLIC_URL: getEnv('NEXT_PUBLIC_KRATOS_PUBLIC_URL'),
  HYDRA_PUBLIC_URL: getEnv('NEXT_PUBLIC_HYDRA_PUBLIC_URL'),
  GATEWAY_URL: getEnv('NEXT_PUBLIC_GATEWAY_URL'),
  APP_URL: getEnv('NEXT_PUBLIC_APP_URL'),
  SIGNOZ_URL: getEnv('NEXT_PUBLIC_SIGNOZ_URL'),
  JAEGER_URL: getEnv('NEXT_PUBLIC_JAEGER_URL'),
  SWAGGER_JSON_URL: getEnv('NEXT_PUBLIC_SWAGGER_JSON_URL'),
  SOCKET_URL: getEnv('NEXT_PUBLIC_SOCKET_URL'),
  
  // Server-side only
  HYDRA_ADMIN_URL: getEnv('HYDRA_ADMIN_URL', false), // Hydra Admin might be internal/optional for some nodes
} as const;
