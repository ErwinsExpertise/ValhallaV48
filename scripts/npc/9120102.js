var vipFaceCoupon = 5152057;
var vipLensCoupon = 5152045;

npc.sendSelection("I'm in charge of the Plastic Surgery here at Showa Shop! I believe your eyes are the most important feature in your body, and with #b#t" + vipFaceCoupon + "##k or #b#t" + vipLensCoupon + "##k, I can prescribe the right kind of plastic surgery and cosmetic lenses for you. Now, what would you like to use?\r\n#L0#Plastic Surgery at Showa (VIP coupon)#l\r\n#L1#Cosmetic Lenses at Showa (VIP coupon)#l");
var selection = npc.selection();

if (selection === 0) {
    var faces = plr.gender() < 1 ? [20020, 20000, 20002, 20004, 20005, 20012] : [21021, 21000, 21002, 21003, 21006, 21012];
    for (var i = 0; i < faces.length; i++) faces[i] = faces[i] + parseInt(plr.face() / 100 % 10) * 100;
    var chosenFace = npc.askAvatar.apply(npc, ["Let's see... for #b#t" + vipFaceCoupon + "##k, you can get a new face. That's right, I can completely transform your face! Wanna give it a shot? Please consider your choice carefully."].concat(faces));
    if (chosenFace < 0 || chosenFace >= faces.length) npc.sendOk("Changed your mind? That's fine. Come back any time.");
    else if (!plr.haveItem(vipFaceCoupon, 1)) npc.sendOk("Hmm ... it looks like you don't have the coupon specifically for this place. Sorry to say this, but without the coupon, there's no plastic surgery for you...");
    else { plr.gainItem(vipFaceCoupon, -1); plr.setFace(faces[chosenFace]); npc.sendOk("Alright, it's all done! Check yourself out in the mirror. Well, aren't you lookin' marvelous? Haha! If you're sick of it, just give me another call, alright?"); }
} else if (selection === 1) {
    var teye = (plr.face() % 100) + (plr.gender() < 1 ? 20000 : 21000);
    var colors = [teye + 100, teye + 200, teye + 300, teye + 400, teye + 500, teye + 600, teye + 700];
    var chosenLens = npc.askAvatar.apply(npc, ["With our specialized machine, you can see the results of your potential treatment in advance. What kind of lens would you like to wear? Choose the style of your liking..."].concat(colors));
    if (chosenLens < 0 || chosenLens >= colors.length) npc.sendOk("Changed your mind? That's fine. Come back any time.");
    else if (!plr.haveItem(vipLensCoupon, 1)) npc.sendOk("I'm sorry, but I don't think you have our cosmetic lens coupon with you right now. Without the coupon, I'm afraid I can't do it for you..");
    else { plr.gainItem(vipLensCoupon, -1); plr.setFace(colors[chosenLens]); npc.sendOk("Enjoy your new and improved cosmetic lenses!"); }
}
