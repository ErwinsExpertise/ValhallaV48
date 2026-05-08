var hairCoupon = 5150019;
var colorCoupon = 5151016;
var couponPrice = 1000000;
var maleHair = [30570, 30450, 30410, 30200, 30050, 30230, 30290, 30300, 30250, 30690];
var femaleHair = [31570, 31480, 31150, 31160, 31020, 31260, 31230, 31220, 31110, 31490];

function buildStyles(source) {
    var out = [];
    var suffix = plr.hair() % 10;
    for (var i = 0; i < source.length; i++) out.push(source[i] + suffix);
    return out;
}

function buildColors() {
    var base = Math.floor(plr.hair() / 10) * 10;
    return [base, base + 1, base + 3, base + 4, base + 5, base + 7];
}

function randomFrom(list) {
    return list[Math.floor(Math.random() * list.length)];
}

var selection = npc.sendMenu(
    "How's it going? I've got some fresh experimental looks if you want to try them. If you have #b#t" + hairCoupon + "##k or #b#t" + colorCoupon + "##k, let me work on your hair.",
    "Buy a coupon",
    "EXP haircut",
    "EXP hair color"
);

if (selection === 0) {
    var buy = npc.sendMenu("Which coupon would you like to buy?", "#t" + hairCoupon + "# (1,000,000 mesos)", "#t" + colorCoupon + "# (1,000,000 mesos)");
    var coupon = buy === 0 ? hairCoupon : colorCoupon;
    if (buy < 0 || buy > 1) npc.sendOk("Come back if you decide you want a coupon.");
    else if (plr.getMesos() < couponPrice) npc.sendOk("You don't have enough mesos to buy that coupon.");
    else if (!plr.canHold(coupon, 1)) npc.sendOk("Make room in your inventory first.");
    else {
        plr.gainMesos(-couponPrice);
        plr.gainItem(coupon, 1);
        npc.sendOk("Enjoy!");
    }
} else if (selection === 1) {
    if (!npc.sendYesNo("If you use the EXP coupon, your hairstyle will change randomly. Do you want to try it?")) npc.sendOk("Take your time and come back when you're ready.");
    else if (!plr.haveItem(hairCoupon, 1)) npc.sendOk("You don't have the right coupon for that service.");
    else {
        plr.gainItem(hairCoupon, -1);
        plr.setHair(randomFrom(buildStyles(plr.gender() < 1 ? maleHair : femaleHair)));
        npc.sendOk("Here's the mirror. Your new look turned out better than I expected.");
    }
} else if (selection === 2) {
    if (!npc.sendYesNo("If you use #b#t" + colorCoupon + "##k, your hair color will change randomly. Do you still want to try it?")) npc.sendOk("Take your time and come back when you're ready.");
    else if (!plr.haveItem(colorCoupon, 1)) npc.sendOk("You don't have the right coupon for that service.");
    else {
        plr.gainItem(colorCoupon, -1);
        plr.setHair(randomFrom(buildColors()));
        npc.sendOk("Here's the mirror. Your new color came out great.");
    }
} else {
    npc.sendOk("Come back if you want a new experimental look.");
}
