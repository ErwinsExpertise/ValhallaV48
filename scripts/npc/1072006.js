if (plr.itemCount(4031013) >= 30) {
    plr.removeAll(4031013);
    npc.sendOk("You brought back all 30 Dark Marbles. Not bad. Take this proof to Athena Pierce and show her your aim held up under pressure.");
    plr.warp(106010000);
} else {
    npc.sendOk("Return with #b30 #t4031013##k. A bowman who loses focus halfway through a test is not ready for the next rank.");
}
