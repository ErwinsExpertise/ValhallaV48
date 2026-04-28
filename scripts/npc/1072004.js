if (plr.itemCount(4031013) >= 30) {
    plr.removeAll(4031013);
    npc.sendOk("You've done well. Take this proof back to Dances with Balrog. He will know you earned it.");
    plr.warp(102020300);
} else {
    npc.sendOk("Come back when you've collected #b30 #t4031013##k. That's the only proof I'll accept.");
}
