var LETTER = 4031008;
var PROOF = 4031012;
var BLACK_CHARM = 4031059;
var STRENGTH_NECKLACE = 4031057;

function rand(min, max) {
    return Math.floor(Math.random() * (max - min + 1)) + min;
}

function useAllRequiredSp(level) {
    return plr.getRemainingSP() <= (level - 30) * 3;
}

function warriorSecondJobMenu() {
    var explain = npc.sendMenu(
        "You made it back safely. Before I advance you, decide which path suits you best.#b",
        "Explain Fighter",
        "Explain Page",
        "Explain Spearman",
        "I want to choose now"
    );

    if (explain === 0) {
        npc.sendOk("Fighters focus on raw offense with swords and axes. They keep pressing forward and overwhelm enemies with power.");
        return;
    }
    if (explain === 1) {
        npc.sendOk("Pages are balanced warriors who favor swords or blunt weapons and develop into knights with elemental charges later on.");
        return;
    }
    if (explain === 2) {
        npc.sendOk("Spearmen specialize in polearms and spears. They support themselves and their parties with long reach and durability.");
        return;
    }

    var pick = npc.sendMenu("Choose your 2nd job.#b", "Fighter", "Page", "Spearman");
    var jobId = pick === 0 ? 110 : pick === 1 ? 120 : 130;
    var name = pick === 0 ? "Fighter" : pick === 1 ? "Page" : "Spearman";

    if (!npc.sendYesNo("Do you want to become a #b" + name + "#k? Once you decide, you cannot go back.")) {
        npc.sendOk("Take your time. Come back when you've made up your mind.");
        return;
    }
    if (!useAllRequiredSp(plr.getLevel())) {
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
    if (jobId === 110) {
        plr.giveMaxHP(rand(300, 350));
    } else {
        plr.giveMaxMP(rand(100, 150));
    }
    plr.increaseSlotSize(2, 4);
    plr.increaseSlotSize(4, 4);
    npc.sendOk("Good. From now on, walk the path of the #b" + name + "#k with pride.");
}

if (plr.job() === 0) {
    npc.sendOk("If you want to become a Warrior, you'll need to be at least #bLevel 10#k with at least #b35 STR#k.");
    if (plr.getLevel() < 10 || plr.getStr() < 35) {
        npc.sendOk("You still need more training before you can become a Warrior.");
    } else if (!npc.sendYesNo("You definitely have the look of a Warrior. Do you want to become a #bWarrior#k?")) {
        npc.sendOk("Come back when you're ready to decide.");
    } else if (plr.getEquipInventoryFreeSlot() < 1) {
        npc.sendOk("Make sure you have at least one free Equip slot. I want to give you a weapon for your first advancement.");
    } else if (!plr.gainItem(1302077, 1)) {
        npc.sendOk("Make sure you have at least one free Equip slot. I want to give you a weapon for your first advancement.");
    } else {
        plr.setJob(100);
        plr.giveMaxHP(rand(200, 250));
        plr.giveSP(1);
        plr.increaseSlotSize(1, 4);
        plr.increaseSlotSize(2, 4);
        plr.increaseSlotSize(3, 4);
        plr.increaseSlotSize(4, 4);
        npc.sendOk("You are a Warrior now. You're much stronger than before, and all of your inventories have grown by a full row.");
    }
} else if (plr.job() === 100) {
    if (plr.getLevel() < 30) {
        npc.sendOk("Keep training. Come back once you've reached the point where you can attempt your second job advancement.");
    } else if (plr.haveItem(LETTER, 1)) {
        npc.sendOk("Take my letter to #b#p1072000##k near #b#m102020300##k. He'll administer the test for me.");
    } else if (plr.haveItem(PROOF, 1)) {
        warriorSecondJobMenu();
    } else if (!npc.sendYesNo("You've grown tremendously. If you want to become stronger, I can test you. Do you want to take the Warrior 2nd job test?")) {
        npc.sendOk("Come back when you're ready.");
    } else if (plr.getEtcInventoryFreeSlot() < 1) {
        npc.sendOk("You need at least one free Etc slot for my letter.");
    } else {
        plr.gainItem(LETTER, 1);
        npc.sendOk("Take this letter to #b#p1072000##k near #b#m102020300##k. He'll explain the rest.");
    }
} else if (plr.job() === 110 || plr.job() === 120 || plr.job() === 130) {
    var state = plr.questData(7500);
    if (state === "s") {
        plr.setQuestData(7500, "p1");
        npc.sendOk("Tylus told me about you. Near the Ant Tunnel, there's a secret passage that only you can enter. Defeat my clone there and bring me #b#t4031059##k.");
    } else if (state === "p1") {
        if (!plr.haveItem(BLACK_CHARM, 1)) {
            npc.sendOk("The secret passage is near the Ant Tunnel. Defeat my clone there and bring back #b#t4031059##k.");
        } else if (plr.getEtcInventoryFreeSlot() < 1) {
            npc.sendOk("Free up at least one Etc slot first.");
        } else {
            plr.gainItem(BLACK_CHARM, -1);
            plr.gainItem(STRENGTH_NECKLACE, 1);
            plr.setQuestData(7500, "p2");
            npc.sendOk("You've proven your strength. Take #b#t4031057##k to #b#p2020008##k in El Nath for the second test.");
        }
    } else if (state === "p2") {
        if (!plr.haveItem(STRENGTH_NECKLACE, 1)) {
            if (plr.getEtcInventoryFreeSlot() < 1) {
                npc.sendOk("Free up at least one Etc slot first.");
            } else {
                plr.gainItem(STRENGTH_NECKLACE, 1);
                npc.sendOk("You lost #b#t4031057##k, so I'm giving you another one. Be careful this time.");
            }
        } else {
            npc.sendOk("Take #b#t4031057##k to #b#p2020008##k in El Nath. He'll handle the next test.");
        }
    } else {
        npc.sendOk("You have chosen wisely.");
    }
} else {
    npc.sendOk("You have chosen wisely.");
}
