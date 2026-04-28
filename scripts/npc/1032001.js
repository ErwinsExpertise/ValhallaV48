if (plr.job() === 0) {
    if (plr.getLevel() < 8 || plr.getInt() < 20) {
        npc.sendOk("Train a bit more and I can show you the way of the #rMagician#k.");
    } else if (!npc.sendYesNo("You have the intellect to pursue magic. If you choose this path, you will study the arcane arts and leave the beginner's road behind. This is a final decision. Do you want to become a #rMagician#k?")) {
        npc.sendOk("Make up your mind and visit me again.");
    } else {
        plr.setJob(200);
        plr.gainItem(1372043, 1);
        npc.sendOk("So be it! Now go, and go with pride.");
    }
} else if (plr.job() === 200) {
    if (plr.getLevel() < 30) {
        npc.sendOk("You have chosen wisely.");
    } else {
        var branch = npc.sendMenu("You have reached the point where your studies may branch into a true specialty.#b", "Wizard (Fire, Poison)", "Wizard (Ice, Lightning)", "Cleric");
        var jobName = branch === 0 ? "Wizard of Fire and Poison" : branch === 1 ? "Wizard of Ice and Lightning" : "Cleric";
        var jobId = branch === 0 ? 210 : branch === 1 ? 220 : 230;
        if (npc.sendYesNo("Do you want to become a #r" + jobName + "#k?")) {
            plr.setJob(jobId);
            npc.sendOk("Then continue your studies as a #b" + jobName + "#k.");
        } else {
            npc.sendOk("Think it over and return when your mind is made up.");
        }
    }
} else {
    npc.sendOk("You have chosen wisely.");
}
