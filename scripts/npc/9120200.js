if (!npc.sendYesNo("Here you are, right in front of the hideout! What? You want to return to #m801000000#?")) {
    npc.sendOk("If you want to return to #m801000000#, then talk to me.");
} else {
    plr.warp(801000000);
}
