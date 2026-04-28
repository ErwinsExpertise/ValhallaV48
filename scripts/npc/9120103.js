npc.sendSelection(
    "Hi, I pretty much shouldn't be doing this, but with a #b#t5152056##k or #b#t5152046##k, I will do it anyway for you. But don't forget, it will be random! Now, what would you like to use?\r\n"
    + "#L0#Plastic Surgery at Showa (REG coupon)#l\r\n"
    + "#L1#Cosmetic Lenses at Showa (REG coupon)#l"
);
var townChoice = npc.selection();

if (townChoice === 0) {
    var baseFaces = plr.gender() < 1 ? [20000, 20016, 20019, 20020, 20021, 20024, 20026] : [21000, 21002, 21009, 21016, 21022, 21025, 21027];
    var variant = parseInt(plr.face() / 100 % 10) * 100;
    var newFace = baseFaces[Math.floor(Math.random() * baseFaces.length)] + variant;
    if (!plr.haveItem(5152056, 1)) {
        npc.sendOk("Hmm ... it looks like you don't have the coupon specifically for this place. Sorry to say this, but without the coupon, there's no plastic surgery for you...");
    } else {
        plr.gainItem(5152056, -1);
        plr.setFace(newFace);
        npc.sendOk("Okay, the surgery's done. Here's a mirror--check it out. What a masterpiece, no? Haha! If you ever get tired of this look, please feel free to come visit me again.");
    }
} else if (townChoice === 1) {
    var teye2 = (plr.face() % 100) + (plr.gender() < 1 ? 20000 : 21000);
    var colorList = [100, 200, 300, 400, 500, 600, 700];
    var lens = teye2 + colorList[Math.floor(Math.random() * colorList.length)];
    if (!plr.haveItem(5152046, 1)) {
        npc.sendOk("I'm sorry, but I don't think you have our cosmetic lens coupon with you right now. Without the coupon, I'm afraid I can't do it for you..");
    } else {
        plr.gainItem(5152046, -1);
        plr.setFace(lens);
        npc.sendOk("Here's the mirror. What do you think? I think they look tailor-made for you. I have to say, you look fabulous. Please come again.");
    }
}
