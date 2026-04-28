var couponCut = 5150007;
var couponDye = 5151007;
var maleHair = [30030, 30020, 30000, 30250, 30190, 30150, 30050, 30280, 30240, 30300, 30160];
var femaleHair = [31040, 31000, 31150, 31280, 31160, 31120, 31290, 31270, 31030, 31230, 31010];

npc.sendSelection("Welcome, welcome, welcome to the Ludibrium Hair Salon! Do you, by any chance, have a #b#t" + couponCut + "##k or a #b#t" + couponDye + "##k? If so, how about letting me take care of your hair? Please choose what you want to do with it...\r\n#L0#Haircut: #i" + couponCut + "##t" + couponCut + "##l\r\n#L1#Dye your hair: #i" + couponDye + "##t" + couponDye + "##l");
var selection = npc.selection();
if (selection === 0) {
    var src = plr.gender() < 1 ? maleHair : femaleHair;
    var styles = [];
    for (var i = 0; i < src.length; i++) styles.push(src[i] + (plr.hair() % 10));
    var choice = npc.askAvatar.apply(npc, ["I can completely change the look of your hair. Aren't you ready for a change? With #b#t" + couponCut + "##k, I'll take care of the rest for you. Choose the style of your liking!"].concat(styles));
    if (choice < 0 || choice >= styles.length) npc.sendOk("Changed your mind? That's fine. Come back any time.");
    else if (!plr.haveItem(couponCut, 1)) npc.sendOk("Hmmm...it looks like you don't have our designated coupon...I'm afraid I can't give you a haircut without it. I'm sorry...");
    else { plr.gainItem(couponCut, -1); plr.setHair(styles[choice]); npc.sendOk("Enjoy your new and improved hairstyle!"); }
} else if (selection === 1) {
    var base = Math.floor(plr.hair() / 10) * 10;
    var colors = [base, base + 1, base + 2, base + 3, base + 4, base + 5, base + 6, base + 7];
    var colorChoice = npc.askAvatar.apply(npc, ["I can completely change the color of your hair. Aren't you ready for a change? With #b#t" + couponDye + "##k, I'll take care of the rest. Choose the color of your liking!"].concat(colors));
    if (colorChoice < 0 || colorChoice >= colors.length) npc.sendOk("Changed your mind? That's fine. Come back any time.");
    else if (!plr.haveItem(couponDye, 1)) npc.sendOk("Hmmm...it looks like you don't have our designated coupon...I'm afraid I can't dye your hair without it. I'm sorry...");
    else { plr.gainItem(couponDye, -1); plr.setHair(colors[colorChoice]); npc.sendOk("Enjoy your new and improved haircolor!"); }
}
