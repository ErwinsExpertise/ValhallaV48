function leavePrompt() {
    if (npc.sendYesNo("Would you like to give up and leave the pirate ship?")) {
        plr.leaveEvent();
    }
}

var STAGE2_IDLE = "0";
var STAGE2_COMPLETE = "end";
var stage2Phases = [
    {
        activeState: "1",
        readyState: "2",
        itemID: 4001120,
        itemCount: 20,
        mobID: 9300114,
        startText: "The pirates will be rushing out here soon. I want you to defeat them and obtain at least #b20 #t4001120##k emblems.",
        missingText: "I do not think you have gathered up #b20 #t4001120##k emblems, yet.",
        clearMessage: "Guon has disabled the portal's first seal.",
        clearText: "I see that you have gathered up all the #b#t4001120##k emblems. Let me know when you are ready to take on more pirates."
    },
    {
        activeState: "3",
        readyState: "4",
        itemID: 4001121,
        itemCount: 20,
        mobID: 9300115,
        startText: "Okay, now I want you to find a way to gather #b20 #t4001121##k emblems. And be careful, the pirates are on their way!",
        missingText: "I do not think you have gathered up #b20 #t4001121##k emblems, yet.",
        clearMessage: "Guon has disabled the portal's second seal.",
        clearText: "I see that you have gathered up all the #b#t4001121##k emblems. Let me know when you are ready to take on more pirates."
    },
    {
        activeState: "5",
        readyState: STAGE2_COMPLETE,
        itemID: 4001122,
        itemCount: 20,
        mobID: 9300116,
        startText: "You're nearly finished. To open the last seal, you'll need #b20 #t4001122##k emblems. Hurry! The pirates are on their way!",
        missingText: "I do not think you have gathered up #b20 #t4001122##k emblems, yet.",
        clearMessage: "Guon has disabled the portal's final seal. Please move to the next spot immediately.",
        clearText: "I see that you have gathered up all the #b#t4001122##k emblems. Great work defeating these pirates. Now move to the portal located on the very right side of the ship."
    }
];

function stage2DisableAllMobs() {
    for (var i = 0; i < stage2Phases.length; i++) {
        map.setMobSpawnEnabled(stage2Phases[i].mobID, false);
    }
}

function stage2ClearMobs() {
    stage2DisableAllMobs();
    map.removeAllMobs();
}

function stage2StartPhase(phase) {
    stage2ClearMobs();
    map.setMobSpawnEnabled(phase.mobID, true);
    plr.setEventProperty("mobGen", phase.activeState);
    npc.sendOk(phase.startText);
}

function stage2CompletePhase(phase) {
    plr.removeItemsByID(phase.itemID, phase.itemCount);
    stage2ClearMobs();
    plr.setEventProperty("mobGen", phase.readyState);
    map.message(phase.clearMessage);
    npc.sendOk(phase.clearText);
}

function stage2Flow() {
    var state = String(plr.getEventProperty("mobGen") || "0");

    if (!plr.isPartyLeader()) {
        if (state === STAGE2_COMPLETE) {
            npc.sendOk("You have succeeded in unsealing the seal created by Lord Pirate. Please move to the next spot immediately.");
        } else {
            npc.sendOk("Watch out! You may see pirates popping up any minute here. That does NOT mean, however, that you can just walk past this area, thanks to Lord Pirate sealing up the portal that sends you to the next stage 3 times!\r\n\r\nTo break the seal, you'll need to acquire the #bPirate Emblem#k, an item that identifies the carrier as a pirate. What you'll need to do is defeat the pirates here, find a way to obtain the #bPirate Emblem#k, and give it to me so I can find a way to disable the seal. Please start this through the leader of your party.");
        }
        return;
    }

    if (state === STAGE2_IDLE) {
        npc.sendNext("Watch out! You may see pirates popping up around here any minute now. Thanks to Lord Pirate sealing up the portal to the next stage 3 times, you cannot just walk past this area!");
        npc.sendNext("To break the seal, you'll need to acquire the #bPirate Emblem#k, an item that identifies the carrier as a pirate. Place the emblem in front of the seal, and the seal will be automatically disarmed. Please defeat the pirates that appear and collect the emblems they drop. Once you have enough #bPirate emblems#k, hand them to me and I will break the seal for you.");
        stage2StartPhase(stage2Phases[0]);
        return;
    }

    for (var i = 0; i < stage2Phases.length; i++) {
        var phase = stage2Phases[i];

        if (state === phase.readyState && phase.readyState !== STAGE2_COMPLETE) {
            stage2StartPhase(stage2Phases[i + 1]);
            return;
        }

        if (state !== phase.activeState) {
            continue;
        }

        if (!plr.haveItem(phase.itemID, phase.itemCount)) {
            npc.sendOk(phase.missingText);
            return;
        }

        stage2CompletePhase(phase);
        return;
    }

    if (state === STAGE2_COMPLETE) {
        npc.sendOk("You have done everything you could here. Please move on to the next spot immediately.");
        return;
    }

    npc.sendOk("You have done everything you could here. Please move on to the next spot immediately.");
}

