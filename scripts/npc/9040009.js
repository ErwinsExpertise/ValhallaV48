function gatekeeperCombo(names, count) {
    var picked = [];
    var available = names.slice();
    while (picked.length < count && available.length > 0) {
        var idx = Math.floor(Math.random() * available.length);
        picked.push(available[idx]);
        available.splice(idx, 1);
    }
    return picked;
}

var leader = plr.getEventProperty("leader");
if (leader == null) {
    plr.warp(990001100);
} else if (leader !== plr.name()) {
    npc.sendOk("I need the leader of your group to speak with me.");
} else {
    var props = map.properties();
    var phase = parseInt(props["stage1phase"] || "1", 10);
    var state = props["stage1status"] || "waiting";

    if (map.hitReactorByName("statuegate")) {
        npc.sendOk("Proceed.");
    } else if (state === "waiting") {
        npc.sendOk(phase === 1
            ? "In this challenge, I shall show a pattern on the statues around me. When I give the word, repeat the pattern to me to proceed."
            : "I shall now present a more difficult puzzle for you. Good luck.");

        var names = [];
        var reactors = map.reactorNames();
        for (var i = 0; i < reactors.length; i++) {
            if (reactors[i] !== "statuegate") names.push(reactors[i]);
        }
        var combo = gatekeeperCombo(names, phase + 3);
        props["stage1phase"] = phase;
        props["stage1status"] = "display";
        props["stage1combo"] = combo.join(",");
        props["stage1guess"] = "";
        map.showEffect("quest/party/wrong_kor");
        map.playSound("Party1/Failed");
        map.revealReactorsByName(combo, 5000, 3500);
    } else if (state === "active") {
        var comboList = (props["stage1combo"] || "").split(",").filter(function(v) { return v.length > 0; });
        var guessList = (props["stage1guess"] || "").split(",").filter(function(v) { return v.length > 0; });
        if (comboList.join(",") === guessList.join(",")) {
            if (phase >= 3) {
                map.hitReactorByName("statuegate");
                map.showEffect("quest/party/clear");
                map.playSound("Party1/Clear");
                props["stage1clear"] = true;
                plr.gainGuildPoints(15);
                npc.sendOk("Excellent work. Please proceed to the next stage.");
            } else {
                props["stage1phase"] = phase + 1;
                props["stage1status"] = "waiting";
                props["stage1guess"] = "";
                npc.sendOk("Very good. You still have more to complete, however. Talk to me again when you're ready.");
            }
        } else {
            props["stage1phase"] = 1;
            props["stage1status"] = "waiting";
            props["stage1guess"] = "";
            npc.sendOk("You have failed this test.");
        }
    } else {
        npc.sendOk("Please wait while the combination is revealed.");
    }
}
