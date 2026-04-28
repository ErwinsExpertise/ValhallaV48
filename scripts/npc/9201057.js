var TICKET_COST = 5000;
var LOCATION_SLOT = "SUBWAY";
var mapId = plr.mapID();

if (mapId === 103000100 || mapId === 600010001) {
    var destination = mapId === 103000100 ? "New Leaf City" : "Kerning City";
    var ticketId = mapId === 103000100 ? 4031711 : 4031713;

    if (!npc.sendYesNo("Travel to " + destination + " will cost you #b" + TICKET_COST + " mesos#k. Would you like to buy a #b#t" + ticketId + "##k?")) {
        npc.sendOk("Come back when you are ready to travel.");
    } else if (plr.getMesos() < TICKET_COST) {
        npc.sendOk("You do not have enough mesos.");
    } else if (!plr.canHold(ticketId, 1)) {
        npc.sendOk("Please make room in your Etc inventory first.");
    } else {
        plr.gainMesos(-TICKET_COST);
        plr.gainItem(ticketId, 1);
        npc.sendOk("Here is your ticket. Have a safe trip.");
    }
} else if (mapId === 600010002 || mapId === 600010004) {
    if (!npc.sendYesNo("Do you want to leave the waiting room? You can, but the ticket is NOT refundable. Are you sure you still want to leave this room?")) {
        npc.sendOk("You will depart soon. Please wait a little longer.");
    } else {
        var fallbackMap = mapId === 600010002 ? 600010001 : 103000100;
        var returnMap = plr.getSavedLocation(LOCATION_SLOT);

        if (returnMap < 0 || returnMap === mapId) {
            returnMap = fallbackMap;
        }

        plr.clearSavedLocation(LOCATION_SLOT);
        plr.warp(returnMap);
    }
} else {
    npc.sendOk("I only handle subway tickets and waiting room exits.");
}
