var pqMaps = [925100000, 925100100, 925100200, 925100201, 925100202, 925100300, 925100301, 925100302, 925100400, 925100500, 925100600];

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
    for (var i = 0; i < pqMaps.length; i++) {
        if (map.playerCount(pqMaps[i], 0) > 0) {
            return true;
        }
    }
    return false;
}

if (plr.mapID() !== 251010404) {
    npc.sendOk("Please gather your party on the deck above the pirate ship first.");
} else {
    npc.sendSelection("What would you like to do?#b\r\n#L0#Enter the Pirate Party Quest.#l\r\n#L1#Tell me about this place.#l#k");
    var sel = npc.selection();

    if (sel === 1) {
        npc.sendOk("Lord Pirate kidnapped Wu Yang and forced the bellflowers under his command. If your party can board the ship and defeat him, Herb Town may finally have peace again.");
    } else if (!plr.inParty()) {
        npc.sendOk("You need a party to challenge Lord Pirate.");
    } else if (!plr.isLeader()) {
        npc.sendOk("Please ask your party leader to speak to me.");
    } else {
        var members = partyOnThisMap();
        if (members.length < 3 || members.length > 6) {
            npc.sendOk("Your party must have between 3 and 6 members here on this map.");
        } else if (pqOccupied()) {
            npc.sendOk("Another party is already on the pirate ship. Please try again later.");
        } else {
            var valid = true;
            for (var i = 0; i < members.length; i++) {
                var level = members[i].level();
                if (level < 55 || level > 100) {
                    valid = false;
                    break;
                }
            }

            if (!valid) {
                npc.sendOk("Everyone in your party must be between Level 55 and 100.");
            } else {
                plr.startPartyQuest("pirate_pq", 0);
            }
        }
    }
}
