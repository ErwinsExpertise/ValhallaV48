var fromMaps = [211000000, 220000000, 240000000];
var toMaps = [211040200, 220050300, 240030000];
var costs = [45000, 25000, 55000];
var location = 0;

for (var i = 0; i < fromMaps.length; i++) {
    if (plr.mapID() === fromMaps[i]) {
        location = i;
        break;
    }
}

npc.sendNext("Hello there! This Bullet Taxi will take you to the danger zone from #m" + plr.mapID() + "# to #b#m" + toMaps[location] + "##k on this Ossyria continent! The transportation fee of #b" + costs[location].toLocaleString() + " mesos#k may seem expensive, but it's worth it when you want to travel through danger zones quickly!");

if (!npc.sendYesNo("Do you want to pay mesos and travel to #b#m" + toMaps[location] + "##k?")) {
    npc.sendOk("Hmm... think it over. This taxi is worth the service! You won't regret it!");
} else if (plr.getMesos() < costs[location]) {
    npc.sendOk("You don't seem to have enough mesos. I am terribly sorry, but I cannot help you unless you pay up. Bring in the mesos by hunting more and come back when you have enough.");
} else {
    plr.warp(toMaps[location]);
    plr.gainMesos(-costs[location]);
}
