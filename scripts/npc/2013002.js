var diaryPages = [4001064, 4001065, 4001066, 4001067, 4001068, 4001069, 4001070, 4001071, 4001072, 4001073];

function rollReward() {
    var r = Math.floor(Math.random() * 251);
    if (r === 0) return [2000004, 10];
    if (r === 1) return [2000002, 100];
    if (r === 2) return [2000003, 100];
    if (r === 3) return [2000006, 50];
    if (r === 4) return [2022000, 50];
    if (r === 5) return [2022003, 50];
    if (r === 6) return [2040002, 1];
    if (r === 7) return [2040402, 1];
    if (r === 8) return [2040502, 1];
    if (r === 9) return [2040505, 1];
    if (r === 10) return [2040602, 1];
    if (r === 11) return [2040802, 1];
    if (r === 12) return [4003000, 70];
    if (r === 13) return [4010000, 20];
    if (r === 14) return [4010001, 20];
    if (r === 15) return [4010002, 20];
    if (r === 16) return [4010003, 20];
    if (r === 17) return [4010004, 20];
    if (r === 18) return [4010005, 20];
    if (r === 19) return [4010006, 15];
    if (r === 20) return [4020000, 20];
    if (r === 21) return [4020001, 20];
    if (r === 22) return [4020002, 20];
    if (r === 23) return [4020003, 20];
    if (r === 24) return [4020004, 20];
    if (r === 25) return [4020005, 20];
    if (r === 26) return [4020006, 20];
    if (r === 27) return [4020007, 10];
    if (r === 28) return [4020008, 10];
    if (r === 29) return [1032013, 1];
    if (r === 30) return [1032011, 1];
    if (r === 31) return [1032014, 1];
    if (r === 32) return [1102021, 1];
    if (r === 33) return [1102022, 1];
    if (r === 34) return [1102023, 1];
    if (r === 35) return [1102024, 1];
    if (r === 36) return [2040803, 1];
    if (r === 37) return [2070011, 1];
    if (r === 38) return [2043001, 1];
    if (r === 39) return [2043101, 1];
    if (r === 40) return [2043201, 1];
    if (r === 41) return [2043301, 1];
    if (r === 42) return [2043701, 1];
    if (r === 43) return [2043801, 1];
    if (r === 44) return [2044001, 1];
    if (r === 45) return [2044101, 1];
    if (r === 46) return [2044201, 1];
    if (r === 47) return [2044301, 1];
    if (r === 48) return [2044401, 1];
    if (r === 49) return [2044501, 1];
    if (r === 50) return [2044601, 1];
    if (r === 51) return [2044701, 1];
    if (r === 52) return [2000004, 35];
    if (r === 53) return [2000002, 80];
    if (r === 54) return [2000003, 80];
    if (r === 55) return [2000006, 35];
    if (r === 56) return [2022000, 35];
    if (r === 57) return [2022003, 35];
    if (r === 58) return [4003000, 75];
    if (r === 59) return [4010000, 18];
    if (r === 60) return [4010001, 18];
    if (r === 61) return [4010002, 18];
    if (r === 62) return [4010003, 18];
    if (r === 63) return [4010004, 18];
    if (r === 64) return [4010005, 18];
    if (r === 65) return [4010006, 12];
    if (r === 66) return [4020000, 18];
    if (r === 67) return [4020001, 18];
    if (r === 68) return [4020002, 18];
    if (r === 69) return [4020003, 18];
    if (r === 70) return [4020004, 18];
    if (r === 71) return [4020005, 18];
    if (r === 72) return [4020006, 18];
    if (r === 73) return [4020007, 7];
    if (r === 74) return [4020008, 7];
    if (r === 75) return [2040001, 1];
    if (r === 76) return [2040004, 1];
    if (r === 77) return [2040301, 1];
    if (r === 78) return [2040401, 1];
    if (r === 79) return [2040501, 1];
    if (r === 80) return [2040504, 1];
    if (r === 81 || r === 82) return [2040601, 1];
    if (r === 83) return [2040701, 1];
    if (r === 84) return [2040704, 1];
    if (r === 85) return [2040707, 1];
    if (r === 86) return [2040801, 1];
    if (r === 87) return [2040901, 1];
    if (r === 88) return [2041001, 1];
    if (r === 89) return [2041004, 1];
    if (r === 90) return [2041007, 1];
    if (r === 91) return [2041010, 1];
    if (r === 92) return [2041013, 1];
    if (r === 93) return [2041016, 1];
    if (r === 94) return [2041019, 1];
    if (r === 95) return [2041022, 1];
    if (r >= 96 && r <= 130) return [2000004, 20];
    if (r >= 131 && r <= 150) return [2000005, 10];
    if (r >= 151 && r <= 180) return [2000002, 100];
    if (r >= 181 && r <= 200) return [2000006, 50];
    return [2000003, 100];
}

function hasAllDiaryPages() {
    for (var i = 0; i < diaryPages.length; i++) {
        if (!plr.haveItem(diaryPages[i], 1)) {
            return false;
        }
    }
    return true;
}

if (plr.mapID() === 920010100) {
    if (!plr.isLeader()) {
        npc.sendOk("Thank you for restoring me. Please follow your party leader to receive the Goddess's blessing.");
    } else if (plr.getEventProperty("completed")) {
        npc.sendOk("The Goddess's blessing has already been granted.");
    } else {
        npc.sendOk("Thank you for restoring the statue and rescuing me. I will now grant your party my blessing.");
        plr.partyGiveExp(23000);
        plr.setEventProperty("completed", true);
        plr.setEventProperty("bonusStarted", true);
        plr.logEvent("final stage cleared; warping party to bonus room");
        plr.warpEventMembersToPortal(920011100, "sp");
    }
} else if (plr.mapID() === 920011300) {
    var canContinue = true;
    if (hasAllDiaryPages() && !plr.haveItem(4161014, 1)) {
        if (npc.sendYesNo("You found all ten pages of my diary. Would you like me to bind them into #b#t4161014##k for you?")) {
            if (!plr.canHold(4161014, 1)) {
                npc.sendOk("Please make space in your inventory first.");
                canContinue = false;
            }
            if (canContinue) {
                for (var i = 0; i < diaryPages.length; i++) {
                    plr.gainItem(diaryPages[i], -1);
                }
                plr.gainItem(4161014, 1);
            }
        }
    }

    if (canContinue) {
        var reward = rollReward();
        if (!plr.canHold(reward[0], reward[1])) {
            npc.sendOk("Please make sure you have enough inventory space before receiving your reward.");
        } else {
            plr.gainItem(reward[0], reward[1]);
            plr.setQuestData(7020, "1");
            plr.logEvent("reward claimed: " + reward[0] + " x" + reward[1]);
            npc.sendOk("Thank you for rescuing me. Please accept this gift.");
            plr.leaveEvent();
        }
    }
}
