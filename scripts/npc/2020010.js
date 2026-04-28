if (!(plr.job() === 310 || plr.job() === 320)) {
    npc.sendOk("A bowman must learn patience, precision, and restraint. Return when you truly walk that path.");
} else if (plr.getLevel() < 70) {
    npc.sendOk("You are not ready yet. Keep training until your aim never wavers.");
} else {
    var branch = npc.sendMenu("You have reached the point where your skill can split into a sharper specialty. Choose your path.#b", "Ranger", "Sniper");
    var jobName = branch === 0 ? "Ranger" : "Sniper";
    var jobId = branch === 0 ? 311 : 321;
    if (plr.job() === jobId) {
        npc.sendOk("Your aim has improved, but there is always a farther mark to hit.");
    } else 
    if (npc.sendYesNo("Do you want to become a #r" + jobName + "#k?")) {
        plr.setJob(jobId);
        plr.giveAP(5);
        npc.sendOk("Then let your aim guide you. You are now a #b" + jobName + "#k.");
    }
}
