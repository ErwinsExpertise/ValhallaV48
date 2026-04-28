package channel

import (
	"fmt"
	"time"

	"github.com/Hucaru/Valhalla/constant"
)

const (
	zakumExpeditionMaxMembers   = 20
	zakumExpeditionSignupWindow = 10 * time.Minute
	zakumExpeditionBossPortal   = "st00"
	zakumExpeditionProp         = "zakumExpedition"
)

type zakumExpeditionMember struct {
	ID   int32
	Name string
}

type zakumExpedition struct {
	LeaderID        int32
	LeaderName      string
	Members         []zakumExpeditionMember
	SignupStartedAt time.Time
	ExpiresAt       time.Time
	Open            bool
	Started         bool
	Timer           *time.Timer
}

func zakumExpeditionResult(ok bool, msg string, handled bool) map[string]interface{} {
	return map[string]interface{}{
		"ok":      ok,
		"message": msg,
		"handled": handled,
	}
}

func (server *Server) sendZakumBlocked(plr *Player) {
	plr.Send(packetMessageRedText("The fight is already active. Please try again later."))
	plr.Send(packetPlayerNoChange())
}

func (server *Server) getZakumBossInstance(instID int) (*fieldInstance, error) {
	field, ok := server.fields[constant.MapBossZakum]
	if !ok {
		return nil, fmt.Errorf("zakum boss field not found")
	}

	inst, err := field.getInstance(instID)
	if err == nil {
		return inst, nil
	}

	return field.getInstance(0)
}

func (server *Server) zakumFightActive(instID int) bool {
	inst, err := server.getZakumBossInstance(instID)
	if err != nil {
		return false
	}

	bossActive, ok := inst.properties["eventActive"].(bool)
	if !ok {
		return false
	}

	return bossActive
}

func (server *Server) stopZakumExpeditionTimer(exp *zakumExpedition) {
	if exp != nil && exp.Timer != nil {
		exp.Timer.Stop()
		exp.Timer = nil
	}
}

func (server *Server) showZakumSignupCountdown(instID int, seconds int32) {
	inst, err := server.getZakumWaitingInstance(instID)
	if err != nil {
		return
	}
	inst.send(packetShowCountdown(seconds))
}

func (server *Server) hideZakumSignupCountdown(instID int) {
	inst, err := server.getZakumWaitingInstance(instID)
	if err != nil {
		return
	}
	inst.send(packetHideCountdown())
}

func (server *Server) getZakumExpeditionRemainingSeconds(instID int) int32 {
	exp := server.getZakumExpedition(instID)
	if exp == nil || !exp.Open {
		return 0
	}

	remaining := int32(time.Until(exp.ExpiresAt).Seconds())
	if remaining < 0 {
		return 0
	}

	return remaining
}

func (server *Server) getZakumWaitingInstance(instID int) (*fieldInstance, error) {
	field, ok := server.fields[constant.MapBossZakumWaiting]
	if !ok {
		return nil, fmt.Errorf("zakum waiting field not found")
	}

	inst, err := field.getInstance(instID)
	if err == nil {
		return inst, nil
	}

	return field.getInstance(0)
}

func (server *Server) getZakumExpedition(instID int) *zakumExpedition {
	inst, err := server.getZakumWaitingInstance(instID)
	if err != nil {
		return nil
	}

	exp, _ := inst.properties[zakumExpeditionProp].(*zakumExpedition)
	if exp != nil && exp.Started && !server.zakumFightActive(instID) {
		server.stopZakumExpeditionTimer(exp)
		delete(inst.properties, zakumExpeditionProp)
		return nil
	}

	return exp
}

func (server *Server) clearZakumExpedition(instID int) {
	inst, err := server.getZakumWaitingInstance(instID)
	if err != nil {
		return
	}

	if exp, ok := inst.properties[zakumExpeditionProp].(*zakumExpedition); ok {
		server.stopZakumExpeditionTimer(exp)
	}

	delete(inst.properties, zakumExpeditionProp)
	inst.send(packetHideCountdown())
}

