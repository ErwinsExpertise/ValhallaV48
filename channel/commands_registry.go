package channel

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Hucaru/Valhalla/constant"
	"github.com/Hucaru/Valhalla/internal"
	"github.com/Hucaru/Valhalla/mnet"
	"github.com/Hucaru/Valhalla/mpacket"
	"github.com/Hucaru/Valhalla/nx"
	"github.com/dop251/goja"
)

type staffRank byte

const (
	staffRankNone staffRank = iota
	staffRankCommunity
	staffRankSupport
	staffRankGM
	staffRankAdmin
)

func rankFromAdminLevel(level int) staffRank {
	switch {
	case level >= int(staffRankAdmin):
		return staffRankAdmin
	case level == int(staffRankGM):
		return staffRankGM
	case level == int(staffRankSupport):
		return staffRankSupport
	case level >= int(staffRankCommunity):
		return staffRankCommunity
	default:
		return staffRankNone
	}
}

func (r staffRank) String() string {
	switch r {
	case staffRankCommunity:
		return "Community"
	case staffRankSupport:
		return "Support"
	case staffRankGM:
		return "Game Master"
	case staffRankAdmin:
		return "Admin"
	default:
		return "None"
	}
}

type staffCommandHandler func(*staffCommandContext, []string) error

type staffCommand struct {
	name    string
	aliases []string
	minRank staffRank
	usage   string
	summary string
	audit   bool
	handler staffCommandHandler
}

type staffCommandContext struct {
	server  *Server
	conn    mnet.Client
	command *staffCommand
	caller  *Player
}

var errCommandUsage = errors.New("command usage")

func (ctx *staffCommandContext) sender() (*Player, error) {
	if ctx.caller != nil {
		return ctx.caller, nil
	}
	plr, err := ctx.server.players.GetFromConn(ctx.conn)
	if err != nil {
		return nil, err
	}
	ctx.caller = plr
	return plr, nil
}

func (ctx *staffCommandContext) sendNotice(msg string) {
	ctx.conn.Send(packetMessageNotice(msg))
}

func (ctx *staffCommandContext) sendError(msg string) {
	ctx.conn.Send(packetMessageRedText(msg))
}

func (ctx *staffCommandContext) sendUsage() {
	if ctx.command != nil && ctx.command.usage != "" {
		ctx.sendError("Usage: " + ctx.command.usage)
	}
}

func (ctx *staffCommandContext) commandRank() staffRank {
	return rankFromAdminLevel(ctx.conn.GetAdminLevel())
}

func (ctx *staffCommandContext) findPlayerByName(name string) (*Player, error) {
	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("player name was not provided")
	}
	return ctx.server.players.GetFromName(name)
}

func (ctx *staffCommandContext) logCommand(args []string, target *Player) {
	caller, err := ctx.sender()
	if err != nil {
		return
	}
	parts := []string{
		fmt.Sprintf("rank=%s", ctx.commandRank()),
		fmt.Sprintf("staff=%s(%d)", caller.Name, caller.ID),
		fmt.Sprintf("account=%d", caller.accountID),
		fmt.Sprintf("channel=%d", ctx.server.id+1),
		fmt.Sprintf("map=%d", caller.mapID),
		fmt.Sprintf("command=/%s", ctx.command.name),
	}
	if len(args) > 0 {
		parts = append(parts, fmt.Sprintf("args=%q", strings.Join(args, " ")))
	}
	if target != nil {
		parts = append(parts, fmt.Sprintf("target=%s(%d)", target.Name, target.ID))
	}
	log.Printf("staff-command %s", strings.Join(parts, " "))
}

func commandSpec(name string, minRank staffRank, usage, summary string, audit bool, handler staffCommandHandler, aliases ...string) staffCommand {
	return staffCommand{
		name:    name,
		aliases: aliases,
		minRank: minRank,
		usage:   usage,
		summary: summary,
		audit:   audit,
		handler: handler,
	}
}

var staffCommands map[string]*staffCommand

func init() {
	staffCommands = buildStaffCommandRegistry()
}

