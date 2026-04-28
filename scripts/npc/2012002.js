if (!npc.sendYesNo("Do you wish to leave the boat?")) {
    npc.sendOk("Please hold on tight. We will arrive soon.");
} else {
    npc.sendOk("All right, see you next time. Take care.");
    if (plr.mapID() === 101000301) {
        plr.warp(101000300);
    } else {
        plr.warp(200000111);
    }
}
