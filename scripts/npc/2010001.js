var couponCut = 5150005;
var couponDye = 5151005;
var maleHair = [30030, 30020, 30000, 30270, 30230, 30260, 30280, 30240, 30290, 30340];
var femaleHair = [31040, 31000, 31250, 31220, 31260, 31240, 31110, 31270, 31030, 31230];

npc.sendSelection("Hello I'm Mino. If you have either a #b#t" + couponCut + "##k or a #b#t" + couponDye + "##k, then please let me take care of your hair. Choose what you want to do with it.\r\n#L0#Haircut: #i" + couponCut + "##t" + couponCut + "##l\r\n#L1#Dye your hair: #i" + couponDye + "##t" + couponDye + "##l");
var selection = npc.selection();
if (selection === 0) {
    var styles = [];
    var src = plr.gender() < 1 ? maleHair : femaleHair;
    for (var i = 0; i < src.length; i++) styles.push(src[i] + (plr.hair() % 10));
    var choice = npc.askAvatar.apply(npc, ["I can totally change up your hairstyle and make it look so good. Why don't you change it up a bit? With #b#t" + couponCut + "##k, I'll take care of the rest for you. Choose the style of your liking!"].concat(styles));
    if (choice < 0 || choice >= styles.length) npc.sendOk("Changed your mind? That's fine. Come back any time.");
    else if (!plr.haveItem(couponCut, 1)) npc.sendOk("Hmmm...it looks like you don't have our designated coupon...I'm afraid I can't give you a haircut without it. I'm sorry...");
    else { plr.gainItem(couponCut, -1); plr.setHair(styles[choice]); npc.sendOk("Enjoy your new and improved hairstyle!"); }
} else if (selection === 1) {
    var base = Math.floor(plr.hair() / 10) * 10;
    var colors = [base, base + 1, base + 2, base + 3, base + 4, base + 5, base + 6, base + 7];
    var colorChoice = npc.askAvatar.apply(npc, ["I can totally change your haircolor and make it look so good. Why don't you change it up a bit? With #b#t" + couponDye + "##k, I'll take care of the rest. Choose the color of your liking!"].concat(colors));
    if (colorChoice < 0 || colorChoice >= colors.length) npc.sendOk("Changed your mind? That's fine. Come back any time.");
    else if (!plr.haveItem(couponDye, 1)) npc.sendOk("Hmmm...it looks like you don't have our designated coupon...I'm afraid I can't dye your hair without it. I'm sorry...");
    else { plr.gainItem(couponDye, -1); plr.setHair(colors[colorChoice]); npc.sendOk("Enjoy your new and improved haircolor!"); }
}
