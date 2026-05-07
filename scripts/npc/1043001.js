var QUEST = 2051;
var ROOT = 4031032;

var status = plr.getQuestStatus(QUEST);

if (plr.getLevel() <= 49) {
    npc.sendOk("Among these herbal bushes, there are roots with a mysterious energy, but the strange aura around them makes them impossible to pick.");
} else if (status === 1) {
    if (npc.sendYesNo("Are you sure you want to take #b#t4031032##k with you?")) {
        if (plr.giveItem(ROOT, 1)) {
            plr.warp(101000000);
        } else {
            npc.sendOk("Your Etc inventory seems full. Please make room before picking up the item.");
        }
    }
} else if (status === 2) {
    if (plr.getEtcInventoryFreeSlot() < 1 || plr.getEquipInventoryFreeSlot() < 1) {
        npc.sendOk("Your Equip and Etc inventories are too full to hold what you might find in these herbal bushes. Free some space and try again.");
    } else {
        var roll = Math.floor(Math.random() * 30) + 1;
        var rewardId;
        var rewardCount;

        if (roll <= 10) {
            rewardId = 4020007;
            rewardCount = 2;
        } else if (roll <= 20) {
            rewardId = 4020008;
            rewardCount = 2;
        } else if (roll <= 29) {
            rewardId = 4010006;
            rewardCount = 2;
        } else {
            rewardId = 1032013;
            rewardCount = 1;
        }

        if (plr.gainItem(rewardId, rewardCount)) {
            plr.warp(101000000);
        } else {
            npc.sendOk("Your Equip and Etc inventories are too full to hold what you found. Free some space and try again.");
        }
    }
} else {
    npc.sendOk("Among these herbal bushes, you'll find roots with a mysterious energy, but since #b#p1061005##k never explained which root to take, there's no way to know what to pick.");
}
