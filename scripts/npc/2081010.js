var exitMap = 221000300;

if (plr.mapID() === exitMap) {
    npc.sendNext("See you next time.");
    plr.warpToPortalName(103000000, "mid00");
} else if (!npc.sendYesNo("Would you like to leave, " + plr.name() + "? Once you leave the map, you'll have to restart the whole quest if you want to try it again, and Juudai will be sad. Do you still want to leave this map?")) {
    npc.sendOk("Stay a little longer if you want to keep trying.");
} else if (plr.getEventProperty("leader") != null) {
    plr.leaveEvent();
} else {
    plr.warp(221000300);
}
