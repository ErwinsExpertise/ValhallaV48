var SNOW_ITEM = 4031875;
var MAX_SNOW = 50000;
var STAGE_STEP = 5000;
var BLOWER_IDS = [9400714, 9400715, 9400716, 9400717, 9400718, 9400719, 9400720, 9400721, 9400722, 9400723, 9400724];
var BOSS_IDS = [9400707, 9400708, 9400709, 9400710];

function props() {
    return map.properties();
}

function currentSnow() {
    var value = props()["wxmasCount"];
    return value ? parseInt(String(value), 10) : 0;
}

function setSnow(value) {
    props()["wxmasCount"] = String(value);
}

function normalizeAfterBoss() {
    if (props()["wxmasBoss"] === "1" && !bossActive()) {
        props()["wxmasBoss"] = "0";
        setSnow(0);
        clearBlowers();
        map.spawnMonster(9400714, 1450, 140);
        map.message("The snow machine calms down and starts over from empty.");
    }
}

function bossActive() {
    for (var i = 0; i < BOSS_IDS.length; i++) {
        if (map.mobCountByID(BOSS_IDS[i]) > 0) {
            return true;
        }
    }
    return false;
}

function clearBlowers() {
    for (var i = 0; i < BLOWER_IDS.length; i++) {
        map.removeMobsByID(BLOWER_IDS[i]);
    }
}

function updateMachine(total) {
    clearBlowers();

    if (total >= MAX_SNOW) {
        if (!bossActive()) {
            props()["wxmasBoss"] = "1";
            map.spawnMonster(9400707, 1250, -422);
            map.spawnMonster(9400708, 710, 60);
            map.message("The snow machine is overfilled! Something has gone wrong in the Extra Frosty Snow Zone.");
        }
        return;
    }

    var stage = Math.floor(total / STAGE_STEP);
    if (stage < 0) {
        stage = 0;
    }
    if (stage >= BLOWER_IDS.length) {
        stage = BLOWER_IDS.length - 1;
    }
    map.spawnMonster(BLOWER_IDS[stage], 1450, 140);
}

var amountOnHand = plr.itemCount(SNOW_ITEM);
normalizeAfterBoss();
var totalSnow = currentSnow();
var remaining = MAX_SNOW - totalSnow;

var action = npc.askMenu(
    "We all want a White Christmas. What would you like to do?#b",
    "How does the snow machine work?",
    "Hand over Snow Powder",
    "Check the current progress"
);

if (action === 0) {
    npc.sendNext("The Extra Frosty Snow Zone depends on #b#t" + SNOW_ITEM + "##k to keep the machine running. Bring any Snow Powder you find to me and I will feed it into the machine.");
    npc.sendOk("Once the machine is filled, the zone will react on its own. Until then, every handful helps.");
} else if (action === 1) {
    if (amountOnHand <= 0) {
        npc.sendOk("You don't have any #b#t" + SNOW_ITEM + "##k with you right now.");
    } else if (remaining <= 0 || bossActive()) {
        npc.sendOk("The snow machine is already at its limit. Stand back and be careful.");
    } else {
        var maxGive = amountOnHand < remaining ? amountOnHand : remaining;
        var give = npc.sendNumber("How much #b#t" + SNOW_ITEM + "##k would you like to contribute?\r\n\r\nCurrent total: #b" + totalSnow + " / " + MAX_SNOW + "#k", maxGive, 1, maxGive);
        if (give < 1 || give > maxGive) {
            npc.sendOk("Come back when you decide how much you want to hand over.");
        } else if (!plr.gainItem(SNOW_ITEM, -give)) {
            npc.sendOk("I couldn't take the Snow Powder from you. Please try again.");
        } else {
            totalSnow += give;
            setSnow(totalSnow);
            updateMachine(totalSnow);

            if (totalSnow >= MAX_SNOW) {
                npc.sendOk("That's enough! The snow machine is overflowing with power now. Stay alert.");
            } else {
                npc.sendOk("Thank you. The snow machine is now at #b" + totalSnow + " / " + MAX_SNOW + "#k.");
            }
        }
    }
} else {
    if (bossActive()) {
        npc.sendOk("The snow machine has already gone out of control. Please help calm the zone down before adding any more powder.");
    } else {
        npc.sendOk("The snow machine currently holds #b" + totalSnow + " / " + MAX_SNOW + "#k Snow Powder.");
    }
}
