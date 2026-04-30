var stage1Map = 670010200;
var stage2Maps = [670010300, 670010301, 670010302];
var stage3Map = 670010400;

function showClear() {
    map.showEffect("quest/party/clear");
    map.playSound("Party1/Clear");
}

function showFail() {
    map.showEffect("quest/party/wrong_kor");
    map.playSound("Party1/Failed");
}

function nextStage2Map() {
    var hour = new Date().getHours();
    if (hour < 8) {
        return 670010302;
    }
    if (hour < 16) {
        return 670010300;
    }
    return 670010301;
}

function chooseStage2Answer() {
    var counts = [0, 0, 0];
    for (var i = 0; i < 5; i++) {
        counts[Math.floor(Math.random() * 3)]++;
    }
    return counts.join(",");
}

function readStage2Answer() {
    var raw = plr.getEventProperty("apqStage2Answer");
    if (raw == null || raw === "") {
        raw = chooseStage2Answer();
        plr.setEventProperty("apqStage2Answer", raw);
    }
    var parts = String(raw).split(",");
    return [Number(parts[0]), Number(parts[1]), Number(parts[2])];
}

function setStage3Answer() {
    var slots = [0, 0, 0, 0, 0, 0, 0, 0, 0];
    var chosen = 0;
    while (chosen < 5) {
        var idx = Math.floor(Math.random() * slots.length);
        if (slots[idx] === 0) {
            slots[idx] = 1;
            chosen++;
        }
    }
    plr.setEventProperty("apqStage3Answer", slots.join(""));
}

function stage1() {
    if (plr.getEventProperty("apqStage1Clear")) {
        if (!plr.isLeader()) {
            npc.sendOk("Ask your party leader to take everyone to the next stage.");
            return;
        }
        plr.warpEventMembersToPortal(nextStage2Map(), "st00");
        return;
    }

    if (!plr.isLeader()) {
        npc.sendOk("Talk to The Glimmer Man below. When your party is ready to move on, have your leader speak to me.");
        return;
    }

    if (!plr.getEventProperty("apqStage1Started")) {
        plr.setEventProperty("apqStage1Started", true);
        npc.sendOk("For the first challenge, head below and speak with my friend, The Glimmer Man. Choose carefully when you descend.");
        return;
    }

    var openIdx = Math.floor(Math.random() * 3);
    map.portalEnabled(openIdx === 0, "go00");
    map.portalEnabled(openIdx === 1, "go01");
    map.portalEnabled(openIdx === 2, "go02");
    npc.sendOk("Think it over and choose carefully. Only one of the three portals will take you where you need to go.");
}

function stage2() {
    if (!plr.getEventProperty("apqStage1Clear")) {
        npc.sendOk("Please clear the first stage before moving on.");
        return;
    }
    if (!plr.isLeader()) {
        npc.sendOk("Ask your party leader to confirm your answer.");
        return;
    }

    var answer = readStage2Answer();
    var a1 = map.playersInArea(1);
    var a2 = map.playersInArea(2);
    var a3 = map.playersInArea(3);

    if (a1 + a2 + a3 !== 5) {
        npc.sendOk("You'll need 5 party members hanging on the ropes.");
        return;
    }
    if (a1 > 2 || a2 > 2 || a3 > 2) {
        npc.sendOk("You can't have more than 2 people on the same rope set.");
        return;
    }

    var correct = 0;
    if (a1 === answer[0]) correct++;
    if (a2 === answer[1]) correct++;
    if (a3 === answer[2]) correct++;

    if (correct === 3) {
        plr.setEventProperty("apqStage2Clear", true);
        map.portalEnabled(true, "next00");
        showClear();
        plr.partyGiveExp(4000);
        npc.sendOk("That's the right answer. Here's the portal to the next stage.");
        return;
    }

    var tries = Number(plr.getEventProperty("apqStage2Try") || 0) + 1;
    plr.setEventProperty("apqStage2Try", tries);
    showFail();

    if (tries >= 7) {
        plr.setEventProperty("apqStage2Try", 0);
        plr.setEventProperty("apqStage2Answer", "");
        for (var i = 0; i < 20; i++) {
            plr.spawnMonster(9400538, -709, -1042);
        }
        npc.sendOk("You have failed too many times. A wave of monsters has been summoned.");
    } else if (tries === 6) {
        npc.sendOk("This is your 6th attempt. You only have one try remaining.");
    } else if (correct === 0) {
        npc.sendOk("This is attempt #" + tries + ". All the rope groups are wrong.");
    } else {
        npc.sendOk("This is attempt #" + tries + ". #b" + correct + "#k rope group(s) are correct.");
    }
}

function stage3() {
    if (!plr.getEventProperty("apqStage2Clear")) {
        npc.sendOk("Please clear the previous stage first.");
        return;
    }
    if (!plr.isLeader()) {
        npc.sendOk("Ask your party leader to confirm the switch combination.");
        return;
    }

    var answer = String(plr.getEventProperty("apqStage3Answer") || "");
    if (answer === "") {
        setStage3Answer();
        answer = String(plr.getEventProperty("apqStage3Answer") || "");
    }

    var current = "";
    var occupied = 0;
    for (var i = 0; i < 9; i++) {
        var count = map.playersInArea(i);
        if (count > 1) {
            npc.sendOk("Only one person can stand on each switch.");
            return;
        }
        occupied += count;
        current += String(count);
    }

    if (occupied !== 5) {
        npc.sendOk("You need exactly 5 party members standing on the switches.");
        return;
    }

    var strike = 0;
    for (var j = 0; j < 9; j++) {
        if (answer.charAt(j) === "1" && current.charAt(j) === "1") {
            strike++;
        }
    }

    if (current === answer) {
        plr.setEventProperty("apqStage3Clear", true);
        map.portalEnabled(true, "next00");
        showClear();
        plr.partyGiveExp(6000);
        npc.sendOk("That's the right answer. Here's the portal to the next stage.");
        return;
    }

    showFail();
    var randomPenalty = [9400540, 9400541, 9400542, 9400543][Math.floor(Math.random() * 4)];
    for (var k = 0; k < strike; k++) {
        plr.spawnMonster(randomPenalty, 1715 + (k * 150), -45);
    }
    if (strike > 0) {
        for (var m = 0; m < 7; m++) {
            plr.spawnMonster(9400539, 1120 + (m * 90), 192 + ((m % 3) * 20));
        }
    }
    npc.sendOk("That's not the right combination. You scored #b" + strike + " strike(s)#k.");
}

if (plr.mapID() === stage1Map) {
    stage1();
} else if (stage2Maps.indexOf(plr.mapID()) !== -1) {
    stage2();
} else if (plr.mapID() === stage3Map) {
    stage3();
} else {
    npc.sendOk("There's nothing for me to handle here.");
}
