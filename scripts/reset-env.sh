#!/bin/bash
set -e

# Colors for better output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}🛑 Stopping and cleaning all containers/volumes...${NC}"
docker compose -f deployments/docker-compose.dev.yaml down -v --remove-orphans

echo -e "${BLUE}🚀 Starting infrastructure containers...${NC}"
docker compose -f deployments/docker-compose.dev.yaml up -d

echo -e "${BLUE}⏳ Waiting for Postgres (rr-postgres) to be ready...${NC}"
until docker exec rr-postgres pg_isready -U user -d postgres > /dev/null 2>&1; do
  echo -n "."
  sleep 1
done
echo -e "\n${GREEN}✅ Postgres is READY!${NC}"

echo -e "${BLUE}🏗️ Running all microservice migrations...${NC}"
task migrate:all

echo -e "${BLUE}🔑 Seeding infrastructure (Admin User & OAuth2 Clients)...${NC}"
task seed:infra

echo -e "${GREEN}✨ Environment reset complete!${NC}"
echo -e "You can now start backend services with: ${YELLOW}task be${NC}"
echo -e "And frontend with: ${YELLOW}task fe${NC}"
