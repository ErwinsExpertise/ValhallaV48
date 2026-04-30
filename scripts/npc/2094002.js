function leavePrompt() {
    if (npc.sendYesNo("Would you like to give up and leave the pirate ship?")) {
        plr.leaveEvent();
    }
}

function stage2Flow() {
    var state = String(plr.getEventProperty("mobGen") || "0");

    if (!plr.isPartyLeader()) {
        if (state === "end") {
            npc.sendOk("The final seal is gone. Move through the portal on the far right.");
        } else {
            npc.sendOk("Please ask your party leader to handle the seals on this stage.");
        }
        return;
    }

    if (state === "0") {
        npc.sendOk("Lord Pirate sealed the next portal three times. Defeat the pirates and bring me #b20 #t4001120##k first.");
        plr.setEventProperty("mobGen", "1");
        map.setMobSpawnEnabled(9300114, true);
        map.setMobSpawnEnabled(9300115, false);
        map.setMobSpawnEnabled(9300116, false);
        return;
    }

    if (state === "1") {
        if (!plr.haveItem(4001120, 20)) {
            npc.sendOk("Bring me #b20 #t4001120##k first.");
            return;
        }
        plr.removeItemsByID(4001120, 20);
        map.setMobSpawnEnabled(9300114, false);
        map.removeMobsByID(9300114);
        map.message("Guon has disabled the portal's first seal.");
        plr.setEventProperty("mobGen", "2");
        npc.sendOk("The first seal is gone. Speak to me again when you are ready for the next wave.");
        return;
    }

    if (state === "2") {
        map.setMobSpawnEnabled(9300115, true);
        plr.setEventProperty("mobGen", "3");
        npc.sendOk("Now bring me #b20 #t4001121##k.");
        return;
    }

    if (state === "3") {
        if (!plr.haveItem(4001121, 20)) {
            npc.sendOk("Bring me #b20 #t4001121##k first.");
            return;
        }
        plr.removeItemsByID(4001121, 20);
        map.setMobSpawnEnabled(9300115, false);
        map.removeMobsByID(9300115);
        map.message("Guon has disabled the portal's second seal.");
        plr.setEventProperty("mobGen", "4");
        npc.sendOk("The second seal is gone. Speak to me again when you are ready for the last wave.");
        return;
    }

    if (state === "4") {
        map.setMobSpawnEnabled(9300116, true);
        plr.setEventProperty("mobGen", "5");
        npc.sendOk("One more time. Bring me #b20 #t4001122##k.");
        return;
    }

    if (state === "5") {
        if (!plr.haveItem(4001122, 20)) {
            npc.sendOk("Bring me #b20 #t4001122##k first.");
            return;
        }
        plr.removeItemsByID(4001122, 20);
        map.setMobSpawnEnabled(9300116, false);
        map.removeMobsByID(9300116);
        map.message("Guon has disabled the final seal. Move to the next stage immediately.");
        plr.setEventProperty("mobGen", "end");
        npc.sendOk("The final seal is broken. Hurry to the next portal.");
        return;
    }

    npc.sendOk("You have already finished this stage. Move on.");
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
