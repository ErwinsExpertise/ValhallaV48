function act() {
    var blowerIds = [9400714, 9400715, 9400716, 9400717, 9400718, 9400719, 9400720, 9400721, 9400722, 9400723, 9400724];
    var bossIds = [9400707, 9400708, 9400709, 9400710];

    for (var i = 0; i < bossIds.length; i++) {
        rm.removeMobsByID(bossIds[i]);
    }
    for (var j = 0; j < blowerIds.length; j++) {
        rm.removeMobsByID(blowerIds[j]);
    }

    rm.setMapProperty("wxmasCount", "0");
    rm.setMapProperty("wxmasBoss", "0");
    rm.spawnMonster(9400714, 1450, 140);
    rm.mapMessage(5, "The snow machine settles down after the commotion.");
}
