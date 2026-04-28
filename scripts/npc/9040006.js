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
    var correct = 0;
    for (var i = a.length - 1; i >= 0; i--) {
        if (a[i] === g[i]) {
            correct++;
            a.splice(i, 1);
            g.splice(i, 1);
        }
    }
    var answerCount = [0, 0, 0, 0];
    var guessCount = [0, 0, 0, 0];
    for (var j = 0; j < a.length; j++) answerCount[a[j]]++;
    for (var k = 0; k < g.length; k++) guessCount[g[k]]++;
    var incorrect = 0;
    var unknown = 0;
    for (var n = 0; n < 4; n++) {
        incorrect += Math.min(answerCount[n], guessCount[n]);
        if (guessCount[n] > answerCount[n]) unknown += guessCount[n] - answerCount[n];
    }
    return [correct, incorrect, unknown];
}

function getGroundItems() {
    var itemInArea = [-1, -1, -1, -1];
    for (var i = 0; i < 4; i++) {
        var items = map.groundItemsInArea(i);
        if (items.length !== 1) return null;
        var id = items[0];
        if (id < 4001027 || id > 4001030) return null;
        itemInArea[i] = id - 4001027;
    }
    return itemInArea[0] * 1000 + itemInArea[1] * 100 + itemInArea[2] * 10 + itemInArea[3];
}

var leader = plr.getEventProperty("leader");
if (leader == null) {
    plr.warp(990001100);
} else if (leader !== plr.name()) {
    npc.sendOk("Please have your leader speak to me.");
} else if (map.hitReactorByName("watergate")) {
    npc.sendOk("You may proceed.");
} else {
    var currentCombo = plr.getEventProperty("stage3combo");
    if (currentCombo == null || currentCombo === "reset") {
        var combo = makeCombo();
        plr.setEventProperty("stage3combo", "" + (combo[0] * 1000 + combo[1] * 100 + combo[2] * 10 + combo[3]));
        plr.setEventProperty("stage3attempt", "1");
        npc.sendOk("This fountain guards the secret passage to the throne room. Offer items in front of the vassals to proceed. They will tell you whether your offerings are accepted, and if not, which vassals are displeased. You have seven attempts. Good luck.");
    } else {
        var attempt = parseInt(plr.getEventProperty("stage3attempt"), 10);
        var comboVal = parseInt(currentCombo, 10);
        var guess = getGroundItems();
        if (guess == null) {
            npc.sendOk("Please make sure your attempt is properly set in front of the vassals and talk to me again.");
        } else if (comboVal === guess) {
            map.hitReactorByName("watergate");
            map.showEffect("quest/party/clear");
            map.playSound("Party1/Clear");
            plr.setEventProperty("stage3clear", true);
            plr.gainGuildPoints(25);
            npc.sendOk("You may proceed.");
        } else if (attempt < 7) {
            var results = compare(parsePattern(comboVal), parsePattern(guess));
            var text = "";
            if (results[0] > 0) text += results[0] + (results[0] === 1 ? " vassal is pleased with their offering.\r\n" : " vassals are pleased with their offerings.\r\n");
            if (results[1] > 0) text += results[1] + (results[1] === 1 ? " vassal received an incorrect offering.\r\n" : " vassals received incorrect offerings.\r\n");
            if (results[2] > 0) text += results[2] + (results[2] === 1 ? " vassal received an unknown offering.\r\n" : " vassals received unknown offerings.\r\n");
            text += "This is your " + attempt + (attempt === 1 ? "st" : attempt === 2 ? "nd" : attempt === 3 ? "rd" : "th") + " attempt.";
            plr.setEventProperty("stage3attempt", "" + (attempt + 1));
            plr.spawnMonster(9300036, -350, 150);
            plr.spawnMonster(9300037, 400, 150);
            npc.sendOk(text);
        } else {
            plr.setEventProperty("stage3combo", "reset");
            for (var i = 0; i < 5; i++) {
                plr.spawnMonster(9300036, -300 + (i * 120), 150);
                plr.spawnMonster(9300037, -240 + (i * 120), 150);
            }
            npc.sendOk("You have failed the test. Please compose yourselves and try again later.");
        }
    }
}
