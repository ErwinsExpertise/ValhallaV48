var couponCut = 5150009;
var couponDye = 5151009;
var maleHair = [30230, 30030, 30260, 30280, 30240, 30290, 30020, 30270, 30340, 30710, 30810];
var femaleHair = [31310, 31300, 31050, 31040, 31160, 31100, 31410, 31030, 31790, 31550];

npc.sendSelection("Welcome to the Showa hair shop. If you have a #b#t" + couponCut + "##k, or a #b#t" + couponDye + "##k, allow me to take care of your hairdo. Please choose the one you want.\r\n#L0#Haircut: #i" + couponCut + "##t" + couponCut + "##l\r\n#L1#Dye your hair: #i" + couponDye + "##t" + couponDye + "##l");
var selection = npc.selection();

if (selection === 0) {
    var src = plr.gender() < 1 ? maleHair : femaleHair;
    var styles = [];
    for (var i = 0; i < src.length; i++) styles.push(src[i] + (plr.hair() % 10));
    var picked = npc.askAvatar.apply(npc, ["I can totally change up your hairstyle and make it look so good. Why don't you change it up a bit? With #b#t" + couponCut + "##k, I'll take care of the rest for you. Choose the style of your liking!"].concat(styles));
    if (picked < 0 || picked >= styles.length) npc.sendOk("Changed your mind? That's fine. Come back any time.");
    else if (!plr.haveItem(couponCut, 1)) npc.sendOk("Hmmm...it looks like you don't have our designated coupon...I'm afraid I can't give you a haircut without it. I'm sorry...");
    else { plr.gainItem(couponCut, -1); plr.setHair(styles[picked]); npc.sendOk("Enjoy your new and improved hairstyle!"); }
} else if (selection === 1) {
    var base = Math.floor(plr.hair() / 10) * 10;
    var colors = [base, base + 1, base + 2, base + 3, base + 4, base + 5, base + 6, base + 7];
    var colorPick = npc.askAvatar.apply(npc, ["I can totally change your haircolor and make it look so good. Why don't you change it up a bit? With #b#t" + couponDye + "##k, I'll take care of the rest. Choose the color of your liking!"].concat(colors));
    if (colorPick < 0 || colorPick >= colors.length) npc.sendOk("Changed your mind? That's fine. Come back any time.");
    else if (!plr.haveItem(couponDye, 1)) npc.sendOk("Hmmm...it looks like you don't have our designated coupon...I'm afraid I can't dye your hair without it. I'm sorry...");
    else { plr.gainItem(couponDye, -1); plr.setHair(colors[colorPick]); npc.sendOk("Enjoy your new and improved haircolor!"); }
}
