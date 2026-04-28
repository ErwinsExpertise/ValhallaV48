var pay = 2000;
var ticket = 4031134;

npc.sendSelection("Have you heard of the beach with a spectacular view of the ocean called #b#m110000000##k, located a little far from #m" + plr.mapID() + "#? I can take you there right now for either #b" + pay + " mesos#k, or if you have #b#t" + ticket + "##k with you, in which case you'll be in for free.\r\n\r\n#L0##bI'll pay " + pay + " mesos.#k#l\r\n#L1##bI have #t" + ticket + "##k#l\r\n#L2##bWhat is #t" + ticket + "#?#k#l");

var selection = npc.selection();
if (selection === 0 || selection === 1) {
    var msg = selection === 0
        ? "You want to pay #b" + pay + " mesos#k and leave for #m110000000#?"
        : "So you have #b#t" + ticket + "##k? You can always head over to #m110000000# with that.";
    if (!npc.sendYesNo(msg + " Okay! Please beware that you may run into some monsters there, so make sure not to get caught off guard. Would you like to head over to #m110000000# right now?")) {
        npc.sendOk("You must have some business to take care of here. You must be tired from all that traveling and hunting. Go take some rest, and if you feel like changing your mind, then come talk to me.");
    } else if (selection === 0) {
        if (plr.getMesos() < pay) {
            npc.sendOk("I think you're lacking mesos. There are many ways to gather up some money, you know, like ... selling your armor ... defeating the monsters ... doing quests ... you know what I'm talking about.");
        } else {
            plr.gainMesos(-pay);
            plr.saveLocation("FLORINA");
            plr.warp(110000000);
        }
    } else if (!plr.haveItem(ticket, 1)) {
        npc.sendOk("Hmmm, so where exactly is #b#t" + ticket + "##k?? Are you sure you have them? Please double-check.");
    } else {
        plr.saveLocation("FLORINA");
        plr.warp(110000000);
    }
} else if (selection === 2) {
    npc.sendNext("You must be curious about #b#t" + ticket + "##k. Yeah, I can see that. #t" + ticket + "# is an item where as long as you have in possession, you may make your way to #m110000000# for free. It's such a rare item that even we had to buy those, but unfortunately I lost mine a few weeks ago during a long weekend.");
    npc.sendNext("I came back without it, and it just feels awful not having it. Hopefully someone picked it up and put it somewhere safe. Anyway, that's my story, and who knows, you may be able to pick it up and put it to good use. If you have any questions, feel free to ask.");
}