func buildStaffCommandRegistry() map[string]*staffCommand {
	defs := []staffCommand{
		commandSpec("help", staffRankCommunity, "/help [command]", "Show available staff commands", false, handleHelpCommand, "commands"),
		commandSpec("online", staffRankCommunity, "/online [channel|world]", "Show online player counts", false, handleOnlineCommand),
		commandSpec("who", staffRankCommunity, "/who [channel|world]", "List online players", false, handleWhoCommand),
		commandSpec("map", staffRankCommunity, "/map", "Show your current map information", false, handleMapCommand, "whereami"),
		commandSpec("where", staffRankSupport, "/where <player>", "Show a player's location in this channel", false, handleWhereCommand),
		commandSpec("search", staffRankCommunity, "/search [type] <query>", "Search items, maps, and quests", false, handleSearchCommand),
		commandSpec("showRates", staffRankCommunity, "/showRates", "Show current rate multipliers", false, handleShowRatesCommand),
		commandSpec("pos", staffRankSupport, "/pos", "Show your current position", false, handlePosCommand),
		commandSpec("notice", staffRankCommunity, "/notice <message>", "Send a notice to the current channel", false, handleChannelNoticeCommand),
		commandSpec("eventNotice", staffRankCommunity, "/eventNotice <message>", "Send an event notice to the current channel", false, handleEventNoticeCommand),
		commandSpec("noticePlayer", staffRankSupport, "/noticePlayer <player> <message>", "Send a private notice to a player in this channel", false, handleNoticePlayerCommand),
		commandSpec("msgBox", staffRankCommunity, "/msgBox <message>", "Show a dialogue box message in the current channel", false, handleMsgBoxCommand),
		commandSpec("header", staffRankCommunity, "/header [message]", "Set or clear the channel scrolling header", false, handleHeaderCommand),
		commandSpec("changeBgm", staffRankCommunity, "/changeBgm [path]", "Set or clear the current map BGM override", false, handleChangeBgmCommand),
		commandSpec("wrong", staffRankCommunity, "/wrong", "Play the failure map effect", false, handleWrongCommand),
		commandSpec("clear", staffRankCommunity, "/clear", "Play the clear map effect", false, handleClearCommand),
		commandSpec("gate", staffRankCommunity, "/gate", "Play the gate portal effect", false, handleGateCommand),
		commandSpec("createInstance", staffRankCommunity, "/createInstance", "Create a new instance for the current map", false, handleCreateInstanceCommand),
		commandSpec("changeInstance", staffRankGM, "/changeInstance <instanceID>", "Move yourself to another instance of the current map", true, handleChangeInstanceCommand),
		commandSpec("mapInfo", staffRankSupport, "/mapInfo", "Show instance details for the current map", false, handleMapInfoCommand),
		commandSpec("eventStart", staffRankCommunity, "/eventStart <script> <instanceID>", "Start an event script for your party or character", false, handleEventStartCommand),
		commandSpec("events", staffRankCommunity, "/events", "List running scripted events", false, handleEventsCommand),
		commandSpec("unstuck", staffRankSupport, "/unstuck [player]", "Warp a player to the nearest safe spawn in the current map", true, handleUnstuckCommand),
		commandSpec("revive", staffRankSupport, "/revive [player]", "Restore a player's HP to full", false, handleReviveCommand),
		commandSpec("questFinish", staffRankSupport, "/questFinish <questID>", "Mark a quest complete", false, handleQuestFinishCommand),
		commandSpec("questUntil", staffRankSupport, "/questUntil <questID> <part>", "Advance a quest record to a specific part", false, handleQuestUntilCommand),
		commandSpec("questReset", staffRankSupport, "/questReset <questID>", "Reset a quest record", false, handleQuestResetCommand),
		commandSpec("clearInstProps", staffRankSupport, "/clearInstProps", "Clear current instance properties", false, handleClearInstPropsCommand),
		commandSpec("properties", staffRankSupport, "/properties", "List current instance properties", false, handlePropertiesCommand),
		commandSpec("kill", staffRankGM, "/kill [player]", "Kill yourself or another player in this channel", true, handleKillCommand),
		commandSpec("warp", staffRankGM, "/warp <mapName|mapID>", "Warp yourself to a map by name or ID", true, handleWarpMapCommand),
		commandSpec("warpTo", staffRankGM, "/warpTo <player>", "Warp yourself to a player's exact location in this channel", true, handleWarpToCommand),
		commandSpec("warpToMe", staffRankGM, "/warpToMe <player>", "Warp a player to your exact location in this channel", true, handleWarpToMeCommand),
		commandSpec("clearDrops", staffRankGM, "/clearDrops", "Clear all drops in the current map instance", true, handleClearDropsCommand),
		commandSpec("removeTimer", staffRankGM, "/removeTimer", "Remove the current field timer", true, handleRemoveTimerCommand),
		commandSpec("killMob", staffRankGM, "/killMob <spawnID>", "Kill a mob by spawn ID", true, handleKillMobCommand),
		commandSpec("killAll", staffRankGM, "/killAll", "Kill all mobs in the current map instance", true, handleKillAllCommand, "killmobs"),
		commandSpec("spawn", staffRankGM, "/spawn <mobID> [count]", "Spawn mobs at your current position", true, handleSpawnCommand, "spawnMob"),
		commandSpec("spawnBoss", staffRankGM, "/spawnBoss <name> [count]", "Spawn a named boss preset", true, handleSpawnBossCommand),
		commandSpec("ban", staffRankGM, "/ban <player> [hours|perm] [reason]", "Ban an online player", true, handleBanCommand),
		commandSpec("unban", staffRankGM, "/unban <player>", "Remove all bans for a player", true, handleUnbanCommand),
		commandSpec("banhistory", staffRankGM, "/banhistory <player>", "Show recent ban history for a player", false, handleBanHistoryCommand),
		commandSpec("enablePortal", staffRankGM, "/enablePortal <name> <true|false>", "Enable or disable a portal in the current instance", true, handleEnablePortalCommand),
		commandSpec("deleteInstance", staffRankGM, "/deleteInstance <instanceID>", "Delete a non-zero instance from the current map", true, handleDeleteInstanceCommand),
		commandSpec("drop", staffRankAdmin, "/drop <itemID> [quantity]", "Drop a specific item at your current position", true, handleDropCommand),
		commandSpec("droptest", staffRankAdmin, "/droptest", "Drop the test loot set at your current position", true, handleDropTestCommand),
		commandSpec("dropr", staffRankGM, "/dropr <dropID>", "Remove a field drop by drop ID", true, handleRemoveDropCommand, "removeDrop"),
		commandSpec("item", staffRankAdmin, "/item <itemID> [quantity]", "Generate a specific item", true, handleItemCommand),
		commandSpec("loadout", staffRankAdmin, "/loadout", "Generate the test loadout", true, handleLoadoutCommand),
		commandSpec("exp", staffRankAdmin, "/exp [player] <amount>", "Set a player's EXP total", true, handleExpCommand),
		commandSpec("gexp", staffRankAdmin, "/gexp [player] <amount>", "Grant EXP to a player", true, handleGrantExpCommand),
		commandSpec("mesos", staffRankAdmin, "/mesos <amount>", "Set your mesos total", true, handleMesosCommand),
		commandSpec("nx", staffRankAdmin, "/nx <amount>", "Add NX to yourself", true, handleNxCommand),
		commandSpec("maplepoints", staffRankAdmin, "/maplepoints <amount>", "Add Maple Points to yourself", true, handleMaplePointsCommand),
		commandSpec("hp", staffRankAdmin, "/hp [player] <amount>", "Set HP", true, handleHPCommand),
		commandSpec("mp", staffRankAdmin, "/mp [player] <amount>", "Set MP", true, handleMPCommand),
		commandSpec("setMaxHP", staffRankAdmin, "/setMaxHP <amount>", "Set your max HP", true, handleSetMaxHPCommand),
		commandSpec("setMaxMP", staffRankAdmin, "/setMaxMP <amount>", "Set your max MP", true, handleSetMaxMPCommand),
		commandSpec("str", staffRankAdmin, "/str [player] <amount>", "Set STR", true, handlePrimaryStatCommand),
		commandSpec("dex", staffRankAdmin, "/dex [player] <amount>", "Set DEX", true, handlePrimaryStatCommand),
		commandSpec("int", staffRankAdmin, "/int [player] <amount>", "Set INT", true, handlePrimaryStatCommand),
		commandSpec("luk", staffRankAdmin, "/luk [player] <amount>", "Set LUK", true, handlePrimaryStatCommand),
		commandSpec("ap", staffRankAdmin, "/ap [player] <amount>", "Set AP", true, handleAPCommand),
		commandSpec("sp", staffRankAdmin, "/sp [player] <amount>", "Set SP", true, handleSPCommand),
		commandSpec("level", staffRankAdmin, "/level [player] <level>", "Set a player's level", true, handleLevelCommand, "setLevel"),
		commandSpec("levelup", staffRankAdmin, "/levelup [player] [amount]", "Increase a player's level", true, handleLevelUpCommand),
		commandSpec("job", staffRankAdmin, "/job <jobName|jobID>", "Set your job", true, handleJobCommand),
		commandSpec("skillLv", staffRankAdmin, "/skillLv <skillID> <level|max>", "Set a player's skill level", true, handleSkillLevelCommand),
		commandSpec("maxSkills", staffRankAdmin, "/maxSkills", "Max all player skills", true, handleMaxSkillsCommand),
		commandSpec("resetSkills", staffRankAdmin, "/resetSkills", "Clear all player skills", true, handleResetSkillsCommand),
		commandSpec("rate", staffRankAdmin, "/rate <exp|drop|mesos> <rate>", "Set world rates", true, handleRateCommand),
		commandSpec("setWorldMessage", staffRankAdmin, "/setWorldMessage <ribbon> [message]", "Set the login world message", true, handleSetWorldMessageCommand),
		commandSpec("partyCreate", staffRankAdmin, "/partyCreate", "Create a party", false, handlePartyCreateCommand),
		commandSpec("guildCreate", staffRankAdmin, "/guildCreate", "Start the guild creation NPC flow", false, handleGuildCreateCommand),
		commandSpec("guildDisband", staffRankAdmin, "/guildDisband", "Disband your guild", true, handleGuildDisbandCommand),
		commandSpec("guildPoints", staffRankAdmin, "/guildPoints <amount>", "Set guild points", true, handleGuildPointsCommand),
		commandSpec("testMob", staffRankAdmin, "/testMob", "Spawn the test mob", true, handleTestMobCommand),
		commandSpec("reloadScripts", staffRankAdmin, "/reloadScripts", "Reload channel script stores from disk", true, handleReloadScriptsCommand),
		commandSpec("saveAll", staffRankAdmin, "/saveAll", "Flush all player state to persistent storage", true, handleSaveAllCommand),
		commandSpec("packet", staffRankAdmin, "/packet <hex>", "Send a raw packet to yourself", true, handlePacketCommand),
	}

	registry := make(map[string]*staffCommand, len(defs)*2)
	for i := range defs {
		def := &defs[i]
		registry[strings.ToLower(def.name)] = def
		for _, alias := range def.aliases {
			registry[strings.ToLower(alias)] = def
		}
	}

	return registry
}

func uniqueAccessibleCommands(rank staffRank) []*staffCommand {
	seen := make(map[string]bool)
	defs := make([]*staffCommand, 0, len(staffCommands))
	for _, def := range staffCommands {
		if seen[def.name] || rank < def.minRank {
			continue
		}
		seen[def.name] = true
		defs = append(defs, def)
	}
	sort.Slice(defs, func(i, j int) bool {
		if defs[i].minRank != defs[j].minRank {
			return defs[i].minRank < defs[j].minRank
		}
		return defs[i].name < defs[j].name
	})
	return defs
}

func (server *Server) gmCommand(conn mnet.Client, msg string) {
	commandText := strings.TrimSpace(strings.TrimPrefix(msg, "/"))
	parts := strings.Fields(commandText)
	if len(parts) == 0 {
		return
	}

	name := strings.ToLower(parts[0])
	def, ok := staffCommands[name]
	if !ok {
		conn.Send(packetMessageRedText("Unknown command. Use /help to list available commands."))
		return
	}

	rank := rankFromAdminLevel(conn.GetAdminLevel())
	if rank < def.minRank {
		conn.Send(packetMessageRedText("You do not have permission to use this command."))
		return
	}

	ctx := &staffCommandContext{server: server, conn: conn, command: def}
	err := def.handler(ctx, parts[1:])
	if err == nil {
		return
	}
	if errors.Is(err, errCommandUsage) {
		ctx.sendUsage()
		return
	}
	ctx.sendError(err.Error())
}

