if (!npc.sendYesNo("Do you want to go to Ellinia?")) {
    npc.sendOk("Until next time...");
} else {
    plr.warp(101000000);
}
