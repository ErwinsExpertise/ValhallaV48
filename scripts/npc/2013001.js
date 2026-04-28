var mapId = plr.mapID();

function stageKey(id) {
    return id + "stageclear";
}

function isStageClear(id) {
    return map.properties()[stageKey(id)] === true || map.properties()[stageKey(id)] === "true";
}

function setStageClear(id) {
    map.properties()[stageKey(id)] = true;
    map.properties()["clear"] = true;
}

function leaderOnly() {
    if (!plr.isLeader()) {
        npc.sendOk("Please ask your party leader to speak with me.");
        return false;
    }
    return true;
}

if (mapId === 920010000) {
    if (isStageClear(1)) {
        plr.warp(920010100);
    } else if (leaderOnly()) {
        plr.partyGiveExp(6000);
        setStageClear(1);
        map.showEffect("quest/party/clear");
        map.playSound("Party1/Clear");
        plr.warpEventMembers(920010100);
        plr.sendMessage("The way to the Door of the Goddess's Tower is now open.");
    }
} else if (mapId === 920010100) {
    var props = map.properties();
    var scars = ["scar1", "scar2", "scar3", "scar4", "scar5", "scar6"];
    var fixed = true;
    for (var i = 0; i < scars.length; i++) {
        if (!props[scars[i]]) {
            fixed = false;
            break;
        }
    }

    if (!leaderOnly()) {
        // already handled
    } else if (fixed) {
        map.showEffect("quest/party/clear");
        map.playSound("Party1/Clear");
        plr.warpEventMembers(920010800);
        plr.sendMessage("The Goddess's statue has been restored. The party moves onward.");
    } else {
        npc.sendOk("Fix the Goddess statue before speaking with me again.");
    }
} else if (mapId === 920010200) {
    if (!leaderOnly()) {
        // already handled
    } else if (isStageClear(3)) {
        npc.sendOk("Please return to the Center Tower and continue the quest.");
    } else if (plr.haveItem(4001050, 30)) {
        plr.gainItem(4001050, -30);
        plr.gainItem(4001044, 1);
        plr.partyGiveExp(7500);
        setStageClear(3);
        map.showEffect("quest/party/clear");
        map.playSound("Party1/Clear");
        plr.sendMessage("The first statue piece has been restored.");
    } else {
        npc.sendOk("This is the walkway of the Tower of Goddess. Bring me #b30 1st Small Pieces#k from the Pixies and I will reassemble the first statue piece.");
    }
} else if (mapId === 920010300) {
    if (!leaderOnly()) {
        // already handled
    } else if (isStageClear(4)) {
        npc.sendOk("Please return to the Center Tower and continue the quest.");
    } else if (plr.haveItem(4001045, 1)) {
        plr.partyGiveExp(7500);
        setStageClear(4);
        map.showEffect("quest/party/clear");
        map.playSound("Party1/Clear");
        plr.sendMessage("The second statue piece has been recovered.");
    } else {
        npc.sendOk("This was once the storage of the Tower of Goddess. Defeat the Cellions and bring back #bStatue of Goddess: 2nd Piece#k.");
    }
} else if (mapId === 920010400) {
    var todayNames = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"];
    var today = new Date().getDay() + 1;
    if (!leaderOnly()) {
        // already handled
    } else if (isStageClear(5)) {
        npc.sendOk("Please return to the Center Tower and continue the quest.");
    } else if (map.reactorStateByName("music") === today) {
        setStageClear(5);
        plr.gainItem(4001046, 1);
        plr.partyGiveExp(7500);
        map.showEffect("quest/party/clear");
        map.playSound("Party1/Clear");
        plr.sendMessage("That was the right song. The spirit of the Goddess has answered you.");
    } else {
        npc.sendOk("This is the lobby where Minerva once listened to music. Today is #b" + todayNames[today - 1] + "#k. Play the correct song and then return to me.");
    }
} else if (mapId === 920010500) {
    if (!leaderOnly()) {
        // already handled
    } else {
        npc.sendOk("This is the sealed room. Solve the platform combination and then speak with me again.");
    }
} else if (mapId === 920010600) {
    if (!leaderOnly()) {
        // already handled
    } else if (isStageClear(7)) {
        npc.sendOk("Please return to the Center Tower and continue the quest.");
    } else if (plr.haveItem(4001052, 40)) {
        plr.gainItem(4001052, -40);
        plr.gainItem(4001048, 1);
        plr.partyGiveExp(7500);
        setStageClear(7);
        map.showEffect("quest/party/clear");
        map.playSound("Party1/Clear");
        plr.sendMessage("The fifth statue piece has been restored.");
    } else {
        npc.sendOk("This is the lounge area of the Tower of Goddess. Search the area and bring me #b40 5th Small Pieces#k.");
    }
} else if (mapId === 920010700) {
    if (!leaderOnly()) {
        // already handled
    } else if (isStageClear(8)) {
        npc.sendOk("Please return to the Center Tower and continue the quest.");
    } else {
        plr.gainItem(4001049, 1);
        plr.partyGiveExp(7500);
        setStageClear(8);
        map.showEffect("quest/party/clear");
        map.playSound("Party1/Clear");
        plr.sendMessage("You've secured another statue piece. Return to the center tower.");
    }
} else if (mapId === 920011000) {
    npc.sendOk("This is the Room of Darkness. You won't find a statue piece here. Search the other rooms instead.");
} else if (mapId === 920011200) {
    npc.sendOk("See you next time!")
    plr.warp(200080101)
} else {
    npc.sendOk("Keep going. The Tower of the Goddess still holds more trials ahead.");
}