func handleHelpCommand(ctx *staffCommandContext, args []string) error {
	rank := ctx.commandRank()
	if len(args) == 1 {
		def, ok := staffCommands[strings.ToLower(args[0])]
		if !ok || rank < def.minRank {
			return fmt.Errorf("Unknown command. Use /help to list available commands.")
		}
		line := fmt.Sprintf("/%s - %s", def.name, def.summary)
		ctx.sendNotice(line)
		ctx.sendNotice(fmt.Sprintf("Minimum rank: %s", def.minRank))
		if def.usage != "" {
			ctx.sendNotice("Usage: " + def.usage)
		}
		if len(def.aliases) > 0 {
			ctx.sendNotice("Aliases: " + strings.Join(def.aliases, ", "))
		}
		return nil
	}
	if len(args) > 1 {
		return errCommandUsage
	}

	defs := uniqueAccessibleCommands(rank)
	groups := []staffRank{staffRankCommunity, staffRankSupport, staffRankGM, staffRankAdmin}
	for _, group := range groups {
		names := make([]string, 0)
		for _, def := range defs {
			if def.minRank == group {
				names = append(names, "/"+def.name)
			}
		}
		if len(names) == 0 {
			continue
		}
		ctx.sendNotice(fmt.Sprintf("%s: %s", group, strings.Join(names, ", ")))
	}
	ctx.sendNotice("Use /help <command> for usage details.")
	return nil
}

func handleOnlineCommand(ctx *staffCommandContext, args []string) error {
	if len(args) > 1 {
		return errCommandUsage
	}
	scope := "channel"
	if len(args) == 1 {
		scope = strings.ToLower(args[0])
	}

	switch scope {
	case "channel":
		ctx.sendNotice(fmt.Sprintf("Online players in channel %d: %d", ctx.server.id+1, ctx.server.players.count()))
	case "world":
		total := 0
		for i, ch := range ctx.server.channels {
			if i == int(ctx.server.id) {
				total += ctx.server.players.count()
				continue
			}
			total += int(ch.Pop)
		}
		ctx.sendNotice(fmt.Sprintf("Online players in world %s: %d", ctx.server.worldName, total))
	default:
		return errCommandUsage
	}

	return nil
}

func handleWhoCommand(ctx *staffCommandContext, args []string) error {
	if len(args) > 1 {
		return errCommandUsage
	}
	scope := "channel"
	if len(args) == 1 {
		scope = strings.ToLower(args[0])
	}

	if scope == "world" {
		ctx.sendNotice("World-wide player lists are not available from channel servers yet. Showing the current channel instead.")
	}
	if scope != "channel" && scope != "world" {
		return errCommandUsage
	}

	entries := make([]string, 0, ctx.server.players.count())
	ctx.server.players.observe(func(plr *Player) {
		entries = append(entries, formatPlayerLocationLine(plr))
	})
	sort.Strings(entries)

	ctx.sendNotice(fmt.Sprintf("Players in channel %d: %d", ctx.server.id+1, len(entries)))
	if len(entries) == 0 {
		return nil
	}
	for _, line := range splitCommandOutput(entries, 3) {
		ctx.sendNotice("Players: " + strings.Join(line, ", "))
	}
	return nil
}

func handleMapCommand(ctx *staffCommandContext, args []string) error {
	if len(args) != 0 {
		return errCommandUsage
	}
	plr, err := ctx.sender()
	if err != nil {
		return err
	}
	ctx.sendNotice(formatMapLine(plr.mapID, plr.inst.id, plr.pos))
	return nil
}

func handleWhereCommand(ctx *staffCommandContext, args []string) error {
	if len(args) != 1 {
		return errCommandUsage
	}
	plr, err := ctx.findPlayerByName(args[0])
	if err != nil {
		return fmt.Errorf("Player %s is not online in this channel.", args[0])
	}
	ctx.sendNotice(formatPlayerLocationLine(plr))
	return nil
}

func handleChannelNoticeCommand(ctx *staffCommandContext, args []string) error {
	if len(args) == 0 {
		return errCommandUsage
	}
	ctx.server.players.broadcast(packetMessageNotice(strings.Join(args, " ")))
	return nil
}

func handleEventNoticeCommand(ctx *staffCommandContext, args []string) error {
	if len(args) == 0 {
		return errCommandUsage
	}
	ctx.server.players.broadcast(packetMessageWhiteBar(strings.Join(args, " ")))
	return nil
}

func handleNoticePlayerCommand(ctx *staffCommandContext, args []string) error {
	if len(args) < 2 {
		return errCommandUsage
	}
	target, err := ctx.findPlayerByName(args[0])
	if err != nil {
		return fmt.Errorf("Player %s is not online in this channel.", args[0])
	}
	target.Send(packetMessageNotice(strings.Join(args[1:], " ")))
	ctx.sendNotice(fmt.Sprintf("Sent a notice to %s.", target.Name))
	return nil
}

func handleUnstuckCommand(ctx *staffCommandContext, args []string) error {
	if len(args) > 1 {
		return errCommandUsage
	}
	target, err := ctx.sender()
	if err != nil {
		return err
	}
	if len(args) == 1 {
		target, err = ctx.findPlayerByName(args[0])
		if err != nil {
			return fmt.Errorf("Player %s is not online in this channel.", args[0])
		}
	}
	field, ok := ctx.server.fields[target.mapID]
	if !ok {
		return fmt.Errorf("Could not find field ID %d", target.mapID)
	}
	inst, err := field.getInstance(target.inst.id)
	if err != nil {
		return err
	}
	portalID, err := inst.calculateNearestSpawnPortalID(target.pos)
	if err != nil {
		portalID = 0
	}
	if int(portalID) >= len(inst.portals) {
		return fmt.Errorf("No valid spawn portal was found for map %d", target.mapID)
	}
	ctx.logCommand(args, target)
	if err := ctx.server.warpPlayerToInstance(target, field, inst.portals[portalID], target.inst.id, true); err != nil {
		return err
	}
	ctx.sendNotice(fmt.Sprintf("Unstuck %s.", target.Name))
	return nil
}

func handleWarpToCommand(ctx *staffCommandContext, args []string) error {
	if len(args) != 1 {
		return errCommandUsage
	}
	caller, err := ctx.sender()
	if err != nil {
		return err
	}
	target, err := ctx.findPlayerByName(args[0])
	if err != nil {
		return fmt.Errorf("Player %s is not online in this channel.", args[0])
	}
	ctx.logCommand(args, target)
	if err := warpPlayerToExactPlayerPosition(ctx.server, caller, target); err != nil {
		return err
	}
	ctx.sendNotice(fmt.Sprintf("Warped to %s.", target.Name))
	return nil
}

func handleWarpToMeCommand(ctx *staffCommandContext, args []string) error {
	if len(args) != 1 {
		return errCommandUsage
	}
	caller, err := ctx.sender()
	if err != nil {
		return err
	}
	target, err := ctx.findPlayerByName(args[0])
	if err != nil {
		return fmt.Errorf("Player %s is not online in this channel.", args[0])
	}
	ctx.logCommand(args, target)
	if err := warpPlayerToExactPlayerPosition(ctx.server, target, caller); err != nil {
		return err
	}
	ctx.sendNotice(fmt.Sprintf("Warped %s to you.", target.Name))
	return nil
}

func warpPlayerToExactPlayerPosition(server *Server, mover *Player, anchor *Player) error {
	if mover == nil || anchor == nil {
		return fmt.Errorf("player not found")
	}
	dstField, ok := server.fields[anchor.mapID]
	if !ok {
		return fmt.Errorf("Invalid map ID")
	}
	dstInst, err := dstField.getInstance(anchor.inst.id)
	if err != nil {
		return err
	}
	portalID, err := dstInst.calculateNearestSpawnPortalID(anchor.pos)
	if err != nil {
		portalID = 0
	}
	if int(portalID) >= len(dstInst.portals) {
		return fmt.Errorf("No valid spawn portal was found for the destination map")
	}
	return server.warpPlayerToInstanceAtPosition(mover, dstField, dstInst.portals[portalID], anchor.inst.id, anchor.pos, true, true)
}

func handleDropCommand(ctx *staffCommandContext, args []string) error {
	if len(args) < 1 || len(args) > 2 {
		return errCommandUsage
	}
	plr, err := ctx.sender()
	if err != nil {
		return err
	}
	field, ok := ctx.server.fields[plr.mapID]
	if !ok {
		return fmt.Errorf("Could not find field ID")
	}
	inst, err := field.getInstance(plr.inst.id)
	if err != nil {
		return err
	}
	item, err := buildCommandItem(args[0], args[1:])
	if err != nil {
		return err
	}
	item.creatorName = plr.Name
	ctx.logCommand(args, nil)
	inst.dropPool.createDrop(dropSpawnNormal, dropFreeForAll, 0, plr.pos, true, true, plr.ID, 0, item)
	ctx.sendNotice(fmt.Sprintf("Dropped item %d x%d.", item.ID, item.amount))
	return nil
}

