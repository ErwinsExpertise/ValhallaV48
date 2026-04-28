var vipFaceCoupon = 5152028
var vipLensCoupon = 5152041
var oneTimeBaseCoupon = 5152100
var maleFaces = [20000, 20001, 20004, 20005, 20006, 20007, 20009, 20012, 20022, 20028, 20031]
var femaleFaces = [21000, 21003, 21005, 21006, 21008, 21009, 21011, 21012, 21023, 21024, 21026]

var mode = npc.sendMenu("Hey, I'm Pata, and I am a renowned plastic surgeon and cosmetic lens expert here in Mu Lung. I believe your face and eyes are the most important features in your body, and with #b#t" + vipFaceCoupon + "##k or #b#t" + vipLensCoupon + "##k, I can prescribe the right kind of facial care and cosmetic lenses for you. Now, what would you like to use?", "Plastic Surgery: #i" + vipFaceCoupon + "##t" + vipFaceCoupon + "#", "Cosmetic Lenses: #i" + vipLensCoupon + "##t" + vipLensCoupon + "#", "One-time Cosmetic Lenses")

function faceVariantBase() {
    return plr.face() % 1000 - (plr.face() % 100)
}

function lensBase() {
    return (plr.face() % 100) + (plr.gender() < 1 ? 20000 : 21000)
}

if (mode === 0) {
    var src = plr.gender() < 1 ? maleFaces : femaleFaces
    var faces = []
    var base = faceVariantBase()
    for (var i = 0; i < src.length; i++) faces.push(src[i] + base)
    var choice = npc.askAvatar.apply(npc, ["I can totally transform your face into something new... how about giving us a try? For #b#t" + vipFaceCoupon + "##k, you can get the face of your liking. Take your time in choosing the face of your preference."].concat(faces))
    if (choice < 0 || choice >= faces.length) npc.sendOk("Changed your mind? That's fine. Come back any time.")
    else if (!plr.haveItem(vipFaceCoupon, 1)) npc.sendOk("I'm sorry, but I don't think you have our plastic surgery coupon with you right now. Without the coupon, I'm afraid I can't do it for you..")
    else { plr.gainItem(vipFaceCoupon, -1); plr.setFace(faces[choice]); npc.sendOk("Enjoy your new and improved face!") }
} else if (mode === 1) {
    var baseLens = lensBase()
    var colors = [baseLens, baseLens + 100, baseLens + 300, baseLens + 500, baseLens + 600, baseLens + 700]
    var lensChoice = npc.askAvatar.apply(npc, ["With our new computer program, you can see yourself after the treatment in advance. What kind of lens would you like to wear? Please choose the style of your liking."].concat(colors))
    if (lensChoice < 0 || lensChoice >= colors.length) npc.sendOk("Changed your mind? That's fine. Come back any time.")
    else if (!plr.haveItem(vipLensCoupon, 1)) npc.sendOk("I'm sorry, but I don't think you have our cosmetic lens coupon with you right now. Without the coupon, I'm afraid I can't do it for you..")
    else { plr.gainItem(vipLensCoupon, -1); plr.setFace(colors[lensChoice]); npc.sendOk("Enjoy your new and improved cosmetic lenses!") }
} else if (mode === 2) {
    var baseOne = lensBase()
    var oneTimeOptions = []
    var couponIds = []
    for (var c = 0; c < 8; c++) {
        if (plr.itemCount(oneTimeBaseCoupon + c) > 0) {
            oneTimeOptions.push(baseOne + (100 * c))
            couponIds.push(oneTimeBaseCoupon + c)
        }
    }
    if (oneTimeOptions.length === 0) {
        npc.sendOk("You don't have any One-time Cosmetic Lenses to use.")
    } else {
        var oneChoice = npc.askAvatar.apply(npc, ["What kind of lens would you like to wear? Please choose the style of your liking."].concat(oneTimeOptions))
        if (oneChoice < 0 || oneChoice >= oneTimeOptions.length) npc.sendOk("Changed your mind? That's fine. Come back any time.")
        else if (plr.itemCount(couponIds[oneChoice]) < 1) npc.sendOk("I'm sorry, but I don't think you have the right one-time cosmetic lens coupon with you right now.")
        else { plr.gainItem(couponIds[oneChoice], -1); plr.setFace(oneTimeOptions[oneChoice]); npc.sendOk("Enjoy your new and improved cosmetic lenses!") }
    }
}
