var QUEST = 2050;
var HERB = 4031020;

var status = plr.getQuestStatus(QUEST);

if (plr.getLevel() <= 24) {
    npc.sendOk("Among these flowers, there are some surrounded by a mysterious aura. You can't pick them because of the energy around them.");
} else if (status === 1) {
    if (npc.sendYesNo("Are you sure you want to take #b#t4031020##k with you?")) {
        if (plr.giveItem(HERB, 1)) {
            plr.warp(101000000);
        } else {
            npc.sendOk("Your Etc inventory seems full. Please make room before picking up the item.");
        }
    }
} else if (status === 2) {
    if (plr.getEtcInventoryFreeSlot() < 1) {
        npc.sendOk("You need at least one free Etc slot to take what you found among the flowers. Free some space and try again.");
    } else {
        var rewardPool = [4010000, 4010001, 4010002, 4010003, 4010004, 4010005, 4020000, 4020001, 4020002, 4020003, 4020004, 4020005, 4020006];
        var reward = rewardPool[Math.floor(Math.random() * rewardPool.length)];
        if (plr.gainItem(reward, 2)) {
            plr.warp(101000000);
        } else {
            npc.sendOk("Your Etc inventory seems full. Please make room before picking up the item.");
        }
    }
} else {
    npc.sendOk("Among all these flowers, you'll find a few surrounded by a mysterious aura, but you can't pick them because #b#p1061005##k never explained which one you're supposed to take.");
}
