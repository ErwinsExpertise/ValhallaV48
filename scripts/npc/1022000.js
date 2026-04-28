if (plr.job() === 0) {
    if (plr.getLevel() < 10 || plr.getStr() < 35) {
        npc.sendOk("Train a bit more and I can show you the way of the #rWarrior#k.");
    } else if (!npc.sendYesNo("You have the build of a warrior. If you choose this path, you will devote yourself to close combat and raw strength. This is a final decision. Do you want to become a #rWarrior#k?")) {
        npc.sendOk("Make up your mind and visit me again.");
    } else {
        plr.setJob(100);
        plr.gainItem(1402001, 1);
        npc.sendOk("So be it! Now go, and go with pride.");
    }
} else if (plr.job() === 100) {
    if (plr.getLevel() < 30) {
        npc.sendOk("You have chosen wisely.");
    } else {
        var branch = npc.sendMenu("You've trained enough to choose your next path.#b", "Fighter", "Page", "Spearman");
        var jobName = branch === 0 ? "Fighter" : branch === 1 ? "Page" : "Spearman";
        var jobId = branch === 0 ? 110 : branch === 1 ? 120 : 130;
        if (npc.sendYesNo("Do you want to become a #r" + jobName + "#k?")) {
            plr.setJob(jobId);
            npc.sendOk("Good. From here on, walk the path of the #b" + jobName + "#k with pride.");
        } else {
            npc.sendOk("Think it over and speak to me again when you're sure.");
        }
    }
} else {
    npc.sendOk("You have chosen wisely.");
}
