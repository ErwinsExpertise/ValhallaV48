var regFaceCoupon = 5152027
var regLensCoupon = 5152042
var maleFaces = [20002, 20005, 20007, 20011, 20014, 20017, 20029]
var femaleFaces = [21001, 21010, 21013, 21018, 21020, 21021, 21030]

var mode = npc.sendMenu("Hey, I'm Noma, and I am assisting Pata in changing faces and applying lenses as part of my internship. With #b#t" + regFaceCoupon + "##k or #b#t" + regLensCoupon + "##k, I can change the way you look. Now, what would you like to use?", "Plastic Surgery: #i" + regFaceCoupon + "##t" + regFaceCoupon + "#", "Cosmetic Lenses: #i" + regLensCoupon + "##t" + regLensCoupon + "#")

if (mode === 0) {
    if (!npc.sendYesNo("If you use the regular coupon, your face may transform into a random new look... do you still want to do it using #b#t" + regFaceCoupon + "##k?")) {
        npc.sendOk("Changed your mind? That's fine. Come back any time.")
    } else if (!plr.haveItem(regFaceCoupon, 1)) {
        npc.sendOk("I'm sorry, but I don't think you have our plastic surgery coupon with you right now. Without the coupon, I'm afraid I can't do it for you..")
    } else {
        var src = plr.gender() < 1 ? maleFaces : femaleFaces
        var variants = plr.face() % 1000 - (plr.face() % 100)
        var options = []
        for (var i = 0; i < src.length; i++) options.push(src[i] + variants)
        plr.gainItem(regFaceCoupon, -1)
        plr.setFace(options[Math.floor(Math.random() * options.length)])
        npc.sendOk("Enjoy your new and improved face!")
    }
} else if (mode === 1) {
    if (!npc.sendYesNo("If you use the regular coupon, you'll be awarded a random pair of cosmetic lenses. Are you going to use a #b#t" + regLensCoupon + "##k and really make the change to your eyes?")) {
        npc.sendOk("Changed your mind? That's fine. Come back any time.")
    } else if (!plr.haveItem(regLensCoupon, 1)) {
        npc.sendOk("I'm sorry, but I don't think you have our cosmetic lens coupon with you right now. Without the coupon, I'm afraid I can't do it for you..")
    } else {
        var base = (plr.face() % 100) + (plr.gender() < 1 ? 20000 : 21000)
        var colors = [base, base + 100, base + 300, base + 500, base + 600, base + 700]
        plr.gainItem(regLensCoupon, -1)
        plr.setFace(colors[Math.floor(Math.random() * colors.length)])
        npc.sendOk("Enjoy your new and improved cosmetic lenses!")
    }
}
