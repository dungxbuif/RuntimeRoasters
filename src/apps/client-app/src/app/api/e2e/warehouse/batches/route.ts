import { NextRequest, NextResponse } from 'next/server';
import fs from 'node:fs';
import path from 'node:path';
import { Pool } from 'pg';

let pool: Pool | null = null;

function loadWarehouseDatabaseUrl(): string {
  const envPath = path.resolve(process.cwd(), '../warehouse-service/.env');
  const fileValues = fs.existsSync(envPath)
    ? Object.fromEntries(
        fs
          .readFileSync(envPath, 'utf8')
          .split('\n')
          .map((line) => line.trim())
          .filter((line) => line && !line.startsWith('#') && line.includes('='))
          .map((line) => {
            const idx = line.indexOf('=');
            return [line.slice(0, idx), line.slice(idx + 1)];
          })
      )
    : {};

  return (
    process.env.WAREHOUSE_DATABASE_URL ||
    fileValues.WAREHOUSE_DATABASE_URL ||
    fileValues.DATABASE_URL ||
    'postgres://user:password@localhost:54321/warehouse_db?sslmode=disable'
  );
}

function getPool(): Pool {
  if (!pool) {
    const connectionString = loadWarehouseDatabaseUrl();
    pool = new Pool({
      connectionString,
    });
  }

  return pool;
}

export async function GET(req: NextRequest) {
  if (process.env.NODE_ENV === 'production') {
    return NextResponse.json({ message: 'not found' }, { status: 404 });
  }

  const harvestId = req.nextUrl.searchParams.get('harvestId');
  if (!harvestId) {
    return NextResponse.json({ message: 'harvestId is required' }, { status: 400 });
  }

  const query = `
    SELECT id, harvest_id, batch_id, status, coffee_type, origin_code, intake_weight, intake_note, created_at, updated_at
    FROM production_batches
    WHERE harvest_id = $1
    ORDER BY created_at DESC
    LIMIT 1
  `;

  try {
    const result = await getPool().query(query, [harvestId]);
    return NextResponse.json({ batch: result.rows[0] ?? null });
  } catch (error) {
    console.error('[E2E][WAREHOUSE] query failed', error);
    return NextResponse.json({ message: 'failed to query warehouse batch' }, { status: 500 });
  }
}
