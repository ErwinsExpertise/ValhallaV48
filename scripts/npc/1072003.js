var rooms = [108000400, 108000401, 108000402];

if (!(plr.job() === 400 && plr.getLevel() >= 30)) {
    npc.sendOk("If you're here for the thief trial, come back when the Dark Lord sends you.");
} else if (!npc.sendYesNo("You've come for the thief trial. Inside, you'll face monsters that won't reward you with experience or items. Bring back #b30 #t4031013##k to prove your speed and precision. Are you ready to enter?")) {
    npc.sendOk("Come back when you're ready to be tested.");
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
