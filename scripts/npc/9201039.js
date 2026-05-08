var freeHairCoupon = 4031528;
var questStateKey = 8861;
var maleHair = [30032, 30020, 30000, 30132, 30192, 30240, 30162, 30270, 30112];
var femaleHair = [31050, 31040, 31030, 31001, 31070, 31310, 31091, 31250, 31150];

if (plr.questData(questStateKey) === "end") {
    npc.sendOk("I've already done your hair once as a trade-for-services. You'll need an EXP Hair coupon from the Cash Shop if you want another change.");
} else if (!npc.sendYesNo("Ready for an awesome hairdo? Just say the word and we'll get started.")) {
    npc.sendOk("All right, I'll give you a minute.");
} else if (!plr.haveItem(freeHairCoupon, 1)) {
    npc.sendOk("Are you sure you have the free coupon I need? I can't do the haircut without it.");
} else {
    var styles = plr.gender() < 1 ? maleHair : femaleHair;
    plr.gainItem(freeHairCoupon, -1);
    plr.setHair(styles[Math.floor(Math.random() * styles.length)]);
    plr.setQuestData(questStateKey, "end");
    npc.sendOk("Not bad, if I do say so myself. I knew those books I studied would come in handy.");
}
