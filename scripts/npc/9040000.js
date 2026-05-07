var GQItems = [1032033, 4001024, 4001025, 4001026, 4001027, 4001028, 4001029, 4001030, 4001031, 4001032, 4001033, 4001034, 4001035, 4001036, 4001037];

function clearItems(target) {
    for (var i = 0; i < GQItems.length; i++) {
        target.removeAll(GQItems[i]);
    }
}

function hasQuestItems(target) {
    for (var i = 0; i < GQItems.length; i++) {
        if (target.itemCount(GQItems[i]) > 0) {
            return true;
        }
    }
    return false;
}

function tryStartGuildQuest() {
    if (!plr.inGuild() || plr.guildRank() >= 3) {
        npc.sendOk("Only a Master or Jr. Master of the guild can start a Guild Quest.");
        return;
    }
    if (plr.guildQuestActive()) {
        npc.sendOk("Your guild already has an active Guild Quest.");
        return;
    }

    var members = plr.guildMembersOnMap();
    if (members.length < 6) {
        npc.sendOk("You need at least 6 guild members from the same guild here to begin.");
        return;
    }
    if (members.length > 30) {
        npc.sendOk("A Guild Quest can only start with up to 30 participants.");
        return;
    }

    for (var i = 0; i < members.length; i++) {
        if (hasQuestItems(members[i])) {
            npc.sendOk("One of your guild members is already carrying Guild Quest items. Please clean up the party first.");
            return;
        }
    }

    for (var j = 0; j < members.length; j++) {
        clearItems(members[j]);
    }

    plr.startGuildQuest("guild_pq", -1);
}

npc.sendSelection("The path to Sharenian starts here. What would you like to do?#b\r\n#L0#Start a Guild Quest#l\r\n#L1#Join your guild's Guild Quest#l");
var selection = npc.selection();

if (selection === 0) {
    tryStartGuildQuest();
} else if (selection === 1) {
    if (!plr.inGuild()) {
        npc.sendOk("You must be in a guild to join a Guild Quest.");
    } else if (!plr.guildQuestActive()) {
        npc.sendOk("Your guild does not currently have an active Guild Quest.");
    } else if (!plr.joinGuildQuest()) {
        npc.sendOk("The door has already opened. You can no longer enter this run.");
    } else {
        clearItems(plr);
        plr.warpToPortalName(990000000, "sp");
    }
}
