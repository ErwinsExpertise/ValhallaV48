npc.sendOk("OK, so are you ready? Now you will try to use the Rock of Evolution to evolve the pet.");
npc.sendOk("However, I don't even know what it will turn out to look like. Aren't you even more excited? Well, if you want to see the evolution in progress, then I suggest you move this window down just a little.");

if (plr.getMesos() < 10000) {
    npc.sendOk("What's this? You don't even have enough mesos! What were you doing that kept you from preparing for this most important moment? Tut tut.");
} else {
    var result = plr.requestPetEvol();
    if (result === 0) {
        plr.takeMesos(10000);
        plr.forceCompleteQuest(8185);
        npc.sendOk("The Dragon has successfully evolved! What do you think? Does it look good or what?");
    } else {
        npc.sendOk("Hmm... something's not right. Check again to make sure you are prepared.");
    }
}
