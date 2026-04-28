var NORMAL_TICKET = 5220000
var REMOTE_TICKET = 5451000
var GACHAPON_NAME = "Sleepywood"

var rewards = [2044102, 2044502, 2041020, 2041017, 2041011, 2041014, 2044702, 2044602, 2043302, 2040302, 2040805, 2040502, 2044402, 2040902, 2040708, 2040402, 2043002, 2044101, 2041022, 2044701, 2040804, 2040702, 2040707, 2043801, 2044001, 2043701, 2048003, 2048000, 4020001, 4020002, 2060001, 2020002, 2012003, 4004002, 4020007, 2000004, 2012001, 2050003, 2020005, 4010006, 2020004, 2002002, 2020012, 2020009, 4010005, 2020003, 4004000, 2000005, 2020013, 2030000, 2030001, 2030002, 2030003, 2030004, 2030005, 2030006, 2030007, 2030019, 2020000, 2012002, 4020005, 4010004, 2020014, 4006001, 4006000, 2050002, 2002003, 1032003, 1302022, 1432009, 1102014, 1102018, 1322023, 1322025, 1032008, 1432008, 1322022, 1442018, 1442039, 1322027, 1032004, 1032026, 1442015, 1032016, 1032018, 1422004, 1422006, 1302021, 1322024, 1322012, 1051017, 1432015, 1032001, 1432018, 1432000, 1402014, 1032000, 1422000, 1032009, 1082145, 1452004, 1452000, 1002162, 1452003, 1040068, 1060057, 1002163, 1452001, 1002161, 1002038, 1002036, 1002013, 1372002, 1372003, 1040044, 1041048, 1041049, 1002150, 1472004, 1002175, 1472003, 1472002, 1041050, 1041047, 1332013, 1060019, 1442009, 1442006, 1422002, 1402003, 1040021, 1442010, 1002009, 1442007, 1442031, 1332002, 1060017, 2100000, 2044902, 2044901, 2044803, 2044804]

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
