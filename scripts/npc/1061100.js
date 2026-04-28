npc.sendNext("Welcome. We're the Sleepywood Hotel. Our hotel works hard to serve you the best at all times. If you are tired and worn out from hunting, how about a relaxing stay at our hotel?");

npc.sendSelection("We offer two kinds of rooms for our service. Please choose the one of your liking.\r\n#b#L0#Regular sauna (499 mesos per use)#l\r\n#L1#VIP sauna (999mesos per use)#l");

var selection = npc.selection();
if (selection !== 0 && selection !== 1) {
    npc.sendOk("We offer other kinds of services, too, so please think carefully and then make your decision.");
} else {
    var regular = selection === 0;
    var cost = regular ? 499 : 999;
    var destination = regular ? 105040401 : 105040402;
    var prompt = regular
        ? "You have chosen the regular sauna. Your HP and MP will recover fast and you can even purchase some items there. Are you sure you want to go in?"
        : "You've chosen the VIP sauna. Your HP and MP will recover even faster than that of the regular sauna and you can even find a special item in there. Are you sure you want to go in?";

    if (!npc.sendYesNo(prompt)) {
        npc.sendOk("We offer other kinds of services, too, so please think carefully and then make your decision.");
    } else if (plr.getMesos() < cost) {
        npc.sendOk("I'm sorry. It looks like you don't have enough mesos. It will cost you at least " + cost + "mesos to stay at our hotel.");
    } else {
        plr.warp(destination);
        plr.gainMesos(-cost);
    }
}
