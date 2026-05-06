var HORN_ITEM = 4031063;
var STATE_KEY = 8845;
var COOLDOWN_KEY = 8846;
var COOLDOWN_MS = 24 * 60 * 60 * 1000;
var rewards = [1012011, 1012012, 1012013, 1012014, 1012015, 1012016, 1012017, 1012018, 1012019, 1012020, 3992024, 3992025, 3992026, 2020012, 2020013];

function cooldownRemaining() {
    var last = parseInt(plr.questData(COOLDOWN_KEY) || "0", 10);
    if (!last) {
        return 0;
    }
    var remaining = COOLDOWN_MS - (Date.now() - last);
    return remaining > 0 ? remaining : 0;
}

function rewardHorn() {
    var reward = rewards[Math.floor(Math.random() * rewards.length)];
    if (!plr.canHold(reward, 1)) {
        npc.sendOk("Please make sure you have at least one free slot for the reward first.");
        return;
    }
    if (!plr.gainItem(HORN_ITEM, -1)) {
        npc.sendOk("I couldn't take the horn from you. Please try again.");
        return;
    }
    plr.gainItem(reward, 1);
    plr.setQuestData(STATE_KEY, "done");
    plr.setQuestData(COOLDOWN_KEY, String(Date.now()));
    npc.sendOk("You found my horn! Thank you so much. Please take this as a reward, and I promise I'll try not to lose it again.");
}

if (plr.level() < 10) {
    npc.sendOk("I need someone a little stronger to help me search for my horn. Please come back after training some more.");
} else {
    var state = plr.questData(STATE_KEY);
    var remaining = cooldownRemaining();

    if (state === "search") {
        if (plr.haveItem(HORN_ITEM, 1)) {
            rewardHorn();
        } else {
            npc.sendOk("My horn is still missing. If you find #b#t" + HORN_ITEM + "##k, please bring it back to me right away.");
        }
    } else if (remaining > 0) {
        var hours = Math.ceil(remaining / (60 * 60 * 1000));
        npc.sendOk("You already helped me today. Please come back in about #b" + hours + " hour(s)#k if I lose my horn again.");
    } else {
        var choice = npc.askMenu(
            "My horn... where is my horn?! What would you like to ask?#b",
            "How can I help?",
            "Where should I look?"
        );

        if (choice === 0) {
            if (npc.sendYesNo("You really want to help me find it? If you bring my lost horn back, I'll make sure it is worth your while.")) {
                plr.setQuestData(STATE_KEY, "search");
                npc.sendOk("Thank you! I lost it somewhere during all the holiday commotion. Please bring it back if you find it.");
            } else {
                npc.sendOk("Oh... alright. Please tell me if you change your mind.");
            }
        } else {
            npc.sendOk("I've heard all kinds of things, but I am sure somebody in the Maple world is carrying it around by now. If you find #b#t" + HORN_ITEM + "##k, please bring it back to me.");
        }
    }
}
