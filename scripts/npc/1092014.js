var maps = [104000000, 102000000, 100000000, 103000000, 101000000];
var regularCosts = [1200, 1000, 1000, 1200, 1000];
var beginnerCosts = [120, 100, 100, 120, 100];

npc.sendNext("How's it going? I drive the Nautilus' Mid-Sized Taxi. If you want to go from town to town safely and fast, then ride our cab. We'll gladly take you to your destination for an affordable price.");

var beginner = plr.job() === 0;
var text = beginner
    ? "We have a special 90% discount for beginners. Choose your destination, for fees will change from place to place.#b"
    : "Choose your destination, for fees will change from place to place.#b";

for (var i = 0; i < maps.length; i++) {
    var price = beginner ? beginnerCosts[i] : regularCosts[i];
    text += "\r\n#L" + i + "##m" + maps[i] + "# (" + price.toLocaleString() + " mesos)#l";
}

npc.sendSelection(text);
var selection = npc.selection();

if (selection >= 0 && selection < maps.length) {
    var cost = beginner ? beginnerCosts[selection] : regularCosts[selection];
    if (!npc.sendYesNo("You don't have anything else to do here, huh? Do you really want to go to #b#m" + maps[selection] + "##k? It'll cost you #b" + cost.toLocaleString() + " mesos#k.")) {
        npc.sendOk("There's a lot to see in this town, too. Come back and find us when you need to go to a different town.");
    } else if (plr.getMesos() < cost) {
        npc.sendOk("You don't have enough mesos.");
    } else {
        plr.gainMesos(-cost);
        plr.warp(maps[selection]);
    }
}