func handleDropTestCommand(ctx *staffCommandContext, args []string) error {
	if len(args) != 0 {
		return errCommandUsage
	}
	plr, err := ctx.sender()
	if err != nil {
		return err
	}
	field, ok := ctx.server.fields[plr.mapID]
	if !ok {
		return fmt.Errorf("Could not find field ID")
	}
	inst, err := field.getInstance(plr.inst.id)
	if err != nil {
		return err
	}
	items := []int32{1372010, 1402005, 1422013, 1412021, 1382016, 1432030, 1442002, 1302023, 1322045, 1312015, 1332027, 1332026, 1462017, 1472033, 1452020, 1092029, 1092025}
	drops := make([]Item, 0, len(items))
	for _, id := range items {
		item, err := createPerfectItemFromID(id, 1)
		if err != nil {
			return err
		}
		item.creatorName = plr.Name
		drops = append(drops, item)
	}
	ctx.logCommand(args, nil)
	inst.dropPool.createDrop(dropSpawnNormal, dropFreeForAll, 1000, plr.pos, true, true, plr.ID, 0, drops...)
	ctx.sendNotice("Dropped the test loot set.")
	return nil
}

func handleReloadScriptsCommand(ctx *staffCommandContext, args []string) error {
	if len(args) != 0 {
		return errCommandUsage
	}
	ctx.logCommand(args, nil)
	stores := []*scriptStore{
		ctx.server.npcScriptStore,
		ctx.server.questScriptStore,
		ctx.server.eventScriptStore,
		ctx.server.portalScriptStore,
		ctx.server.reactorScriptStore,
	}
	for _, store := range stores {
		if store == nil {
			continue
		}
		store.scripts = make(map[string]*goja.Program)
		if err := store.loadScripts(); err != nil {
			return err
		}
	}
	ctx.sendNotice("Reloaded channel scripts.")
	return nil
}

func handleSaveAllCommand(ctx *staffCommandContext, args []string) error {
	if len(args) != 0 {
		return errCommandUsage
	}
	ctx.logCommand(args, nil)
	ctx.server.players.Flush()
	ctx.sendNotice("Flushed all online player state.")
	return nil
}

func buildCommandItem(itemArg string, quantityArgs []string) (Item, error) {
	itemID64, err := strconv.ParseInt(itemArg, 10, 32)
	if err != nil || itemID64 <= 0 {
		return Item{}, fmt.Errorf("Item ID must be a positive number")
	}
	itemID := int32(itemID64)

	item, err := CreateItemFromID(itemID, 1)
	if err != nil {
		return Item{}, fmt.Errorf("Invalid item ID %d", itemID)
	}

	quantity := int16(1)
	if len(quantityArgs) == 1 {
		qty64, err := strconv.ParseInt(quantityArgs[0], 10, 16)
		if err != nil || qty64 <= 0 {
			return Item{}, fmt.Errorf("Quantity must be a positive number")
		}
		quantity = int16(qty64)
	}

	if item.invID == constant.InventoryEquip || item.pet {
		item.amount = 1
		if len(quantityArgs) == 1 && quantity != 1 {
			return Item{}, fmt.Errorf("Equips and pets can only be dropped one at a time")
		}
		return item, nil
	}

	slotMax := getItemSlotMax(itemID)
	if slotMax <= 0 {
		slotMax = constant.MaxItemStack
	}
	if quantity > slotMax {
		return Item{}, fmt.Errorf("Quantity %d exceeds the stack limit of %d for item %d", quantity, slotMax, itemID)
	}
	item.amount = quantity
	return item, nil
}

func formatPlayerLocationLine(plr *Player) string {
	return fmt.Sprintf("%s(ch %d, %s, pos %s)", plr.Name, plr.ChannelID+1, formatMapLabel(plr.mapID), plr.pos.String())
}

func formatMapLine(mapID int32, instanceID int, p pos) string {
	return fmt.Sprintf("%s, instance %d, pos %s", formatMapLabel(mapID), instanceID, p.String())
}

func formatMapLabel(mapID int32) string {
	mapName := fmt.Sprintf("map %d", mapID)
	if data, err := nx.GetMap(mapID); err == nil {
		name := strings.TrimSpace(data.MapName)
		street := strings.TrimSpace(data.StreetName)
		switch {
		case street != "" && name != "":
			mapName = fmt.Sprintf("map %d (%s - %s)", mapID, street, name)
		case name != "":
			mapName = fmt.Sprintf("map %d (%s)", mapID, name)
		}
	}
	return mapName
}

func splitCommandOutput(values []string, perLine int) [][]string {
	if perLine <= 0 || len(values) == 0 {
		return nil
	}
	result := make([][]string, 0, (len(values)+perLine-1)/perLine)
	for i := 0; i < len(values); i += perLine {
		end := i + perLine
		if end > len(values) {
			end = len(values)
		}
		result = append(result, values[i:end])
	}
	return result
}

func withAudit(ctx *staffCommandContext, args []string, target *Player) {
	if ctx.command != nil && ctx.command.audit {
		ctx.logCommand(args, target)
	}
}

func parseIntArg(text, name string) (int, error) {
	val, err := strconv.Atoi(text)
	if err != nil {
		return 0, fmt.Errorf("%s must be a number", name)
	}
	return val, nil
}

func parseTargetAndAmount(ctx *staffCommandContext, args []string) (*Player, int, error) {
	var (
		target *Player
		amount int
		err    error
	)
	if len(args) == 1 {
		target, err = ctx.sender()
		if err != nil {
			return nil, 0, err
		}
		amount, err = parseIntArg(args[0], "Amount")
		return target, amount, err
	}
	if len(args) == 2 {
		target, err = ctx.findPlayerByName(args[0])
		if err != nil {
			return nil, 0, err
		}
		amount, err = parseIntArg(args[1], "Amount")
		return target, amount, err
	}
	return nil, 0, errCommandUsage
}

func currentFieldInstance(ctx *staffCommandContext) (*Player, *field, *fieldInstance, error) {
	plr, err := ctx.sender()
	if err != nil {
		return nil, nil, nil, err
	}
	fld, ok := ctx.server.fields[plr.mapID]
	if !ok {
		return nil, nil, nil, fmt.Errorf("Could not find field ID")
	}
	inst, err := fld.getInstance(plr.inst.id)
	if err != nil {
		return nil, nil, nil, err
	}
	return plr, fld, inst, nil
}

func handleSearchCommand(ctx *staffCommandContext, args []string) error {
	if len(args) < 1 {
		return errCommandUsage
	}
	searchType := "all"
	queryStart := 0
	if len(args) > 1 {
		switch strings.ToLower(args[0]) {
		case "map", "maps":
			searchType = "map"
			queryStart = 1
		case "quest", "quests":
			searchType = "quest"
			queryStart = 1
		default:
			if itemType, ok := nx.NormalizeItemSearchCategory(args[0]); ok {
				searchType = itemType
				queryStart = 1
			}
		}
	}
	if len(args) <= queryStart {
		return fmt.Errorf("Search query was not provided")
	}
	query := strings.Join(args[queryStart:], " ")
	const perTypeLimit = 5
	sent := 0
	sendMatches := func(label string, matches []nx.StringMatch) {
		if len(matches) == 0 {
			return
		}
		ctx.sendError(label + ":")
		for _, match := range matches {
			line := fmt.Sprintf("  [%d] %s", match.ID, match.Name)
			if match.Extra != "" {
				line += fmt.Sprintf(" (%s)", match.Extra)
			}
			ctx.sendError(line)
			sent++
		}
	}
	switch searchType {
	case "map":
		sendMatches("Maps", nx.SearchMapsByName(query, perTypeLimit))
	case "quest":
		sendMatches("Quests", nx.SearchQuestsByName(query, perTypeLimit))
	default:
		if searchType == "all" {
			sendMatches("Items", nx.SearchItemsByName(query, perTypeLimit))
			sendMatches("Maps", nx.SearchMapsByName(query, perTypeLimit))
			sendMatches("Quests", nx.SearchQuestsByName(query, perTypeLimit))
		} else {
			sendMatches(searchLabelForItemType(searchType), nx.SearchItemsByCategory(query, searchType, perTypeLimit))
		}
	}
	if sent == 0 {
		ctx.sendError("No WZ matches found for: " + query)
	}
	return nil
}

func handleShowRatesCommand(ctx *staffCommandContext, args []string) error {
	if len(args) != 0 {
		return errCommandUsage
	}
	ctx.sendNotice(fmt.Sprintf("Exp: x%.2f, Drop: x%.2f, Mesos: x%.2f", ctx.server.rates.exp, ctx.server.rates.drop, ctx.server.rates.mesos))
	return nil
}

