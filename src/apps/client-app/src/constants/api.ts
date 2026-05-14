export const API_ENDPOINTS = {
  AUTH: {
    LOGIN_ACCEPT: '/v1/auth/login/accept',
    CONSENT_ACCEPT: '/api/auth/consent/accept',
    CALLBACK: '/api/auth/callback',
  },
  DEMO: {
    PING: '/v1/demo/ping',
  },
  FARM: {
    FARMS: '/v1/farms',
    HARVESTS: '/v1/harvests',
  },
  WAREHOUSE: {
    BATCHES: '/v1/warehouse/batches',
    ROAST_RUNS: '/v1/warehouse/runs',
    INVENTORY: '/v1/warehouse/inventory',
  },
} as const;
