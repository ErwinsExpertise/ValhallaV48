var regularLensCoupon = 5152025;
var vipLensCoupon = 5152026;

function lensOptions() {
    var base = (plr.face() % 100) + (plr.gender() < 1 ? 20000 : 21000);
    return [base, base + 100, base + 300, base + 400, base + 600, base + 700];
}

var selection = npc.sendMenu(
    "I'm Dr. Roberts, the contact lens specialist here in Amoria. With #b#t" + regularLensCoupon + "##k or #b#t" + vipLensCoupon + "##k, I can prescribe the right lenses for you.",
    "Amoria cosmetic lenses (regular coupon)",
    "Amoria cosmetic lenses (VIP coupon)"
);

var options = lensOptions();
if (selection === 0) {
    if (!npc.sendYesNo("If you use #b#t" + regularLensCoupon + "##k, you'll receive a random pair of lenses. Do you still want to try it?")) npc.sendOk("That's understandable. Come back if you change your mind.");
    else if (!plr.haveItem(regularLensCoupon, 1)) npc.sendOk("You don't have the right coupon for that service.");
    else {
        plr.gainItem(regularLensCoupon, -1);
        plr.setFace(options[Math.floor(Math.random() * options.length)]);
        npc.sendOk("Take a look. Those lenses fit you perfectly.");
    }
} else if (selection === 1) {
    var choice = npc.askAvatar.apply(npc, ["With our special machine, you can preview the result in advance. Choose the lens style you like."].concat(options));
    if (choice < 0 || choice >= options.length) npc.sendOk("Changed your mind? Come back any time.");
    else if (!plr.haveItem(vipLensCoupon, 1)) npc.sendOk("You don't have the right coupon for that service.");
    else {
        plr.gainItem(vipLensCoupon, -1);
        plr.setFace(options[choice]);
        npc.sendOk("Take a look. Those lenses fit you perfectly.");
    }
} else {
    npc.sendOk("Come back if you decide you want a new pair of lenses.");
}
