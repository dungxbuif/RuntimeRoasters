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
  TRACE: {
    BY_ID: (id: string) => `/v1/trace/${id}`,
    DOCUMENT: (id: string) => `/v1/trace/${id}/document`,
  },
  LOGISTICS: {
    SHIPMENTS: '/v1/logistics/shipments',
    LOCATIONS: '/v1/logistics/locations',
    UPDATE_LOCATION: '/v1/logistics/drivers/location',
  },
  PAYMENT: {
    PAYMENTS: '/v1/payments',
    BY_ORDER: (orderId: string) => `/v1/payments/orders/${orderId}`,
  },
  AUDIT: {
    BY_PARTITION: (partition: string) => `/v1/audit/${partition}`,
  },
} as const;
