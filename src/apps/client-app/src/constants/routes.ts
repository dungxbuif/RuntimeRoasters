export const APP_ROUTES = {
  HOME: '/',
  LOGIN: '/login',
  CONSENT: '/consent',
  DASHBOARD: {
    USERS: '/dashboard/users',
    FARMS: '/dashboard/farm-ops/registry',
    HARVESTS: '/dashboard/farm-ops/harvests',
    FARM_TELEMETRY: '/dashboard/farm-ops/telemetry',
    BATCHES: '/dashboard/batches',
    LOGISTICS: '/dashboard/logistics',
    WAREHOUSE: '/dashboard/warehouse',
    RETAIL: '/dashboard/retail',
    RETAIL_ORDERS: '/dashboard/retail/orders',
    TOPOLOGY: '/dashboard/topology-mesh',
    EXPLORER: '/dashboard/explorer',
    RESILIENCY: '/dashboard/resiliency',
    TRACEABILITY: '/dashboard/traceability',
    FINANCE: '/dashboard/finance',
    AUDIT: '/dashboard/audit',
    PROFILE: '/dashboard/profile',
  }
} as const;
