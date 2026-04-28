if (plr.job() === 0) {
    if (plr.getLevel() < 10 || plr.getDex() < 25) {
        npc.sendOk("Train a bit more and I can show you the way of the #rBowman#k.");
    } else if (!npc.sendYesNo("Your eyes and your footing are ready for the bow. If you choose this path, your future will depend on precision and control. This is a final decision. Do you want to become a #rBowman#k?")) {
        npc.sendOk("Make up your mind and visit me again.");
    } else {
        plr.setJob(300);
        plr.gainItem(1452002, 1);
        plr.gainItem(2060000, 1000);
        npc.sendOk("So be it! Now go, and go with pride.");
    }
} else if (plr.job() === 300) {
    if (plr.getLevel() < 30) {
        npc.sendOk("You have chosen wisely.");
    } else {
        var branch = npc.sendMenu("You've trained enough to choose your next path.#b", "Hunter", "Crossbowman");
        var jobName = branch === 0 ? "Hunter" : "Crossbowman";
        var jobId = branch === 0 ? 310 : 320;
        if (npc.sendYesNo("Do you want to become a #r" + jobName + "#k?")) {
            plr.setJob(jobId);
            npc.sendOk("Then let your aim guide you. You are now a #b" + jobName + "#k.");
        } else {
            npc.sendOk("Think it over and speak to me again when you're sure.");
        }
    }
} else {
    npc.sendOk("You have chosen wisely.");
}
