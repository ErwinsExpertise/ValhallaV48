var TICKET_ID = 4031713;
var WAITING_ROOM = 600010002;
var LOCATION_SLOT = "SUBWAY";
var props = map.properties();
var boardingOpen = ("canBoard" in props) && props["canBoard"];

if (!plr.haveItem(TICKET_ID, 1)) {
    npc.sendOk("You need a subway ticket before I can let you through.");
} else if (!npc.sendYesNo("Please have your ticket ready. I will send you to the waiting room for the train to Kerning City. Do you want to go in now?")) {
    npc.sendOk("Come back when you are ready to board.");
} else if (!boardingOpen) {
    npc.sendOk("The subway is not boarding right now. Please wait for the next departure window.");
} else if (!plr.gainItem(TICKET_ID, -1)) {
    npc.sendOk("I could not take your ticket. Please try again.");
} else {
    plr.saveLocation(LOCATION_SLOT);
    plr.warp(WAITING_ROOM);
}
