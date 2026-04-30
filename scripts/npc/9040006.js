function makeCombo() {
    var combo = [];
    while (combo.length < 4) {
        combo.push(Math.floor(Math.random() * 4));
    }
    return combo;
}

function parsePattern(value) {
    var items = [];
    var temp = parseInt(value, 10);
    for (var i = 0; i < 4; i++) {
        items.push(Math.floor(temp / Math.pow(10, 3 - i)));
        temp = temp % Math.pow(10, 3 - i);
    }
    return items;
}

function compare(answer, guess) {
    var a = answer.slice();
    var g = guess.slice();
    var strike = 0;
    for (var i = a.length - 1; i >= 0; i--) {
        if (a[i] === g[i]) {
            strike++;
            a.splice(i, 1);
            g.splice(i, 1);
        }
    }
    var answerCount = [0, 0, 0, 0];
    var guessCount = [0, 0, 0, 0];
    for (var j = 0; j < a.length; j++) answerCount[a[j]]++;
    for (var k = 0; k < g.length; k++) guessCount[g[k]]++;
    var ball = 0;
    var unknown = 0;
    for (var n = 0; n < 4; n++) {
        ball += Math.min(answerCount[n], guessCount[n]);
        if (guessCount[n] > answerCount[n]) unknown += guessCount[n] - answerCount[n];
    }
    return [strike, ball, unknown];
}

function getGroundItems() {
    var placed = [-1, -1, -1, -1];
    for (var i = 1; i <= 4; i++) {
        var items = map.groundItemsInArea(i);
        if (items.length !== 1) return null;
        var id = items[0];
        if (id < 4001027 || id > 4001030) return null;
        placed[i - 1] = id - 4001027;
    }
    return placed[0] * 1000 + placed[1] * 100 + placed[2] * 10 + placed[3];
}

var leader = plr.getEventProperty("leader");
if (leader == null) {
    plr.warp(990001100);
} else if (leader !== plr.name()) {
    npc.sendOk("Please have the registered leader speak to me.");
} else if (plr.getEventProperty("watergateopen") === "yes" || map.reactorStateByName("watergate") === 1) {
    npc.sendOk("The secret passage is already open.");
} else {
    var currentCombo = plr.getEventProperty("stage3combo");
    if (currentCombo == null || currentCombo === "reset") {
        var combo = makeCombo();
        plr.setEventProperty("stage3combo", "" + (combo[0] * 1000 + combo[1] * 100 + combo[2] * 10 + combo[3]));
        plr.setEventProperty("stage3attempt", "0");
        npc.sendOk("Offer the four royal gifts to the vassals. They will tell you how close you are, but you only have 7 attempts.");
    } else {
        var attempt = parseInt(plr.getEventProperty("stage3attempt") || "0", 10) + 1;
        var comboVal = parseInt(currentCombo, 10);
        var guess = getGroundItems();
        if (guess == null) {
            npc.sendOk("Place exactly one valid offering in front of each vassal and speak to me again.");
        } else if (comboVal === guess) {
            map.setReactorStateByName("watergate", 1);
            plr.setEventProperty("stage3clear", true);
            plr.setEventProperty("watergateopen", "yes");
            map.showEffect("quest/party/clear");
            map.playSound("Party1/Clear");
            plr.logEvent("gpq stage3 clear attempts=" + attempt);
            npc.sendOk("The vassals approve. The secret passage is now open.");
        } else {
            plr.setEventProperty("stage3attempt", "" + attempt);
            var results = compare(parsePattern(comboVal), parsePattern(guess));
            if (attempt >= 7) {
                plr.setEventProperty("stage3combo", "reset");
                plr.setEventProperty("stage3attempt", "0");
                for (var i = 0; i < 4; i++) {
                    plr.spawnMonster(9300036, -184 + (i * 136), 140);
                    plr.spawnMonster(9300037, -116 + (i * 136), 140);
                }
                npc.sendOk("You have enraged the vassals. Defeat the summoned guardians and try again.");
            } else {
                plr.spawnMonster(9300036, -476, 140);
                plr.spawnMonster(9300037, 552, 140);
                var text = "This is attempt #" + attempt + ".\r\n";
                if (results[0] > 0) text += results[0] + " vassal(s) received the correct offering.\r\n";
                if (results[1] > 0) text += results[1] + " vassal(s) received the wrong offering.\r\n";
                if (results[2] > 0) text += results[2] + " vassal(s) do not recognize their offering.\r\n";
                if (attempt === 6) text += "You only have one chance left.";
                npc.sendOk(text);
            }
        }
    }
}
