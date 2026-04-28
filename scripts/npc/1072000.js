var rooms = [108000300, 108000301, 108000302];

if (!(plr.job() === 100 && plr.getLevel() >= 30)) {
    npc.sendOk("If you're here for the warrior trial, come back when Dances with Balrog sends you.");
} else if (!npc.sendYesNo("You've come for the warrior trial. Inside, you'll face monsters that won't reward you with experience or items. Your only goal is to collect #b30 #t4031013##k. Are you ready to enter?")) {
    npc.sendOk("Come back when you're ready to prove your strength.");
} else {
    var room = -1;
    for (var i = 0; i < rooms.length; i++) {
        if (map.playerCount(rooms[i], 0) === 0) {
            room = rooms[i];
            break;
        }
    }
    if (room < 0) {
        npc.sendOk("All test chambers are currently occupied. Please wait until one becomes available.");
    } else {
        plr.removeAll(4031013);
        plr.removeAll(4031012);
        plr.warp(room);
    }
}
