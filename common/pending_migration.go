package common

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

const (
	MigrationTypeChannel  = "channel"
	MigrationTypeCashShop = "cashshop"
	defaultMigrationTTL   = 30 * time.Second
)

var ErrPendingMigrationIPMismatch = errors.New("pending migration ip mismatch")

type PendingMigration struct {
	ID              int64
	AccountID       int32
	CharacterID     int32
	WorldID         byte
	DestinationType string
	DestinationID   int
	ClientIP        string
	Nonce           string
	CreatedAt       int64
	ExpiresAt       int64
	ConsumedAt      int64
}

func CleanupExpiredPendingMigrations() {
	if DB == nil {
		return
	}

	now := time.Now().UnixMilli()
	if _, err := DB.Exec("DELETE FROM pending_migrations WHERE expiresAt < ? OR (consumedAt > 0 AND consumedAt < ?)", now, now-int64((5*time.Minute)/time.Millisecond)); err != nil {
		log.Printf("cleanup pending migrations failed: %v", err)
	}
}

func CreatePendingMigration(accountID, characterID int32, worldID byte, destinationType string, destinationID int, clientIP string, ttl time.Duration) (*PendingMigration, error) {
	if DB == nil {
		return nil, errors.New("database unavailable")
	}
	if ttl <= 0 {
		ttl = defaultMigrationTTL
	}

	nonce, err := randomHex(16)
	if err != nil {
		return nil, err
	}

	now := time.Now().UnixMilli()
	expiresAt := now + ttl.Milliseconds()
	clientIP = strings.TrimSpace(clientIP)

	if _, err := DB.Exec("DELETE FROM pending_migrations WHERE accountID=? OR characterID=?", accountID, characterID); err != nil {
		return nil, err
	}

	res, err := DB.Exec(`INSERT INTO pending_migrations
		(accountID, characterID, worldID, destinationType, destinationID, clientIP, nonce, createdAt, expiresAt, consumedAt)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 0)`,
		accountID, characterID, worldID, destinationType, destinationID, clientIP, nonce, now, expiresAt)
	if err != nil {
		return nil, err
	}

	id, _ := res.LastInsertId()
	return &PendingMigration{
		ID:              id,
		AccountID:       accountID,
		CharacterID:     characterID,
		WorldID:         worldID,
		DestinationType: destinationType,
		DestinationID:   destinationID,
		ClientIP:        clientIP,
		Nonce:           nonce,
		CreatedAt:       now,
		ExpiresAt:       expiresAt,
	}, nil
}

func ConsumePendingMigration(characterID int32, destinationType string, destinationID int, clientIP string) (*PendingMigration, error) {
	if DB == nil {
		return nil, errors.New("database unavailable")
	}

	tx, err := DB.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if tx != nil {
			_ = tx.Rollback()
		}
	}()

	now := time.Now().UnixMilli()
	var row PendingMigration
	err = tx.QueryRow(`SELECT id, accountID, characterID, worldID, destinationType, destinationID, clientIP, nonce, createdAt, expiresAt, consumedAt
		FROM pending_migrations
		WHERE characterID=? AND destinationType=? AND destinationID=? AND consumedAt=0 AND expiresAt>=?
		ORDER BY createdAt DESC LIMIT 1 FOR UPDATE`, characterID, destinationType, destinationID, now).
		Scan(&row.ID, &row.AccountID, &row.CharacterID, &row.WorldID, &row.DestinationType, &row.DestinationID, &row.ClientIP, &row.Nonce, &row.CreatedAt, &row.ExpiresAt, &row.ConsumedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	clientIP = strings.TrimSpace(clientIP)
	if row.ClientIP != "" && clientIP != "" && !strings.EqualFold(row.ClientIP, clientIP) {
		return nil, ErrPendingMigrationIPMismatch
	}

	res, err := tx.Exec("UPDATE pending_migrations SET consumedAt=? WHERE id=? AND consumedAt=0", now, row.ID)
	if err != nil {
		return nil, err
	}
	if affected, _ := res.RowsAffected(); affected != 1 {
		return nil, fmt.Errorf("pending migration %d already consumed", row.ID)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	tx = nil
	row.ConsumedAt = now
	return &row, nil
}

func DeletePendingMigrationsForAccount(accountID int32) {
	if DB == nil || accountID <= 0 {
		return
	}
	if _, err := DB.Exec("DELETE FROM pending_migrations WHERE accountID=?", accountID); err != nil {
		log.Printf("delete pending migrations for account %d failed: %v", accountID, err)
	}
}

func DeletePendingMigrationForCharacter(characterID int32) {
	if DB == nil || characterID <= 0 {
		return
	}
	if _, err := DB.Exec("DELETE FROM pending_migrations WHERE characterID=?", characterID); err != nil {
		log.Printf("delete pending migration for character %d failed: %v", characterID, err)
	}
}

func ReconcileAccountLoginState(accountID int32) (bool, error) {
	if DB == nil || accountID <= 0 {
		return false, nil
	}

	now := time.Now().UnixMilli()
	var activeCharacters int
	if err := DB.QueryRow("SELECT COUNT(*) FROM characters WHERE accountID=? AND (channelID != -1 OR inCashShop = 1)", accountID).Scan(&activeCharacters); err != nil {
		return false, err
	}

	var activeMigrations int
	if err := DB.QueryRow("SELECT COUNT(*) FROM pending_migrations WHERE accountID=? AND consumedAt=0 AND expiresAt>=?", accountID, now).Scan(&activeMigrations); err != nil {
		return false, err
	}

	if activeCharacters > 0 || activeMigrations > 0 {
		return true, nil
	}

	if _, err := DB.Exec("UPDATE accounts SET isLogedIn=0 WHERE accountID=?", accountID); err != nil {
		return false, err
	}
	if _, err := DB.Exec("UPDATE characters SET migrationID=-1, inCashShop=0 WHERE accountID=?", accountID); err != nil {
		return false, err
	}
	DeletePendingMigrationsForAccount(accountID)
	return false, nil
}

func randomHex(size int) (string, error) {
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
