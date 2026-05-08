function tryEvolution(questId) {
    npc.sendOk("OK, so are you ready? Now you will try to use the Rock of Evolution to evolve the pet.");

    if (questId === 8185) {
        npc.sendOk("However, I don't even know what it will turn out to look like. Aren't you even more excited? Well, if you want to see the evolution in progress, then I suggest you move this window down just a little.");
    }

    if (plr.getMesos() < 10000) {
        npc.sendOk("What's this? You don't even have enough mesos! What were you doing that kept you from preparing for this most important moment? Tut tut.");
        return;
    }

    var result = plr.requestPetEvol();
    if (result === 0) {
        plr.takeMesos(10000);
        plr.forceCompleteQuest(questId);
        npc.sendOk("The Dragon has successfully evolved! What do you think? Does it look good or what?");
    } else {
        npc.sendOk("Hmm... something's not right. Check again to make sure you are prepared.");
    }
}

function startQuest(questId, prompt, decline, failure) {
    if (!npc.sendYesNo(prompt)) {
        npc.sendOk(decline);
        return;
    }
    if (!plr.startQuest(questId)) {
        npc.sendOk(failure);
    }
}

if (plr.questStarted(8185)) {
    tryEvolution(8185);
} else if (plr.questStarted(8189)) {
    tryEvolution(8189);
} else if (!plr.questCompleted(8184)) {
    if (plr.questStarted(8184)) {
        if (!npc.sendYesNo("Can you prove your love for the Baby Dragon?\r\n\r\nGreat job. It's very evident that you really care for your pet.")) {
            npc.sendOk("You still haven't gotten rid of those pests? Be strong.");
        } else if (!plr.completeQuest(8184)) {
            npc.sendOk("You still haven't gotten rid of those pests? Be strong.");
        }
    } else {
        startQuest(
            8184,
            "The Dragon is a special pet that has an ability to evolve. However, before evolving it you must give your Dragon assurance that it is loved by its master. Prove your love for your pet.\r\n\r\nI heard that recently the number of monsters hunting Dragon eggs is increasing. First things first, those monsters need to be dealt with. I want you to bring me 50 #b#t4000023#s#k and 50 #b#t4000029#s#k.",
            "People are such foolish animals. They say the appearance of an evolved Dragon is quite spectacular, but if you say you don't have time, then... goodbye.",
            "People are such foolish animals. They say the appearance of an evolved Dragon is quite spectacular, but if you say you don't have time, then... goodbye."
        );
    }
} else if (!plr.questCompleted(8185)) {
    startQuest(
        8185,
        "I can see that you truly care for your Dragon. Now, shall I help you with your Baby Dragon's evolution?\r\n\r\nIn order for you to help with the evolution of your Baby Dragon, I will need an item called #t5380000#. Someone told me it is being sold in the Cash Shop.\r\n\r\nThe Rock of Evolution, coupled with 10,000 mesos, should be enough to evolve your Baby Dragon. What do you think? I think it's a reasonable deal.",
        "Your heart's not ready yet? That's disappointing...",
        "I thought I told you that you can't do it without the Rock of Evolution. Weren't you listening? Hurry up and find it!"
    );
} else if (!plr.questCompleted(8189)) {
    startQuest(
        8189,
        "If you are bored with your pet's current appearance, then try evolving it once more.\r\n\r\nFor re-evolution you will need 1 #t5380000# just like the first time. I know you saw it in the Cash Shop before, so you know how to get it right?\r\n\r\nThe cost is the same as well, 10,000 mesos. Can't you just see how much better your Dragon will look?",
        "Your heart's not ready yet? That's disappointing...",
        "I thought I told you that you can't do it without the Rock of Evolution. Weren't you listening? Hurry up and find it!"
    );
} else {
    npc.sendOk("Hi, I'm Garnox the Pet Scientist. Have you heard of the evolution of special pets?");
}
