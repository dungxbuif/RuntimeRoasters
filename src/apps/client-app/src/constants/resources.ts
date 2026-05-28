/**
 * Centralized authorization resources and actions for Runtime Roasters.
 * Aligned with backend Casbin policies in default_policies.csv.
 */

export const AUTH_RESOURCES = {
  // Domain Resources
  FARM: 'farm',
  HARVEST: 'harvest',
  ORDER: '/v1/orders',
  STORE: '/v1/stores',
  PAYMENT: '/v1/payments',
  WAREHOUSE_INVENTORY: '/v1/warehouse/inventory',
  LOGISTICS_SHIPMENT: '/v1/logistics/shipments',
  
  // System Resources
  SYSTEM_SERVICE: '/runtime.system.v1.SystemService/*',
  FARM_SERVICE: '/runtime.farm.v1.FarmService/*',
  AUDIT_LOGS: '/v1/audit/.*',
  TRACEABILITY: '/v1/trace/.*',
  TOPOLOGY: '/v1/traces/topology/.*',
  REALTIME: '/v1/realtime/.*',
} as const;

export const AUTH_ACTIONS = {
  READ: 'read',
  WRITE: 'write',
  DELETE: 'delete',
  MANAGE: '.*',
} as const;

export type AuthResource = typeof AUTH_RESOURCES[keyof typeof AUTH_RESOURCES];
export type AuthAction = typeof AUTH_ACTIONS[keyof typeof AUTH_ACTIONS];
