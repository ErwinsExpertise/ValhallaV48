var NORMAL_TICKET = 5220000
var REMOTE_TICKET = 5451000
var GACHAPON_NAME = "New Leaf City"

var commonRewards = [
    [2000004, 1, 50],
    [2020012, 1, 50],
    [2000005, 1, 50],
    [2030007, 1, 50],
    [2022027, 1, 50]
]

var rareRewards = [2040001, 2041002, 2040805, 2040702, 2043802, 2040402, 2043702, 1302022, 1322021, 1322026, 1302026, 1442017, 1082147, 1102043, 1442016, 1402012, 1302027, 1322027, 1322025, 1312012, 1062000, 1332020, 1302028, 1372002, 1002033, 1092022, 1302021, 1102041, 1102042, 1322024, 1082148, 1002012, 1322012, 1322022, 1002020, 1302013, 1082146, 1442014, 1002096, 1302017, 1442012, 1322010, 1442011, 1442018, 1092011, 1092014, 1302003, 1432001, 1312011, 1002088, 1041020, 1322015, 1442004, 1422008, 1302056, 1432000, 1382001, 1041053, 1060014, 1050053, 1051032, 1050073, 1061036, 1002253, 1002034, 1050067, 1051052, 1002072, 1002144, 1051054, 1050069, 1372007, 1050056, 1050074, 1002254, 1002274, 1002218, 1051055, 1382010, 1002246, 1050039, 1382007, 1372000, 1002013, 1050072, 1002036, 1002243, 1372008, 1382008, 1382011, 1092021, 1051034, 1050047, 1040019, 1041031, 1051033, 1002153, 1002252, 1051024, 1050068, 1382003, 1382006, 1050055, 1051031, 1050025, 1002155, 1002245, 1452004, 1452023, 1060057, 1040071, 1002137, 1462009, 1452017, 1040025, 1041027, 1452005, 1452007, 1061057, 1472006, 1472019, 1060084, 1472028, 1002179, 1082074, 1332015, 1432001, 1060071, 1472007, 1472002, 1051009, 1061037, 1332016, 1332034, 1472020, 1102084, 1102086, 1032026, 1082149]

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
    var rewardId = 0

    if (Math.random() < 0.30) {
        rewardId = rareRewards[Math.floor(Math.random() * rareRewards.length)]
    } else {
        var picked = commonRewards[Math.floor(Math.random() * commonRewards.length)]
        rewardId = picked[0]
    }

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
