export const AUTH_PARAMS = {
  ACCESS_TOKEN: 'access_token',
  LOGIN_CHALLENGE: 'login_challenge',
  CONSENT_CHALLENGE: 'consent_challenge',
} as const;

export const OIDC_CONFIG = {
  CLIENT_ID: process.env.NEXT_PUBLIC_OIDC_CLIENT_ID || 'client-app',
  CLIENT_SECRET: process.env.OIDC_CLIENT_SECRET || 'client-secret', 
  GRANT_TYPE: 'authorization_code',
} as const;

/**
 * Danh sách các OAuth2 Client được tin cậy hoàn toàn.
 * Khi Hydra yêu cầu Consent cho những Client này,
 * server sẽ tự động Accept mà không hiển thị UI.
 */
export const TRUSTED_CLIENTS: readonly string[] = [
  OIDC_CONFIG.CLIENT_ID,
] as const;
