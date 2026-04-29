package login

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	"strings"
)

func generateCredentialSalt() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	return hex.EncodeToString(buf), nil
}

func hashSaltedMD5(value, salt string) string {
	sum := md5.Sum([]byte(salt + value))
	return hex.EncodeToString(sum[:])
}

func hashLegacyMD5(value string) string {
	sum := md5.Sum([]byte(value))
	return hex.EncodeToString(sum[:])
}

func hashLegacySHA512(value string) string {
	hasher := sha512.New()
	hasher.Write([]byte(value))
	return hex.EncodeToString(hasher.Sum(nil))
}

func makeStoredCredential(value string) (string, string, error) {
	salt, err := generateCredentialSalt()
	if err != nil {
		return "", "", err
	}

	return hashSaltedMD5(value, salt), salt, nil
}

func verifyStoredCredential(value, stored string, salt sql.NullString) bool {
	if salt.Valid && salt.String != "" {
		return strings.EqualFold(stored, hashSaltedMD5(value, salt.String))
	}

	if strings.EqualFold(stored, hashLegacyMD5(value)) {
		return true
	}

	if strings.EqualFold(stored, value) {
		return true
	}

	// Support existing unsalted SHA-512 rows during migration to salted hashes.
	return strings.EqualFold(stored, hashLegacySHA512(value))
}
