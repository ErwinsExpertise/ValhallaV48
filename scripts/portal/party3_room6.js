if (plr.getEventProperty("stage6_clear") === "1") {
    portal.block("There is nothing left to do in this room.");
} else {
    if (plr.getEventProperty("room6_ans") !== "1") {
        var slots = ["1", "2", "3", "4", "5"];
        for (var i = slots.length - 1; i > 0; i--) {
            var j = Math.floor(Math.random() * (i + 1));
            var tmp = slots[i];
            slots[i] = slots[j];
            slots[j] = tmp;
        }
        plr.setEventProperty("stage6_ans1", slots[0]);
        plr.setEventProperty("stage6_ans2", slots[1]);
        plr.setEventProperty("stage6_wans1", slots[2]);
        plr.setEventProperty("stage6_wans2", slots[3]);
        plr.setEventProperty("stage6_wans3", slots[4]);
        plr.setEventProperty("room6_ans", "1");
    }

    if (!plr.isLeader()) {
        if (map.playerCountInMap(920010700) > 0) {
            portal.warp(920010700, "st00");
        } else {
            portal.block("You may only enter the room your party leader is already in.");
        }
    } else {
        plr.sendMessage("Your party leader entered On the Way Up.");
        portal.warp(920010700, "st00");
    }
}
