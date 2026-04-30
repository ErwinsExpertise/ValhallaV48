var hatThresholds = [5, 10, 15, 20];
var hatItems = [1002571, 1002572, 1002573, 1002574];

function pirateHatOwned(plrRef) {
    for (var i = 0; i < hatItems.length; i++) {
        if (plrRef.haveItem(hatItems[i], 1) || plrRef.isWearingItem(hatItems[i])) {
            return true;
        }
    }
    return false;
}

function grantOrUpgradeHat(clearCount, tier) {
    if (tier >= hatItems.length) {
        npc.sendOk("You have already received every reward I can offer. Thank you again for saving us.");
        return;
    }

    if (clearCount < hatThresholds[tier]) {
        npc.sendOk("You have defeated Lord Pirate #b" + clearCount + " time(s)#k. Keep going if you want the next reward.");
        return;
    }

    if (plr.getEquipInventoryFreeSlot() < 1) {
        npc.sendOk("Please make space in your Equip inventory first.");
        return;
    }

    if (tier === 0) {
        if (pirateHatOwned(plr)) {
            npc.sendOk("You already have the first Lord Pirate hat.");
            return;
        }
        if (plr.giveItem(hatItems[0], 1)) {
            plr.setQuestData(7041, "1");
            npc.sendOk("Please accept #b#t1002571##k as a token of our gratitude.");
            return;
        }
        npc.sendOk("I could not give you the hat. Please make sure your Equip inventory has space.");
        return;
    }

    var oldHat = hatItems[tier - 1];
    var newHat = hatItems[tier];
    if (plr.isWearingItem(oldHat)) {
        npc.sendOk("Please unequip your current Lord Pirate hat before asking for an upgrade.");
        return;
    }
    if (!plr.haveItem(oldHat, 1)) {
        npc.sendOk("Bring me your current Lord Pirate hat and I will upgrade it.");
        return;
    }
    if (plr.haveItem(newHat, 1) || plr.isWearingItem(newHat)) {
        npc.sendOk("You already have that upgraded hat.");
        return;
    }
    if (plr.inventoryExchange(oldHat, 1, newHat, 1)) {
        plr.setQuestData(7041, String(tier + 1));
        npc.sendOk("Your Lord Pirate hat has been upgraded.");
        return;
    }
    npc.sendOk("I could not upgrade your hat. Please check your inventory and try again.");
}

if (plr.mapID() === 925100500) {
    if (!plr.isPartyLeader()) {
        npc.sendOk("Please ask your party leader to speak to me.");
    } else {
        var exp = 42000;
        var over70 = parseInt(plr.getEventProperty("over70") || "0", 10);
        var avgLevel = parseInt(plr.getEventProperty("avgLevel") || "0", 10);

        if (over70 > 0) {
            if (avgLevel <= 70) {
                exp = 35000;
            } else if (avgLevel <= 80) {
                exp = 28000;
            } else if (avgLevel <= 90) {
                exp = 20000;
            } else {
                exp = 10000;
            }
        }

        if (plr.getEventProperty("bossRewarded") !== true) {
            var members = plr.partyMembersOnMap();
            for (var i = 0; i < members.length; i++) {
                var clears = parseInt(members[i].questData(7040) || "0", 10);
                if (clears < 500) {
                    members[i].setQuestData(7040, String(clears + 1));
                }
            }

            plr.partyGiveExp(exp);
            plr.setEventProperty("completed", true);
            plr.setEventProperty("bossRewarded", true);
        }

        npc.sendOk("Thank you for rescuing me from Lord Pirate. I will take your party to a safe place so we can properly thank you.");
        plr.warpEventMembersToPortal(925100600, "st00");
    }
} else if (plr.mapID() === 925100600) {
    if (plr.questData(7040) === "") {
        plr.setQuestData(7040, "0");
    }
    if (plr.questData(7041) === "") {
        plr.setQuestData(7041, "0");
    }

    var clearCount = parseInt(plr.questData(7040) || "0", 10);
    var tier = parseInt(plr.questData(7041) || "0", 10);

    npc.sendSelection("Thank you for saving us from Lord Pirate. How can I help you?#b\r\n#L0#Check my Lord Pirate clear rewards.#l\r\n#L1#Reset my Lord Pirate record.#l\r\n#L2#Leave this place.#l#k");
    var sel = npc.selection();

    if (sel === 0) {
        grantOrUpgradeHat(clearCount, tier);
    } else if (sel === 1) {
        if (tier < 1 || clearCount < hatThresholds[0]) {
            npc.sendOk("You cannot reset your Lord Pirate record yet.");
        } else if (pirateHatOwned(plr)) {
            npc.sendOk("Please remove every Lord Pirate hat from your inventory and equipment before resetting your record.");
        } else if (npc.sendYesNo("Reset your Lord Pirate clear count back to 0?")) {
            plr.setQuestData(7040, "0");
            plr.setQuestData(7041, "0");
            npc.sendOk("Your Lord Pirate record has been reset.");
        }
    } else {
        plr.leaveEvent();
    }
}
