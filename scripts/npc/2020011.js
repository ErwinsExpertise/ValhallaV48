if (!(plr.job() === 410 || plr.job() === 420)) {
    npc.sendOk("The shadows do not open for just anyone. Return when you truly understand the thief's road.");
} else if (plr.getLevel() < 70) {
    npc.sendOk("Not yet. Keep training until your hands and your head are both ready.");
} else {
    var branch = npc.sendMenu("You've made it this far because your instincts are good. Choose the road that suits your talents best.#b", "Hermit", "Chief Bandit");
    var jobName = branch === 0 ? "Hermit" : "Chief Bandit";
    var jobId = branch === 0 ? 411 : 421;
    if (plr.job() === jobId) {
        npc.sendOk("A thief survives by staying sharp. Keep honing your instincts.");
    } else 
    if (npc.sendYesNo("Do you want to become a #r" + jobName + "#k?")) {
        plr.setJob(jobId);
        plr.giveAP(5);
        npc.sendOk("Good. From here on, you are a #b" + jobName + "#k.");
    }
}
