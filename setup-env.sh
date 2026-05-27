#!/usr/bin/env bash
set -e

echo "Cleaning up stray processes..."
pkill -9 -f next || true
pkill -9 -f node || true
pkill -9 -f "go run" || true
killall auth-service farm-service retail-service logistics-service payment-service trace-service audit-service warehouse-service socket-service 2>/dev/null || true

echo "Starting Docker Compose Environment..."
cd /Users/dungxbuif/workspace/RuntimeRoasters
docker compose -f deployments/docker-compose.dev.yaml up -d

echo "Running Database Migrations..."
task migrate:all

echo "Starting Go Backend Services..."
cd src

SERVICES=("auth-service" "farm-service" "warehouse-service" "retail-service" "payment-service" "logistics-service" "trace-service" "audit-service" "socket-service")

mkdir -p bin
for s in "${SERVICES[@]}"; do
    if [ -d "apps/$s" ]; then
        echo "Building $s..."
        go build -o bin/$s ./apps/$s/cmd
    fi
done

for s in "${SERVICES[@]}"; do
    if [ -f "bin/$s" ] && [ -d "apps/$s" ]; then
        echo "Starting $s..."
        ( cd apps/$s && ../../bin/$s > "../../bin/$s.log" 2>&1 ) &
    fi
done

echo "Building NextJS..."
cd apps/client-app
npm run build

echo "Starting NextJS in background..."
npm start > next-prod.log 2>&1 &
NEXTJS_PID=$!

echo "Waiting 10 seconds for services to fully start..."
sleep 10

echo "Running E2E Specs..."
npx playwright test e2e/saga-branches.spec.ts e2e/manual-live-evidence.spec.ts

echo "E2E complete. Shutting down..."
kill -9 $NEXTJS_PID
cd ../../
for s in "${SERVICES[@]}"; do
    killall $s 2>/dev/null || true
done
