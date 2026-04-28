var NORMAL_TICKET = 5220000
var REMOTE_TICKET = 5451000
var GACHAPON_NAME = "Perion"

var rewards = [2022176, 2022113, 2043202, 2043201, 2044102, 2044101, 2040602, 2040601, 2043302, 2043301, 2040002, 2040001, 2044402, 2002017, 1402010, 1312014, 1442017, 1002063, 1060062, 1050018, 1002392, 1040037, 1002160, 1060005, 1332009, 1332008, 1442009, 1302004, 1312006, 1002154, 1002175, 1060064, 1061088, 1402012, 1002024, 1312005, 1432002, 1302050, 1002048, 1040061, 1041067, 1002131, 1072263, 1332001, 1312027, 1322015, 1432006, 1041088, 1061087, 1402013, 1302051, 1002023, 1402006, 1322000, 1372002, 1442001, 1422004, 1432003, 1040088, 1002100, 1041004, 1061047, 1322022, 1040021, 1061091, 1102012, 1050006, 1060018, 1041044, 1041024, 1041087, 1082146, 1332043, 1062001, 1051014, 1402030, 1432004, 1060060, 1432018, 1002096, 1442010, 1422003, 1472014, 1002021, 1060060, 1442031, 1402000, 1040089, 1432005, 2040402, 2022130, 4130014, 2000004, 2000005, 2022113, 1322008, 1302021, 1322022, 1302013, 1051010, 1060079, 1002005, 1002023, 1002085, 1332017, 1322010, 1051031, 1002212, 1002117, 1040081, 1051037, 1472026, 1332015, 1041060, 1472003, 1060086, 1060087, 1472009, 1060051, 1041080, 1041106, 1092018]

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
