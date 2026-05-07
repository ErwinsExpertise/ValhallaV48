var QUEST_DIET_MED = 2050;
var QUEST_AGING_MED = 2051;

var dietStatus = plr.getQuestStatus(QUEST_DIET_MED);
var agingStatus = plr.getQuestStatus(QUEST_AGING_MED);

if (agingStatus === 0) {
    if (dietStatus === 1) {
        var firstEntryCost = plr.getLevel() * 100;
        if (npc.sendYesNo("So you came here at the request of #b#p1061005##k to gather medicinal herbs? Well... I inherited this land from my father, and I can't just let strangers in for free. But for #r" + firstEntryCost + "#k mesos, that's a different story. Do you want to pay the entrance fee?")) {
            if (plr.getMesos() < firstEntryCost) {
                npc.sendOk("Are you short on money? Make sure you have at least #r" + firstEntryCost + "#k mesos with you. Don't expect any discounts from me.");
            } else {
                plr.gainMesos(-firstEntryCost);
                plr.warp(101000100);
            }
        } else {
            npc.sendOk("I understand... but you should understand my side too. I can't let you in here for free.");
        }
    } else if (dietStatus === 2) {
        npc.sendOk("It's you from the other day... Is #p1061005# working hard on the diet medicine? I was honestly surprised you made it through this place without much trouble. As a reward, I'll let you in for free for a while. You might even find some interesting things inside while you're at it.");
        plr.warp(101000100);
    } else {
        npc.sendOk("Do you want to enter this place? I'm sure you've heard there are precious medicinal herbs in here, but I can't let some stranger wander around property he doesn't even know I own. Sorry, but that's all there is to it.");
    }
} else if (agingStatus === 1) {
    var secondEntryCost = plr.getLevel() * 200;
    if (npc.sendYesNo("It's you from the other day... did #b#p1061005##k send you here again? What? You need to stay longer this time? Hmm... It's very dangerous in there, but for #r" + secondEntryCost + "#k mesos, I'll let you search as much as you want. So, are you going to pay the entrance fee?")) {
        if (plr.getMesos() < secondEntryCost) {
            npc.sendOk("Are you short on money? Make sure you have at least #r" + secondEntryCost + "#k mesos with you. Don't expect any discounts from me.");
        } else {
            plr.gainMesos(-secondEntryCost);
            plr.warp(101000102);
        }
    } else {
        npc.sendOk("I understand... but you should understand my side too. I can't let you in here for free.");
    }
} else if (agingStatus === 2) {
    npc.sendOk("It's you from the other day... Is #p1061005# working hard on the anti-aging medicine? I was honestly surprised you made it through this place without much trouble. As a reward, I'll let you in for free for a while. You might even find some interesting things inside while you're at it.");
    npc.sendOk("By the way... #p1032100# from this town tried to sneak in earlier. I caught her, but in the process she dropped something in there. I tried looking for it, but I have no idea where it went. Why don't you take a look for it while you're inside?");
    plr.warp(101000102);
} else {
    npc.sendOk("Do you want to enter this place? I'm sure you've heard there are precious medicinal herbs in here, but I can't let some stranger wander around property he doesn't even know I own. Sorry, but that's all there is to it.");
}
