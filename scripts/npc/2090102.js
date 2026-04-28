var coupon = 5153006
var skins = [0, 1, 2, 3, 4]

npc.sendNext("Well, hello! Welcome to the Mu Lung Skin-Care! Would you like to have a firm, tight, healthy looking skin like mine? With #b#t" + coupon + "##k, you can let us take care of the rest and have the kind of skin you've always wanted~!")

var choice = npc.sendAvatar("With our specialized machine, you can see the way you'll look after the treatment prior to the procedure. What kind of skin-treatment would you like to do? Choose the style of your liking...", skins)

if (choice < 0 || choice >= skins.length) {
    npc.sendOk("Changed your mind? That's fine. Come back any time.")
} else if (plr.itemCount(coupon) < 1) {
    npc.sendOk("Um... you don't have the skin-care coupon you need to receive the treatment. Sorry, but we can't do it for you...")
} else {
    plr.removeItemsByID(coupon, 1)
    plr.setSkinColor(skins[choice])
    npc.sendOk("Enjoy your new and improved skin!")
}
