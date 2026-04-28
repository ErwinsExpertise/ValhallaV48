var mapId = plr.mapID();

if (mapId !== 800000000) {
    var cost = plr.job() === 0 ? 300 : 3000;
    npc.sendNext("If you're tired of the monotonous daily life, how about getting out for a change? there's nothing quite like soaking up a new culture, learning something new by the minute! It's time for you to get out and travel. We, at the Maple Travel Agency recommend you going on a #bWorld Tour#k! Are you worried about the travel expense? You shouldn't be! We, the #bMaple Travel Agency#k, have carefully come up with a plan to let you travel for ONLY #b" + cost.toLocaleString() + " mesos#k!");
    npc.sendSelection("We currently offer this destination for adventurous travelers: #bMushroom Shrine of Japan#k. I will be there to serve as your travel guide, and more destinations may be added over time. Would you like to head to Mushroom Shrine now?\r\n#L0##bYes, take me to Mushroom Shrine (Japan).#k#l");
    if (npc.selection() === 0 && npc.sendYesNo("Would you like to travel to #bMushroom Shrine of Japan#k? If you want to experience the essence of Japan, there is nothing quite like visiting the Shrine, a cultural melting pot with a rich history.")) {
        if (plr.getMesos() < cost) {
            npc.sendOk("Please check and make sure you have enough mesos for the trip.");
        } else {
            plr.gainMesos(-cost);
            plr.saveLocation("WORLDTOUR");
            plr.warp(800000000);
        }
    }
} else {
    var returnMap = plr.getSavedLocation("WORLDTOUR");
    if (returnMap < 0) {
        returnMap = 100000000;
    }
    npc.sendSelection("How's the traveling? Are you enjoying it?\r\n#L0##bYes, I'm done with travelling. Can I go back to #m" + returnMap + "#?#k#l\r\n#L1##bNo, I'd like to continue exploring this place.#k#l");
    if (npc.selection() === 0) {
        plr.warp(returnMap);
    } else {
        npc.sendOk("OK. If you ever change your mind, please let me know.");
    }
}
