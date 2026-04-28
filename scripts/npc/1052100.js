var couponCut = 5150003;
var couponDye = 5151003;
var maleHair = [30030, 30020, 30000, 30130, 30350, 30190, 30110, 30180, 30050, 30040, 30160, 30780];
var femaleHair = [31050, 31040, 31000, 31060, 31090, 31020, 31130, 31120, 31140, 31330, 31010];

npc.sendSelection("Hello! I'm Don Giovanni, head of the beauty salon! If you have either #b#t" + couponCut + "##k or #b#t" + couponDye + "##k, why don't you let me take care of the rest? Decide what you want to do with your hair...\r\n#L0#Haircut: #i" + couponCut + "##t" + couponCut + "##l\r\n#L1#Dye your hair: #i" + couponDye + "##t" + couponDye + "##l");
var selection = npc.selection();
if (selection === 0) {
    var src = plr.gender() < 1 ? maleHair : femaleHair;
    var styles = [];
    for (var i = 0; i < src.length; i++) styles.push(src[i] + (plr.hair() % 10));
    var choice = npc.askAvatar.apply(npc, ["I can totally change up your hairstyle and make it look so good. Why don't you change it up a bit? If you have #b#t" + couponCut + "##k I'll change it for you. Choose the one to your liking~."].concat(styles));
    if (choice < 0 || choice >= styles.length) npc.sendOk("Changed your mind? That's fine. Come back any time.");
    else if (!plr.haveItem(couponCut, 1)) npc.sendOk("Hmmm...it looks like you don't have our designated coupon...I'm afraid I can't give you a haircut without it. I'm sorry...");
    else { plr.gainItem(couponCut, -1); plr.setHair(styles[choice]); npc.sendOk("Enjoy your new and improved hairstyle!"); }
} else if (selection === 1) {
    var base = Math.floor(plr.hair() / 10) * 10;
    var colors = [base, base + 1, base + 2, base + 3, base + 4, base + 5, base + 6, base + 7];
    var colorChoice = npc.askAvatar.apply(npc, ["I can totally change your haircolor and make it look so good. Why don't you change it up a bit? With #b#t" + couponDye + "##k I'll change it for you. Choose the one to your liking."].concat(colors));
    if (colorChoice < 0 || colorChoice >= colors.length) npc.sendOk("Changed your mind? That's fine. Come back any time.");
    else if (!plr.haveItem(couponDye, 1)) npc.sendOk("Hmmm...it looks like you don't have our designated coupon...I'm afraid I can't dye your hair without it. I'm sorry...");
    else { plr.gainItem(couponDye, -1); plr.setHair(colors[colorChoice]); npc.sendOk("Enjoy your new and improved haircolor!"); }
}
