var COUPON = 4001106;
var GOAL_MAP = 809050015;
var REWARD_MAP = 809050016;
var FAIL_MAP = 809050017;

if (!plr.partyQuestActive()) {
    npc.sendOk("The maze is no longer active. Please talk to Rolly outside the maze if you need to leave.");
} else if (!plr.eventMembersOnMap(GOAL_MAP)) {
    npc.sendOk("I don't believe all members of your party are present at the moment. Let me know when everyone is ready.");
} else if (!plr.isLeader()) {
    npc.sendOk("Great job escaping the maze! Please tell #byour party leader#k to speak to me after gathering all the coupons from the party members.");
} else {
    var count = plr.itemCount(COUPON);

    if (!npc.sendYesNo("So you have gathered up #b" + count + " coupons#k with your collective effort. Are these all that your party has collected?")) {
        npc.sendOk("Please check once more, and let me know when you're ready.");
    } else if (!npc.sendYesNo("Great work! If you gather up 30 Maze Coupons, then you'll receive a cool prize! Would you like to head to the exit?")) {
        npc.sendOk("I am guessing you'd like to collect more coupons. Let me know if you wish to enter the Exit stage.");
    } else {
        plr.logEvent("lmpq: goal accepted coupons=" + count);
        plr.removeAll(COUPON);
        plr.eventGiveExp(50 * count);
        plr.warpEventMembers(count < 30 ? FAIL_MAP : REWARD_MAP);
        plr.finishEvent();
    }
}
