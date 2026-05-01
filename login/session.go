package login

import (
	"fmt"

	"github.com/Hucaru/Valhalla/common"
	"github.com/Hucaru/Valhalla/mnet"
)

type sessionStage byte

const (
	sessionStageAwaitLogin sessionStage = iota
	sessionStageAwaitEULA
	sessionStageAwaitPin
	sessionStageAwaitPinRegistration
	sessionStageAwaitWorldSelect
	sessionStageAwaitChannelSelect
	sessionStageAwaitCharacterSelect
	sessionStageMigrating
)

type session struct {
	stage sessionStage

	username      string
	accountID     int32
	worldID       byte
	channelID     byte
	viewAll       bool
	onlineMarked  bool
	migrationChar int32
}

func (server *Server) sessionFor(conn mnet.Client) *session {
	if server.sessions == nil {
		server.sessions = make(map[mnet.Client]*session)
	}

	if sess, ok := server.sessions[conn]; ok {
		return sess
	}

	sess := &session{stage: sessionStageAwaitLogin}
	server.sessions[conn] = sess
	return sess
}

func (server *Server) closeSession(conn mnet.Client) *session {
	if server.sessions == nil {
		return nil
	}

	sess := server.sessions[conn]
	delete(server.sessions, conn)
	delete(server.migrating, conn)
	return sess
}

func clientRemoteIP(conn fmt.Stringer) string {
	return common.RemoteIPFromConn(conn)
}

func (server *Server) validWorld(worldID byte) bool {
	return int(worldID) >= 0 && int(worldID) < len(server.worlds) && server.worlds[worldID].Conn != nil
}

func (server *Server) validChannel(worldID, channelID byte) bool {
	if !server.validWorld(worldID) {
		return false
	}
	return int(channelID) >= 0 && int(channelID) < len(server.worlds[worldID].Channels)
}

func (server *Server) firstAvailableChannel(worldID byte) (byte, bool) {
	if !server.validWorld(worldID) {
		return 0, false
	}

	for i, ch := range server.worlds[worldID].Channels {
		if ch.Port != 0 && len(ch.IP) == 4 {
			return byte(i), true
		}
	}

	return 0, false
}
