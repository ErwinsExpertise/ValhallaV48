if (!npc.sendYesNo("                          #e<Ariant PQ>#n\r\n\r\nWould you like to go to the #bAriant Coliseum#k?\r\nYou should be between level #e20 and 30#n to participate.")) {
    npc.sendOk("Talk to me again if you decide to participate.");
} else if ((plr.getLevel() >= 20 && plr.getLevel() < 31) || plr.isGM()) {
    plr.saveLocation("ARIANT_PQ");
    plr.warp(980010000);
} else {
    npc.sendOk("You are not between level 20 and 30. Sorry, you cannot participate.");
}
