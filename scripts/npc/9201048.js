var maps = [670010200, 670010300, 670010301, 670010302, 670010400, 670010500, 670010600, 670010700, 670010750, 670010800, 670011000];

function partyOnThisMap() {
    var all = plr.partyMembersOnMap();
    var here = [];
    for (var i = 0; i < all.length; i++) {
        if (all[i].mapID() === plr.mapID()) {
            here.push(all[i]);
        }
    }
    return here;
}

function pqOccupied() {
    for (var i = 0; i < maps.length; i++) {
        if (map.getMap(maps[i], 0).playerCount(maps[i], 0) > 0) {
            return true;
        }
    }
    return false;
}

if (plr.mapID() !== 670010100) {
    npc.sendOk("Please wait for your party in the entrance room first.");
} else {
    npc.sendSelection("Okay. What would you like to do?#b\r\n#L0#I'd like to start the Party Quest.#l\r\n#L1#Please get us out of here!#l#k");
    var sel = npc.selection();
    if (sel === 1) {
        plr.warp(670010000);
    } else if (!plr.inParty()) {
        npc.sendOk("Form a party of 6 before talking to me.");
    } else if (!plr.isLeader()) {
        npc.sendOk("Please ask your party leader to talk to me.");
    } else {
        var members = partyOnThisMap();
        if (members.length !== 6) {
            npc.sendOk("Your party must have exactly 6 members on this map.");
        } else {
            var badLevel = false;
            for (var i = 0; i < members.length; i++) {
                if (members[i].level() < 40) {
                    badLevel = true;
                    break;
                }
            }

            if (badLevel) {
                npc.sendOk("Someone in your party is below Level 40. Please check again.");
            } else if (pqOccupied()) {
                npc.sendOk("Another party is already taking on the Amorian Challenge. Please try again later.");
            } else {
                plr.startPartyQuest("amoria_pq", 0);
            }
        }
    }
}
