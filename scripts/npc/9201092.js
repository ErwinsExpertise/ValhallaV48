var DAILY_TIME_KEY = 8200;
var DAILY_COUNT_KEY = 8201;
var GRUBBER_INTRO_KEY = 8204;
var GRUBBER_FULL_KEY = 8205;
var STOCK_ITEM = 4031880;
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
        plr.setQuestData(GRUBBER_FULL_KEY, "");
    }
}

function getDailyCount() {
    var raw = plr.questData(DAILY_COUNT_KEY);
    return raw ? parseInt(raw, 10) : 0;
}

function rollMesos(bundleCount) {
    var roll = Math.random();
    if (bundleCount === 25) {
        if (roll < 0.65) return 10000;
        if (roll < 0.90) return 25000;
        return 50000;
    }
    if (bundleCount === 50) {
        if (roll < 0.55) return 25000;
        if (roll < 0.90) return 50000;
        return 150000;
    }
    if (bundleCount === 75) {
        if (roll < 0.55) return 50000;
        if (roll < 0.90) return 200000;
        return 500000;
    }
    if (roll < 0.55) return 100000;
    if (roll < 0.90) return 500000;
    return 1000000;
}

resetDailyIfNeeded();

if (plr.level() < 10) {
    npc.sendOk("Come back when you're strong enough to do real business, kid.");
} else if (plr.questData(GRUBBER_INTRO_KEY) !== "end") {
    if (!npc.sendYesNo("You look like a smart kid... and a smart kid knows when opportunity is knocking. Want in on a lucrative holiday business deal?")) {
        npc.sendOk("Then get out of my way until you're ready to think like a winner.");
    } else {
        npc.sendNext("Maple Claws and O-Pongo are sentimental fools. Me? I pay cash. Bring me stolen presents in bundles of #b25, 50, 75, or 100#k and I'll buy them.");
        npc.sendNext("For every #b25 presents#k you sell me, I'll give you one #b#t" + STOCK_ITEM + "##k from Grubber Industries, plus one shot at a bag full of mesos.");
        plr.setQuestData(GRUBBER_INTRO_KEY, "end");
    }
} else {
    var presentItem = getPresentItem();
    var dailyCount = getDailyCount();
    var bundle = npc.askMenu(
        "How many presents are you selling to me today?#b",
        "25 presents",
        "50 presents",
        "75 presents",
        "100 presents"
    );

    var amounts = [25, 50, 75, 100];
    var chosen = amounts[bundle];
    if (!chosen) {
        npc.sendOk("Come back when you're ready to make a deal.");
    } else if (chosen === 100 && plr.questData(GRUBBER_FULL_KEY) === "end") {
        npc.sendOk("You've already pulled the full #b100 present#k deal for today. Pick a smaller bundle if you want more cash.");
    } else if (dailyCount + chosen > 100) {
        npc.sendOk("You can only sell me #b100 presents#k per day. You've already sold #b" + dailyCount + "#k today.");
    } else if (!plr.haveItem(presentItem, chosen)) {
        npc.sendOk("You don't have enough #b#t" + presentItem + "##k for that deal.");
    } else {
        var stockCount = chosen / 25;
        var mesos = rollMesos(chosen);
        if (!npc.sendYesNo("You'll sell me #b" + chosen + "#k presents for #b" + mesos + " mesos#k and #b#t" + STOCK_ITEM + "##k x " + stockCount + "#k. Deal?")) {
            npc.sendOk("That's fine. Come back when you want real money.");
        } else if (!plr.canHold(STOCK_ITEM, stockCount)) {
            npc.sendOk("Make room in your Etc inventory first. I don't hand out stock certificates without a proper filing system.");
        } else if (!plr.gainItem(presentItem, -chosen)) {
            npc.sendOk("I couldn't take the presents from you. Try again.");
        } else {
            plr.gainMesos(mesos);
            plr.gainItem(STOCK_ITEM, stockCount);
            if (!plr.questData(DAILY_TIME_KEY)) {
                plr.setQuestData(DAILY_TIME_KEY, String(Date.now()));
            }
            plr.setQuestData(DAILY_COUNT_KEY, String(dailyCount + chosen));
            if (chosen === 100) {
                plr.setQuestData(GRUBBER_FULL_KEY, "end");
            }
            npc.sendOk("Pleasure doing business. You pulled in #b" + mesos + " mesos#k this time.");
        }
    }
}
