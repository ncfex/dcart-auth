package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ncfex/dcart-auth/internal/domain/shared"
	"github.com/ncfex/dcart-auth/internal/domain/user"
)

type EventMetadata struct {
	ID            string          `json:"id"`
	AggregateID   string          `json:"aggregate_id"`
	AggregateType string          `json:"aggregate_type"`
	EventType     string          `json:"event_type"`
	Version       int             `json:"version"`
	Timestamp     time.Time       `json:"timestamp"`
	Payload       json.RawMessage `json:"payload"`
}

type PostgresEventStore struct {
	db *sql.DB
}

func NewPostgresEventStore(db *sql.DB) *PostgresEventStore {
	return &PostgresEventStore{db: db}
}

func (s *PostgresEventStore) SaveEvents(ctx context.Context, aggregateID string, events []shared.Event) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		SELECT 1 FROM events 
		WHERE aggregate_id = $1 
		FOR UPDATE`,
		aggregateID)
	if err != nil {
		return fmt.Errorf("lock aggregate events: %w", err)
	}

	var latestVersion int
	err = tx.QueryRowContext(ctx, `
		SELECT COALESCE(MAX(version), 0) 
		FROM events 
		WHERE aggregate_id = $1`,
		aggregateID).Scan(&latestVersion)
	if err != nil {
		return fmt.Errorf("get latest version: %w", err)
	}

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO events (
			id, 
			aggregate_id, 
			aggregate_type, 
			event_type, 
			version, 
			timestamp, 
			payload
		) VALUES ($1, $2, $3, $4, $5, $6, $7)`)
	if err != nil {
		return fmt.Errorf("prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, event := range events {
		payload, err := json.Marshal(event.GetPayload())
		if err != nil {
			return fmt.Errorf("marshal payload: %w", err)
		}

		expectedVersion := latestVersion + 1
		if event.GetVersion() != expectedVersion {
			return fmt.Errorf("concurrent modification detected: expected version %d, got %d",
				expectedVersion, event.GetVersion())
		}

		_, err = stmt.ExecContext(ctx,
			fmt.Sprintf("%s-%d", aggregateID, event.GetVersion()),
			event.GetAggregateID(),
			event.GetAggregateType(),
			event.GetEventType(),
			event.GetVersion(),
			event.GetTimestamp(),
			payload)
		if err != nil {
			return fmt.Errorf("insert event: %w", err)
		}

		latestVersion = event.GetVersion()
	}

	return tx.Commit()
}

func (s *PostgresEventStore) GetEvents(ctx context.Context, aggregateID string) ([]shared.Event, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT 
			aggregate_id, 
			aggregate_type, 
			event_type, 
			version, 
			timestamp, 
			payload 
		FROM events 
		WHERE aggregate_id = $1 
		ORDER BY version ASC`,
		aggregateID)
	if err != nil {
		return nil, fmt.Errorf("query events: %w", err)
	}
	defer rows.Close()

	return s.scanEvents(rows)
}

func (s *PostgresEventStore) GetEventsByType(ctx context.Context, eventType string) ([]shared.Event, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT 
			aggregate_id, 
			aggregate_type, 
			event_type, 
			version, 
			timestamp, 
			payload 
		FROM events 
		WHERE event_type = $1 
		ORDER BY timestamp ASC`,
		eventType)
	if err != nil {
		return nil, fmt.Errorf("query events by type: %w", err)
	}
	defer rows.Close()

	return s.scanEvents(rows)
}

func (s *PostgresEventStore) scanEvents(rows *sql.Rows) ([]shared.Event, error) {
	var events []shared.Event

	for rows.Next() {
		var metadata EventMetadata
		err := rows.Scan(
			&metadata.AggregateID,
			&metadata.AggregateType,
			&metadata.EventType,
			&metadata.Version,
			&metadata.Timestamp,
			&metadata.Payload,
		)
		if err != nil {
			return nil, fmt.Errorf("scan event: %w", err)
		}

		// todo improve payload marshal approicah
		var payload interface{}
		switch metadata.EventType {
		case "USER_REGISTERED":
			var p user.UserRegisteredEventPayload
			if err := json.Unmarshal(metadata.Payload, &p); err != nil {
				return nil, fmt.Errorf("unmarshal USER_REGISTERED payload: %w", err)
			}
			payload = p
		case "USER_PASSWORD_CHANGED":
			var p user.UserPasswordChangedEventPayload
			if err := json.Unmarshal(metadata.Payload, &p); err != nil {
				return nil, fmt.Errorf("unmarshal USER_PASSWORD_CHANGED payload: %w", err)
			}
			payload = p
		default:
			return nil, fmt.Errorf("unknown event type: %s", metadata.EventType)
		}

		event := shared.BaseEvent{
			AggregateID:   metadata.AggregateID,
			AggregateType: metadata.AggregateType,
			EventType:     metadata.EventType,
			Version:       metadata.Version,
			Timestamp:     metadata.Timestamp,
			Payload:       payload,
		}

		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %w", err)
	}

	return events, nil
}
