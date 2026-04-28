if (plr.job() === 0) {
    if (plr.getLevel() < 10 || plr.getDex() < 25) {
        npc.sendOk("Train a bit more and I can show you the way of the #rThief#k.");
    } else if (!npc.sendYesNo("You have the quick hands and agility a thief needs. If you choose this path, you'll fight with speed, timing, and precision from the shadows. This is a final decision. Do you want to become a #rThief#k?")) {
        npc.sendOk("Make up your mind and visit me again.");
    } else {
        plr.setJob(400);
        plr.gainItem(1332063, 1);
        npc.sendOk("So be it! Now go, and go with pride.");
    }
} else if (plr.job() === 400) {
    if (plr.getLevel() < 30) {
        npc.sendOk("You have chosen wisely.");
    } else {
        var branch = npc.sendMenu("You've gotten fast enough for the next step. Choose the path that suits you best.#b", "Assassin", "Bandit");
        var jobName = branch === 0 ? "Assassin" : "Bandit";
        var jobId = branch === 0 ? 410 : 420;
        if (npc.sendYesNo("Do you want to become a #r" + jobName + "#k?")) {
            plr.setJob(jobId);
            npc.sendOk("Good. From this point on, you walk the road of the #b" + jobName + "#k.");
        } else {
            npc.sendOk("Think it over and speak to me again when you're sure.");
        }
    }
} else {
    npc.sendOk("You have chosen wisely.");
}
