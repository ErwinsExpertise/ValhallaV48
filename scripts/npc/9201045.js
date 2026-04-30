function showClear() {
    map.showEffect("quest/party/clear");
    map.playSound("Party1/Clear");
}

function clearPQItems(member) {
    var items = [4031592, 4031593, 4031594, 4031595, 4031596, 4031597];
    for (var i = 0; i < items.length; i++) {
        member.removeAll(items[i]);
    }
}

function eventMembersHere() {
    var all = plr.partyMembersOnMap();
    var here = [];
    for (var i = 0; i < all.length; i++) {
        if (all[i].mapID() === plr.mapID()) {
            here.push(all[i]);
        }
    }
    return here;
}

function eventMembers() {
    return plr.partyMembersOnMap();
}

function setCouplePairs() {
    var members = eventMembers();
    var seen = {};
    var pairs = [];
    for (var i = 0; i < members.length; i++) {
        var member = members[i];
        if (!member.isMarried() || seen[member.name()]) {
            continue;
        }
        var partner = member.partnerName();
        if (!partner || seen[partner]) {
            continue;
        }
        if (member.gender() === 0) {
            pairs.push(member.name() + "=" + partner);
            seen[member.name()] = true;
            seen[partner] = true;
        }
    }
    plr.setEventProperty("apqCouplePairs", pairs.join(";"));
    plr.setEventProperty("apqCoupleEligible", pairs.length === 3);
}

function getCouplePairs() {
    var raw = String(plr.getEventProperty("apqCouplePairs") || "");
    if (raw === "") {
        return [];
    }
    return raw.split(";");
}

function openVault() {
    var members = eventMembers();
    for (var i = 0; i < members.length; i++) {
        clearPQItems(members[i]);
    }
    plr.warpEventMembersToPortal(670010800, "st00");
}

function rewardCupidPieces() {
    var count = plr.itemCount(4031597);
    if (count < 1) {
        npc.sendOk("You don't have any #b#t4031597##k for me.");
        return;
    }

    var exp = 0;
    if (count <= 10) exp = 800;
    else if (count <= 20) exp = 1300;
    else if (count <= 30) exp = 2000;
    else if (count < 35) exp = 2800;
    else if (count === 35) exp = 4000;
    else {
        npc.sendOk("I only accept a maximum of 35 Cupid Code Pieces from one person.");
        return;
    }

    if (!plr.removeItemsByID(4031597, count)) {
        npc.sendOk("Please check your inventory and try again.");
        return;
    }
    plr.giveEXP(exp);
    npc.sendOk("You've earned some bonus EXP for those Cupid Code Pieces.");
}

function rewardSpecialCape() {
    var capes = plr.gender() === 0 ? [1102101, 1102102, 1102103] : [1102104, 1102105, 1102106];
    var cape = capes[Math.floor(Math.random() * capes.length)];
    if (plr.haveItem(cape, 1)) {
        rewardCupidPieces();
        return;
    }
    if (!plr.removeItemsByID(4031597, 35)) {
        npc.sendOk("You don't have the 35 Cupid Code Pieces I need.");
        return;
    }
    if (!plr.giveItem(cape, 1)) {
        plr.giveItem(4031597, 35);
        npc.sendOk("Please make room in your Equip inventory and try again.");
        return;
    }
    npc.sendOk("Congratulations! You've received a special cape in memory of Elias and Cecelia.");
}

function stage4() {
    if (!plr.getEventProperty("apqStage3Clear")) {
        npc.sendOk("Please clear the previous stage first.");
    } else if (!plr.isLeader()) {
        npc.sendOk("Ask your party leader to talk to me.");
    } else if (plr.itemCount(4031597) < 50) {
        npc.sendOk("You'll need 50 #b#t4031597##k to unlock the door.");
    } else if (!plr.removeItemsByID(4031597, 50)) {
        npc.sendOk("Please check your inventory and try again.");
    } else {
        plr.setEventProperty("apqStage4Clear", true);
        plr.partyGiveExp(8000);
        showClear();
        plr.warpEventMembersToPortal(670010600, "st00");
    }
}

function stage5() {
    if (!plr.getEventProperty("apqStage4Clear")) {
        npc.sendOk("Please clear the previous stage first.");
        return;
    }

    if (!plr.getEventProperty("apqStage5Clear")) {
        if (!plr.isLeader()) {
            npc.sendOk("Ask your party leader to talk to me.");
            return;
        }
        if (map.playersInArea(0) !== 6) {
            npc.sendOk("I'm sorry, but it seems like not all of your party members are here yet.");
            return;
        }

        plr.setEventProperty("apqStage5Clear", true);
        plr.setEventProperty("apqStage5PerfectGender", map.malePlayersInArea(1) === 3 && map.femalePlayersInArea(1) === 3);
        setCouplePairs();
        plr.partyGiveExp(9000);
        showClear();
        npc.sendOk("Magnificent speed! The end is near. Talk to me again when your party is ready for the final stage.");
        return;
    }

    if (!plr.isLeader()) {
        npc.sendOk("Ask your party leader to move the party to the last stage.");
        return;
    }

    var members = eventMembers();
    for (var i = 0; i < members.length; i++) {
        members[i].removeAll(4031597);
    }
    plr.warpEventMembersToPortal(670010700, "st00");
}

