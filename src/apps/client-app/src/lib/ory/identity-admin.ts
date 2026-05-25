import { readFileSync } from 'node:fs';
import path from 'node:path';

function readInternalSecret(): string | undefined {
  if (process.env.INTERNAL_SECRET) {
    return process.env.INTERNAL_SECRET;
  }

  try {
    const envPath = path.resolve(process.cwd(), '../auth-service/.env');
    const envText = readFileSync(envPath, 'utf8');
    const match = envText.match(/^INTERNAL_SECRET=(.+)$/m);
    return match?.[1]?.trim();
  } catch {
    return undefined;
  }
}

export interface IdentityClaims {
  email: string;
  role: string;
  org_id: string;
  store_ids: string[];
  warehouse_ids: string[];
}

export async function resolveIdentityClaims(subject?: string | null): Promise<IdentityClaims> {
  const fallback: IdentityClaims = {
    email: '',
    role: 'GUEST',
    org_id: 'org-root-001',
    store_ids: [],
    warehouse_ids: [],
  };
  if (!subject) {
    return fallback;
  }

  const headers = new Headers();
  const internalSecret = readInternalSecret();
  if (internalSecret) {
    headers.set('X-Internal-Secret', internalSecret);
  }

  try {
    const response = await fetch(
      `${process.env.KRATOS_ADMIN_URL || 'http://localhost:4434'}/admin/identities/${subject}`,
      {
        headers,
        cache: 'no-store',
      }
    );

    if (!response.ok) {
      return fallback;
    }

    const identity = (await response.json()) as {
      traits?: {
        email?: string;
        role?: string;
        org_id?: string;
        store_ids?: string[];
        warehouse_ids?: string[];
      };
    };
    return {
      email: identity.traits?.email || '',
      role: identity.traits?.role || 'GUEST',
      org_id: identity.traits?.org_id || fallback.org_id,
      store_ids: Array.isArray(identity.traits?.store_ids) ? identity.traits.store_ids : [],
      warehouse_ids: Array.isArray(identity.traits?.warehouse_ids) ? identity.traits.warehouse_ids : [],
    };
  } catch {
    return fallback;
  }
}

export async function resolveIdentityRole(subject?: string | null): Promise<string> {
  const claims = await resolveIdentityClaims(subject);
  return claims.role;
}
