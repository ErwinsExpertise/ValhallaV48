var QUEST = 7500;
var STRENGTH_NECKLACE = 4031057;
var WISDOM_NECKLACE = 4031058;

function confirmAdvance() {
    if (!npc.sendYesNo("You've passed the tests. Make sure you've used all of the SP you've earned up to level 70 first. Do you want to become a 3rd job Warrior now?")) {
        npc.sendOk("Come back when you're ready.");
        return;
    }
    if (plr.getRemainingSP() > (plr.getLevel() - 70) * 3) {
        npc.sendOk("You still have too much unused SP. Use the SP from your 1st and 2nd job skills first, then come back.");
        return;
    }

    var jobId = plr.job() === 110 ? 111 : plr.job() === 120 ? 121 : 131;
    plr.setJob(jobId);
    plr.giveSP(1);
    plr.giveAP(5);
    npc.sendOk("Excellent. You have taken the next step as a Warrior.");
}

if (plr.job() !== 110 && plr.job() !== 120 && plr.job() !== 130 && plr.job() !== 111 && plr.job() !== 121 && plr.job() !== 131) {
    npc.sendOk("I'm Tylus, chief of all Warriors. Speak to your own class chief if you need guidance.");
} else if (plr.job() === 111 || plr.job() === 121 || plr.job() === 131) {
    npc.sendOk("You already passed my tests. Keep training and prove yourself worthy of the title you now bear.");
} else if (plr.getLevel() < 70) {
    npc.sendOk("You need to be at least level 70 before you can attempt your 3rd job advancement.");
} else if (plr.getQuestStatus(QUEST) === 0) {
    if (!npc.sendYesNo("I am Tylus, chief of all Warriors. Do you want to take the test for the 3rd job advancement?")) {
        npc.sendOk("Come back when you're truly ready.");
    } else {
        plr.startQuest(QUEST);
        plr.setQuestData(QUEST, "s");
        npc.sendOk("Go to #b#p1022000##k in Perion and complete the physical test first. Bring back #b#t4031057##k when you have passed it.");
    }
} else if (plr.questData(QUEST) === "s" || plr.questData(QUEST) === "p1") {
    npc.sendOk("You need to pass Dances with Balrog's physical test first. Bring me #b#t4031057##k.");
} else if (plr.questData(QUEST) === "p2") {
    if (!plr.haveItem(STRENGTH_NECKLACE, 1)) {
        npc.sendOk("Go back to Dances with Balrog and get #b#t4031057##k first.");
    } else {
        plr.gainItem(STRENGTH_NECKLACE, -1);
        plr.setQuestData(QUEST, "end1");
        npc.sendOk("Good. Now go to the #bHoly Stone#k in the snowfield and offer a #b#t4005004##k. Pass its test and bring me #b#t4031058##k.");
    }
} else if (plr.questData(QUEST) === "end1") {
    if (!plr.haveItem(WISDOM_NECKLACE, 1)) {
        npc.sendOk("Pass the Holy Stone's test and bring me #b#t4031058##k.");
    } else {
        plr.gainItem(WISDOM_NECKLACE, -1);
        plr.forceCompleteQuest(QUEST);
        confirmAdvance();
    }
} else if (plr.getQuestStatus(QUEST) === 2) {
    confirmAdvance();
} else {
    npc.sendOk("You need to finish the tests before I can advance you.");
}
