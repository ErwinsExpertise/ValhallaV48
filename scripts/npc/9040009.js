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
    npc.sendOk("I need the registered leader to speak with me.");
} else {
    var props = map.properties();
    var phase = parseInt(props["stage1phase"] || "1", 10);
    var status = props["stage1status"] || "waiting";

    if (map.reactorStateByName("statuegate") === 1 || props["stage1clear"] === true || props["stage1clear"] === "true") {
        npc.sendOk("You have already passed my trial. Proceed.");
    } else if (status === "waiting") {
        var names = map.reactorNames().filter(function(name) { return name !== "statuegate"; });
        var combo = gatekeeperCombo(names, phase + 3);
        props["stage1phase"] = phase;
        props["stage1status"] = "display";
        props["stage1combo"] = combo.join(",");
        props["stage1guess"] = "";
        props["stage1displaycount"] = 0;

        npc.sendOk(phase === 1
            ? "Watch the statues closely and remember the order in which they shine. Strike them in that same order, then return to me."
            : "You succeeded once, but I require a harder answer now. Watch carefully.");

        map.revealReactorsByName(combo, 5000, 3500);
    } else if (status === "display") {
        npc.sendOk("Wait until all of the statues finish revealing the pattern.");
    } else {
        var comboList = (props["stage1combo"] || "").split(",").filter(function(v) { return v.length > 0; });
        var guessList = (props["stage1guess"] || "").split(",").filter(function(v) { return v.length > 0; });
        if (guessList.length < comboList.length) {
            npc.sendOk("You have not struck enough statues yet.");
        } else if (comboList.join(",") === guessList.join(",")) {
            if (phase >= 3) {
                map.setReactorStateByName("statuegate", 1);
                props["stage1clear"] = true;
                props["stage1status"] = "done";
                props["statuegateopen"] = "yes";
                map.showEffect("quest/party/clear");
                map.playSound("Party1/Clear");
                plr.logEvent("gpq stage1 clear");
                npc.sendOk("Brilliant. I will open the way to the castle.");
            } else {
                props["stage1phase"] = phase + 1;
                props["stage1status"] = "waiting";
                props["stage1guess"] = "";
                props["stage1displaycount"] = 0;
                npc.sendOk("Correct. Speak to me again when you are ready for the next sequence.");
            }
        } else {
            props["stage1phase"] = 1;
            props["stage1status"] = "waiting";
            props["stage1guess"] = "";
            props["stage1displaycount"] = 0;
            npc.sendOk("Incorrect. You must begin again from the first test.");
        }
    }
}
