var pqItems = [4001015, 4001016, 4001018];

npc.sendNext("Tough luck there eh? You can always come back if your'e prepared... but anyway, I will take all the items you obtained from the PQ :)");

if (!plr.isGM()) {
    for (var i = 0; i < pqItems.length; i++) {
        plr.removeAll(pqItems[i]);
    }
}

plr.warp(211042300);
