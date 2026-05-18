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

export async function resolveIdentityRole(subject?: string | null): Promise<string> {
  if (!subject) {
    return 'GUEST';
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
      return 'GUEST';
    }

    const identity = (await response.json()) as { traits?: { role?: string } };
    return identity.traits?.role || 'GUEST';
  } catch {
    return 'GUEST';
  }
}