function stage6() {
    if (!plr.getEventProperty("apqStage5Clear")) {
        npc.sendOk("Please clear the previous stage first.");
        return;
    }

    if (!plr.getEventProperty("apqBossSpawned")) {
        if (!plr.isLeader()) {
            npc.sendOk("Ask your party leader to talk to me.");
            return;
        }
        if (map.playersInArea(0) !== 6) {
            npc.sendOk("It seems like not all of your party members are present at the moment.");
            return;
        }
        if (npc.sendYesNo("Are you ready to descend and face Geist Balrog?")) {
            plr.setEventProperty("apqBossSpawned", true);
            plr.warpEventMembersToPortal(670010700, "st01");
            plr.spawnMonster(9400536, 933, 432);
        }
        return;
    }

    if (!plr.getEventProperty("apqStage6Clear")) {
        if (!plr.isLeader()) {
            npc.sendOk("Ask your party leader to talk to me.");
            return;
        }
        if (!plr.haveItem(4031594, 1)) {
            npc.sendOk("You must bring me a #b#t4031594##k to prove your victory over Geist Balrog.");
            return;
        }
        if (!plr.removeItemsByID(4031594, 1)) {
            npc.sendOk("Please check your inventory and try again.");
            return;
        }
        plr.setEventProperty("apqStage6Clear", true);
        plr.partyGiveExp(11000);
        showClear();
    }

    if (!plr.isLeader()) {
        npc.sendOk("Your leader will decide where the party goes next.");
    } else if (plr.getEventProperty("apqCoupleEligible") && plr.getEventProperty("apqStage5PerfectGender")) {
        plr.warpEventMembersToPortal(670010750, "st00");
    } else {
        openVault();
    }
}

function couplesStage() {
    var pairs = getCouplePairs();
    var eligible = !!plr.getEventProperty("apqCoupleEligible");

    if (!plr.getEventProperty("apqCoupleClear")) {
        if (!eligible) {
            if (!plr.isLeader()) {
                npc.sendOk("Ask your party leader what to do next.");
                return;
            }
            npc.sendSelection("It doesn't look like every pair in your party is a married couple.#b\r\n#L0#Wait for everyone.#l\r\n#L1#Take us to the vault instead.#l#k");
            if (npc.selection() === 1) {
                openVault();
            } else {
                npc.sendOk("I'll wait for your party members.");
            }
            return;
        }

        if (map.playersInArea(2) !== 6) {
            if (!plr.isLeader()) {
                npc.sendOk("Wait for the rest of your party to gather here.");
                return;
            }
            npc.sendSelection("What would you like to do?#b\r\n#L0#We'll wait for the rest of the party.#l\r\n#L1#We've lost some members. Take us to the vault.#l#k");
            if (npc.selection() === 1) {
                openVault();
            } else {
                npc.sendOk("I'll wait.");
            }
            return;
        }

        if (plr.gender() === 0) {
            for (var i = 0; i < pairs.length; i++) {
                var maleName = pairs[i].split("=")[0];
                if (plr.name() === maleName) {
                    plr.warpToPortalName(670010750, "st0" + (i + 1));
                    return;
                }
            }
            npc.sendOk("Men should enter first.");
            return;
        }

        if (map.malePlayersInArea(0) !== 0) {
            npc.sendOk("Please let the men enter their rooms first.");
            return;
        }

        var partner = plr.partnerName();
        for (var j = 0; j < pairs.length; j++) {
            var split = pairs[j].split("=");
            if (split[0] === partner) {
                plr.warpToPortalName(670010750, "st0" + (j + 1));
                return;
            }
        }
        npc.sendOk("I couldn't match you with your spouse. Please have your leader talk to me.");
        return;
    }

    npc.sendSelection("Hello there. What would you like to do?#b\r\n#L0#I brought Cupid Code Pieces. Reward me.#l\r\n#L1#Take us to the vault.#l#k");
    var sel = npc.selection();
    if (sel === 1) {
        if (!plr.isLeader()) {
            npc.sendOk("Ask your party leader to take everyone to the vault.");
            return;
        }
        openVault();
        return;
    }

    var first = String(plr.getEventProperty("apqFirstClaimer") || "");
    var firstMate = String(plr.getEventProperty("apqFirstMate") || "");
    var hadFullTurnIn = !!plr.getEventProperty("apqFirstClaimFull35");
    if (first === "" && plr.itemCount(4031597) > 0) {
        plr.setEventProperty("apqFirstClaimer", plr.name());
        plr.setEventProperty("apqFirstMate", plr.partnerName());
        first = plr.name();
        firstMate = plr.partnerName();
    }

    if ((plr.name() === first || plr.name() === firstMate) && plr.itemCount(4031597) === 35) {
        if (plr.name() === first) {
            plr.setEventProperty("apqFirstClaimFull35", true);
        }
        if (Math.floor(Math.random() * 51) === 0) {
            rewardSpecialCape();
        } else {
            rewardCupidPieces();
        }
        return;
    }

    if (plr.name() !== firstMate && plr.name() !== first && hadFullTurnIn && plr.itemCount(4031597) === 35) {
        rewardCupidPieces();
        return;
    }

    rewardCupidPieces();
}

if (plr.mapID() === 670010500) {
    stage4();
} else if (plr.mapID() === 670010600) {
    stage5();
} else if (plr.mapID() === 670010700) {
    stage6();
} else if (plr.mapID() === 670010750) {
    couplesStage();
} else {
    npc.sendOk("There's nothing for me to handle here.");
}
