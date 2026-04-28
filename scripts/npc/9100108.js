var NORMAL_TICKET = 5220000
var REMOTE_TICKET = 5451000
var GACHAPON_NAME = "Ludibrium"

var rewards = [2048000, 2040601, 2041019, 2041007, 2041016, 2041022, 2041001, 2041010, 2041013, 2041004, 2044701, 2043301, 2040301, 2048004, 2048001, 2040901, 2040701, 2040704, 2040707, 2040602, 2041020, 2041008, 2041017, 2041023, 2041002, 2041011, 2041014, 2041005, 2044702, 2043302, 2040302, 2040002, 2044402, 2048005, 2048002, 2040702, 2040705, 2040708, 2044302, 2043802, 2040402, 2043702, 2044811, 2000004, 2000005, 4006000, 4006001, 1032003, 1432009, 1302022, 1302029, 1102014, 1102018, 1312014, 1302026, 1102015, 1032011, 1312013, 1032008, 1032019, 1032007, 1332030, 1032020, 1032004, 1302027, 1032022, 1312012, 1032021, 1032006, 1302028, 1322003, 1032016, 1032015, 1302024, 1092008, 1032018, 1302021, 1032014, 1332021, 1322012, 1032005, 1032013, 1102012, 1302025, 1302013, 1032002, 1032001, 1032012, 1302017, 1032010, 1402014, 1102017, 1102013, 1442021, 1032009, 1422011, 1402017, 1422005, 1002023, 1332016, 1432005, 1002037, 1002034, 1002064, 1002038, 1002013, 1002036, 1382011, 1002035, 1002065, 1382014, 1372001, 1452026, 1002162, 1002164, 1462018, 1002165, 1452014, 1002163, 1452012, 1002161, 1452009, 1462007, 1332022, 1002175, 1002172, 1002174, 1040096, 1472033, 1002173, 1332054, 1472054, 1002171, 1332016, 1002646, 2040805, 1002419, 1442018]

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
