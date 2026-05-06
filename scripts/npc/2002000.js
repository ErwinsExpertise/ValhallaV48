if (npc.sendYesNo("You look like you are done with Happyville for now. Would you like me to send you back where you came from?")) {
    // Clean up temporary Happyville items before returning the player.
    var removeIds = [1472063, 2060005, 2060006];
    for (var i = 0; i < removeIds.length; i++) {
        var count = plr.itemCount(removeIds[i]);
        if (count > 0) {
            plr.gainItem(removeIds[i], -count);
        }
    }

    var location = plr.getSavedLocation("HV_MAP");
    if (location < 0) {
        location = 100000000;
    }
    plr.clearSavedLocation("HV_MAP");
    plr.warp(location);
} else {
    npc.sendOk("Take your time. Happyville will still be here when you are ready to leave.");
}
