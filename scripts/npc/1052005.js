var coupon = 5152000;
var maleFaces = [20005, 20006, 20007, 20008, 20009];
var femaleFaces = [21005, 21006, 21007, 21008, 21009];

npc.sendSelection("Hey there! If you have a #b#t" + coupon + "##k, I can change your face for you!\r\n#L0#I already have a Coupon!#l");
if (npc.selection() === 0) {
    var src = plr.gender() < 1 ? maleFaces : femaleFaces;
    var newFace = src[Math.floor(Math.random() * src.length)] + (plr.face() % 1000);
    if (!plr.haveItem(coupon, 1)) npc.sendOk("Hmm... looks like you don't have a #b#t" + coupon + "##k.");
    else { plr.gainItem(coupon, -1); plr.setFace(newFace); npc.sendOk("All done! What do you think of your new look?"); }
}
