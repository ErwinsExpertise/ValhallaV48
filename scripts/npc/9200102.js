var regCoupon = 5152012;
var vipCoupon = 5152015;

function currentLensBase() {
    return (plr.face() % 100) + (plr.gender() < 1 ? 20000 : 21000);
}

npc.sendSelection(
    "Um... hi, I'm Dr. Bosch, and I am a cosmetic lens expert here at the Ludibrium Plastic Surgery Shop. I believe your eyes are the most important feature in your body, and with #b#t" + regCoupon + "##k or #b#t" + vipCoupon + "##k, I can prescribe the right kind of cosmetic lenses for you. Now, what would you like to use?\r\n"
    + "#L0#Cosmetic Lenses: #i" + regCoupon + "##t" + regCoupon + "##l\r\n"
    + "#L1#Cosmetic Lenses: #i" + vipCoupon + "##t" + vipCoupon + "##l"
);

var selection = npc.selection();
if (selection === 0) {
    if (!npc.sendYesNo("If you use the regular coupon, I'll have to warn you that you'll be awarded a random pair of cosmetic lenses. Are you going to use #b#t" + regCoupon + "##k and really make the change to your eyes?")) {
        npc.sendOk("Changed your mind? That's fine. Come back any time.");
    } else if (!plr.haveItem(regCoupon, 1)) {
        npc.sendOk("I'm sorry, but I don't think you have our cosmetic lens coupon with you right now. Without the coupon, I'm afraid I can't do it for you..");
    } else {
        var base = currentLensBase();
        var colors = [base, base + 100, base + 200, base + 300, base + 400, base + 500, base + 600, base + 700];
        plr.gainItem(regCoupon, -1);
        plr.setFace(colors[Math.floor(Math.random() * colors.length)]);
        npc.sendOk("Here's the mirror. What do you think? I think they look tailor-made for you. I have to say, you look fabulous. Please come again.");
    }
} else if (selection === 1) {
    var base2 = currentLensBase();
    var options = [base2, base2 + 100, base2 + 200, base2 + 300, base2 + 400, base2 + 500, base2 + 600, base2 + 700];
    var choice = npc.askAvatar.apply(npc, ["With our specialized machine, you can see yourself after the treatment in advance. What kind of lens would you like to wear? Choose the style of your liking..."].concat(options));
    if (choice < 0 || choice >= options.length) {
        npc.sendOk("Changed your mind? That's fine. Come back any time.");
    } else if (!plr.haveItem(vipCoupon, 1)) {
        npc.sendOk("I'm sorry, but I don't think you have our cosmetic lens coupon with you right now. Without the coupon, I'm afraid I can't do it for you..");
    } else {
        plr.gainItem(vipCoupon, -1);
        plr.setFace(options[choice]);
        npc.sendOk("Here's the mirror. What do you think? I think they look tailor-made for you. I have to say, you look fabulous. Please come again.");
    }
}