func (server *Server) getZakumExpeditionSnapshot(instID int) map[string]interface{} {
	exp := server.getZakumExpedition(instID)
	if exp == nil {
		return map[string]interface{}{
			"exists": false,
		}
	}

	members := make([]map[string]interface{}, 0, len(exp.Members))
	for _, member := range exp.Members {
		members = append(members, map[string]interface{}{
			"id":       member.ID,
			"name":     member.Name,
			"isLeader": member.ID == exp.LeaderID,
		})
	}

	remaining := int32(time.Until(exp.ExpiresAt).Seconds())
	if remaining < 0 {
		remaining = 0
	}

	return map[string]interface{}{
		"exists":           true,
		"leaderID":         exp.LeaderID,
		"leaderName":       exp.LeaderName,
		"members":          members,
		"open":             exp.Open,
		"started":          exp.Started,
		"signupStartedAt":  exp.SignupStartedAt.UnixMilli(),
		"expiresAt":        exp.ExpiresAt.UnixMilli(),
		"remainingSeconds": remaining,
	}
}

func (server *Server) createZakumExpedition(plr *Player) map[string]interface{} {
	if server.zakumFightActive(plr.inst.id) {
		server.sendZakumBlocked(plr)
		return zakumExpeditionResult(false, "", true)
	}
	waitingInst, err := server.getZakumWaitingInstance(plr.inst.id)
	if err != nil {
		return zakumExpeditionResult(false, "The Zakum waiting map could not be found.", false)
	}

	if server.getZakumExpedition(plr.inst.id) != nil {
		return zakumExpeditionResult(false, "A Zakum expedition has already been registered.", false)
	}

	now := time.Now()
	exp := &zakumExpedition{
		LeaderID:        plr.ID,
		LeaderName:      plr.Name,
		Members:         []zakumExpeditionMember{{ID: plr.ID, Name: plr.Name}},
		SignupStartedAt: now,
		ExpiresAt:       now.Add(zakumExpeditionSignupWindow),
		Open:            true,
	}

	exp.Timer = time.AfterFunc(zakumExpeditionSignupWindow, func() {
		server.dispatch <- func() {
			server.startZakumExpedition(plr.inst.id, 0, true)
		}
	})

	waitingInst.properties[zakumExpeditionProp] = exp
	server.showZakumSignupCountdown(plr.inst.id, int32(zakumExpeditionSignupWindow/time.Second))

	return zakumExpeditionResult(true, "The Zakum expedition has been registered. You have 10 minutes to gather up to 20 members.", false)
}

func (server *Server) joinZakumExpedition(plr *Player) map[string]interface{} {
	if server.zakumFightActive(plr.inst.id) {
		server.sendZakumBlocked(plr)
		return zakumExpeditionResult(false, "", true)
	}

	exp := server.getZakumExpedition(plr.inst.id)
	if exp == nil {
		return zakumExpeditionResult(false, "There is no Zakum expedition accepting signups right now.", false)
	}

	if !exp.Open {
		return zakumExpeditionResult(false, "Registrations for the current expedition have already closed.", false)
	}

	for _, member := range exp.Members {
		if member.ID == plr.ID {
			return zakumExpeditionResult(false, "You are already registered for this expedition.", false)
		}
	}

	if len(exp.Members) >= zakumExpeditionMaxMembers {
		return zakumExpeditionResult(false, "The Zakum expedition is already full.", false)
	}

	exp.Members = append(exp.Members, zakumExpeditionMember{ID: plr.ID, Name: plr.Name})

	return zakumExpeditionResult(true, "You have joined the Zakum expedition. Please wait for the leader's signal.", false)
}