func handlePosCommand(ctx *staffCommandContext, args []string) error {
	if len(args) != 0 {
		return errCommandUsage
	}
	plr, err := ctx.sender()
	if err != nil {
		return err
	}
	ctx.sendNotice(plr.pos.String())
	return nil
}

func handleMsgBoxCommand(ctx *staffCommandContext, args []string) error {
	if len(args) == 0 {
		return errCommandUsage
	}
	ctx.server.players.broadcast(packetMessageDialogueBox(strings.Join(args, " ")))
	return nil
}

func handleHeaderCommand(ctx *staffCommandContext, args []string) error {
	ctx.server.header = strings.Join(args, " ")
	ctx.server.players.broadcast(packetMessageScrollingHeader(ctx.server.header))
	return nil
}

func handleChangeBgmCommand(ctx *staffCommandContext, args []string) error {
	if len(args) > 1 {
		return errCommandUsage
	}
	_, _, inst, err := currentFieldInstance(ctx)
	if err != nil {
		return err
	}
	bgm := ""
	if len(args) == 1 {
		bgm = args[0]
	}
	inst.changeBgm(bgm)
	return nil
}

func handleWrongCommand(ctx *staffCommandContext, args []string) error {
	if len(args) != 0 {
		return errCommandUsage
	}
	plr, err := ctx.sender()
	if err != nil {
		return err
	}
	plr.inst.send(packetShowScreenEffect("quest/party/wrong_kor"))
	plr.inst.send(packetPlaySound("Party1/Failed"))
	return nil
}

func handleClearCommand(ctx *staffCommandContext, args []string) error {
	if len(args) != 0 {
		return errCommandUsage
	}
	plr, err := ctx.sender()
	if err != nil {
		return err
	}
	plr.inst.send(packetShowScreenEffect("quest/party/clear"))
	plr.inst.send(packetPlaySound("Party1/Clear"))
	return nil
}

func handleGateCommand(ctx *staffCommandContext, args []string) error {
	if len(args) != 0 {
		return errCommandUsage
	}
	plr, err := ctx.sender()
	if err != nil {
		return err
	}
	plr.inst.send(packetPortalEffectt(2, "gate"))
	return nil
}

func handleCreateInstanceCommand(ctx *staffCommandContext, args []string) error {
	if len(args) != 0 {
		return errCommandUsage
	}
	plr, fld, _, err := currentFieldInstance(ctx)
	if err != nil {
		return err
	}
	_ = plr
	id := fld.createInstance(&ctx.server.rates, ctx.server)
	ctx.sendNotice("Created instance: " + strconv.Itoa(id))
	return nil
}

func handleChangeInstanceCommand(ctx *staffCommandContext, args []string) error {
	if len(args) != 1 {
		return errCommandUsage
	}
	plr, fld, _, err := currentFieldInstance(ctx)
	if err != nil {
		return err
	}
	instanceID, err := parseIntArg(args[0], "Instance ID")
	if err != nil {
		return err
	}
	withAudit(ctx, args, plr)
	if err := fld.changePlayerInstance(plr, instanceID); err != nil {
		return err
	}
	ctx.sendNotice("Changed instance to " + strconv.Itoa(instanceID))
	return nil
}

func handleMapInfoCommand(ctx *staffCommandContext, args []string) error {
	if len(args) != 0 {
		return errCommandUsage
	}
	plr, _, _, err := currentFieldInstance(ctx)
	if err != nil {
		return err
	}
	fld, ok := ctx.server.fields[plr.mapID]
	if !ok {
		return fmt.Errorf("Could not find field ID")
	}
	for i, v := range fld.instances {
		ctx.sendNotice("instance " + strconv.Itoa(i) + ":")
		ctx.sendNotice(v.String())
	}
	return nil
}

func handleEventStartCommand(ctx *staffCommandContext, args []string) error {
	if len(args) != 2 {
		return errCommandUsage
	}
	program, ok := ctx.server.eventScriptStore.scripts[args[0]]
	if !ok {
		return fmt.Errorf("Could not find event script: %s", args[0])
	}
	instanceID, err := parseIntArg(args[1], "Instance ID")
	if err != nil {
		return err
	}
	plr, err := ctx.sender()
	if err != nil {
		return err
	}
	ids := []int32{}
	if plr.party != nil {
		for i, id := range plr.party.PlayerID {
			if plr.mapID == plr.party.MapID[i] && plr.party.players[i] != nil && plr.inst.id == plr.party.players[i].inst.id {
				ids = append(ids, id)
			}
		}
	} else {
		ids = append(ids, plr.ID)
	}
	event, err := createEvent(plr.ID, instanceID, ids, ctx.server, program)
	if err != nil {
		return err
	}
	ctx.server.events[plr.ID] = event
	event.start()
	return nil
}

func handleEventsCommand(ctx *staffCommandContext, args []string) error {
	if len(args) != 0 {
		return errCommandUsage
	}
	ctx.sendNotice("There are currently " + strconv.Itoa(len(ctx.server.events)) + " events running")
	for id, event := range ctx.server.events {
		ctx.sendNotice("id: " + strconv.Itoa(int(id)))
		info := "player ids:"
		for _, v := range event.playerIDs {
			info += " " + strconv.Itoa(int(v))
		}
		ctx.sendNotice(info)
		ctx.sendNotice("remaining time: " + time.Until(event.endTime).String())
	}
	return nil
}

func handleReviveCommand(ctx *staffCommandContext, args []string) error {
	if len(args) > 1 {
		return errCommandUsage
	}
	target, err := ctx.sender()
	if err != nil {
		return err
	}
	if len(args) == 1 {
		target, err = ctx.findPlayerByName(args[0])
		if err != nil {
			return err
		}
	}
	target.setHP(target.maxHP)
	return nil
}

func handleQuestFinishCommand(ctx *staffCommandContext, args []string) error {
	if len(args) != 1 {
		return errCommandUsage
	}
	qid64, err := strconv.ParseInt(args[0], 10, 16)
	if err != nil {
		return fmt.Errorf("Quest ID must be a number")
	}
	plr, err := ctx.sender()
	if err != nil {
		return err
	}
	questID := int16(qid64)
	nowMs := time.Now().UnixMilli()
	plr.quests.complete(questID, nowMs)
	setQuestCompleted(plr.ID, questID, nowMs)
	plr.Send(packetQuestComplete(questID))
	ctx.sendNotice(fmt.Sprintf("Quest %d completed", questID))
	return nil
}

func handleQuestUntilCommand(ctx *staffCommandContext, args []string) error {
	if len(args) != 2 {
		return errCommandUsage
	}
	qid64, err := strconv.ParseInt(args[0], 10, 16)
	if err != nil {
		return fmt.Errorf("Quest ID must be a number")
	}
	part, err := parseIntArg(args[1], "Part")
	if err != nil || part < 0 {
		return fmt.Errorf("Part must be a non-negative number")
	}
	plr, err := ctx.sender()
	if err != nil {
		return err
	}
	questID := int16(qid64)
	record := fmt.Sprintf("p%d", part)
	plr.quests.add(questID, record)
	upsertQuestRecord(plr.ID, questID, record)
	plr.Send(packetQuestUpdate(questID, record))
	ctx.sendNotice(fmt.Sprintf("Quest %d progressed to %s", questID, record))
	return nil
}

func handleQuestResetCommand(ctx *staffCommandContext, args []string) error {
	if len(args) != 1 {
		return errCommandUsage
	}
	qid64, err := strconv.ParseInt(args[0], 10, 16)
	if err != nil {
		return fmt.Errorf("Quest ID must be a number")
	}
	plr, err := ctx.sender()
	if err != nil {
		return err
	}
	questID := int16(qid64)
	delete(plr.quests.inProgress, questID)
	delete(plr.quests.completed, questID)
	delete(plr.quests.mobKills, questID)
	deleteQuest(plr.ID, questID)
	clearQuestMobKills(plr.ID, questID)
	plr.Send(packetQuestRemove(questID))
	ctx.sendNotice(fmt.Sprintf("Quest %d has been reset", questID))
	return nil
}

func handleClearInstPropsCommand(ctx *staffCommandContext, args []string) error {
	if len(args) != 0 {
		return errCommandUsage
	}
	plr, err := ctx.sender()
	if err != nil {
		return err
	}
	for i := range plr.inst.properties {
		delete(plr.inst.properties, i)
	}
	return nil
}

func handlePropertiesCommand(ctx *staffCommandContext, args []string) error {
	if len(args) != 0 {
		return errCommandUsage
	}
	plr, err := ctx.sender()
	if err != nil {
		return err
	}
	for i, v := range plr.inst.properties {
		ctx.sendError(fmt.Sprintf("prop: %s , value: %v", i, v))
	}
	return nil
}

