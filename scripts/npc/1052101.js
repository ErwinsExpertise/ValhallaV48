var couponCut = 5150002;
var couponDye = 5151002;
var maleHair = [30000, 30020, 30030, 30040, 30050, 30110, 30130, 30160, 30180, 30190, 30350, 30610, 30440, 30400];
var femaleHair = [31000, 31010, 31020, 31040, 31050, 31060, 31090, 31120, 31130, 31140, 31330, 31700, 31620, 31610];

npc.sendSelection("Hello! I'm the assistant at the Kerning salon. Got #b#t" + couponCut + "##k or #b#t" + couponDye + "##k? Let's give you a new look!\r\n#L0#Haircut: #i" + couponCut + "##t" + couponCut + "##l\r\n#L1#Dye your hair: #i" + couponDye + "##t" + couponDye + "##l");
var selection = npc.selection();
if (selection === 0) {
    var src = plr.gender() < 1 ? maleHair : femaleHair;
    var newHair = src[Math.floor(Math.random() * src.length)] + (plr.hair() % 10);
    if (!plr.haveItem(couponCut, 1)) npc.sendOk("Looks like you're missing a #b#t" + couponCut + "##k. Sorry, I can't do it without that!");
    else { plr.gainItem(couponCut, -1); plr.setHair(newHair); npc.sendOk("All done! What do you think? A fresh new cut for a fresh new start!"); }
} else if (selection === 1) {
    var base = Math.floor(plr.hair() / 10) * 10;
    var colors = [base, base + 1, base + 2, base + 3, base + 4, base + 5, base + 6, base + 7];
    var newColor = colors[Math.floor(Math.random() * colors.length)];
    if (!plr.haveItem(couponDye, 1)) npc.sendOk("No #b#t" + couponDye + "##k, no color. Sorry, friend!");
    else { plr.gainItem(couponDye, -1); plr.setHair(newColor); npc.sendOk("Your new color is poppin'! Come back anytime!"); }
}
