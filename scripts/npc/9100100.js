var NORMAL_TICKET = 5220000
var REMOTE_TICKET = 5451000
var GACHAPON_NAME = "Henesys"

var rewards = [2040317, 3010013, 2000005, 2022113, 2043201, 2044001, 2041038, 2041039, 2041036, 2041037, 2041040, 2041041, 2041026, 2041027, 2044600, 2043301, 2040308, 2040309, 2040304, 2040305, 2040810, 2040811, 2040812, 2040813, 2040814, 2040815, 2040008, 2040009, 2040010, 2040011, 2040012, 2040013, 2040510, 2040511, 2040508, 2040509, 2040518, 2040519, 2040520, 2040521, 2044401, 2040900, 2040902, 2040908, 2040909, 2044301, 2040406, 2040407, 1302026, 1061054, 1452003, 1382037, 1302063, 1041067, 1372008, 1432006, 1332053, 1432016, 1302021, 1002393, 1051009, 1082148, 1102082, 1061043, 1452005, 1051016, 1442012, 1372017, 1332000, 1050026, 1041062]

if (!plr.haveItem(NORMAL_TICKET, 1) && !plr.haveItem(REMOTE_TICKET, 1)) {
    var info = npc.sendMenu("Welcome to the " + GACHAPON_NAME + " Gachapon. How may I help you?", "What is Gachapon?", "Where can I buy Gachapon tickets?")

    if (info === 0) {
        npc.sendNext("Play Gachapon to earn rare scrolls, equipment, chairs, mastery books, and other cool items! All you need is a #bGachapon Ticket#k to be the winner of a random mix of items.")
        npc.sendOk("You'll find a variety of items from the " + GACHAPON_NAME + " Gachapon, but you'll most likely find items and scrolls related to " + GACHAPON_NAME + ".")
    } else if (info === 1) {
        npc.sendNext("Gachapon Tickets are available in the #rCash Shop#k and can be purchased using NX or Maple Points.")
        npc.sendOk("Click on the red SHOP at the lower right hand corner of the screen to visit the #rCash Shop#k where you can purchase tickets.")
    }
} else if (!npc.sendYesNo("You may use the " + GACHAPON_NAME + " Gachapon. Would you like to use your Gachapon ticket?")) {
    npc.sendOk("See you next time, when you decide to try your luck.")
} else {
    var rewardId = rewards[Math.floor(Math.random() * rewards.length)]
    var ticketId = plr.haveItem(NORMAL_TICKET, 1) ? NORMAL_TICKET : REMOTE_TICKET

    if (!plr.canHold(rewardId, 1)) {
        npc.sendOk("Please check your inventory and try again.")
    } else if (!plr.gainItem(ticketId, -1)) {
        npc.sendOk("I could not take your ticket. Please try again.")
    } else if (!plr.gainItem(rewardId, 1)) {
        npc.sendOk("Please check your inventory and try again.")
    } else {
        npc.dispose()
        npc.sendOk("You have obtained #b#t" + rewardId + "##k.")
    }
}
