if (plr.itemCount(4031013) >= 30) {
    plr.removeAll(4031013);
    npc.sendOk("So you collected all 30 Dark Marbles. Good. Take this proof back to the Dark Lord and show him your hands stayed sharp all the way through.");
    plr.warp(102040000);
} else {
    npc.sendOk("Bring me #b30 #t4031013##k. If you can't finish this cleanly, you're not ready for the next step.");
}
