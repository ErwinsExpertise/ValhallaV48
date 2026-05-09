var LOCATION_SLOT = "SUBWAY";
var props = map.properties();
var boardingOpen = ("canBoard" in props) && props["canBoard"];

var choices = [];

function canBoardNLC() {
    if (!boardingOpen) {
        npc.sendOk("The subway is not boarding right now. Please wait for the next departure window.");
        return false;
    }
    return true;
}

function addChoice(type, label, ticketId, mapId) {
    if (plr.haveItem(ticketId, 1)) {
        choices.push({ type: type, label: label, ticketId: ticketId, mapId: mapId });
    }
}

addChoice("construction", "Construction Site B1", 4031036, 103000900);
addChoice("construction", "Construction Site B2", 4031037, 103000903);
addChoice("construction", "Construction Site B3", 4031038, 103000906);
addChoice("nlc", "New Leaf City (for Beginners)", 4031710, 600010004);
addChoice("nlc", "New Leaf City (Regular)", 4031711, 600010004);

if (choices.length === 0) {
    npc.sendOk("This is the ticket reader. You won't be allowed to go in without a ticket.");
} else if (choices.length === 1) {
    var only = choices[0];
    if (only.type === "construction") {
        if (npc.sendYesNo("This is the ticket reader. Will you use #b#t" + only.ticketId + "##k? If you do, you'll be moved inside right away.")) {
            if (!plr.gainItem(only.ticketId, -1)) npc.sendOk("Please insert #b#t" + only.ticketId + "##k into the ticket reader.");
            else plr.warp(only.mapId);
        }
    } else if (npc.sendYesNo("Please have your ticket ready. I will send you to the waiting room for the train to New Leaf City. Do you want to go in now?")) {
        if (canBoardNLC() && plr.gainItem(only.ticketId, -1)) {
            plr.saveLocation(LOCATION_SLOT);
            plr.warp(only.mapId);
        } else if (boardingOpen) {
            npc.sendOk("Please insert #b#t" + only.ticketId + "##k into the ticket reader.");
        }
    }
} else {
    var menu = "This is the ticket reader. You'll be moved inside right away. Which ticket would you like to use?";
    for (var i = 0; i < choices.length; i++) {
        menu += "\r\n#b#L" + i + "# " + choices[i].label + "#l#k";
    }
    var sel = npc.sendMenu(menu);
    if (sel >= 0 && sel < choices.length) {
        var choice = choices[sel];
        if (choice.type === "nlc" && !canBoardNLC()) {
        } else if (!plr.gainItem(choice.ticketId, -1)) {
            npc.sendOk("Please insert #b#t" + choice.ticketId + "##k into the ticket reader.");
        } else {
            if (choice.type === "nlc") plr.saveLocation(LOCATION_SLOT);
            plr.warp(choice.mapId);
        }
    }
}
