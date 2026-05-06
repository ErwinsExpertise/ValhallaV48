var DAILY_TIME_KEY = 8200;
var DAILY_COUNT_KEY = 8201;
var FULL_BUNDLE_KEY = 8202;
var STAR_ITEM = 4031878;
var DAY_MS = 24 * 60 * 60 * 1000;
var MAPLEMAS_CHOICE_KEY = 4997;
var VERSALMAS_CHOICE_KEY = 4998;
var DELIVERY_ITEM = 4031486;
var RETURN_GIFT_ITEM = 4031519;
var RETURN_GIFT_STATE_KEY = 8831;
var DELIVERY_PRE_KEYS = [8839, 8840, 8841, 8842];
var DELIVERY_STATE_KEYS = [8835, 8836, 8837, 8838];
var RETURN_GIFT_REWARDS = [2000000, 2000003, 4020003, 1322000, 2060000, 4010004, 2000006, 4011006, 2000001, 2022120, 4010003, 4010005, 2050004, 2000005, 2000004, 1072103, 2000002, 2002010, 1040044, 4010006, 2002004, 4004000, 2041013, 2041016, 2041019, 2041022];

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

function getDailyCount() {
    var raw = plr.questData(DAILY_COUNT_KEY);
    return raw ? parseInt(raw, 10) : 0;
}

function resetDailyIfNeeded() {
    var last = parseInt(plr.questData(DAILY_TIME_KEY) || "0", 10);
    if (!last) {
        return;
    }
    if (Date.now() - last >= DAY_MS) {
        plr.setQuestData(DAILY_TIME_KEY, "");
        plr.setQuestData(DAILY_COUNT_KEY, "0");
        plr.setQuestData(FULL_BUNDLE_KEY, "");
    }
}

function rewardBundle(bundle) {
    var rewards = [[STAR_ITEM, bundle.stars]];
    for (var i = 0; i < bundle.rewards.length; i++) {
        rewards.push(bundle.rewards[i]);
    }
    return rewards;
}

function bundleText(bundle) {
    var text = "Here is what I can offer for #b" + bundle.count + "#k presents:\r\n";
    text += "#b#t" + STAR_ITEM + "##k x " + bundle.stars;
    for (var i = 0; i < bundle.rewards.length; i++) {
        text += "\r\n#b#t" + bundle.rewards[i][0] + "##k x " + bundle.rewards[i][1];
    }
    return text;
}

function activeDeliveryIndex() {
    for (var i = 0; i < DELIVERY_PRE_KEYS.length; i++) {
        if (plr.questData(DELIVERY_PRE_KEYS[i]) === "ing") {
            return i;
        }
    }
    return -1;
}

function giveReplacementBox(index) {
    if (plr.haveItem(DELIVERY_ITEM, 1)) {
        npc.sendOk("You still have the present box with you. Please deliver it first.");
    } else if (!plr.canHold(DELIVERY_ITEM, 1)) {
        npc.sendOk("Please make room in your Etc inventory first.");
    } else {
        plr.gainItem(DELIVERY_ITEM, 1);
        npc.sendOk("You seem to have misplaced the present box, so I prepared another one for you. Please try not to lose it this time.");
    }
}

function startDelivery() {
    var current = activeDeliveryIndex();
    if (current >= 0) {
        if (plr.questData(DELIVERY_STATE_KEYS[current]) === "end") {
            if (!plr.canHold(RETURN_GIFT_ITEM, 1)) {
                npc.sendOk("Please make room in your Etc inventory first so I can hand you your returned gift.");
                return;
            }
            plr.gainItem(RETURN_GIFT_ITEM, 1);
            plr.setQuestData(DELIVERY_PRE_KEYS[current], "end");
            plr.setQuestData(RETURN_GIFT_STATE_KEY, "ing");
            npc.sendOk("Wonderful! You delivered the present safely. Here is a returned gift for you. Ask someone to help you open it when you're ready.");
            return;
        }
        giveReplacementBox(current);
        return;
    }

    var choice = npc.askMenu(
        "I still need to send presents to a few friends in other towns. Who should receive one for me?#b",
        "Rowen the Fairy in Ellinia",
        "Ayan in Perion",
        "Ericsson in Orbis",
        "Porter in Omega Sector"
    );

    if (!plr.canHold(DELIVERY_ITEM, 1)) {
        npc.sendOk("Please make room in your Etc inventory first.");
        return;
    }

    plr.gainItem(DELIVERY_ITEM, 1);
    plr.setQuestData(DELIVERY_PRE_KEYS[choice], "ing");
    plr.setQuestData(DELIVERY_STATE_KEYS[choice], "");
    npc.sendOk("Please bring this present box to my friend and keep trying if they are too busy to accept it at first.");
}

function openReturnedGift() {
    if (plr.questData(RETURN_GIFT_STATE_KEY) !== "ing") {
        npc.sendOk("You do not have a returned holiday gift waiting to be opened right now.");
        return;
    }
    if (!plr.haveItem(RETURN_GIFT_ITEM, 1)) {
        npc.sendOk("You seem to have misplaced the returned gift. Please finish another delivery and come back.");
        return;
    }
    var reward = RETURN_GIFT_REWARDS[Math.floor(Math.random() * RETURN_GIFT_REWARDS.length)];
    if (!plr.canHold(reward, 1)) {
        npc.sendOk("Please make room in your inventory before opening the gift.");
        return;
    }
    if (!plr.gainItem(RETURN_GIFT_ITEM, -1)) {
        npc.sendOk("I couldn't take the gift package from you. Please try again.");
        return;
    }
    plr.gainItem(reward, 1);
    plr.setQuestData(RETURN_GIFT_STATE_KEY, "done");
    npc.sendOk("Let's see what was inside... here you go!");
}

