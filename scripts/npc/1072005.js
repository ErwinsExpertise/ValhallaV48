if (plr.itemCount(4031013) >= 30) {
    plr.removeAll(4031013);
    npc.sendOk("You gathered every last Dark Marble. Good. Take this proof back to Grendel and show him your control was not a matter of luck.");
    plr.warp(101020000);
} else {
    npc.sendOk("Come back when you've collected #b30 #t4031013##k. A magician who cannot finish a trial should not be asking about the next rank.");
}
