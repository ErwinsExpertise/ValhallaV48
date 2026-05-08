var vipFaceCoupon = 5152022;
var base = Math.floor(plr.face() / 100 % 10) * 100;
var faces = plr.gender() < 1
    ? [20018, 20019, 20000, 20001, 20003, 20004, 20005, 20006, 20007, 20008]
    : [21018, 21019, 21001, 21002, 21003, 21004, 21005, 21006, 21007, 21012];

for (var i = 0; i < faces.length; i++) faces[i] += base;

var choice = npc.askAvatar.apply(npc, ["Ready to look like a million mesos? With #b#t" + vipFaceCoupon + "##k, I can give you the exact face you want."].concat(faces));
if (choice < 0 || choice >= faces.length) npc.sendOk("Changed your mind? Come back any time.");
else if (!plr.haveItem(vipFaceCoupon, 1)) npc.sendOk("You don't have the right coupon for that service.");
else {
    plr.gainItem(vipFaceCoupon, -1);
    plr.setFace(faces[choice]);
    npc.sendOk("Your new face is a true work of art. Come back any time you want another look.");
}
