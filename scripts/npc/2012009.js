var coupon = 5152004;
var maleFaces = [20025, 20026, 20027, 20028, 20029];
var femaleFaces = [21025, 21026, 21027, 21028, 21029];

npc.sendSelection("Hey there! If you have a #b#t" + coupon + "##k, I can change your face for you!\r\n#L0#I already have a Coupon!#l");
if (npc.selection() === 0) {
    var src = plr.gender() < 1 ? maleFaces : femaleFaces;
    var newFace = src[Math.floor(Math.random() * src.length)] + (plr.face() % 1000);
    if (!plr.haveItem(coupon, 1)) npc.sendOk("Hmm... looks like you don't have a #b#t" + coupon + "##k.");
    else { plr.gainItem(coupon, -1); plr.setFace(newFace); npc.sendOk("All done! What do you think of your new look?"); }
}
