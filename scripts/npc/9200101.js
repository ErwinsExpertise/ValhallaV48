var regCoupon = 5152011;
var vipCoupon = 5152014;

function currentLensBase() {
    return (plr.face() % 100) + (plr.gender() < 1 ? 20000 : 21000);
}

npc.sendSelection(
    "Hello, I'm Dr. Rhomes, head of the cosmetic lens department here at the Orbis Plastic Surgery Shop. My goal here is to add personality to everyone's eyes through the wonders of cosmetic lenses, and with #b#t" + regCoupon + "##k or #b#t" + vipCoupon + "##k, I can do the same for you, too! Now, what would you like to use?\r\n"
    + "#L0#Cosmetic Lenses: #i" + regCoupon + "##t" + regCoupon + "##l\r\n"
    + "#L1#Cosmetic Lenses: #i" + vipCoupon + "##t" + vipCoupon + "##l"
);

var selection = npc.selection();
if (selection === 0) {
    if (!npc.sendYesNo("If you use the regular coupon, you'll be awarded a random pair of cosmetic lenses. Are you going to use a #b#t" + regCoupon + "##k and really make the change to your eyes?")) {
        npc.sendOk("I see... take your time and see if you really want it. Let me know when you've decided.");
    } else if (!plr.haveItem(regCoupon, 1)) {
        npc.sendOk("I'm sorry, but I don't think you have our cosmetic lens coupon with you right now. Without the coupon, I'm afraid I can't do it for you..");
    } else {
        var base = currentLensBase();
        var colors = [base, base + 100, base + 200, base + 300, base + 400, base + 500, base + 600, base + 700];
        plr.gainItem(regCoupon, -1);
        plr.setFace(colors[Math.floor(Math.random() * colors.length)]);
        npc.sendOk("Enjoy your new and improved cosmetic lenses!");
    }
} else if (selection === 1) {
    var base2 = currentLensBase();
    var options = [base2, base2 + 100, base2 + 200, base2 + 300, base2 + 400, base2 + 500, base2 + 600, base2 + 700];
    var choice = npc.askAvatar.apply(npc, ["With our new computer program, you can see yourself after the treatment in advance. What kind of lens would you like to wear? Please choose the style of your liking."].concat(options));
    if (choice < 0 || choice >= options.length) {
        npc.sendOk("Changed your mind? That's fine. Come back any time.");
    } else if (!plr.haveItem(vipCoupon, 1)) {
        npc.sendOk("I'm sorry, but I don't think you have our cosmetic lens coupon with you right now. Without the coupon, I'm afraid I can't do it for you..");
    } else {
        plr.gainItem(vipCoupon, -1);
        plr.setFace(options[choice]);
        npc.sendOk("Enjoy your new and improved cosmetic lenses!");
    }
}
