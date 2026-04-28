var couponCut = 5150004;
var couponDye = 5151004;
var maleHair = [30000, 30020, 30030, 30230, 30240, 30260, 30270, 30280, 30290, 30340, 30610, 30440, 30400];
var femaleHair = [31000, 31030, 31040, 31110, 31220, 31230, 31240, 31250, 31260, 31270, 31320, 31700, 31620, 31610];

npc.sendSelection("Welcome to Orbis Hair! If you have #b#t" + couponCut + "##k or #b#t" + couponDye + "##k, I can take care of your hair!\r\n#L0#Haircut: #i" + couponCut + "##t" + couponCut + "##l\r\n#L1#Dye your hair: #i" + couponDye + "##t" + couponDye + "##l");
var selection = npc.selection();
if (selection === 0) {
    var src = plr.gender() < 1 ? maleHair : femaleHair;
    var newHair = src[Math.floor(Math.random() * src.length)] + (plr.hair() % 10);
    if (!plr.haveItem(couponCut, 1)) npc.sendOk("It looks like you don't have a #b#t" + couponCut + "##k. Sorry, can't help without it!");
    else { plr.gainItem(couponCut, -1); plr.setHair(newHair); npc.sendOk("There we go! What do you think? A new style for a new you!"); }
} else if (selection === 1) {
    var base = Math.floor(plr.hair() / 10) * 10;
    var colors = [base, base + 1, base + 2, base + 3, base + 4, base + 5, base + 6, base + 7];
    var newColor = colors[Math.floor(Math.random() * colors.length)];
    if (!plr.haveItem(couponDye, 1)) npc.sendOk("No #b#t" + couponDye + "##k, no dye job. Sorry!");
    else { plr.gainItem(couponDye, -1); plr.setHair(newColor); npc.sendOk("Beautiful! Your new color shines brighter than the clouds themselves!"); }
}
