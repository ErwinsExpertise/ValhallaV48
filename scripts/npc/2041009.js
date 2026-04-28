var couponCut = 5150012;
var couponDye = 5151006;
var maleHair = [30250, 30190, 30150, 30050, 30280, 30240, 30300, 30160, 30650, 30540, 30640, 30680];
var femaleHair = [31150, 31280, 31160, 31120, 31290, 31270, 31030, 31230, 31010, 31640, 31540, 31680, 31600];

npc.sendSelection("Hi, I'm the assistant here. Don't worry, I'm plenty good enough for this. If you have #b#t" + couponCut + "##k or #b#t" + couponDye + "##k by any chance, then allow me to take care of the rest, alright?\r\n#L0#Haircut: #i" + couponCut + "##t" + couponCut + "##l\r\n#L1#Dye your hair: #i" + couponDye + "##t" + couponDye + "##l");
var selection = npc.selection();
if (selection === 0) {
    if (!npc.sendYesNo("If you use the EXP coupon your hair will change RANDOMLY with a chance to obtain a new experimental style that I came up with. Are you going to use #b#t" + couponCut + "##k and really change your hairstyle?")) npc.sendOk("Changed your mind? That's fine. Come back any time.");
    else if (!plr.haveItem(couponCut, 1)) npc.sendOk("Hmmm...it looks like you don't have our designated coupon...I'm afraid I can't give you a haircut without it. I'm sorry...");
    else {
        var src = plr.gender() < 1 ? maleHair : femaleHair;
        var styles = [];
        for (var i = 0; i < src.length; i++) styles.push(src[i] + (plr.hair() % 10));
        plr.gainItem(couponCut, -1);
        plr.setHair(styles[Math.floor(Math.random() * styles.length)]);
        npc.sendOk("Enjoy your new and improved hairstyle!");
    }
} else if (selection === 1) {
    if (!npc.sendYesNo("If you use a regular coupon your hair will change RANDOMLY. Do you still want to use #b#t" + couponDye + "##k and change it up?")) npc.sendOk("Changed your mind? That's fine. Come back any time.");
    else if (!plr.haveItem(couponDye, 1)) npc.sendOk("Hmmm...it looks like you don't have our designated coupon...I'm afraid I can't dye your hair without it. I'm sorry...");
    else {
        var base = Math.floor(plr.hair() / 10) * 10;
        var colors = [base, base + 1, base + 2, base + 3, base + 4, base + 5, base + 6, base + 7];
        plr.gainItem(couponDye, -1);
        plr.setHair(colors[Math.floor(Math.random() * colors.length)]);
        npc.sendOk("Enjoy your new and improved haircolor!");
    }
}
