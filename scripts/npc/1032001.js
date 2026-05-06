var LETTER = 4031009;
var PROOF = 4031012;
var BLACK_CHARM = 4031059;
var STRENGTH_NECKLACE = 4031057;

function rand(min, max) {
    return Math.floor(Math.random() * (max - min + 1)) + min;
}

function spentEnoughSp() {
    return plr.getRemainingSP() <= (plr.getLevel() - 30) * 3;
}

function magicianSecondJobMenu() {
    var explain = npc.sendMenu(
        "Before I advance you, choose the direction of your studies.#b",
        "Explain Fire/Poison Wizard",
        "Explain Ice/Lightning Wizard",
        "Explain Cleric",
        "I want to choose now"
    );

    if (explain === 0) {
        npc.sendOk("Fire/Poison Wizards specialize in aggressive elemental magic and poison-based control.");
        return;
    }
    if (explain === 1) {
        npc.sendOk("Ice/Lightning Wizards specialize in freezing and stunning enemies with elemental magic.");
        return;
    }
    if (explain === 2) {
        npc.sendOk("Clerics use holy magic and healing. They support allies and excel against undead monsters.");
        return;
    }

    var pick = npc.sendMenu("Choose your 2nd job.#b", "Wizard (Fire, Poison)", "Wizard (Ice, Lightning)", "Cleric");
    var jobId = pick === 0 ? 210 : pick === 1 ? 220 : 230;
    var name = pick === 0 ? "Wizard of Fire and Poison" : pick === 1 ? "Wizard of Ice and Lightning" : "Cleric";

    if (!npc.sendYesNo("Do you want to become a #b" + name + "#k? Once you decide, you cannot go back.")) {
        npc.sendOk("Come back after you've thought it through.");
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
    plr.giveMaxMP(rand(450, 500));
    plr.increaseSlotSize(4, 4);
    npc.sendOk("From this point on, continue your studies as a #b" + name + "#k.");
}

if (plr.job() === 0) {
    npc.sendOk("If you want to become a Magician, you need to be at least #bLevel 8#k with at least #b20 INT#k.");
    if (plr.getLevel() < 8 || plr.getInt() < 20) {
        npc.sendOk("Train a bit more before asking to become a Magician.");
    } else if (!npc.sendYesNo("You definitely have the look of a Magician. Do you want to become a #bMagician#k?")) {
        npc.sendOk("Come back when you've made up your mind.");
    } else if (!plr.canHold(1372005, 1)) {
        npc.sendOk("Make sure you have at least one free Equip slot. I want to give you a weapon for your first advancement.");
    } else if (!plr.gainItem(1372005, 1)) {
        npc.sendOk("Make sure you have at least one free Equip slot. I want to give you a weapon for your first advancement.");
    } else {
        plr.setJob(200);
        plr.giveMaxMP(rand(100, 150));
        plr.giveSP(1);
        npc.sendOk("You are a Magician now. Keep studying and grow stronger every day.");
    }
} else if (plr.job() === 200) {
    if (plr.getLevel() < 30) {
        npc.sendOk("Keep training. Come back when you're ready for your second job advancement.");
    } else if (plr.haveItem(LETTER, 1)) {
        npc.sendOk("Take my letter to #b#p1072001##k near #b#m101020000##k. He'll administer the test for me.");
    } else if (plr.haveItem(PROOF, 1)) {
        magicianSecondJobMenu();
    } else if (!npc.sendYesNo("You seem qualified for the next step. Do you want to take the Magician 2nd job test?")) {
        npc.sendOk("Come back when you're ready.");
    } else if (plr.getEtcInventoryFreeSlot() < 1) {
        npc.sendOk("You need at least one free Etc slot for my letter.");
    } else {
        plr.gainItem(LETTER, 1);
        npc.sendOk("Take this letter to #b#p1072001##k near #b#m101020000##k. He'll explain the rest.");
    }
} else if (plr.job() === 210 || plr.job() === 220 || plr.job() === 230) {
    var state = plr.questData(7500);
    if (state === "s") {
        plr.setQuestData(7500, "p1");
        npc.sendOk("Robeira told me about you. Near Ellinia's forest, there's a secret passage only you can enter. Defeat my clone there and bring me #b#t4031059##k.");
    } else if (state === "p1") {
        if (!plr.haveItem(BLACK_CHARM, 1)) {
            npc.sendOk("The secret passage is near Ellinia's forest. Defeat my clone there and bring back #b#t4031059##k.");
        } else if (plr.getEtcInventoryFreeSlot() < 1) {
            npc.sendOk("Free up at least one Etc slot first.");
        } else {
            plr.gainItem(BLACK_CHARM, -1);
            plr.gainItem(STRENGTH_NECKLACE, 1);
            plr.setQuestData(7500, "p2");
            npc.sendOk("You've proven your strength. Take #b#t4031057##k to #b#p2020009##k in El Nath for the second test.");
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
            npc.sendOk("Take #b#t4031057##k to #b#p2020009##k in El Nath. She'll handle the next test.");
        }
    } else {
        npc.sendOk("Keep refining your magic.");
    }
} else {
    npc.sendOk("Keep refining your magic.");
}
