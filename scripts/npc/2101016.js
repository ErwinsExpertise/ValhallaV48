var jewelry = plr.itemCount(4031868);

if (jewelry <= 5) {
    plr.removeAll(4031868);
    npc.sendNext("                                  #e<APQ>#n\r\n\r\nBring more #eJewelry#n next time if you want to earn more #eexperience#n.");
    plr.warp(980010020);
} else {
    plr.removeAll(4031868);
    npc.sendNext("                                  #e<APQ>#n\r\n\r\nThanks for the #b#eJewelry#k#n.");
    plr.giveEXP(100 * jewelry);
    plr.warp(980010020);
}
