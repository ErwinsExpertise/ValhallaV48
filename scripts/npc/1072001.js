var LETTER = 4031009;
var MARBLE = 4031013;
var rooms = [108000200, 108000201, 108000202];

function openRoom() {
    for (var i = 0; i < rooms.length; i++) {
        if (map.playerCount(rooms[i], 0) === 0) {
            return rooms[i];
        }
    }
    return 0;
}

if (plr.job() !== 200 || plr.getLevel() < 30) {
    npc.sendOk("If Grendel did not send you here, then this test is not for you.");
} else if (!plr.haveItem(LETTER, 1)) {
    npc.sendOk("Go back to Grendel the Really Old first.");
} else if (plr.itemCount(MARBLE) > 0) {
    if (!npc.sendYesNo("So you gave up once already. That's fine. I can let you try again, but I have to remove the Dark Marbles you collected first. Do you want to re-enter the test?")) {
        npc.sendOk("Come back when you're ready.");
    } else {
        var retryRoom = openRoom();
        if (retryRoom === 0) {
            npc.sendOk("Someone else is already taking the test. Please come back later.");
        } else {
            plr.removeAll(MARBLE);
            plr.warp(retryRoom);
        }
    }
} else if (!npc.sendYesNo("I'll send you to a hidden map filled with unusual monsters. They don't give EXP or items, but they do drop #b#t4031013##k. Collect 30 and talk to my colleague inside. Do you want to begin now?")) {
    npc.sendOk("Come back when you're fully prepared.");
} else {
    var room = openRoom();
    if (room === 0) {
        npc.sendOk("Someone else is already taking the test. Please come back later.");
    } else {
        npc.sendOk("Defeat the monsters, collect 30 Dark Marbles, and speak to my colleague inside for your proof.");
        plr.warp(room);
    }
}
