var coupon = 5152005;
var maleFaces = [20000, 20001, 20002, 20003, 20004, 20005, 20006, 20007, 20008, 20012, 20014];
var femaleFaces = [21000, 21001, 21002, 21003, 21004, 21005, 21006, 21007, 21008, 21012, 21014];

npc.sendSelection("Well well well, welcome to the Orbis Plastic Surgery! Would you like to transform your face into something new? With a #b#t" + coupon + "##k, you can let us take care of the rest and have the face you've always wanted~!\r\n#L0#I already have a Coupon!#l");
if (npc.selection() === 0) {
    var variants = plr.face() % 1000 - (plr.face() % 100);
    var src = plr.gender() < 1 ? maleFaces : femaleFaces;
    var faces = [];
    for (var i = 0; i < src.length; i++) faces.push(src[i] + variants);
    var choice = npc.askAvatar.apply(npc, ["I can totally transform your face into something new... how about giving us a try? For #b#t" + coupon + "##k, you can get the face of your liking...take your time in choosing the face of your preference."].concat(faces));
    if (choice < 0 || choice >= faces.length) npc.sendOk("Changed your mind? That's fine. Come back any time.");
    else if (!plr.haveItem(coupon, 1)) npc.sendOk("Hmm ... it looks like you don't have the coupon specifically for this place. Sorry to say this, but without the coupon, there's no plastic surgery for you...");
    else { plr.gainItem(coupon, -1); plr.setFace(faces[choice]); npc.sendOk("Enjoy your new and improved face!"); }
}
