# Kafka No-ZooKeeper Final Evidence

Date: 2026-05-23

Scope: verify the development Kafka stack no longer depends on ZooKeeper and can run as a single-node Apache Kafka broker.

## Configuration Evidence

- Compose file: `deployments/docker-compose.dev.yaml`
- Kafka image: `apache/kafka:4.3.0`
- Mode: single-node Apache Kafka broker with an internal controller listener
- External app port remains compatible: `localhost:9094`
- Internal Docker port remains compatible: `kafka:9092`
- No `KAFKA_ZOOKEEPER_CONNECT` references remain in the compose file, and Kafka has no ZooKeeper dependency.
- SigNoz/ClickHouse may define its own coordination service; that service is observability-only and not part of Kafka.

## Runtime Evidence

Command:

```bash
docker compose -f deployments/docker-compose.dev.yaml up -d kafka kafka-ui
```

Result:

```text
Container rr-kafka Running
Container rr-kafka-ui Running
Container rr-kafka Healthy
```

Command:

```bash
docker ps --format 'table {{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}'
```

Result excerpt:

```text
rr-kafka     apache/kafka:4.3.0     Up ... (healthy)     0.0.0.0:9094->9094/tcp
rr-kafka-ui  provectuslabs/kafka-ui:latest
```

Command:

```bash
docker compose -f deployments/docker-compose.dev.yaml config --services
```

Result excerpt:

```text
kafka
kafka-ui
krakend
hydra-migrate
postgres
hydra
kratos-migrate
kratos
identity
elasticsearch
kibana
seeder
valkey
cassandra
signoz-zookeeper-1
signoz-clickhouse
otel-collector
signoz
```

There is no Kafka ZooKeeper service or `rr-zookeeper` container. SigNoz/ClickHouse coordination, if enabled, is separate from Kafka.

## Kafka CLI Evidence

Command:

```bash
docker exec rr-kafka /opt/kafka/bin/kafka-topics.sh --bootstrap-server localhost:9092 --list
```

Result:

```text
__consumer_offsets
farm.harvest.events
payment.completed
```

Command:

```bash
docker exec rr-kafka /opt/kafka/bin/kafka-topics.sh --bootstrap-server localhost:9092 --create --if-not-exists --topic rr.final.kafka.evidence --partitions 1 --replication-factor 1
```

Result:

```text
Created topic rr.final.kafka.evidence.
```

Command:

```bash
docker exec rr-kafka /opt/kafka/bin/kafka-verifiable-producer.sh --bootstrap-server localhost:9092 --topic rr.final.kafka.evidence --max-messages 1
```

Result excerpt:

```text
"producer_send_success"
"sent":1
"acked":1
```

Command:

```bash
docker exec rr-kafka /opt/kafka/bin/kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic rr.final.kafka.evidence --from-beginning --max-messages 1
```

Result:

```text
0
Processed a total of 1 messages
```

## Screenshot Evidence

- Kafka UI screenshot: `docs/business/sprint-emergency-final-demo/evidence/kafka-no-zookeeper-ui.png`

## Verdict

PASS. The local development Kafka stack now runs on Apache Kafka 4.3.0 without ZooKeeper, and the broker accepts topic creation, message production, and message consumption.
