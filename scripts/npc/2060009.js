var vipTicket = 4031242;
var unknownMap = 230030200;
var unknownPortal = "st00";
var herbTownMap = 251000100;

var beginner = plr.job() === 0;
var unknownFee = beginner ? 100 : 1000;
var herbTownFee = beginner ? 1000 : 10000;

var text = "The oceans are all connected to each other. Places you can't reach by foot can be reached quickly by sea. What do you think about taking the #bDolphin Taxi#k today?\r\n#b";

if (plr.haveItem(vipTicket, 1)) {
    text += "#L0#Use #t4031242# to travel to #m" + unknownMap + "##l\r\n";
    text += "#L1#Travel to #m" + herbTownMap + "# for #b" + herbTownFee + " mesos#k.#l";
} else {
    if (beginner) {
        text += "#L0#Travel to #m" + unknownMap + "# for #b" + unknownFee + " mesos#k. Apprentices receive a 90% discount.#l\r\n";
    } else {
        text += "#L0#Travel to #m" + unknownMap + "# for #b" + unknownFee + " mesos#k.#l\r\n";
    }
    text += "#L1#Travel to #m" + herbTownMap + "# for #b" + herbTownFee + " mesos#k.#l";
}

npc.sendSelection(text + "#k");
var sel = npc.selection();

if (sel === 0) {
    if (plr.haveItem(vipTicket, 1)) {
        plr.gainItem(vipTicket, -1);
        plr.warpToPortalName(unknownMap, unknownPortal);
    } else if (plr.mesos() < unknownFee) {
        npc.sendOk("I don't think you have enough mesos...");
    } else {
        plr.takeMesos(unknownFee);
        plr.warpToPortalName(unknownMap, unknownPortal);
    }
} else if (sel === 1) {
    if (plr.mesos() < herbTownFee) {
        npc.sendOk("I don't think you have enough mesos...");
    } else {
        plr.takeMesos(herbTownFee);
        plr.warp(herbTownMap);
    }
}
