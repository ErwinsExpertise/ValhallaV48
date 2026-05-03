package common

import (
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

// DB object used for queries
var DB *sql.DB

type DBTLSConfig struct {
	Mode       string
	CAFile     string
	CertFile   string
	KeyFile    string
	ServerName string
}

type DBConfig struct {
	Address  string
	Port     string
	User     string
	Password string
	Database string
	TLS      DBTLSConfig
}

// ConnectToDB - connect to a MySQL instance
func ConnectToDB(dbConfig DBConfig) error {
	var err error
	cfg := mysql.NewConfig()
	cfg.User = dbConfig.User
	cfg.Passwd = dbConfig.Password
	cfg.Net = "tcp"
	cfg.Addr = dbConfig.Address + ":" + dbConfig.Port
	cfg.DBName = dbConfig.Database
	cfg.ParseTime = true
	cfg.Loc = time.UTC
	cfg.Collation = "utf8mb4_0900_ai_ci"
	cfg.Params = map[string]string{
		"charset":   "utf8mb4",
		"time_zone": "'+00:00'",
	}

	tlsMode, err := configureMySQLTLS(dbConfig.TLS)
	if err != nil {
		return err
	}
	if tlsMode != "" {
		cfg.TLSConfig = tlsMode
	}

	DB, err = sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		return err
	}

	err = DB.Ping()

	if err != nil {
		return err
	}

	DB.SetMaxIdleConns(10)

	return nil
}

func configureMySQLTLS(cfg DBTLSConfig) (string, error) {
	mode := strings.TrimSpace(strings.ToLower(cfg.Mode))
	if mode == "" || mode == "false" || mode == "off" || mode == "disable" || mode == "disabled" || mode == "none" {
		return "", nil
	}

	switch mode {
	case "preferred":
		return "preferred", nil
	case "skip-verify", "insecure", "insecure-skip-verify":
		return "skip-verify", nil
	case "true", "required", "require", "verify-ca", "verify-full", "custom":
		if cfg.CAFile == "" && cfg.CertFile == "" && cfg.KeyFile == "" && cfg.ServerName == "" {
			return "true", nil
		}
		return registerCustomMySQLTLS(cfg)
	default:
		return "", fmt.Errorf("unsupported MySQL TLS mode %q", cfg.Mode)
	}
}

func registerCustomMySQLTLS(cfg DBTLSConfig) (string, error) {
	tlsConfig := &tls.Config{MinVersion: tls.VersionTLS12}
	if serverName := strings.TrimSpace(cfg.ServerName); serverName != "" {
		tlsConfig.ServerName = serverName
	}

	if caFile := strings.TrimSpace(cfg.CAFile); caFile != "" {
		pem, err := os.ReadFile(caFile)
		if err != nil {
			return "", fmt.Errorf("load MySQL TLS CA file %q: %w", caFile, err)
		}
		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(pem) {
			return "", fmt.Errorf("parse MySQL TLS CA file %q: no certificates found", caFile)
		}
		tlsConfig.RootCAs = pool
	}

	if certFile := strings.TrimSpace(cfg.CertFile); certFile != "" || strings.TrimSpace(cfg.KeyFile) != "" {
		keyFile := strings.TrimSpace(cfg.KeyFile)
		if certFile == "" || keyFile == "" {
			return "", fmt.Errorf("MySQL TLS client certificate requires both cert_file and key_file")
		}
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return "", fmt.Errorf("load MySQL TLS client certificate/key: %w", err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	sum := sha1.Sum([]byte(strings.Join([]string{cfg.Mode, cfg.CAFile, cfg.CertFile, cfg.KeyFile, cfg.ServerName}, "|")))
	name := fmt.Sprintf("valhalla-custom-%x", sum[:6])
	if err := mysql.RegisterTLSConfig(name, tlsConfig); err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			return "", fmt.Errorf("register MySQL TLS config: %w", err)
		}
	}
	return name, nil
}
