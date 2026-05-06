package login

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	"strings"
)

type credentialScheme int

const (
	credentialSchemeUnknown credentialScheme = iota
	credentialSchemeSaltedSHA512
	credentialSchemeSaltedMD5
	credentialSchemeLegacyMD5
	credentialSchemePlaintext
	credentialSchemeLegacySHA512
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

func hashSaltedSHA512(value, salt string) string {
	hasher := sha512.New()
	hasher.Write([]byte(salt + value))
	return hex.EncodeToString(hasher.Sum(nil))
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

	return hashSaltedSHA512(value, salt), salt, nil
}

func verifyStoredCredential(value, stored string, salt sql.NullString) (bool, credentialScheme) {
	if salt.Valid && salt.String != "" {
		if strings.EqualFold(stored, hashSaltedSHA512(value, salt.String)) {
			return true, credentialSchemeSaltedSHA512
		}

		if strings.EqualFold(stored, hashSaltedMD5(value, salt.String)) {
			return true, credentialSchemeSaltedMD5
		}

		return false, credentialSchemeUnknown
	}

	if strings.EqualFold(stored, hashLegacyMD5(value)) {
		return true, credentialSchemeLegacyMD5
	}

	if strings.EqualFold(stored, value) {
		return true, credentialSchemePlaintext
	}

	// Support existing unsalted SHA-512 rows during migration to salted hashes.
	if strings.EqualFold(stored, hashLegacySHA512(value)) {
		return true, credentialSchemeLegacySHA512
	}

	return false, credentialSchemeUnknown
}
