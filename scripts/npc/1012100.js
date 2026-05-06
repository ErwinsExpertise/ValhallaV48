var LETTER = 4031010;
var PROOF = 4031012;
var BLACK_CHARM = 4031059;
var STRENGTH_NECKLACE = 4031057;

function rand(min, max) {
    return Math.floor(Math.random() * (max - min + 1)) + min;
}

function spentEnoughSp() {
    return plr.getRemainingSP() <= (plr.getLevel() - 30) * 3;
}

function bowmanSecondJobMenu() {
    var explain = npc.sendMenu(
        "Before I advance you, choose which style of bowman suits you best.#b",
        "Explain Hunter",
        "Explain Crossbowman",
        "I want to choose now"
    );

    if (explain === 0) {
        npc.sendOk("Hunters emphasize fast bow attacks and consistent pressure from range.");
        return;
    }
    if (explain === 1) {
        npc.sendOk("Crossbowmen strike more slowly, but with heavier hits and stronger control at range.");
        return;
    }

    var pick = npc.sendMenu("Choose your 2nd job.#b", "Hunter", "Crossbowman");
    var jobId = pick === 0 ? 310 : 320;
    var name = pick === 0 ? "Hunter" : "Crossbowman";

    if (!npc.sendYesNo("Do you want to become a #b" + name + "#k? Once you decide, you cannot go back.")) {
        npc.sendOk("Come back when you've made your decision.");
        return;
    }
    if (!spentEnoughSp()) {
        npc.sendOk("You still have too much unused SP. Use the SP you've earned up to level 30 first, then come back.");
        return;
    }
    if (!plr.haveItem(PROOF, 1)) {
        npc.sendOk("You need #b#t4031012##k from my instructor before I can advance you.");
        return;
    }

    plr.gainItem(PROOF, -1);
    plr.setJob(jobId);
    plr.giveSP(1);
    plr.giveMaxHP(rand(300, 350));
    plr.giveMaxMP(rand(150, 200));
    plr.increaseSlotSize(4, 4);
    npc.sendOk("Good. From this point on, you walk the road of the #b" + name + "#k.");
}

if (plr.job() === 0) {
    npc.sendOk("If you want to become a Bowman, you need to be at least #bLevel 10#k with at least #b25 DEX#k.");
    if (plr.getLevel() < 10 || plr.getDex() < 25) {
        npc.sendOk("Train more before asking to become a Bowman.");
    } else if (!npc.sendYesNo("You have the sharp eyes and steady hands of a Bowman. Do you want to become a #bBowman#k?")) {
        npc.sendOk("Come back when you've made up your mind.");
    } else if (plr.getEquipInventoryFreeSlot() < 1 || plr.getUseInventoryFreeSlot() < 3) {
        npc.sendOk("Make sure you have at least one free Equip slot and three free Use slots first.");
    } else if (!plr.gainItem(1452051, 1) || !plr.gainItem(2060000, 6000)) {
        npc.sendOk("Make sure you have enough free Equip and Use inventory space first.");
    } else {
        plr.increaseSlotSize(1, 4);
        plr.increaseSlotSize(2, 4);
        plr.setJob(300);
        plr.giveMaxHP(rand(100, 150));
        plr.giveMaxMP(rand(30, 50));
        plr.giveSP(1);
        npc.sendOk("You are a Bowman now. Your Equip and Use inventories have grown, and you're much stronger than before.");
    }
} else if (plr.job() === 300) {
    if (plr.getLevel() < 30) {
        npc.sendOk("Keep training. Come back when you're ready for your second job advancement.");
    } else if (plr.haveItem(LETTER, 1)) {
        npc.sendOk("Take my letter to #b#p1072002##k near #b#m106010000##k. She'll administer the test for me.");
    } else if (plr.haveItem(PROOF, 1)) {
        bowmanSecondJobMenu();
    } else if (!npc.sendYesNo("You've grown into a fine Bowman. Do you want to take the 2nd job test?")) {
        npc.sendOk("Come back when you're ready.");
    } else if (plr.getEtcInventoryFreeSlot() < 1) {
        npc.sendOk("You need at least one free Etc slot for my letter.");
    } else {
        plr.gainItem(LETTER, 1);
        npc.sendOk("Take this letter to #b#p1072002##k near #b#m106010000##k. She'll explain the rest.");
    }
} else if (plr.job() === 310 || plr.job() === 320) {
    var state = plr.questData(7500);
    if (state === "s") {
        plr.setQuestData(7500, "p1");
        npc.sendOk("Rene told me about you. In Sleepy Dungeon, there's a secret passage only you can enter. Defeat my clone there and bring me #b#t4031059##k.");
    } else if (state === "p1") {
        if (!plr.haveItem(BLACK_CHARM, 1)) {
            npc.sendOk("The secret passage is in Sleepy Dungeon. Defeat my clone there and bring back #b#t4031059##k.");
        } else if (plr.getEtcInventoryFreeSlot() < 1) {
            npc.sendOk("Free up at least one Etc slot first.");
        } else {
            plr.gainItem(BLACK_CHARM, -1);
            plr.gainItem(STRENGTH_NECKLACE, 1);
            plr.setQuestData(7500, "p2");
            npc.sendOk("You've proven your strength. Take #b#t4031057##k to #b#p2020010##k in El Nath for the second test.");
        }
    } else if (state === "p2") {
        if (!plr.haveItem(STRENGTH_NECKLACE, 1)) {
            if (plr.getEtcInventoryFreeSlot() < 1) {
                npc.sendOk("Free up at least one Etc slot first.");
            } else {
                plr.gainItem(STRENGTH_NECKLACE, 1);
                npc.sendOk("You lost #b#t4031057##k, so I'm giving you another one. Keep it safe this time.");
            }
        } else {
            npc.sendOk("Take #b#t4031057##k to #b#p2020010##k in El Nath. She'll handle the next test.");
        }
    } else {
        npc.sendOk("Let your aim keep improving.");
    }
} else {
    npc.sendOk("Let your aim keep improving.");
}
