package cassandra

import (
	"context"
	"time"

	"github.com/dungxbuif/RuntimeRoasters/apps/audit-service/internal/domain"
	"github.com/gocql/gocql"
)

type Store struct {
	session *gocql.Session
}

func NewStore(hosts []string, keyspace string) (*Store, error) {
	if len(hosts) == 0 {
		hosts = []string{"127.0.0.1:9042"}
	}
	if keyspace == "" {
		keyspace = "runtime_roasters_audit"
	}

	setup := gocql.NewCluster(hosts...)
	setup.Consistency = gocql.Quorum
	setup.Timeout = 10 * time.Second
	setup.ConnectTimeout = 10 * time.Second
	setupSession, err := setup.CreateSession()
	if err != nil {
		return nil, err
	}
	defer setupSession.Close()
	if err := setupSession.Query(`CREATE KEYSPACE IF NOT EXISTS ` + keyspace + ` WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}`).Exec(); err != nil {
		return nil, err
	}

	cluster := gocql.NewCluster(hosts...)
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum
	cluster.Timeout = 10 * time.Second
	cluster.ConnectTimeout = 10 * time.Second
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	store := &Store{session: session}
	if err := store.ensureSchema(); err != nil {
		session.Close()
		return nil, err
	}
	return store, nil
}

func (s *Store) Close() {
	if s != nil && s.session != nil {
		s.session.Close()
	}
}

func (s *Store) Insert(ctx context.Context, log domain.AuditLog) error {
	return s.session.Query(`
		INSERT INTO audit_logs (
			partition_key, occurred_at, id, message_id, topic, store_id, payload, previous_hash, current_hash, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		log.PartitionKey,
		log.OccurredAt,
		log.ID,
		log.MessageID,
		log.Topic,
		log.StoreID,
		log.Payload,
		log.PreviousHash,
		log.CurrentHash,
		log.CreatedAt,
	).WithContext(ctx).Exec()
}

func (s *Store) ListByPartition(ctx context.Context, partition string) ([]domain.AuditLog, error) {
	iter := s.session.Query(`
		SELECT id, partition_key, message_id, topic, store_id, payload, previous_hash, current_hash, occurred_at, created_at
		FROM audit_logs
		WHERE partition_key = ?`, partition).WithContext(ctx).Iter()
	var logs []domain.AuditLog
	var log domain.AuditLog
	for iter.Scan(&log.ID, &log.PartitionKey, &log.MessageID, &log.Topic, &log.StoreID, &log.Payload, &log.PreviousHash, &log.CurrentHash, &log.OccurredAt, &log.CreatedAt) {
		logs = append(logs, log)
	}
	return logs, iter.Close()
}

func (s *Store) ensureSchema() error {
	return s.session.Query(`
		CREATE TABLE IF NOT EXISTS audit_logs (
			partition_key text,
			occurred_at timestamp,
			id text,
			message_id text,
			topic text,
			store_id text,
			payload text,
			previous_hash text,
			current_hash text,
			created_at timestamp,
			PRIMARY KEY ((partition_key), occurred_at, id)
		) WITH CLUSTERING ORDER BY (occurred_at ASC, id ASC)`).Exec()
}
