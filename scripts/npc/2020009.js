if (!(plr.job() === 210 || plr.job() === 220 || plr.job() === 230)) {
    npc.sendOk("Arcane study is not for the impatient. Return when you truly walk the path of magic.");
} else if (plr.getLevel() < 70) {
    npc.sendOk("You are not ready yet. Continue your studies and return when your foundation is complete.");
} else {
    var branch = npc.sendMenu("You are ready to devote yourself to a higher branch of magic. Choose your path.#b", "Fire/Poison Mage", "Ice/Lightning Mage", "Priest");
    var jobName = branch === 0 ? "Fire/Poison Mage" : branch === 1 ? "Ice/Lightning Mage" : "Priest";
    var jobId = branch === 0 ? 211 : branch === 1 ? 221 : 231;
    if (plr.job() === jobId) {
        npc.sendOk("A mage's journey never truly ends. Keep refining your craft.");
    } else 
    if (npc.sendYesNo("Do you want to become a #r" + jobName + "#k?")) {
        plr.setJob(jobId);
        plr.giveAP(5);
        npc.sendOk("Then take the next step. From this point on, you are a #b" + jobName + "#k.");
    }
}
