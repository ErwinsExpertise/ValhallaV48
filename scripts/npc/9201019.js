var regularFaceCoupon = 5152021;
var base = Math.floor(plr.face() / 100 % 10) * 100;
var faces = plr.gender() < 1
    ? [20018, 20019, 20000, 20001, 20003, 20004, 20005, 20006, 20007, 20008]
    : [21018, 21019, 21001, 21002, 21003, 21004, 21005, 21006, 21007, 21012];

for (var i = 0; i < faces.length; i++) faces[i] += base;

if (!npc.sendYesNo("I can give you a lovely new look with #b#t" + regularFaceCoupon + "##k. Your new face will be chosen at random. Do you want to try it?")) {
    npc.sendOk("Think it over and come back if you decide you want the change.");
} else if (!plr.haveItem(regularFaceCoupon, 1)) {
    npc.sendOk("You don't have the right coupon for that service.");
} else {
    plr.gainItem(regularFaceCoupon, -1);
    plr.setFace(faces[Math.floor(Math.random() * faces.length)]);
    npc.sendOk("There we go. Your new face came out wonderfully.");
}
