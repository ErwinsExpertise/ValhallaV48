var LETTER = 4031011;
var PROOF = 4031012;
var BLACK_CHARM = 4031059;
var STRENGTH_NECKLACE = 4031057;

function rand(min, max) {
    return Math.floor(Math.random() * (max - min + 1)) + min;
}

function spentEnoughSp() {
    return plr.getRemainingSP() <= (plr.getLevel() - 30) * 3;
}

function thiefSecondJobMenu() {
    var explain = npc.sendMenu(
        "Before I advance you, decide which side of a thief's craft suits you best.#b",
        "Explain Assassin",
        "Explain Bandit",
        "I want to choose now"
    );

    if (explain === 0) {
        npc.sendOk("Assassins specialize in claws and throwing stars. They strike quickly and precisely from range.");
        return;
    }
    if (explain === 1) {
        npc.sendOk("Bandits specialize in daggers. They fight up close with speed, control, and multi-hit attacks.");
        return;
    }

    var pick = npc.sendMenu("Choose your 2nd job.#b", "Assassin", "Bandit");
    var jobId = pick === 0 ? 410 : 420;
    var name = pick === 0 ? "Assassin" : "Bandit";

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
    plr.increaseSlotSize(2, 4);
    npc.sendOk("Good. From now on, walk the path of the #b" + name + "#k.");
}

if (plr.job() === 0) {
    npc.sendOk("If you want to become a Thief, you need to be at least #bLevel 10#k with at least #b25 DEX#k.");
    if (plr.getLevel() < 10 || plr.getDex() < 25) {
        npc.sendOk("You're still just an apprentice. Train more and come back later.");
    } else if (!npc.sendYesNo("You look like you could become one of us. Do you want to become a #bThief#k?")) {
        npc.sendOk("Choosing your path matters. Come back when you've decided.");
    } else if (plr.getEquipInventoryFreeSlot() < 2 || plr.getUseInventoryFreeSlot() < 3) {
        npc.sendOk("Make sure you have at least two free Equip slots and three free Use slots first.");
    } else if (!plr.gainItem(1472061, 1) || !plr.gainItem(1332063, 1) || !plr.gainItem(2070015, 3000)) {
        npc.sendOk("Make sure you have enough free Equip and Use inventory space first.");
    } else {
        plr.increaseSlotSize(1, 4);
        plr.increaseSlotSize(4, 4);
        plr.setJob(400);
        plr.giveMaxHP(rand(100, 150));
        plr.giveMaxMP(rand(30, 50));
        plr.giveSP(1);
        npc.sendOk("You're one of us now. Your Equip and Etc inventories have grown, and you're much stronger than before.");
    }
} else if (plr.job() === 400) {
    if (plr.getLevel() < 30) {
        npc.sendOk("Keep training. Come back when you're ready for your second job advancement.");
    } else if (plr.haveItem(LETTER, 1)) {
        npc.sendOk("Take my letter to #b#p1072003##k near #b#m102040000##k. He'll administer the test for me.");
    } else if (plr.haveItem(PROOF, 1)) {
        thiefSecondJobMenu();
    } else if (!npc.sendYesNo("You seem ready for the next step. Do you want to take the Thief 2nd job test?")) {
        npc.sendOk("Come back when you're ready.");
    } else if (plr.getEtcInventoryFreeSlot() < 1) {
        npc.sendOk("You need at least one free Etc slot for my letter.");
    } else {
        plr.gainItem(LETTER, 1);
        npc.sendOk("Take this letter to #b#p1072003##k near #b#m102040000##k. He'll explain the rest.");
    }
} else if (plr.job() === 410 || plr.job() === 420) {
    var state = plr.questData(7500);
    if (state === "s") {
        plr.setQuestData(7500, "p1");
        npc.sendOk("Arec told me about you. In Monkey Swamp II, there is a crack that leads to another dimension. Defeat my clone there and bring me #b#t4031059##k.");
    } else if (state === "p1") {
        if (!plr.haveItem(BLACK_CHARM, 1)) {
            npc.sendOk("Go to Monkey Swamp II, enter the crack, defeat my clone, and bring back #b#t4031059##k.");
        } else if (plr.getEtcInventoryFreeSlot() < 1) {
            npc.sendOk("Free up at least one Etc slot first.");
        } else {
            plr.gainItem(BLACK_CHARM, -1);
            plr.gainItem(STRENGTH_NECKLACE, 1);
            plr.setQuestData(7500, "p2");
            npc.sendOk("You've proven your strength. Take #b#t4031057##k to #b#p2020011##k in El Nath for the second test.");
        }
    } else if (state === "p2") {
        if (!plr.haveItem(STRENGTH_NECKLACE, 1)) {
            if (plr.getEtcInventoryFreeSlot() < 1) {
                npc.sendOk("Free up at least one Etc slot first.");
            } else {
                plr.gainItem(STRENGTH_NECKLACE, 1);
                npc.sendOk("You lost #b#t4031057##k, so I'm giving you another one. Don't lose it again.");
            }
        } else {
            npc.sendOk("Take #b#t4031057##k to #b#p2020011##k in El Nath. He'll handle the next test.");
        }
    } else {
        npc.sendOk("Stay sharp. A thief who gets careless doesn't last long.");
    }
} else {
    npc.sendOk("Stay sharp. A thief who gets careless doesn't last long.");
}
