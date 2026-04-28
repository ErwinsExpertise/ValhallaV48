var recommendationLetter = 4031801;

if (!npc.sendYesNo("Take this ship and you'll head off to a bigger continent. For #e150 mesos#n, I'll take you to #bVictoria Island#k. The thing is, once you leave this place, you can't ever come back. What do you think? Do you want to go to Victoria Island?")) {
    npc.sendOk("Hmm... I guess you still have things to do here?");
} else if (plr.haveItem(recommendationLetter, 1)) {
    plr.gainItem(recommendationLetter, -1);
    plr.warpToPortalName(104000000, "maple00");
} else if (plr.getLevel() <= 6) {
    npc.sendOk("Let's see... I don't think you are strong enough. You'll have to be at least Level 7 to go to Victoria Island.");
} else {
    if (plr.getMesos() < 150) {
        npc.sendOk("What? You're telling me you wanted to go without any money? You're one weirdo...");
    } else {
        plr.gainMesos(-150);
        plr.warpToPortalName(104000000, "maple00");
    }
}
