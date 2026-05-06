var LETTER = 4031011;
var MARBLE = 4031013;
var rooms = [108000400, 108000401, 108000402];

function openRoom() {
    for (var i = 0; i < rooms.length; i++) {
        if (map.playerCount(rooms[i], 0) === 0) {
            return rooms[i];
        }
    }
    return 0;
}

if (plr.job() !== 400 || plr.getLevel() < 30) {
    npc.sendOk("If the Dark Lord did not send you here, then this test is not for you.");
} else if (!plr.haveItem(LETTER, 1)) {
    npc.sendOk("Go back to the Dark Lord first.");
} else if (plr.itemCount(MARBLE) > 0) {
    if (!npc.sendYesNo("So you gave up once already. That's fine. I can let you try again, but I have to remove the Dark Marbles you collected first. Do you want to go back in?")) {
        npc.sendOk("Come back when you're fully prepared.");
    } else {
        var retryRoom = openRoom();
        if (retryRoom === 0) {
            npc.sendOk("Someone else is already taking the test. Please come back later.");
        } else {
            plr.removeAll(MARBLE);
            plr.warp(retryRoom);
        }
    }
} else if (!npc.sendYesNo("I'll send you to a hidden map where special monsters drop #b30 #t4031013##k. Collect them all and speak to my colleague inside to receive your proof. Do you want to begin now?")) {
    npc.sendOk("Come back when you're ready.");
} else {
    var room = openRoom();
    if (room === 0) {
        npc.sendOk("Someone else is already taking the test. Please come back later.");
    } else {
        npc.sendOk("Defeat the monsters, collect 30 Dark Marbles, and speak to my colleague inside for your proof.");
        plr.warp(room);
    }
}
