var QUEST_DIET = 2050;
var QUEST_AGING = 2051;
var HERB_DIET = 4031020;
var HERB_AGING = 4031032;

var dietStatus = plr.getQuestStatus(QUEST_DIET);
var agingStatus = plr.getQuestStatus(QUEST_AGING);

if (agingStatus === 1) {
    if (!plr.haveItem(HERB_AGING, 1)) {
        npc.sendOk("I found the place where the herb I need grows. Go to Shane's forest near Ellinia and bring me #b#t4031032##k. It's red, and the root splits in two, so remember what you're looking for.");
    } else if (!plr.completeQuest(QUEST_AGING)) {
        npc.sendOk("Make sure you have enough room for the reward before I finish this quest for you.");
    } else {
        npc.sendOk("Ohhh ... this is IT! With #b#t4031032##k, I can finally make the anti-aging medicine!!! Hahaha, if you ever become old and weak, find me. By then I may have a special medicine for just that!");
        npc.sendOk("Oh, I almost forgot. Since you helped me out, I should thank you for your hard work ... #b#t4021009##k is something I found at the very bottom of a valley a long time ago in the middle of a journey. It'll probably help you down the road. I also boosted your fame level and from here on out, #p1032003# may let you in for free. Well, so long...");
    }
} else if (agingStatus === 2) {
    npc.sendOk("It's you!! Thanks to the herbs you got me, the medicine is well on its way. It should be done pretty soon. Thanks again for your help.");
} else if (dietStatus === 1) {
    if (!plr.haveItem(HERB_DIET, 1)) {
        npc.sendOk("Actually, I found a place where there are good medicinal herbs. It's in a forest not too far from here. There are plenty of obstacles along the way, but if you make it through, you should find what we need. Bring me #bPink Anthurium#k from there. It's green grass with a small pink flower, so remember it well.");
    } else if (!plr.completeQuest(QUEST_DIET)) {
        npc.sendOk("Make sure you have enough room for the reward before I finish this quest for you.");
    } else {
        npc.sendOk("Ohhh ... this is IT! With #b#t4031020##k, I can finally make the diet medicine!! Hahaha, if you ever feel like you've gained weight, feel free to find me, because by then, I may have a special medicine in place just for that!");
        npc.sendOk("Oh, I almost forgot. Since you helped me out, I should thank you for your hard work. Here, take this scroll. My brother made this for me a while back, and it improves the defensive ability of your overall armor. And from here on out, #p1032003# will let you in for free. Thanks for your help...");
    }
} else if (dietStatus === 2 && plr.getLevel() >= 50) {
    if (!npc.sendYesNo("Ohhh, you're the traveler that helped me out a lot the other day! I made the diet medicine with the herbs you got me and made some money ... and this time, I'd like to make a different kind of medicine. What do you think? Do you want to help me out one more time?")) {
        npc.sendOk("Alright. Come back if you change your mind.");
    } else if (!plr.startQuest(QUEST_AGING)) {
        npc.sendOk("You can't start that quest right now.");
    } else {
        npc.sendOk("I ran into a place where I can get the herb I need for this new medicine. Go to Shane's forest near Ellinia and bring me #b#t4031032##k. The root splits in two, and it's red, so remember that well.");
    }
} else if (dietStatus === 2) {
    npc.sendOk("Lots of medicinal herbs can be found in this forest. Nothing makes me happier than discovering a new herb here!");
} else if (plr.getLevel() < 25) {
    npc.sendOk("Lots of medicinal herbs can be found in this forest. Nothing makes me happier than discovering a new herb here!");
} else if (!npc.sendYesNo("Wait, hold on one second. I am a herb collector who travels around the world searching for herbs. I've been looking for useful medicinal herbs around this area, but they are harder to find these days. Actually, I found a place where good medicinal herbs grow in abundance. It's in a forest not too far from here. What do you think? Do you want to go there in my place?")) {
    npc.sendOk("That's too bad. Come back if you change your mind.");
} else if (!plr.startQuest(QUEST_DIET)) {
    npc.sendOk("You can't start that quest right now.");
} else {
    npc.sendOk("Actually, I found a place where there are good medicinal herbs. It's in a forest not too far from here. There are plenty of obstacles along the way, but if you make it through, you should find what we need. Bring me #bPink Anthurium#k from there. It's green grass with a small pink flower, so remember it well.");
}
