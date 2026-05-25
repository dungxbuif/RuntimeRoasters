#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
COMPOSE_FILE="$ROOT_DIR/deployments/docker-compose.dev.yaml"

topics=(
  retail.order.created
  payment.intent.created
  payment.completed
  payment.simulated_completed
  payment.failed
  payment.refunded
  warehouse.stock.reserved
  warehouse.stock.reservation_failed
  warehouse.inventory.updated
  logistics.delivery.assigned
  logistics.delivery.completed
  logistics.gps.updated
)

for topic in "${topics[@]}"; do
  docker exec rr-kafka kafka-topics --bootstrap-server localhost:9092 --delete --topic "$topic" >/dev/null 2>&1 || true
done

truncate_table() {
  local db="$1"
  local table="$2"
  docker exec rr-postgres psql -U user -d "$db" -c "TRUNCATE TABLE $table RESTART IDENTITY CASCADE" >/dev/null 2>&1 || true
}

truncate_table retail_db orders
truncate_table retail_db outbox_events
truncate_table retail_db inbox_events
truncate_table payment_db payments
truncate_table payment_db inbox_events
truncate_table logistics_db shipments
truncate_table logistics_db inbox_events
truncate_table trace_db trace_events
truncate_table trace_db trace_documents
truncate_table audit_db audit_logs

curl -fsS -X DELETE "http://localhost:9200/coffee_traceability" >/dev/null 2>&1 || true
docker exec rr-cassandra cqlsh -e "DROP KEYSPACE IF EXISTS runtime_roasters_audit" >/dev/null 2>&1 || true

docker compose -f "$COMPOSE_FILE" restart krakend >/dev/null 2>&1 || true

echo "RuntimeRoasters demo state reset."
