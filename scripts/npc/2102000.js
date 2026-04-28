var TICKET_ID = 4031045;

if (!npc.sendYesNo("It seems that there is still plenty of room on this ride. Please have your ticket ready so I can let you on. The journey will be long, but you will get to your destination safely. What do you think? Do you want to take this ride?")) {
    npc.sendOk("You must have some business to take care of here, right?");
} else if (plr.haveItem(TICKET_ID, 1)) {
    plr.gainItem(TICKET_ID, -1);
    plr.warp(200000100);
} else {
    npc.sendOk("Oh, no... It looks like you do not have a ticket with you. I cannot let you on without it. Please buy the ticket at the ticket sales guide.");
}
