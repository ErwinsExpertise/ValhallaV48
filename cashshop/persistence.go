package cashshop

import (
	"database/sql"
	"fmt"

	"github.com/Hucaru/Valhalla/channel"
	"github.com/Hucaru/Valhalla/common"
)

func loadWishList(characterID int32) ([]int32, error) {
	rows, err := common.DB.Query(`
		SELECT slot, sn
		FROM cashshop_wishlist
		WHERE characterID=?
		ORDER BY slot ASC`, characterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	wish := make([]int32, 10)
	for rows.Next() {
		var slot int
		var sn int32
		if err := rows.Scan(&slot, &sn); err != nil {
			return nil, err
		}
		if slot >= 0 && slot < len(wish) {
			wish[slot] = sn
		}
	}

	return wish, rows.Err()
}

func saveWishList(characterID int32, sns []int32) error {
	tx, err := common.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if _, err = tx.Exec("DELETE FROM cashshop_wishlist WHERE characterID=?", characterID); err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO cashshop_wishlist(characterID, slot, sn) VALUES(?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for slot, sn := range sns {
		if sn == 0 {
			continue
		}
		if _, err = stmt.Exec(characterID, slot, sn); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func ensureCashShopStorageHeader(tx *sql.Tx, accountID int32) error {
	_, err := tx.Exec(`
		INSERT INTO account_cashshop_storage(accountID, slots)
		VALUES(?, ?)
		AS new ON DUPLICATE KEY UPDATE slots=account_cashshop_storage.slots`, accountID, cashShopStorageMinSlots)
	return err
}

func giftCashItemToAccount(accountID int32, item channel.Item) error {
	tx, err := common.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if err = ensureCashShopStorageHeader(tx, accountID); err != nil {
		return err
	}

	var maxSlots int
	if err = tx.QueryRow("SELECT slots FROM account_cashshop_storage WHERE accountID=?", accountID).Scan(&maxSlots); err != nil {
		return err
	}

	rows, err := tx.Query("SELECT slotNumber FROM account_cashshop_storage_items WHERE accountID=? ORDER BY slotNumber ASC", accountID)
	if err != nil {
		return err
	}

	used := make(map[int]bool)
	for rows.Next() {
		var slot int
		if scanErr := rows.Scan(&slot); scanErr != nil {
			rows.Close()
			return scanErr
		}
		used[slot] = true
	}
	if err = rows.Err(); err != nil {
		rows.Close()
		return err
	}
	rows.Close()

	freeSlot := 0
	for slot := 1; slot <= maxSlots; slot++ {
		if !used[slot] {
			freeSlot = slot
			break
		}
	}
	if freeSlot == 0 {
		return fmt.Errorf("cash shop locker full for account %d", accountID)
	}

	if err = item.SaveToCashShopStorage(tx, accountID, int16(freeSlot)); err != nil {
		return err
	}

	return tx.Commit()
}
