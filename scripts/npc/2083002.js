var pqItems = [4001087, 4001088, 4001089, 4001090, 4001091, 4001092, 4001093];

if (plr.mapID() !== 240050500) {
    if (!npc.sendYesNo("Do you wish to go out from here? You will have to start from scratch again next time...")) {
        npc.sendOk("Ok, keep persevering!");
    } else {
        plr.leaveEvent();
        plr.warp(240050500);
    }
} else {
    npc.sendNext("Tough luck there, eh? You can always come back if you're prepared... but anyway, I will take all the items you obtained from the PQ :)");
    for (var i = 0; i < pqItems.length; i++) {
        plr.removeAll(pqItems[i]);
    }
    plr.warp(240040700);
}
