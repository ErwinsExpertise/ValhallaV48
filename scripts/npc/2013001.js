var mapId = plr.mapID();

function getStage(key) {
    var value = plr.getEventProperty(key);
    return value == null ? "" : String(value);
}

function setStage(key, value) {
    plr.setEventProperty(key, value);
}

function completeStage(key, exp, logName) {
    setStage(key, "1");
    map.properties()["clear"] = true;
    map.showEffect("quest/party/clear");
    map.playSound("Party1/Clear");
    plr.partyGiveExp(exp);
    if (logName) {
        plr.logEvent(logName);
    }
}

function leaderOnly() {
    if (!plr.isLeader()) {
        npc.sendOk("Please ask your party leader to speak with me.");
        return false;
    }
    return true;
}

if (mapId === 920010000) {
    if (!leaderOnly()) {
    } else if (getStage("prestage_clear") !== "1") {
        if (!plr.haveItem(4001063, 20)) {
            npc.sendOk("Please collect #b20 Cloud Pieces#k from the clouds around Wonky so I can restore my body.");
        } else {
            plr.gainItem(4001063, -20);
            completeStage("prestage_clear", 6000, "prestage cleared");
            npc.sendOk("Thank you for restoring my body. I will guide your party to the tower entrance now.");
            plr.warpEventMembersToPortal(920010100, "st00");
        }
    } else {
        npc.sendOk("I'll guide your party to the Center Tower again.");
        plr.warpEventMembersToPortal(920010100, "st00");
    }
} else if (mapId === 920010100) {
    if (!leaderOnly()) {
    } else {
        var statueState = map.reactorStateByName("minerva");
        var centerStage = getStage("stage0_clear");

        if (statueState === 6) {
            setStage("stage0_clear", "1");
            npc.sendOk("The statue has been restored. The #bGrass of Life#k can only be cultivated in the garden. I'll send the party there now.");
            plr.logEvent("statue restored; warping to garden");
            plr.warpEventMembersToPortal(920010800, "in00");
        } else if (statueState === 7) {
            npc.sendOk("Goddess Minerva has already been saved. Please speak with her.");
        } else if (centerStage === "1") {
            npc.sendOk("If you have found the #bGrass of Life#k, place it in the center of the statue.");
        } else {
            if (centerStage === "") {
                setStage("stage0_clear", "s");
            }
            npc.sendOk("Restore the six broken statue pieces in the Center Tower, then bring back the #bGrass of Life#k from the garden.");
        }
    }
} else if (mapId === 920010200) {
    if (!leaderOnly()) {
    } else if (getStage("stage1_clear") === "1") {
        npc.sendOk("There is nothing left to do here. Please continue to another room.");
    } else if (getStage("stage1_clear") === "") {
        setStage("stage1_clear", "s");
        npc.sendOk("Collect #b30 1st Small Pieces#k from the Pixies in this walkway so I can restore the first statue piece.");
    } else if (!plr.haveItem(4001050, 30)) {
        npc.sendOk("Have you collected #b30 1st Small Pieces#k yet?");
    } else if (!plr.canHold(4001044, 1)) {
        npc.sendOk("Please make space in your inventory for the statue piece.");
    } else {
        plr.gainItem(4001050, -30);
        plr.gainItem(4001044, 1);
        completeStage("stage1_clear", 7500, "stage 1 cleared");
        npc.sendOk("Excellent work. I have restored #bStatue of Goddess: 1st Piece#k.");
    }
} else if (mapId === 920010300) {
    if (!leaderOnly()) {
    } else if (getStage("stage2_clear") === "1") {
        npc.sendOk("There is nothing left to do here. Please continue to another room.");
    } else if (getStage("stage2_clear") === "") {
        setStage("stage2_clear", "s");
        npc.sendOk("Search the storage room and recover #bStatue of Goddess: 2nd Piece#k.");
    } else if (!plr.haveItem(4001045, 1)) {
        npc.sendOk("Have you recovered #bStatue of Goddess: 2nd Piece#k yet?");
    } else {
        completeStage("stage2_clear", 7500, "stage 2 cleared");
        npc.sendOk("You found the second statue piece. Return to the Center Tower when you are ready.");
    }
} else if (mapId === 920010400) {
    var todayNames = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"];
    var today = new Date().getDay() + 1;
    if (!leaderOnly()) {
    } else if (getStage("stage3_clear") === "1") {
        npc.sendOk("There is nothing left to do here. Please continue to another room.");
    } else if (map.reactorStateByName("music") === 0) {
        npc.sendOk("This is Minerva's music hall. Play the song that matches today's mood, then talk to me again.");
    } else if (map.reactorStateByName("music") !== today) {
        npc.sendOk("That is not the correct song for #b" + todayNames[today - 1] + "#k.");
    } else {
        map.hitReactorByName("stone3");
        completeStage("stage3_clear", 7500, "stage 3 cleared");
        npc.sendOk("Yes, that's the song Minerva loved to hear. Something appeared on the altar.");
    }
} else if (mapId === 920010500) {
    if (!leaderOnly()) {
    } else if (getStage("stage4_clear") === "") {
        setStage("stage4_clear", "s");
        setStage("try", "0");
        npc.sendOk("Five party members must stand on the three platforms in the correct combination. Talk to me after setting the positions.");
    } else if (getStage("stage4_clear") === "1") {
        npc.sendOk("There is nothing left to do here. Please continue to another room.");
    } else {
        var ans1 = getStage("ans1");
        var ans2 = getStage("ans2");
        var ans3 = getStage("ans3");
        if (ans1 === "" || ans2 === "" || ans3 === "") {
            var a = Math.floor(Math.random() * 6);
            var b = Math.floor(Math.random() * (6 - a));
            var c = 5 - a - b;
            var perms = [
                [a, b, c], [a, c, b], [b, a, c],
                [b, c, a], [c, a, b], [c, b, a]
            ];
            var pick = perms[Math.floor(Math.random() * perms.length)];
            setStage("ans1", String(pick[0]));
            setStage("ans2", String(pick[1]));
            setStage("ans3", String(pick[2]));
            ans1 = String(pick[0]);
            ans2 = String(pick[1]);
            ans3 = String(pick[2]);
        }

        var area1 = map.playersInArea(1);
        var area2 = map.playersInArea(2);
        var area3 = map.playersInArea(3);

        if (area1 + area2 + area3 !== 5) {
            npc.sendOk("You need exactly five party members standing on the platforms.");
        } else {
            var correct = 0;
            if (area1 === parseInt(ans1, 10)) correct++;
            if (area2 === parseInt(ans2, 10)) correct++;
            if (area3 === parseInt(ans3, 10)) correct++;

            if (correct < 3) {
                var tries = parseInt(getStage("try") || "0", 10) + 1;
                setStage("try", String(tries));
                if (tries === 6) {
                    npc.sendOk("This is your 6th attempt. You only have one try left, so be careful.");
                } else if (tries >= 7) {
                    setStage("try", "0");
                    setStage("stage4_clear", "");
                    plr.logEvent("stage 4 failed after seven tries");
                    npc.sendOk("You failed the platform test and will be sent back to the Center Tower.");
                    plr.warpEventMembersToPortal(920010100, "in03");
                } else if (correct === 0) {
                    npc.sendOk("That was attempt #" + tries + ". None of the platform weights are correct yet.");
                } else {
                    npc.sendOk("That was attempt #" + tries + ". #b" + correct + "#k platform weight(s) are correct.");
                }
            } else {
                map.hitReactorByName("stone4");
                completeStage("stage4_clear", 7500, "stage 4 cleared");
                npc.sendOk("Correct. Something appeared on the altar.");
            }
        }
    }
} else if (mapId === 920010600) {
    if (!leaderOnly()) {
    } else if (getStage("stage5_clear") === "1") {
        npc.sendOk("There is nothing left to do here. Please continue to another room.");
    } else if (getStage("stage5_clear") === "") {
        setStage("stage5_clear", "s");
        npc.sendOk("Search the lounge and bedrooms for #b40 5th Small Pieces#k so I can restore the fifth statue piece.");
    } else if (!plr.haveItem(4001052, 40)) {
        npc.sendOk("Have you collected #b40 5th Small Pieces#k yet?");
    } else if (!plr.canHold(4001048, 1)) {
        npc.sendOk("Please make space in your inventory for the statue piece.");
    } else {
        plr.gainItem(4001052, -40);
        plr.gainItem(4001048, 1);
        completeStage("stage5_clear", 7500, "stage 5 cleared");
        npc.sendOk("Excellent work. I have restored #bStatue of Goddess: 5th Piece#k.");
    }
} else if (mapId === 920010700) {
    if (!leaderOnly()) {
    } else if (getStage("stage6_clear") === "1") {
        npc.sendOk("There is nothing left to do here. Please continue to another room.");
    } else if (getStage("stage6_clear") === "") {
        setStage("stage6_clear", "s");
        npc.sendOk("At the top of this tower are five levers. Exactly two of them must be switched on.");
    } else {
        var ansA = getStage("stage6_ans1");
        var ansB = getStage("stage6_ans2");
        var wrong1 = getStage("stage6_wans1");
        var wrong2 = getStage("stage6_wans2");
        var wrong3 = getStage("stage6_wans3");
        var ok = map.reactorStateByName(ansA) === 1 && map.reactorStateByName(ansB) === 1 &&
            map.reactorStateByName(wrong1) === 0 && map.reactorStateByName(wrong2) === 0 &&
            map.reactorStateByName(wrong3) === 0;

        if (!ok) {
            npc.sendOk("Nothing happened. It looks like the wrong levers were switched.");
        } else {
            map.hitReactorByName("stone6");
            completeStage("stage6_clear", 7500, "stage 6 cleared");
            npc.sendOk("You moved the correct levers. Something appeared on the altar.");
        }
    }
} else if (mapId === 920010800) {
    if (!leaderOnly()) {
        npc.sendOk("Please follow your party leader.");
    } else if (plr.haveItem(4001055, 1)) {
        npc.sendOk("You found the #bGrass of Life#k. Return to the Center Tower and place it in the statue.");
    } else {
        npc.sendOk("Cultivate the #bGrass of Life#k here. Be careful, because Papa Pixie may appear.");
    }
} else if (mapId === 920010900) {
    npc.sendOk("This is the Room of the Guilty. There is no statue piece here, but something may still be hidden.");
} else if (mapId === 920011000) {
    npc.sendOk("This is the Room of Darkness. There is no statue piece here, but you may find diary pages from the Goddess.");
} else if (mapId === 920011200) {
    npc.sendOk("See you next time.");
    plr.warp(200080101);
} else if (mapId === 920011100) {
    npc.sendOk("Break the treasure boxes quickly before the Goddess's blessing fades.");
} else {
    npc.sendOk("Please continue onward.");
}
