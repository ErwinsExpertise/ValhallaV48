var minLevel = 51;
var maxLevel = 70;
var minPlayers = 6;
var maxPlayers = 6;
var foodItems = [2020001, 2022001, 2020004, 2022003, 2001001, 2020028, 2001002];

if (plr.mapID() === 200080101) {
    npc.sendSelection("I'm Wonky the Fairy. What do you want?#b\r\n#L0#Register for the Orbis Party Quest.#l\r\n#L1#Tell me about the Tower of Goddess.#l\r\n#L2#Give Wonky something to eat.#l#k");
    var sel = npc.selection();

    if (sel === 0) {
        if (!plr.inParty()) {
            npc.sendOk("You need to be in a party to enter the Tower of Goddess.");
        } else if (!plr.isLeader()) {
            npc.sendOk("Please ask your party leader to speak with me.");
        } else if (plr.partyQuestActive()) {
            npc.sendOk("Your party is already in the middle of a party quest.");
        } else {
            var members = plr.partyMembersOnMap();
            if (members.length < minPlayers || members.length > maxPlayers) {
                npc.sendOk("Your party must have exactly 6 members here to enter.");
            } else {
                var ok = true;
                for (var i = 0; i < members.length; i++) {
                    var level = members[i].level();
                    if (level < minLevel || level > maxLevel) {
                        ok = false;
                        break;
                    }
                }

                if (!ok) {
                    npc.sendOk("Everyone in the party must be between level 51 and 70.");
                } else {
                    plr.startPartyQuest("orbis_pq", -1);
                    npc.sendOk("Please make your way inside and rescue Goddess Minerva.");
                }
            }
        }
    } else if (sel === 1) {
        npc.sendOk("A mysterious tower appeared behind the statue of Goddess Minerva in Orbis. Gather six statue pieces, restore the statue, find the #bGrass of Life#k, and save the Goddess.");
    } else if (sel === 2) {
        var text = "Hey, did you bring me something tasty?#b";
        for (var j = 0; j < foodItems.length; j++) {
            text += "\r\n#L" + j + "##t" + foodItems[j] + "##l";
        }
        npc.sendSelection(text + "#k");
        var food = npc.selection();
        if (food >= 0 && food < foodItems.length && plr.haveItem(foodItems[food], 1)) {
            plr.gainItem(foodItems[food], -1);
            npc.sendOk("Mmm, that helps... but I'm still hungry.");
        } else {
            npc.sendOk("What? You don't have it!");
        }
    }
} else if (npc.sendYesNo("Would you like to leave the Orbis Party Quest? You will have to start over next time.")) {
    plr.leavePartyQuest();
}
