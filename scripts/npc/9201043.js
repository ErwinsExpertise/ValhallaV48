var ticketQuest = 8883;
var cooldownQuest = 8884;
var entranceTicket = 4031592;
var lipLockKey = 4031593;
var outsideMap = 670010000;
var waitingMap = 670010100;
var cooldownMs = 6 * 60 * 60 * 1000;

function lastAttemptMs() {
    var raw = plr.questData(cooldownQuest);
    if (raw == null || raw === "") {
        return 0;
    }
    var parsed = Number(raw);
    return isNaN(parsed) ? 0 : parsed;
}

function remainingMinutes() {
    var remaining = cooldownMs - (Date.now() - lastAttemptMs());
    if (remaining <= 0) {
        return 0;
    }
    return Math.ceil(remaining / 60000);
}

if (plr.mapID() !== outsideMap) {
    npc.sendOk("Please meet me in Amos' Training Ground.");
} else {
    var state = plr.questData(ticketQuest);
    if (plr.itemCount(entranceTicket) > 0) {
        state = "end";
    }

    if (state === "end") {
        if (!plr.haveItem(entranceTicket, 1)) {
            plr.setQuestData(ticketQuest, "");
            npc.sendOk("Don't you have the Entrance Ticket? I'm sorry, but you'll have to bring me 10 #b#t4031593##k for another one.");
        } else if (!plr.removeItemsByID(entranceTicket, 1)) {
            npc.sendOk("Please make room in your inventory and talk to me again.");
        } else {
            plr.setQuestData(ticketQuest, "");
            plr.warp(waitingMap);
        }
    } else if (state === "ing") {
        if (!plr.isMarried()) {
            npc.sendOk("You must be married before you can take on the Amorian Challenge.");
        } else if (plr.level() < 40) {
            npc.sendOk("You'll need to be at least Level 40 to enter my hunting ground.");
        } else if (plr.itemCount(lipLockKey) < 10) {
            npc.sendOk("Bring me #b10 #t4031593##k and I'll hand you an #b#t4031592##k.");
        } else if (!plr.inventoryExchange(lipLockKey, 10, entranceTicket, 1)) {
            npc.sendOk("Please check your inventory space and try again.");
        } else {
            plr.setQuestData(ticketQuest, "end");
            plr.setQuestData(cooldownQuest, String(Date.now()));
            npc.sendOk("A worthy warrior deserves a chance. Take this ticket and head to the entrance.");
        }
    } else if (!plr.isMarried()) {
        npc.sendOk("I admire your bravery, but you must be married to brave the dangers of the Amorian Challenge.");
    } else if (remainingMinutes() > 0) {
        npc.sendOk("Easy there. You can only claim another ticket every 6 hours. Come back in about #b" + remainingMinutes() + " minute(s)#k.");
    } else if (plr.level() < 40) {
        npc.sendOk("You'll need to be at least Level 40 to enter my hunting ground.");
    } else if (npc.sendYesNo("I am Amos the Strong! I built this hunting ground for those strong enough to protect the ones they love. Do you want to take on the Amorian Challenge?")) {
        plr.setQuestData(ticketQuest, "ing");
        npc.sendOk("Good. Bring me #b10 #t4031593##k and I'll give you an #b#t4031592##k.");
    } else {
        npc.sendOk("Come back when you're ready. I'll be waiting.");
    }
}