func handleKillCommand(ctx *staffCommandContext, args []string) error {
	if len(args) > 1 {
		return errCommandUsage
	}
	target, err := ctx.sender()
	if err != nil {
		return err
	}
	if len(args) == 1 {
		target, err = ctx.findPlayerByName(args[0])
		if err != nil {
			return err
		}
	}
	withAudit(ctx, args, target)
	target.setHP(0)
	return nil
}

func handleWarpMapCommand(ctx *staffCommandContext, args []string) error {
	if len(args) != 1 {
		return errCommandUsage
	}
	var id int32
	if val, err := strconv.Atoi(args[0]); err == nil {
		id = int32(val)
	} else if mapped, ok := convertMapNameToID(args[0]); ok {
		id = mapped
	} else {
		return fmt.Errorf("Unknown map destination: %s", args[0])
	}
	if _, err := nx.GetMap(id); err != nil {
		return err
	}
	plr, err := ctx.sender()
	if err != nil {
		return err
	}
	dstField, ok := ctx.server.fields[id]
	if !ok {
		return fmt.Errorf("Invalid map ID")
	}
	inst, err := dstField.getInstance(plr.inst.id)
	if err != nil {
		inst, err = dstField.getInstance(0)
		if err != nil {
			return err
		}
	}
	portal, err := inst.getRandomSpawnPortal()
	if err != nil {
		return err
	}
	withAudit(ctx, args, plr)
	return ctx.server.warpPlayer(plr, dstField, portal, true)
}

func handleClearDropsCommand(ctx *staffCommandContext, args []string) error {
	if len(args) != 0 {
		return errCommandUsage
	}
	_, _, inst, err := currentFieldInstance(ctx)
	if err != nil {
		return err
	}
	withAudit(ctx, args, nil)
	inst.dropPool.clearDrops()
	return nil
}

func handleRemoveTimerCommand(ctx *staffCommandContext, args []string) error {
	if len(args) != 0 {
		return errCommandUsage
	}
	_, _, inst, err := currentFieldInstance(ctx)
	if err != nil {
		return err
	}
	withAudit(ctx, args, nil)
	inst.fieldTimer.Reset(0)
	return nil
}

func handleKillMobCommand(ctx *staffCommandContext, args []string) error {
	if len(args) != 1 {
		return errCommandUsage
	}
	spawnID, err := parseIntArg(args[0], "Spawn ID")
	if err != nil {
		return err
	}
	plr, _, inst, err := currentFieldInstance(ctx)
	if err != nil {
		return err
	}
	withAudit(ctx, args, nil)
	inst.lifePool.mobDamaged(int32(spawnID), plr, int32(^uint32(0)>>1))
	return nil
}

func handleKillAllCommand(ctx *staffCommandContext, args []string) error {
	if len(args) != 0 {
		return errCommandUsage
	}
	plr, _, inst, err := currentFieldInstance(ctx)
	if err != nil {
		return err
	}
	withAudit(ctx, args, nil)
	for spawnID, mob := range inst.lifePool.mobs {
		inst.lifePool.mobDamaged(spawnID, plr, mob.hp)
	}
	return nil
}

func handleSpawnCommand(ctx *staffCommandContext, args []string) error {
	if len(args) < 1 || len(args) > 2 {
		return errCommandUsage
	}
	mobID, err := parseIntArg(args[0], "Mob ID")
	if err != nil {
		return err
	}
	count := 1
	if len(args) == 2 {
		count, err = parseIntArg(args[1], "Count")
		if err != nil {
			return err
		}
	}
	plr, _, inst, err := currentFieldInstance(ctx)
	if err != nil {
		return err
	}
	withAudit(ctx, args, nil)
	for i := 0; i < count; i++ {
		if err := inst.lifePool.spawnMobFromID(int32(mobID), plr.pos, false, true, true, constant.MobSummonTypeInstant, plr.ID); err != nil {
			return err
		}
	}
	return nil
}

func handleSpawnBossCommand(ctx *staffCommandContext, args []string) error {
	if len(args) < 1 || len(args) > 2 {
		return errCommandUsage
	}
	mobIDs, err := covnertMobNameToID(args[0])
	if err != nil {
		return err
	}
	count := 1
	if len(args) == 2 {
		count, err = parseIntArg(args[1], "Count")
		if err != nil {
			return err
		}
	}
	plr, _, inst, err := currentFieldInstance(ctx)
	if err != nil {
		return err
	}
	withAudit(ctx, args, nil)
	for i := 0; i < count; i++ {
		for _, id := range mobIDs {
			if err := inst.lifePool.spawnMobFromID(id, plr.pos, false, true, true, constant.MobSummonTypeInstant, plr.ID); err != nil {
				return err
			}
		}
	}
	return nil
}

func handleBanCommand(ctx *staffCommandContext, args []string) error {
	if len(args) < 1 {
		return errCommandUsage
	}
	if ctx.server.ac == nil {
		return fmt.Errorf("anti-cheat is not available")
	}
	target, err := ctx.findPlayerByName(args[0])
	if err != nil {
		return fmt.Errorf("Player not found")
	}
	hours := 168
	reason := "Banned by GM"
	if len(args) >= 2 {
		if args[1] == "perm" {
			hours = 0
		} else if h, err := strconv.Atoi(args[1]); err == nil {
			hours = h
		}
	}
	if len(args) >= 3 {
		reason = strings.Join(args[2:], " ")
	}
	withAudit(ctx, args, target)
	if err := ctx.server.ac.IssueBan(target.Conn.GetAccountID(), hours, reason, "", ""); err != nil {
		return err
	}
	ctx.sendError(fmt.Sprintf("Banned %s for %d hours", target.Name, hours))
	return nil
}

func handleUnbanCommand(ctx *staffCommandContext, args []string) error {
	if len(args) != 1 {
		return errCommandUsage
	}
	if ctx.server.ac == nil {
		return fmt.Errorf("anti-cheat is not available")
	}
	withAudit(ctx, args, nil)
	if err := ctx.server.ac.Unban(args[0]); err != nil {
		return err
	}
	ctx.sendError(fmt.Sprintf("Unbanned player %s", args[0]))
	return nil
}

func handleBanHistoryCommand(ctx *staffCommandContext, args []string) error {
	if len(args) != 1 {
		return errCommandUsage
	}
	if ctx.server.ac == nil {
		return fmt.Errorf("anti-cheat is not available")
	}
	history, err := ctx.server.ac.GetBanHistory(args[0], 10)
	if err != nil {
		return err
	}
	if len(history) == 0 {
		ctx.sendError("No ban history")
		return nil
	}
	for _, entry := range history {
		ctx.sendError(entry)
	}
	return nil
}

func handleEnablePortalCommand(ctx *staffCommandContext, args []string) error {
	if len(args) != 2 {
		return errCommandUsage
	}
	plr, err := ctx.sender()
	if err != nil {
		return err
	}
	port, err := plr.inst.getPortalFromName(args[0])
	if err != nil {
		return err
	}
	port.enabled = args[1] == "true"
	withAudit(ctx, args, nil)
	ctx.sendError(fmt.Sprintf("portal %s has been set to %v", port.name, port.enabled))
	return nil
}

func handleDeleteInstanceCommand(ctx *staffCommandContext, args []string) error {
	if len(args) != 1 {
		return errCommandUsage
	}
	instanceID, err := parseIntArg(args[0], "Instance ID")
	if err != nil {
		return err
	}
	if instanceID < 1 {
		return fmt.Errorf("Cannot delete instance 0")
	}
	plr, fld, _, err := currentFieldInstance(ctx)
	if err != nil {
		return err
	}
	if plr.inst.id == instanceID {
		return fmt.Errorf("Cannot delete the same instance you are in")
	}
	withAudit(ctx, args, nil)
	if err := fld.deleteInstance(instanceID); err != nil {
		return err
	}
	ctx.sendNotice("Deleted")
	return nil
}

func handleRemoveDropCommand(ctx *staffCommandContext, args []string) error {
	if len(args) != 1 {
		return errCommandUsage
	}
	dropID, err := parseIntArg(args[0], "Drop ID")
	if err != nil {
		return err
	}
	_, _, inst, err := currentFieldInstance(ctx)
	if err != nil {
		return err
	}
	withAudit(ctx, args, nil)
	inst.dropPool.removeDrop(0, int32(dropID))
	return nil
}