function treasureStage(reactorName, introKey, clearKey, mobA, mobB) {
    if (plr.getEventProperty(introKey) !== true) {
        plr.setEventProperty(introKey, true);
        npc.sendOk("Defeat every monster in this room, then speak to me again.");
        return;
    }

    if (map.mobCountByID(mobA) + map.mobCountByID(mobB) > 0) {
        npc.sendOk("There are still monsters left in this room.");
        return;
    }

    if (plr.getEventProperty(clearKey) !== true) {
        map.setReactorStateByName(reactorName, 1);
        plr.setEventProperty(clearKey, true);
        npc.sendOk("The treasure chest is ready, but it will still need a key.");
        return;
    }

    npc.sendOk("The treasure chest is ready. Use a key if you have one.");
}

var mapID = plr.mapID();

if (mapID === 925100700) {
    plr.warp(251010404);
} else if (mapID === 925100000) {
    npc.sendSelection("What would you like to do?#b\r\n#L0#Listen to Guon's story.#l\r\n#L1#Leave the pirate ship.#l#k");
    if (npc.selection() === 0) {
        npc.sendOk("Lord Pirate kidnapped Wu Yang and forced Herb Town's bellflowers to obey him. Please rescue her before the pirate ship escapes.");
    } else {
        leavePrompt();
    }
} else if (mapID === 925100100) {
    npc.sendSelection("What would you like to do?#b\r\n#L0#Continue the quest.#l\r\n#L1#Leave the pirate ship.#l#k");
    if (npc.selection() === 0) {
        stage2Flow();
    } else {
        leavePrompt();
    }
} else if (mapID === 925100200 || mapID === 925100300) {
    npc.sendSelection("What would you like to do?#b\r\n#L0#Listen to Guon's story.#l\r\n#L1#Leave the pirate ship.#l#k");
    if (npc.selection() === 0) {
        npc.sendOk("Lord Pirate is ahead. Defeat his crew and keep moving until you find Wu Yang.");
    } else {
        leavePrompt();
    }
} else if (mapID === 925100201) {
    treasureStage("treasure1", "d201_clear", "clear3hd", 9300112, 9300113);
} else if (mapID === 925100301) {
    treasureStage("treasure2", "d301_clear", "clear4hd", 9300112, 9300113);
} else if (mapID === 925100202 || mapID === 925100302) {
    npc.sendSelection("What would you like to do?#b\r\n#L0#Listen to Guon's story.#l\r\n#L1#Leave the pirate ship.#l#k");
    if (npc.selection() === 0) {
        npc.sendOk("This room is guarded by Lord Pirate's most trusted servants. They may be hiding something valuable.");
    } else {
        leavePrompt();
    }
} else if (mapID === 925100400) {
    npc.sendSelection("What would you like to do?#b\r\n#L0#Listen to Guon's story.#l\r\n#L1#Leave the pirate ship.#l#k");
    if (npc.selection() === 0) {
        npc.sendOk("The pirates keep pouring in through these doors. Defeat them, collect #b#t4001117##k, and seal every entrance.");
    } else {
        leavePrompt();
    }
} else if (mapID === 925100500) {
    leavePrompt();
}
