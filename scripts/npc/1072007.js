var LETTER = 4031011;
var PROOF = 4031012;
var MARBLE = 4031013;

if (plr.job() !== 400 || plr.getLevel() < 30) {
    npc.sendOk("You shouldn't be here. I'll send you back.");
    plr.warp(102040000);
} else if (plr.itemCount(MARBLE) >= 30) {
    plr.removeAll(MARBLE);
    plr.gainItem(LETTER, -1);
    if (!plr.gainItem(PROOF, 1)) {
        npc.sendOk("Something is wrong. Make sure you still have the letter and at least one free Etc slot.");
    } else {
        npc.sendOk("Excellent. You passed. Take #b#t4031012##k back to the Dark Lord in Kerning City.");
        plr.warp(102040000);
    }
} else if (!npc.sendYesNo("You still haven't collected 30 #b#t4031013##k. If you give up now, you can leave and try again later. Do you want to give up and leave?")) {
    npc.sendOk("Then keep fighting and come back when you've collected 30 Dark Marbles.");
} else {
    npc.sendOk("I'll send you out. Don't give up; you can always try again.");
    plr.warp(102040000);
}
