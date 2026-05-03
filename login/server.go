package login

import (
	"log"
	"time"

	"github.com/Hucaru/Valhalla/anticheat"
	"github.com/Hucaru/Valhalla/common"
	"github.com/Hucaru/Valhalla/internal"
	"github.com/Hucaru/Valhalla/mnet"
)

// Server state
type Server struct {
	migrating map[mnet.Client]bool
	sessions  map[mnet.Client]*session
	// db        *sql.DB
	worlds       []internal.World
	withPin      bool
	autoRegister bool
	ac           *anticheat.AntiCheat
}

// Initialise the server
func (server *Server) Initialise(dbConfig common.DBConfig, withpin bool, autoRegister bool) {
	server.migrating = make(map[mnet.Client]bool)
	server.sessions = make(map[mnet.Client]*session)
	server.withPin = withpin
	server.autoRegister = autoRegister

	err := common.ConnectToDB(dbConfig)

	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Connected to database")
	common.CleanupExpiredPendingMigrations()
	server.CleanupDB()

	log.Println("Cleaned up the database")

	// Initialize anti-cheat
	server.ac = anticheat.New(common.DB, nil)
	server.ac.StartCleanup()
	log.Println("Anti-cheat initialized")
}

// CleanupDB sets all accounts isLogedIn to 0
func (server *Server) CleanupDB() {
	now := time.Now().UnixMilli()
	if _, err := common.DB.Exec("UPDATE characters SET migrationID=-1 WHERE migrationID != -1 AND channelID = -1 AND inCashShop = 0"); err != nil {
		log.Println(err)
	}
	res, err := common.DB.Exec(`UPDATE accounts a
		SET a.isLogedIn = 0
		WHERE a.isLogedIn = 1
		AND NOT EXISTS (
			SELECT 1 FROM characters c
			WHERE c.accountID = a.accountID AND (c.channelID != -1 OR c.inCashShop = 1)
		)
		AND NOT EXISTS (
			SELECT 1 FROM pending_migrations pm
			WHERE pm.accountID = a.accountID AND pm.consumedAt = 0 AND pm.expiresAt >= ?
		)`, now)
	if err != nil {
		log.Fatal(err)
	}
	amount, _ := res.RowsAffected()
	log.Printf("Set %d isLogedin rows to 0.", amount)
}

// ServerDisconnected handler
func (server *Server) ServerDisconnected(conn mnet.Server) {
	for i, v := range server.worlds {
		if v.Conn == conn {
			log.Println(v.Name, "disconnected")
			server.worlds[i].Conn = nil
			server.worlds[i].Channels = []internal.Channel{}
			break
		}
	}
}

// ClientDisconnected from server
func (server *Server) ClientDisconnected(conn mnet.Client) {
	isMigrating := server.migrating[conn]
	sess := server.closeSession(conn)
	if isMigrating {
		delete(server.migrating, conn)
	} else if sess != nil && sess.onlineMarked {
		common.DeletePendingMigrationsForAccount(sess.accountID)
		_, err := common.DB.Exec("UPDATE accounts SET isLogedIn=0 WHERE accountID=?", sess.accountID)

		if err != nil {
			log.Println("Unable to complete logout for ", sess.accountID)
		}
	} else if conn.GetAccountID() > 0 {
		common.DeletePendingMigrationsForAccount(conn.GetAccountID())
		_, err := common.DB.Exec("UPDATE accounts SET isLogedIn=0 WHERE accountID=?", conn.GetAccountID())

		if err != nil {
			log.Println("Unable to complete logout for ", conn.GetAccountID())
		}
	}

	conn.Cleanup()
}
