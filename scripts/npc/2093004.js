var cost = 10000;

if (!npc.sendYesNo("Will you move to #b#m230000000##k now? The price is #b" + cost + " mesos#k.")) {
    npc.sendOk("Hmmm ... too busy to do it right now? If you feel like doing it, though, come back and find me.");
} else if (plr.getMesos() < cost) {
    npc.sendOk("I don't think you have enough money...");
} else {
    plr.gainMesos(-cost);
    plr.warp(230000000);
}
