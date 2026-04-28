var cost = 20000;
var ticket = 4031731;
var props = map.properties();
var boardingOpen = ("canBoard" in props) && props["canBoard"];

npc.sendSelection("Hello there~ I am #p9270041# from Singapore Airport. I was transferred to #m103000000# to celebrate new opening of our service! How can i help you?\r\n#L0##bI would like to buy a plane ticket to Singapore#k#l\r\n#L1##bLet me go in to the departure point.#k#l");

var selection = npc.selection();
if (selection === 0) {
    if (!npc.sendYesNo("The ticket will cost you 20,000 mesos. Will you purchase the ticket?")) {
        npc.sendOk("I am here for a long time. Please talk to me again when you change your mind.");
    } else if (!plr.canHold(ticket, 1) || plr.getMesos() < cost) {
        npc.sendOk("I don't think you have enough meso or empty slot in your ETC inventory. Please check and talk to me again.");
    } else {
        plr.gainMesos(-cost);
        plr.gainItem(ticket, 1);
    }
} else if (selection === 1) {
    if (!npc.sendYesNo("Would you like to go in now? You will lose your ticket once you go in~ Thank you for choosing Wizet Airline.")) {
        npc.sendOk("Please confirm the departure time you wish to leave. Thank you.");
    } else if (!plr.haveItem(ticket, 1)) {
        npc.sendOk("Please do purchase the ticket first. Thank you~");
    } else if (!boardingOpen) {
        npc.sendOk("We are sorry but the gate is closed 1 minute before the departure.");
    } else {
        plr.gainItem(ticket, -1);
        plr.warp(540010100);
    }
}
