if (!plr.questCompleted(3034)) {
    npc.sendOk("Go away, I'm trying to meditate.");
} else if (!npc.sendYesNo("You've been so much of a help to me... If you have any Dark Crystal Ore, I can refine it for you for only #b500000 meso#k each.")) {
    npc.sendOk("Use it wisely.");
} else {
    var qty = npc.sendNumber("Okay, so how many do you want me to make?", 1, 1, 100);
    var totalCost = 500000 * qty;
    var totalOre = 10 * qty;

    if (plr.getMesos() < totalCost) {
        npc.sendOk("I'm sorry, but I am NOT doing this for free.");
    } else if (!plr.haveItem(4004004, totalOre)) {
        npc.sendOk("I need that ore to refine the Crystal. No exceptions..");
    } else if (!plr.canHold(4005004, qty)) {
        npc.sendOk("Are you having trouble with no empty slots on your inventory? Sort that out first!");
    } else {
        plr.gainItem(4004004, -totalOre);
        plr.gainMesos(-totalCost);
        plr.gainItem(4005004, qty);
        npc.sendOk("Use it wisely.");
    }
}
