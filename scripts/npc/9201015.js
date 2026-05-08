var hairCoupon = 5150020;
var colorCoupon = 5151017;
var membershipCoupon = 5420000;
var couponPrice = 1000000;
var maleHair = [30580, 30590, 30280, 30670, 30410, 30200, 30050, 30230, 30290, 30300, 30250];
var femaleHair = [31580, 31590, 31310, 31200, 31150, 31160, 31020, 31260, 31230, 31220, 31110];

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

var selection = npc.sendMenu(
    "Welcome! My name is Julius Styleman. If you have #b#t" + hairCoupon + "##k, #b#t" + colorCoupon + "##k, or #b#t" + membershipCoupon + "##k, I'll make your hair unforgettable.",
    "Buy a coupon",
    "VIP haircut",
    "VIP hair color",
    "VIP membership haircut"
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
} else if (selection === 1 || selection === 3) {
    var couponId = selection === 1 ? hairCoupon : membershipCoupon;
    var styles = buildStyles(plr.gender() < 1 ? maleHair : femaleHair);
    var choice = npc.askAvatar.apply(npc, [selection === 1
        ? "Choose the VIP hairstyle you like."
        : "Choose the hairstyle you want to start with for your VIP membership coupon."].concat(styles));
    if (choice < 0 || choice >= styles.length) npc.sendOk("Changed your mind? Come back any time.");
    else if (!plr.haveItem(couponId, 1)) npc.sendOk("You don't have the right coupon for that service.");
    else {
        plr.gainItem(couponId, -1);
        plr.setHair(styles[choice]);
        npc.sendOk("Take a look. That's a work of art if I do say so myself.");
    }
} else if (selection === 2) {
    var colors = buildColors();
    var colorChoice = npc.askAvatar.apply(npc, ["Choose the VIP hair color you like."].concat(colors));
    if (colorChoice < 0 || colorChoice >= colors.length) npc.sendOk("Changed your mind? Come back any time.");
    else if (!plr.haveItem(colorCoupon, 1)) npc.sendOk("You don't have the right coupon for that service.");
    else {
        plr.gainItem(colorCoupon, -1);
        plr.setHair(colors[colorChoice]);
        npc.sendOk("That color suits you nicely.");
    }
} else {
    npc.sendOk("Come back if you want a new style.");
}
