if (!(plr.job() === 110 || plr.job() === 120 || plr.job() === 130)) {
    npc.sendOk("Only those who have already walked the warrior's path belong here.");
} else if (plr.getLevel() < 70) {
    npc.sendOk("You are not ready yet. Train until both your body and your technique can support the next step.");
} else {
    var branch = npc.sendMenu("You've trained long enough to take on a greater title. Choose the road ahead.#b", "Crusader", "White Knight", "Dragon Knight");
    var jobName = branch === 0 ? "Crusader" : branch === 1 ? "White Knight" : "Dragon Knight";
    var jobId = branch === 0 ? 111 : branch === 1 ? 121 : 131;
    if (plr.job() === jobId) {
        npc.sendOk("A warrior's training never ends. Keep sharpening your body and your will.");
    } else 
    if (npc.sendYesNo("Do you wish to become a #r" + jobName + "#k?")) {
        plr.setJob(jobId);
        plr.giveAP(5);
        npc.sendOk("Good. Then carry that title with pride. You are now a #b" + jobName + "#k.");
    }
}
