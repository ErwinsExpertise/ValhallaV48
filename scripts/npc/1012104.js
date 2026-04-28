var couponCut = 5150000;
var couponDye = 5151000;
var maleHair = [30000, 30020, 30030, 30060, 30120, 30140, 30150, 30200, 30210, 30310, 30330, 30410];
var femaleHair = [31000, 31030, 31040, 31050, 31080, 31070, 31100, 31150, 31160, 31300, 31310, 31410];

npc.sendSelection("Hi there! I'm the assistant here in Henesys. If you have #b#t" + couponCut + "##k or #b#t" + couponDye + "##k, I can help you change up your look!\r\n#L0#Haircut: #i" + couponCut + "##t" + couponCut + "##l\r\n#L1#Dye your hair: #i" + couponDye + "##t" + couponDye + "##l");
var selection = npc.selection();
if (selection === 0) {
    var src = plr.gender() < 1 ? maleHair : femaleHair;
    var newHair = src[Math.floor(Math.random() * src.length)] + (plr.hair() % 10);
    if (!plr.haveItem(couponCut, 1)) npc.sendOk("It seems you don't have a #b#t" + couponCut + "##k. Sorry, I can't cut your hair without it.");
    else { plr.gainItem(couponCut, -1); plr.setHair(newHair); npc.sendOk("Take a look! I think this new style suits you perfectly!"); }
} else if (selection === 1) {
    var base = Math.floor(plr.hair() / 10) * 10;
    var colors = [base, base + 1, base + 2, base + 3, base + 4, base + 5, base + 6, base + 7];
    var newColor = colors[Math.floor(Math.random() * colors.length)];
    if (!plr.haveItem(couponDye, 1)) npc.sendOk("You don't have a #b#t" + couponDye + "##k. Sorry, I can't dye your hair without it.");
    else { plr.gainItem(couponDye, -1); plr.setHair(newColor); npc.sendOk("Your new color looks great! Come back if you ever want another change!"); }
}
