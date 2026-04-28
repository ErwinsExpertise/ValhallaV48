var GQItems = [4001024, 4001025, 4001026, 4001027, 4001028, 4001031, 4001032, 4001033, 4001034, 4001035, 4001037];

npc.sendSelection("The path to Sharenian starts here. What would you like to do?#b\r\n#L0#Start a Guild Quest#l\r\n#L1#Join your guild's Guild Quest#l");
var selection = npc.selection();

if (selection === 0) {
    if (!plr.inGuild() || plr.guildRank() >= 3) {
        npc.sendOk("Only a Master or Jr. Master of the guild can start an instance.");
    } else {
        for (var i = 0; i < GQItems.length; i++) {
            plr.removeAll(GQItems[i]);
        }
		plr.startGuildQuest("guild_pq", 0);
    }
} else if (selection === 1) {
    if (!plr.inGuild()) {
        npc.sendOk("You must be in a guild to join an instance.");
    } else if (!plr.joinGuildQuest()) {
        npc.sendOk("Your guild is currently not registered for an instance.");
    } else {
        if (plr.getEventProperty("canEnter") === false || plr.getEventProperty("canEnter") === "false") {
            npc.sendOk("I'm sorry, but the guild has gone on without you. Try again later.");
            return;
        }
        for (var j = 0; j < GQItems.length; j++) {
            plr.removeAll(GQItems[j]);
        }
        plr.warp(990000000);
    }
}
