npc.sendSelection("Have you heard of the beach with a spectacular view of the ocean called #bFlorina Beach#k, located near Lith Harbor? I can take you there right now for either #b1500 mesos#k, or if you have a #bVIP Ticket to Florina Beach#k with you, in which case you'll be there for free.\r\n\r\n#L0##b I'll pay 1500 mesos.#l\r\n#L1# I have a VIP Ticket to Florina Beach.#l\r\n#L2# What is a VIP Ticket to Florina Beach#k?#l");

var selection = npc.selection();
if (selection === 0) {
    if (!npc.sendYesNo("So you want to pay #b1500 mesos#k and leave for Florina Beach? Alright, then, but just be aware that you may be running into some monsters around there, too. Okay, would you like to head over to Florina Beach right now?")) {
        npc.sendOk("You must have some business to take care of here. You must be tired from all that traveling and hunting. Go take some rest, and if you feel like changing your mind, then come talk to me.");
    } else if (plr.getMesos() < 1500) {
        npc.sendOk("I think you're lacking mesos. There are many ways to gather up some money, you know, like... selling your armor... defeating monsters... doing quests... you know what I'm talking about.");
    } else {
        plr.saveLocation("FLORINA");
        plr.gainMesos(-1500);
        plr.warp(110000000);
    }
} else if (selection === 1) {
    if (!npc.sendYesNo("So you have a #bVIP Ticket to Florina Beach#k? You can always head over to Florina Beach with that. Alright then, but just be aware that you may be running into some monsters there too. Okay, would you like to head over to Florina Beach right now?")) {
        npc.sendOk("You must have some business to take care of here. You must be tired from all that traveling and hunting. Go take some rest, and if you feel like changing your mind, then come talk to me.");
    } else if (!plr.haveItem(4031134, 1)) {
        npc.sendOk("Hmmm, so where exactly is your #bVIP Ticket to Florina\r\nBeach#k? Are you sure you have one? Please double-check.");
    } else {
        plr.saveLocation("FLORINA");
        plr.warp(110000000);
    }
} else if (selection === 2) {
    npc.sendNext("You must be curious about a #bVIP Ticket to Florina Beach#k. Haha, that's very understandable. A VIP Ticket to Florina Beach is an item where as long as you have in possession, you may make your way to Florina Beach for free. It's such a rare item that even we had to buy those, but unfortunately I lost mine a few weeks ago during my precious summer break.");
    npc.sendNext("I came back without it, and it just feels awful not having it. Hopefully someone picked it up and put it somewhere safe. Anyway, this is my story and who knows, you may be able to pick it up and put it to good use. If you have any questions, feel free to ask.");
}