func handleItemCommand(ctx *staffCommandContext, args []string) error {
	if len(args) < 1 || len(args) > 2 {
		return errCommandUsage
	}
	itemID, err := parseIntArg(args[0], "Item ID")
	if err != nil {
		return err
	}
	amount := 1
	if len(args) == 2 {
		amount, err = parseIntArg(args[1], "Quantity")
		if err != nil {
			return err
		}
	}
	item, err := CreateItemFromID(int32(itemID), int16(amount))
	if err != nil {
		return err
	}
	plr, err := ctx.sender()
	if err != nil {
		return err
	}
	item.creatorName = plr.Name
	withAudit(ctx, args, nil)
	_, err = plr.GiveItem(item)
	return err
}

func handleLoadoutCommand(ctx *staffCommandContext, args []string) error {
	if len(args) != 0 {
		return errCommandUsage
	}
	plr, err := ctx.sender()
	if err != nil {
		return err
	}
	withAudit(ctx, args, nil)
	equips := []int32{1372010, 1402005, 1422013, 1412021, 1382016, 1432030, 1442002, 1302023, 1322045, 1312015, 1332027, 1332026, 1462017, 1472033, 1452020, 1092029, 1092025}
	for _, v := range equips {
		item, err := createPerfectItemFromID(v, 1)
		if err != nil {
			return err
		}
		item.creatorName = plr.Name
		if _, err = plr.GiveItem(item); err != nil {
			return err
		}
	}
	etc := []int32{4006001, 4006000, 4001017, 4031179, 4031059}
	for _, v := range etc {
		item, err := createPerfectItemFromID(v, 100)
		if err != nil {
			return err
		}
		item.creatorName = plr.Name
		if _, err = plr.GiveItem(item); err != nil {
			return err
		}
	}
	return nil
}

func handleExpCommand(ctx *staffCommandContext, args []string) error {
	target, amount, err := parseTargetAndAmount(ctx, args)
	if err != nil {
		return err
	}
	withAudit(ctx, args, target)
	target.setEXP(int32(amount))
	return nil
}

func handleGrantExpCommand(ctx *staffCommandContext, args []string) error {
	target, amount, err := parseTargetAndAmount(ctx, args)
	if err != nil {
		return err
	}
	withAudit(ctx, args, target)
	target.giveEXP(int32(amount), false, false)
	return nil
}

func handleMesosCommand(ctx *staffCommandContext, args []string) error {
	if len(args) != 1 {
		return errCommandUsage
	}
	amount, err := parseIntArg(args[0], "Amount")
	if err != nil {
		return err
	}
	plr, err := ctx.sender()
	if err != nil {
		return err
	}
	withAudit(ctx, args, nil)
	plr.setMesos(int32(amount))
	return nil
}

func handleNxCommand(ctx *staffCommandContext, args []string) error {
	if len(args) != 1 {
		return errCommandUsage
	}
	amount, err := parseIntArg(args[0], "Amount")
	if err != nil {
		return err
	}
	plr, err := ctx.sender()
	if err != nil {
		return err
	}
	withAudit(ctx, args, nil)
	plr.nx += int32(amount)
	return nil
}

func handleMaplePointsCommand(ctx *staffCommandContext, args []string) error {
	if len(args) != 1 {
		return errCommandUsage
	}
	amount, err := parseIntArg(args[0], "Amount")
	if err != nil {
		return err
	}
	plr, err := ctx.sender()
	if err != nil {
		return err
	}
	withAudit(ctx, args, nil)
	plr.maplepoints += int32(amount)
	return nil
}

func handleHPCommand(ctx *staffCommandContext, args []string) error {
	target, amount, err := parseTargetAndAmount(ctx, args)
	if err != nil {
		return err
	}
	if int16(amount) > target.maxHP {
		target.setMaxHP(int16(amount))
	}
	withAudit(ctx, args, target)
	target.setHP(int16(amount))
	return nil
}

func handleMPCommand(ctx *staffCommandContext, args []string) error {
	target, amount, err := parseTargetAndAmount(ctx, args)
	if err != nil {
		return err
	}
	if int16(amount) > target.maxMP {
		target.setMaxMP(int16(amount))
	}
	withAudit(ctx, args, target)
	target.setMP(int16(amount))
	return nil
}

func handleSetMaxHPCommand(ctx *staffCommandContext, args []string) error {
	if len(args) != 1 {
		return errCommandUsage
	}
	amount, err := parseIntArg(args[0], "Amount")
	if err != nil {
		return err
	}
	if amount < 1 {
		return fmt.Errorf("Max HP must be at least 1")
	}
	plr, err := ctx.sender()
	if err != nil {
		return err
	}
	withAudit(ctx, args, nil)
	plr.setMaxHP(int16(amount))
	ctx.sendNotice(fmt.Sprintf("Set Max HP to %d", amount))
	return nil
}

func handleSetMaxMPCommand(ctx *staffCommandContext, args []string) error {
	if len(args) != 1 {
		return errCommandUsage
	}
	amount, err := parseIntArg(args[0], "Amount")
	if err != nil {
		return err
	}
	if amount < 0 {
		return fmt.Errorf("Max MP cannot be negative")
	}
	plr, err := ctx.sender()
	if err != nil {
		return err
	}
	withAudit(ctx, args, nil)
	plr.setMaxMP(int16(amount))
	ctx.sendNotice(fmt.Sprintf("Set Max MP to %d", amount))
	return nil
}

func handlePrimaryStatCommand(ctx *staffCommandContext, args []string) error {
	target, amount, err := parseTargetAndAmount(ctx, args)
	if err != nil {
		return err
	}
	if amount < 0 {
		return fmt.Errorf("Stat cannot be negative")
	}
	withAudit(ctx, args, target)
	switch ctx.command.name {
	case "str":
		target.str = int16(amount)
		target.MarkDirty(DirtyStr, time.Millisecond*300)
		target.Send(packetPlayerStatChange(false, constant.StrID, int32(target.str)))
	case "dex":
		target.dex = int16(amount)
		target.MarkDirty(DirtyDex, time.Millisecond*300)
		target.Send(packetPlayerStatChange(false, constant.DexID, int32(target.dex)))
	case "int":
		target.intt = int16(amount)
		target.MarkDirty(DirtyInt, time.Millisecond*300)
		target.Send(packetPlayerStatChange(false, constant.IntID, int32(target.intt)))
	case "luk":
		target.luk = int16(amount)
		target.MarkDirty(DirtyLuk, time.Millisecond*300)
		target.Send(packetPlayerStatChange(false, constant.LukID, int32(target.luk)))
	}
	return nil
}

func handleAPCommand(ctx *staffCommandContext, args []string) error {
	target, amount, err := parseTargetAndAmount(ctx, args)
	if err != nil {
		return err
	}
	withAudit(ctx, args, target)
	target.setAP(int16(amount))
	return nil
}

func handleSPCommand(ctx *staffCommandContext, args []string) error {
	target, amount, err := parseTargetAndAmount(ctx, args)
	if err != nil {
		return err
	}
	withAudit(ctx, args, target)
	target.setSP(int16(amount))
	return nil
}

func handleLevelCommand(ctx *staffCommandContext, args []string) error {
	target, amount, err := parseTargetAndAmount(ctx, args)
	if err != nil {
		return err
	}
	withAudit(ctx, args, target)
	target.setLevel(byte(amount))
	return nil
}

func handleLevelUpCommand(ctx *staffCommandContext, args []string) error {
	var target *Player
	amount := 1
	var err error
	switch len(args) {
	case 0:
		target, err = ctx.sender()
	case 1:
		target, err = ctx.sender()
		if err == nil {
			amount, err = parseIntArg(args[0], "Amount")
		}
	case 2:
		target, err = ctx.findPlayerByName(args[0])
		if err == nil {
			amount, err = parseIntArg(args[1], "Amount")
		}
	default:
		return errCommandUsage
	}
	if err != nil {
		return err
	}
	withAudit(ctx, args, target)
	target.giveLevel(byte(amount))
	return nil
}

func handleJobCommand(ctx *staffCommandContext, args []string) error {
	if len(args) != 1 {
		return errCommandUsage
	}
	jobID := int16(0)
	if val, err := strconv.Atoi(args[0]); err == nil {
		jobID = int16(val)
	} else {
		jobID = convertJobNameToID(args[0])
	}
	plr, err := ctx.sender()
	if err != nil {
		return err
	}
	withAudit(ctx, args, nil)
	plr.setJob(jobID)
	return nil
}

