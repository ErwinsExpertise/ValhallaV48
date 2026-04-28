var mapId = plr.mapID();

if (mapId === 670010200) {
    if (!plr.isLeader()) {
        npc.sendOk("Ask your party leader to talk to me.");
    } else if (plr.countMonster() === 0) {
        npc.sendOk("Kill this fairy I'm about to spawn, and drop the Hammer onto the mirror to break it.");
        plr.spawnMonster(9400518, -5, 150);
    } else {
        npc.sendOk("Please kill all the monsters in the map before talking to me.");
    }
} else if (mapId === 670011000) {
    plr.warp(670010000);
}
