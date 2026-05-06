var DAILY_TIME_KEY = 8200;
var DAILY_COUNT_KEY = 8201;
var VERSAL_FULL_KEY = 8203;
var MAPLEMAS_CHOICE_KEY = 4997;
var VERSALMAS_CHOICE_KEY = 4998;
var CHEER_ITEM = 4031879;
var DAY_MS = 24 * 60 * 60 * 1000;

function getPresentItem() {
    var level = plr.level();
    if (level <= 20) {
        return 4031443;
    }
    if (level <= 30) {
        return 4031440;
    }
    if (level <= 40) {
        return 4031441;
    }
    if (level <= 60) {
        return 4031439;
    }
    return 4031442;
}

function resetDailyIfNeeded() {
    var last = parseInt(plr.questData(DAILY_TIME_KEY) || "0", 10);
    if (last && Date.now() - last >= DAY_MS) {
        plr.setQuestData(DAILY_TIME_KEY, "");
        plr.setQuestData(DAILY_COUNT_KEY, "0");
    }
}

function getDailyCount() {
    var raw = plr.questData(DAILY_COUNT_KEY);
    return raw ? parseInt(raw, 10) : 0;
}

function buildChoices() {
    return [
        { count: 10, cheer: 1, rewards: [[2000002, 25]] },
        { count: 20, cheer: 2, rewards: [[4003000, 10], [4003001, 10], [4011001, 2]] },
        { count: 30, cheer: 3, rewards: [
            [[4020007, 10]], [[4020008, 10]], [[2020014, 5]], [[2002022, 5]]
        ][Math.floor(Math.random() * 4)] },
        { count: 40, cheer: 4, rewards: [[2022272, 3]] },
        { count: 60, cheer: 6, rewards: [
            [[2020014, 5], [2002022, 5]], [[2022275, 5]], [[2022274, 3]], [[2022273, 3]], [[2022277, 3]], [[1082228, 1]]
        ][Math.floor(Math.random() * 6)] },
        { count: 90, cheer: 9, rewards: [
            [[2022275, 10]], [[2022274, 5]], [[2022273, 3]], [[2022277, 3]], [[2040805, 1]], [[1082228, 1]], [[1442061, 1]]
        ][Math.floor(Math.random() * 7)] },
        { count: 100, cheer: 10, rewards: [[1442061, 1]] }
    ];
}

function describeRewards(bundle) {
    var text = "For #b" + bundle.count + "#k presents, I can offer #b#t" + CHEER_ITEM + "##k x " + bundle.cheer;
    for (var i = 0; i < bundle.rewards.length; i++) {
        text += "\r\n#b#t" + bundle.rewards[i][0] + "##k x " + bundle.rewards[i][1];
    }
    return text;
}

resetDailyIfNeeded();

if (plr.level() < 10) {
    npc.sendOk("You need to be at least level 10 before you can help me spread Versalmas cheer.");
} else if (plr.questData(MAPLEMAS_CHOICE_KEY) === "2") {
    npc.sendOk("You decided to celebrate #bMaplemas#k this season, so Maple Claws is the one who should receive your presents.");
} else {
    var presentItem = getPresentItem();
    var dailyCount = getDailyCount();
    var choseVersalmas = plr.questData(VERSALMAS_CHOICE_KEY) === "2";
    var choices = buildChoices();

    var action = npc.askMenu(
        "O-Hoy! Happy Versalmas to you! What would you like to do?#b",
        "Tell me about Versalmas.",
        "Which present colors should I collect?",
        "Turn in presents for Versalmas rewards",
        "Check my daily progress"
    );

    if (action === 0) {
        npc.sendNext("Versalmas is all about sharing bright gifts and brighter cheer. Bring me the stolen presents and I will Versalize them into proper holiday rewards.");
        npc.sendOk("You may turn in bundles of #b10, 20, 30, 40, 60, 90, or 100#k presents, up to #b100 per day#k.");
    } else if (action === 1) {
        npc.sendOk("If you are level 10-20, collect the #bRed and Green#k presents. Level 21-30 should collect #bRed and White#k, level 31-40 should collect #bRed and Blue#k, level 41-60 should collect #bBlue and White#k, and level 61+ should collect #bGreen and White#k presents.");
    } else if (action === 3) {
        npc.sendOk("You have turned in #b" + dailyCount + " / 100#k presents today.");
    } else {
        var picked = npc.askMenu(
            "How many presents would you like to Versalize?#b",
            "10 presents",
            "20 presents",
            "30 presents",
            "40 presents",
            "60 presents",
            "90 presents",
            "100 presents"
        );
        var bundle = choices[picked];

        if (!bundle) {
            npc.sendOk("Come back when you have a proper bundle ready.");
        } else if (bundle.count >= 30 && !choseVersalmas) {
            npc.sendOk("Before I can accept bundles of #b30 or more#k, you need to help #bLittle Suzy#k in New Leaf City choose #bVersalmas#k.");
        } else if (bundle.count === 100 && plr.questData(VERSAL_FULL_KEY) === "end") {
            npc.sendOk("I can only give the #bVersalmas Cactus#k once. You can still turn in #b90 or fewer#k presents for more rewards.");
        } else if (dailyCount + bundle.count > 100) {
            npc.sendOk("You may only turn in up to #b100 presents#k per day, and you have already turned in #b" + dailyCount + "#k.");
        } else if (!plr.haveItem(presentItem, bundle.count)) {
            npc.sendOk("You don't have enough #b#t" + presentItem + "##k for that bundle.");
        } else {
            var fullRewards = [[CHEER_ITEM, bundle.cheer]].concat(bundle.rewards);
            if (!npc.sendYesNo(describeRewards(bundle) + "\r\n\r\nWould you like to make the trade?")) {
                npc.sendOk("Come back when you are ready.");
            } else if (!plr.canHoldAll(fullRewards)) {
                npc.sendOk("Please make enough room in your inventory first.");
            } else if (!plr.gainItem(presentItem, -bundle.count)) {
                npc.sendOk("I couldn't take the presents from you. Please try again.");
            } else {
                for (var i = 0; i < fullRewards.length; i++) {
                    plr.gainItem(fullRewards[i][0], fullRewards[i][1]);
                }
                if (!plr.questData(DAILY_TIME_KEY)) {
                    plr.setQuestData(DAILY_TIME_KEY, String(Date.now()));
                }
                plr.setQuestData(DAILY_COUNT_KEY, String(dailyCount + bundle.count));
                if (bundle.count === 100) {
                    plr.setQuestData(VERSAL_FULL_KEY, "end");
                }
                npc.sendOk("O-Hoy! Thank you for helping me spread Versalmas cheer!");
            }
        }
    }
}
