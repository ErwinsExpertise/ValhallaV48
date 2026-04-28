var TICKET_ID = 4031331;

if (!npc.sendYesNo("It seems that there is still plenty of room for this ride. Please have your ticket ready so I can let you on. The journey will be long, but you will get to your destination safely. What do you think? Do you want take this ride?")) {
    npc.sendOk("You must have some business to take care of here, right?");
} else if (plr.haveItem(TICKET_ID, 1)) {
    plr.gainItem(TICKET_ID, -1);
    plr.warp(240000100);
} else {
    npc.sendOk("Oh, no... It looks like you don't have a ticket with you. I can't let you on without it. Please buy the ticket from the ticket guide.");
}
