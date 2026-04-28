var towns = [104000000, 102000000, 101000000, 100000000, 120000000];
var basePrices = [1000, 800, 1200, 1000, 1000];

npc.sendNext("What's up? I drive the Regular Cab. If you want to go from town to town safely and fast, then ride our cab. We'll glady take you to your destination with an affordable price.");

var beginner = plr.job() === 0;
var text = beginner
    ? "We have a special 90% discount for beginners. Choose your destination, for fees will change from place to place.#b"
    : "Choose your destination, for fees will change from place to place.#b";

for (var i = 0; i < towns.length; i++) {
    var price = beginner ? Math.floor(basePrices[i] * 0.10) : basePrices[i];
    text += "\r\n#L" + i + "##m" + towns[i] + "# (" + price.toLocaleString() + " mesos)#l";
}

npc.sendSelection(text);

var selection = npc.selection();
if (selection < 0 || selection >= towns.length) {
    npc.sendOk("There's a lot to see in this town, too. Come back and find us when you need to go to a different town.");
} else {
    var finalCost = beginner ? Math.floor(basePrices[selection] * 0.10) : basePrices[selection];
    if (!npc.sendYesNo("You don't have anything else to do here, huh? Do you really want to go to #b#m" + towns[selection] + "##k? It'll cost you #b" + finalCost.toLocaleString() + " mesos#k.")) {
        npc.sendOk("There's a lot to see in this town, too. Come back and find us when you need to go to a different town.");
    } else if (plr.getMesos() < finalCost) {
        npc.sendOk("You don't have enough mesos. Sorry to say this, but without them, you won't be able to ride the cab.");
    } else {
        plr.gainMesos(-finalCost);
        plr.warp(towns[selection]);
    }
}