func (server *Server) leaveZakumExpedition(plr *Player) map[string]interface{} {
	exp := server.getZakumExpedition(plr.inst.id)
	if exp == nil {
		return zakumExpeditionResult(false, "There is no Zakum expedition to leave.", false)
	}

	if exp.LeaderID == plr.ID {
		return zakumExpeditionResult(false, "As the leader, you must disband the expedition instead.", false)
	}

	for i, member := range exp.Members {
		if member.ID == plr.ID {
			exp.Members = append(exp.Members[:i], exp.Members[i+1:]...)
			return zakumExpeditionResult(true, "You have been removed from the Zakum expedition signup list.", false)
		}
	}

	return zakumExpeditionResult(false, "You are not registered for the current expedition.", false)
}

func (server *Server) kickZakumExpeditionMember(plr *Player, memberID int32) map[string]interface{} {
	exp := server.getZakumExpedition(plr.inst.id)
	if exp == nil {
		return zakumExpeditionResult(false, "There is no Zakum expedition to manage.", false)
	}

	if exp.LeaderID != plr.ID {
		return zakumExpeditionResult(false, "Only the expedition leader may remove members.", false)
	}

	if exp.LeaderID == memberID {
		return zakumExpeditionResult(false, "The expedition leader cannot be removed.", false)
	}

	for i, member := range exp.Members {
		if member.ID != memberID {
			continue
		}

		exp.Members = append(exp.Members[:i], exp.Members[i+1:]...)
		if kicked, err := server.players.GetFromID(memberID); err == nil {
			kicked.Send(packetMessageRedText("You have been removed from the Zakum expedition."))
		}
		return zakumExpeditionResult(true, fmt.Sprintf("%s has been removed from the expedition.", member.Name), false)
	}

	return zakumExpeditionResult(false, "That player is not registered for the expedition.", false)
}

func (server *Server) terminateZakumExpedition(plr *Player) map[string]interface{} {
	waitingInst, err := server.getZakumWaitingInstance(plr.inst.id)
	if err != nil {
		return zakumExpeditionResult(false, "The Zakum waiting map could not be found.", false)
	}

	exp := server.getZakumExpedition(plr.inst.id)
	if exp == nil {
		return zakumExpeditionResult(false, "There is no Zakum expedition to disband.", false)
	}

	if exp.LeaderID != plr.ID {
		return zakumExpeditionResult(false, "Only the expedition leader may disband the expedition.", false)
	}

	server.stopZakumExpeditionTimer(exp)
	delete(waitingInst.properties, zakumExpeditionProp)
	waitingInst.send(packetHideCountdown())

	for _, member := range exp.Members {
		if target, err := server.players.GetFromID(member.ID); err == nil {
			target.Send(packetHideCountdown())
			target.Send(packetMessageRedText("The Zakum expedition has been disbanded."))
		}
	}

	return zakumExpeditionResult(true, "The Zakum expedition has been disbanded.", false)
}