resetDailyIfNeeded();

if (plr.level() < 10) {
    npc.sendOk("I need a helper who is at least level 10 before I can trust them with Maplemas work.");
} else if (plr.questData(VERSALMAS_CHOICE_KEY) === "2") {
    npc.sendOk("You decided to celebrate #bVersalmas#k this season, so you should bring those presents to #bO-Pongo#k in New Leaf City instead.");
} else {
    var presentItem = getPresentItem();
    var dailyCount = getDailyCount();
    var choseMaplemas = plr.questData(MAPLEMAS_CHOICE_KEY) === "2";
    var choices = [
        {
            count: 10,
            stars: 1,
            rewards: [[2000002, 25]]
        },
        {
            count: 20,
            stars: 2,
            rewards: [[2000006, 30]]
        },
        {
            count: 30,
            stars: 3,
            rewards: [[[2022195, 5], [2022190, 5], [2002020, 5], [2002021, 5]][Math.floor(Math.random() * 4)]]
        },
        {
            count: 40,
            stars: 4,
            rewards: [[2022271, 3]]
        },
        {
            count: 60,
            stars: 6,
            rewards: [
                [[2002020, 5], [2002021, 5]],
                [[2041014, 1]],
                [[2041017, 1]],
                [[2041020, 1]],
                [[2041023, 1]],
                [[1302080, 1]]
            ][Math.floor(Math.random() * 6)]
        },
        {
            count: 90,
            stars: 9,
            rewards: [
                [[2002023, 10]],
                [[2022182, 1]],
                [[1432015, 1]],
                [[2040805, 1]],
                [[1302080, 1]],
                [[2022276, 5]],
                [[1432046, 1]]
            ][Math.floor(Math.random() * 7)]
        },
        {
            count: 100,
            stars: 10,
            rewards: [[1432046, 1]]
        }
    ];

    var action = npc.askMenu(
        "Hi! I'm Maple Claws, and I still need help recovering the stolen Maplemas presents. What would you like to do?#b",
        "How does this event work?",
        "Which present colors should I collect?",
        "Turn in my recovered presents",
        "Check my daily progress",
        "Help deliver presents to my friends",
        "Open a returned holiday gift"
    );

    if (action === 0) {
        npc.sendNext("Monsters all over the world ran off with Maplemas presents. Bring the correct present type for your level back to me in bundles, and I'll reward you for helping me recover them.");
        npc.sendNext("You can turn in #b10, 20, 30, 40, 60, 90, or 100#k presents at a time, and you may recover up to #b100 presents per day#k with me.");
    } else if (action === 1) {
        npc.sendOk("If you are level 10-20, collect the #bRed and Green#k presents. Level 21-30 should collect #bRed and White#k, level 31-40 should collect #bRed and Blue#k, level 41-60 should collect #bBlue and White#k, and level 61+ should collect #bGreen and White#k presents.");
    } else if (action === 3) {
        npc.sendOk("You have turned in #b" + dailyCount + " / 100#k presents today.");
    } else if (action === 4) {
        startDelivery();
    } else if (action === 5) {
        openReturnedGift();
    } else {
        var labels = [];
        for (var i = 0; i < choices.length; i++) {
            labels.push(choices[i].count + " presents");
        }

        var pickedIndex = npc.askMenu("Choose the bundle you want to turn in.#b", labels[0], labels[1], labels[2], labels[3], labels[4], labels[5], labels[6]);
        var bundle = choices[pickedIndex];

        if (!bundle) {
            npc.sendOk("Come back when you have a proper bundle ready.");
        } else if (bundle.count >= 30 && !choseMaplemas) {
            npc.sendOk("Before I can accept bundles of #b30 or more#k, you need to help #bLittle Suzy#k in New Leaf City decide whether she wants #bMaplemas#k or #bVersalmas#k. If she picks Maplemas, come back to me.");
        } else if (dailyCount + bundle.count > 100) {
            npc.sendOk("You can only turn in up to #b100 presents#k per day. You have already turned in #b" + dailyCount + "#k today.");
        } else if (!plr.haveItem(presentItem, bundle.count)) {
            npc.sendOk("You don't have enough #b#t" + presentItem + "##k for that bundle.");
        } else {
            var rewards = rewardBundle(bundle);
            var confirm = npc.sendYesNo(bundleText(bundle) + "\r\n\r\nWould you like to make the trade?");
            if (!confirm) {
                npc.sendOk("Come back when you are ready to trade them in.");
            } else if (!plr.canHoldAll(rewards)) {
                npc.sendOk("Please make enough room in your inventory before trading those presents in.");
            } else if (!plr.gainItem(presentItem, -bundle.count)) {
                npc.sendOk("I couldn't take the presents from you. Please try again.");
            } else {
                for (var j = 0; j < rewards.length; j++) {
                    plr.gainItem(rewards[j][0], rewards[j][1]);
                }
                if (!plr.questData(DAILY_TIME_KEY)) {
                    plr.setQuestData(DAILY_TIME_KEY, String(Date.now()));
                }
                plr.setQuestData(DAILY_COUNT_KEY, String(dailyCount + bundle.count));
                if (bundle.count === 100) {
                    plr.setQuestData(FULL_BUNDLE_KEY, "done");
                }
                npc.sendOk("Wonderful! Every recovered present helps keep Maplemas alive. Here is your reward.");
            }
        }
    }
}
