var aquaRoadMap = 230000000;
var fee = plr.job() === 0 ? 1000 : 10000;

if (npc.sendYesNo("Do you want to head to #b#m" + aquaRoadMap + "##k now? The fare is #b" + fee + " mesos#k.")) {
    if (plr.mesos() < fee) {
        npc.sendOk("I don't think you have enough mesos...");
    } else {
        plr.takeMesos(fee);
        plr.warp(aquaRoadMap);
    }
} else {
    npc.sendOk("Hmm... too busy right now? Come back when you need a ride.");
}
