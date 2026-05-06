var LETTER = 4031008;
var MARBLE = 4031013;
var rooms = [108000300, 108000301, 108000302];

function openRoom() {
    for (var i = 0; i < rooms.length; i++) {
        if (map.playerCount(rooms[i], 0) === 0) {
            return rooms[i];
        }
    }
    return 0;
}

if (plr.job() !== 100 || plr.getLevel() < 30) {
    npc.sendOk("If Dances with Balrog did not send you here, then you do not belong here.");
} else if (!plr.haveItem(LETTER, 1)) {
    npc.sendOk("Go back to Dances with Balrog in Perion first.");
} else if (plr.itemCount(MARBLE) > 0) {
    if (!npc.sendYesNo("So you gave up once already. That's fine. I can let you try again, but I'll have to take away the Dark Marbles you collected first. Do you want to re-enter the test?")) {
        npc.sendOk("Come back when you are fully prepared.");
    } else {
        var retryRoom = openRoom();
        if (retryRoom === 0) {
            npc.sendOk("Someone else is already taking the test. Please come back later.");
        } else {
            plr.removeAll(MARBLE);
            plr.warp(retryRoom);
        }
    }
} else if (!npc.sendYesNo("I'll send you to a hidden map filled with special monsters. They won't give EXP or items, but they will drop #b30 #t4031013##k. Collect 30 and speak to my colleague inside. Once you enter, you must finish or give up. Do you want to start now?")) {
    npc.sendOk("Come back when you're ready.");
} else {
    var room = openRoom();
    if (room === 0) {
        npc.sendOk("Someone else is already taking the test. Please come back later.");
    } else {
        npc.sendOk("Defeat the monsters inside, collect 30 Dark Marbles, and speak to my colleague there for your proof.");
        plr.warp(room);
    }
}
