var TICKET_ID = 4031045;
var PRICE = 6000;

if (!npc.sendYesNo("Hello, I'm in charge of selling tickets for the ship ride to Orbis Station of Ossyria. The ride to Orbis takes off every 10 minutes, beginning on the hour, and it'll cost you #b" + PRICE + " mesos#k. Are you sure you want to purchase #b#t" + TICKET_ID + "##k?")) {
    npc.sendOk("You must have some business to take care of here, right?");
} else if (plr.getMesos() < PRICE || !plr.canHold(TICKET_ID, 1)) {
    npc.sendOk("Are you sure you have #b" + PRICE + " mesos#k? If so, then I urge you to check your etc. inventory, and see if it's full or not.");
} else if (!plr.gainItem(TICKET_ID, 1)) {
    npc.sendOk("Are you sure you have #b" + PRICE + " mesos#k? If so, then I urge you to check your etc. inventory, and see if it's full or not.");
} else {
    plr.gainMesos(-PRICE);
}
