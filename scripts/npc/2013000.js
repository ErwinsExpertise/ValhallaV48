var minLevel = 51;
var maxLevel = 70;
var minPlayers = 5;
var maxPlayers = 6;
var foodArray = [2001001, 2001002, 2020000, 2020003, 2020001, 0];
var blessingArray = [2022090, 2022091, 2022092, 2022093];

function hasBalancedParty(members) {
    var warrior = 0;
    var mage = 0;
    var bowman = 0;
    var thief = 0;

    for (var i = 0; i < members.length; i++) {
        var job = members[i].job();
        var branch = Math.floor(job / 100);
        if (branch === 1) warrior++;
        else if (branch === 2) mage++;
        else if (branch === 3) bowman++;
        else if (branch === 4) thief++;
    }

    return warrior > 0 && mage > 0 && bowman > 0 && thief > 0;
}

if (plr.mapID() === 200080101) {
    npc.sendSelection("Hello, I am Wonky the Fairy. What would you like to do today?#b\r\n#L0#Apply for entrance.#l\r\n#L1#Give Wonky something to eat.#l");
    var sel = npc.selection();
    if (sel === 0) {
        var gmSolo = plr.isGM();
        if (!plr.inParty() && !gmSolo) {
            npc.sendOk("Talk to me after you've formed a party. And also #r#eGIVE ME FOOD!!#n");
        } else if (!plr.isLeader() && !gmSolo) {
            npc.sendOk("Please ask your party leader to talk to me.");
        } else {
            var members = gmSolo ? [plr] : plr.partyMembersOnMap();
            var ok = gmSolo || (members.length >= minPlayers && members.length <= maxPlayers);
            if (ok && !gmSolo) {
                for (var i = 0; i < members.length; i++) {
                    var level = members[i].getLevel();
                    if (level < minLevel || level > maxLevel) {
                        ok = false;
                        break;
                    }
                }
            }
            if (!ok) {
                npc.sendOk("Either your party members are not all in the map, or they are not in the right level range.");
            } else {
                if (gmSolo) {
                    plr.startPartyQuest("orbis_pq", 0);
                } else {
                    plr.startPartyQuest("orbis_pq", 0);
                }
                if (!gmSolo && hasBalancedParty(members)) {
                    for (var j = 0; j < members.length; j++) {
                        members[j].gainItem(blessingArray[Math.floor(Math.random() * blessingArray.length)], 1);
                    }
                }
            }
        }
    } else if (sel === 1) {
        npc.sendSelection("Aww cool what have you got for me?#b\r\n#L0#Ice Cream Pop#l\r\n#L1#Red Bean Sundae#l\r\n#L2#Salad#l\r\n#L3#Pizza#l\r\n#L4#Fried Chicken#l\r\n#L5#Nothing...#l#k");
        var food = npc.selection();
        if (food >= 0 && food <= 4 && plr.haveItem(foodArray[food], 1)) {
            plr.gainItem(foodArray[food], -1);
            npc.sendOk("Thank you for feeding me, but I'M STILL HUNGRY!!!");
        } else if (food === 5) {
            npc.sendOk("What? Are you playing tricks on me?!");
        } else {
            npc.sendOk("What?! Where's the food?!");
        }
    }
} else if (plr.mapID() === 920010000) {
    if (npc.sendYesNo("Would you like to exit the Party Quest?\r\nYou will have to start again next time...")) {
        plr.leavePartyQuest();
    }
}
