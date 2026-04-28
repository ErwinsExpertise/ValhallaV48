var coupon = 5153002;
var skins = [0, 1, 2, 3, 4];

npc.sendNext("Oh, hello! Welcome to the Ludibrium Skin-Care! Are you interested in getting tanned and looking sexy? How about a beautiful, snow-white skin? If you have #b#t" + coupon + "##k, you can let us take care of the rest and have the kind of skin you've always dreamed of!");

var choice = npc.sendAvatar.apply(
    npc,
    ["With our specialized machine, you can see the way you'll look after the treatment in advance. What kind of skin-treatment would you like to do? Choose the style of your liking..."].concat(skins)
);

if (choice < 0 || choice >= skins.length) {
    npc.sendOk("Changed your mind? That's fine. Come back any time.");
} else if (!plr.haveItem(coupon, 1)) {
    npc.sendOk("It looks like you don't have the coupon you need to receive the treatment. I'm sorry but it looks like we cannot do it for you.");
} else {
    plr.gainItem(coupon, -1);
    plr.setSkinColor(skins[choice]);
    npc.sendOk("Here's the mirror, check it out! Doesn't your skin look beautiful and glowing like mine? Hehe, it's wonderful. Please come again!");
}
