npc.sendNext("Hi there! This cab is for VIP customers only. Instead of just taking you to different towns like the regular cabs, we offer a much better service worthy of VIP class. It's a bit pricey, but... for only 10,000 mesos, we'll take you safely to the \r\n#bAnt Tunnel#k.");

var beginner = plr.job() === 0;
var cost = beginner ? 1000 : 10000;
var prompt = beginner
    ? "We have a special 90% discount for beginners. The Ant Tunnel is located deep inside in the dungeon that's placed at the center of the Victoria Island, where the 24 Hr Mobile Store is. Would you like to go there for #b1,000 mesos#k?"
    : "The regular fee applies for all non-beginners. The Ant Tunnel is located deep inside in the dungeon that's placed at the center of the Victoria Island, where 24 Hr Mobile Store is. Would you like to go there for #b10,000 mesos#k?";

if (!npc.sendYesNo(prompt)) {
    npc.sendOk("This town also has a lot to offer. Find us if and when you feel the need to go to the Ant Tunnel Park.");
} else if (plr.getMesos() < cost) {
    npc.sendOk("It looks like you don't have enough mesos. Sorry but you won't be able to use this without it.");
} else {
    plr.gainMesos(-cost);
    plr.warp(105070001);
}
