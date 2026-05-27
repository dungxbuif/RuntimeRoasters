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
    INTAKES: '/v1/warehouse/intakes',
    PICKUP_REQUESTS: '/v1/warehouse/pickup-requests',
    DISPATCH_REQUESTS: '/v1/warehouse/dispatch-requests',
    BATCHES: '/v1/warehouse/batches',
    ROAST_RUNS: '/v1/warehouse/runs',
    INVENTORY: '/v1/warehouse/inventory',
  },
  TRACE: {
    BY_ID: (id: string) => `/v1/trace/${id}`,
    DOCUMENT: (id: string) => `/v1/trace/${id}/document`,
    PUBLIC_TOPOLOGY_CONFIG: '/v1/traces/public/topology/config',
    PUBLIC_TOPOLOGY_HISTORY: '/v1/traces/public/topology/history',
    TOPOLOGY_CONFIG: '/v1/traces/topology/config',
    TOPOLOGY_HISTORY: '/v1/traces/topology/history',
  },
  LOGISTICS: {
    SHIPMENTS: '/v1/logistics/shipments',
    LOCATIONS: '/v1/logistics/locations',
    UPDATE_LOCATION: '/v1/logistics/drivers/location',
    DRIVERS: '/v1/logistics/drivers',
    VEHICLES: '/v1/logistics/vehicles',
    AVAILABLE_VEHICLES: '/v1/logistics/vehicles/available',
  },
  PAYMENT: {
    PAYMENTS: '/v1/payments',
    BY_ORDER: (orderId: string) => `/v1/payments/orders/${orderId}`,
    STRIPE_WEBHOOK_KEY: '/v1/payments/demo/stripe-webhook-key',
    STRIPE_WEBHOOK: '/v1/webhooks/stripe',
  },
  AUDIT: {
    BY_PARTITION: (partition: string) => `/v1/audit/${partition}`,
  },
  REALTIME: {
    TICKETS: '/v1/realtime/tickets',
  },
} as const;