func (server *Server) startZakumExpedition(instID int, leaderID int32, fromTimer bool) map[string]interface{} {
	if server.zakumFightActive(instID) {
		if leaderID != 0 {
			if leader, err := server.players.GetFromID(leaderID); err == nil {
				server.sendZakumBlocked(leader)
			}
		}

		exp := server.getZakumExpedition(instID)
		if exp != nil && fromTimer {
			for _, member := range exp.Members {
				if target, err := server.players.GetFromID(member.ID); err == nil && target.mapID == constant.MapBossZakumWaiting && target.inst != nil && target.inst.id == instID {
					server.sendZakumBlocked(target)
				}
			}
			server.clearZakumExpedition(instID)
		}

		return zakumExpeditionResult(false, "", true)
	}

	exp := server.getZakumExpedition(instID)
	if exp == nil {
		return zakumExpeditionResult(false, "There is no Zakum expedition ready to depart.", false)
	}

	if leaderID != 0 && exp.LeaderID != leaderID {
		return zakumExpeditionResult(false, "Only the expedition leader may start the expedition.", false)
	}

	server.stopZakumExpeditionTimer(exp)
	exp.Open = false
	exp.Started = true
	members := append([]zakumExpeditionMember(nil), exp.Members...)

	bossField, ok := server.fields[constant.MapBossZakum]
	if !ok {
		return zakumExpeditionResult(false, "The Zakum boss map could not be found.", false)
	}

	bossInst, err := bossField.getInstance(instID)
	if err != nil {
		bossInst, err = bossField.getInstance(0)
		if err != nil {
			return zakumExpeditionResult(false, "The Zakum boss instance could not be prepared.", false)
		}
	}

	portal, err := bossInst.getPortalFromName(zakumExpeditionBossPortal)
	if err != nil {
		portal, err = bossInst.getRandomSpawnPortal()
		if err != nil {
			return zakumExpeditionResult(false, "A valid Zakum entry portal could not be found.", false)
		}
	}

	playersToWarp := make([]*Player, 0, len(members))
	for _, member := range members {
		target, err := server.players.GetFromID(member.ID)
		if err != nil || target.inst == nil {
			continue
		}
		if target.mapID != constant.MapBossZakumWaiting || target.inst.id != instID {
			continue
		}
		playersToWarp = append(playersToWarp, target)
	}

	if len(playersToWarp) == 0 {
		server.clearZakumExpedition(instID)
		return zakumExpeditionResult(false, "There were no registered members left in the waiting area to enter Zakum.", false)
	}

	go manageSummonedBoss(bossInst, constant.MobZakum1Body, server)
	server.hideZakumSignupCountdown(instID)

	warped := 0
	for _, target := range playersToWarp {
		target.Send(packetHideCountdown())
		if err := server.warpPlayer(target, bossField, portal, true); err == nil {
			warped++
		}
	}

	return zakumExpeditionResult(true, fmt.Sprintf("The expedition is departing with %d member(s).", warped), false)
}

func (ctrl *scriptPlayerWrapper) GetZakumExpedition() map[string]interface{} {
	if ctrl.plr == nil || ctrl.plr.inst == nil {
		return map[string]interface{}{"exists": false}
	}

	return ctrl.server.getZakumExpeditionSnapshot(ctrl.plr.inst.id)
}

func (ctrl *scriptPlayerWrapper) CreateZakumExpedition() map[string]interface{} {
	if ctrl.plr == nil || ctrl.plr.inst == nil {
		return zakumExpeditionResult(false, "You are not currently in a valid map instance.", false)
	}

	return ctrl.server.createZakumExpedition(ctrl.plr)
}

func (ctrl *scriptPlayerWrapper) JoinZakumExpedition() map[string]interface{} {
	if ctrl.plr == nil || ctrl.plr.inst == nil {
		return zakumExpeditionResult(false, "You are not currently in a valid map instance.", false)
	}

	return ctrl.server.joinZakumExpedition(ctrl.plr)
}

func (ctrl *scriptPlayerWrapper) LeaveZakumExpedition() map[string]interface{} {
	if ctrl.plr == nil || ctrl.plr.inst == nil {
		return zakumExpeditionResult(false, "You are not currently in a valid map instance.", false)
	}

	return ctrl.server.leaveZakumExpedition(ctrl.plr)
}

func (ctrl *scriptPlayerWrapper) KickZakumExpeditionMember(memberID int32) map[string]interface{} {
	if ctrl.plr == nil || ctrl.plr.inst == nil {
		return zakumExpeditionResult(false, "You are not currently in a valid map instance.", false)
	}

	return ctrl.server.kickZakumExpeditionMember(ctrl.plr, memberID)
}

func (ctrl *scriptPlayerWrapper) StartZakumExpedition() map[string]interface{} {
	if ctrl.plr == nil || ctrl.plr.inst == nil {
		return zakumExpeditionResult(false, "You are not currently in a valid map instance.", false)
	}

	return ctrl.server.startZakumExpedition(ctrl.plr.inst.id, ctrl.plr.ID, false)
}

func (ctrl *scriptPlayerWrapper) TerminateZakumExpedition() map[string]interface{} {
	if ctrl.plr == nil || ctrl.plr.inst == nil {
		return zakumExpeditionResult(false, "You are not currently in a valid map instance.", false)
	}

	return ctrl.server.terminateZakumExpedition(ctrl.plr)
}