func handleSkillLevelCommand(ctx *staffCommandContext, args []string) error {
	var target *Player
	var err error
	useArgs := args
	switch len(args) {
	case 2:
		target, err = ctx.sender()
	case 3:
		target, err = ctx.findPlayerByName(args[0])
		useArgs = args[1:]
	default:
		return errCommandUsage
	}
	if err != nil {
		return err
	}
	skillID64, err := strconv.ParseInt(useArgs[0], 10, 32)
	if err != nil {
		return fmt.Errorf("Skill must be a numeric ID")
	}
	skillID := int32(skillID64)
	levels, err := nx.GetPlayerSkill(skillID)
	if err != nil || len(levels) == 0 {
		return fmt.Errorf("Unknown or invalid skill ID: %d", skillID)
	}
	var level byte
	if strings.EqualFold(useArgs[1], "max") {
		level = byte(len(levels))
	} else {
		lv64, err := strconv.ParseInt(useArgs[1], 10, 32)
		if err != nil {
			return fmt.Errorf("Level must be a number or 'max'")
		}
		if lv64 < 0 || int(lv64) > len(levels) {
			return fmt.Errorf("Invalid level, max for skill %d is %d", skillID, len(levels))
		}
		level = byte(lv64)
	}
	withAudit(ctx, args, target)
	if level == 0 {
		delete(target.skills, skillID)
		target.MarkDirty(DirtySkills, time.Millisecond*300)
		ctx.sendNotice(fmt.Sprintf("Removed skill %d from %s", skillID, target.Name))
		return nil
	}
	ps, err := createPlayerSkillFromData(skillID, level)
	if err != nil {
		return err
	}
	if target.skills == nil {
		target.skills = make(map[int32]playerSkill, 8)
	}
	target.skills[skillID] = ps
	target.MarkDirty(DirtySkills, time.Millisecond*300)
	ctx.sendNotice(fmt.Sprintf("Set %s's skill %d to level %d", target.Name, skillID, ps.Level))
	return nil
}

func handleMaxSkillsCommand(ctx *staffCommandContext, args []string) error {
	if len(args) != 0 {
		return errCommandUsage
	}
	target, err := ctx.sender()
	if err != nil {
		return err
	}
	if target.skills == nil {
		target.skills = make(map[int32]playerSkill, 1024)
	}
	withAudit(ctx, args, target)
	jobFamilies := []int{int(constant.BeginnerJobID), int(constant.WarriorJobID), int(constant.FighterJobID), int(constant.CrusaderJobID), int(constant.PageJobID), int(constant.WhiteKnightJobID), int(constant.SpearmanJobID), int(constant.DragonKnightJobID), int(constant.MagicianJobID), int(constant.FirePoisonWizardJobID), int(constant.FirePoisonMageJobID), int(constant.IceLightWizardJobID), int(constant.IceLightMageJobID), int(constant.ClericJobID), int(constant.PriestJobID), int(constant.BowmanJobID), int(constant.HunterJobID), int(constant.RangerJobID), int(constant.CrossbowmanJobID), int(constant.SniperJobID), int(constant.ThiefJobID), int(constant.AssassinJobID), int(constant.HermitJobID), int(constant.BanditJobID), int(constant.ChiefBanditJobID), int(constant.GmJobID), int(constant.SuperGmJobID)}
	count := 0
	for _, job := range jobFamilies {
		base := job * 10000
		for idx := 0; idx <= 1999; idx++ {
			skillID := int32(base + idx)
			levels, err := nx.GetPlayerSkill(skillID)
			if err != nil || len(levels) == 0 {
				continue
			}
			ps, err := createPlayerSkillFromData(skillID, byte(len(levels)))
			if err != nil {
				continue
			}
			target.skills[skillID] = ps
			count++
		}
	}
	target.MarkDirty(DirtySkills, time.Millisecond*300)
	ctx.sendNotice(fmt.Sprintf("Maxed %d skills across all classes for %s", count, target.Name))
	return nil
}

func handleResetSkillsCommand(ctx *staffCommandContext, args []string) error {
	if len(args) != 0 {
		return errCommandUsage
	}
	target, err := ctx.sender()
	if err != nil {
		return err
	}
	withAudit(ctx, args, target)
	for k := range target.skills {
		delete(target.skills, k)
	}
	target.MarkDirty(DirtySkills, time.Millisecond*300)
	ctx.sendNotice(fmt.Sprintf("Reset all skills for %s", target.Name))
	return nil
}

func handleRateCommand(ctx *staffCommandContext, args []string) error {
	if len(args) != 2 {
		return errCommandUsage
	}
	rates := map[string]func(rate float32) mpacket.Packet{"exp": internal.PacketChangeExpRate, "drop": internal.PacketChangeDropRate, "mesos": internal.PacketChangeMesosRate}
	mFunc, ok := rates[args[0]]
	if !ok {
		return fmt.Errorf("Choose between exp/drop/mesos rates")
	}
	rate, err := strconv.ParseFloat(args[1], 32)
	if err != nil {
		return fmt.Errorf("<rate> should be a number")
	}
	withAudit(ctx, args, nil)
	ctx.server.world.Send(mFunc(float32(rate)))
	return nil
}

func handleSetWorldMessageCommand(ctx *staffCommandContext, args []string) error {
	if len(args) < 1 {
		return errCommandUsage
	}
	ribbon, err := parseIntArg(args[0], "Ribbon number")
	if err != nil || ribbon < 0 {
		return fmt.Errorf("Invalid ribbon number")
	}
	message := ""
	if len(args) > 1 {
		message = strings.Join(args[1:], " ")
	}
	withAudit(ctx, args, nil)
	ctx.server.world.Send(internal.PacketUpdateLoginInfo(byte(ribbon), message))
	return nil
}

func handlePartyCreateCommand(ctx *staffCommandContext, args []string) error {
	if len(args) != 0 {
		return errCommandUsage
	}
	plr, err := ctx.sender()
	if err != nil {
		return err
	}
	ctx.server.world.Send(internal.PacketChannelPartyCreateRequest(plr.ID, ctx.server.id, plr.mapID, int32(plr.job), int32(plr.level), plr.Name))
	return nil
}

func handleGuildCreateCommand(ctx *staffCommandContext, args []string) error {
	if len(args) != 0 {
		return errCommandUsage
	}
	plr, err := ctx.sender()
	if err != nil {
		return err
	}
	program, ok := ctx.server.npcScriptStore.scripts["2010007"]
	if !ok {
		return fmt.Errorf("Unable to find guild npc script")
	}
	controller, err := createNpcChatController(2010007, "2010007", ctx.conn, program, plr, ctx.server)
	if err != nil {
		return err
	}
	ctx.server.npcChat[ctx.conn] = controller
	ctx.server.updateNPCInteractionMetric(1)
	if controller.run() {
		delete(ctx.server.npcChat, ctx.conn)
		ctx.server.updateNPCInteractionMetric(-1)
	}
	return nil
}

func handleGuildDisbandCommand(ctx *staffCommandContext, args []string) error {
	if len(args) != 0 {
		return errCommandUsage
	}
	plr, err := ctx.sender()
	if err != nil {
		return err
	}
	if plr.guild == nil {
		return fmt.Errorf("Not in guild, cannot disband")
	}
	withAudit(ctx, args, nil)
	ctx.server.world.Send(internal.PacketGuildDisband(plr.guild.id))
	return nil
}

func handleGuildPointsCommand(ctx *staffCommandContext, args []string) error {
	if len(args) != 1 {
		return errCommandUsage
	}
	plr, err := ctx.sender()
	if err != nil {
		return err
	}
	if plr.guild == nil {
		return fmt.Errorf("Not in guild, cannot disband")
	}
	points, err := parseIntArg(args[0], "Points")
	if err != nil {
		return err
	}
	withAudit(ctx, args, nil)
	ctx.server.world.Send(internal.PacketGuildPointsUpdate(plr.guild.id, int32(points)))
	return nil
}

func handleTestMobCommand(ctx *staffCommandContext, args []string) error {
	if len(args) != 0 {
		return errCommandUsage
	}
	plr, _, inst, err := currentFieldInstance(ctx)
	if err != nil {
		return err
	}
	withAudit(ctx, args, nil)
	return inst.lifePool.spawnMobFromID(5100001, plr.pos, true, true, true, constant.MobSummonTypeInstant, plr.ID)
}

func handlePacketCommand(ctx *staffCommandContext, args []string) error {
	if len(args) != 1 {
		return errCommandUsage
	}
	data, err := hex.DecodeString(args[0])
	if err != nil {
		log.Println("Error in decoding string for gm command packet:", args[0])
		return nil
	}
	withAudit(ctx, args, nil)
	log.Println("Sent packet:", hex.EncodeToString(data))
	ctx.conn.Send(append(make([]byte, 4), data...))
	return nil
}
