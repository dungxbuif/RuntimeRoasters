export const API_ENDPOINTS = {
  AUTH: {
    LOGIN_ACCEPT: '/v1/auth/login/accept',
    CONSENT_ACCEPT: '/api/auth/consent/accept',
    CALLBACK: '/api/auth/callback',
  },
  DEMO: {
    PING: '/v1/demo/ping',
  },
} as const;
