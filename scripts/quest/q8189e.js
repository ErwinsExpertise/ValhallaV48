npc.sendOk("OK, so are you ready? Now you will try to use the Rock of Evolution to evolve the pet.");

if (plr.getMesos() < 10000) {
    npc.sendOk("What's this? You don't even have enough mesos! What were you doing that kept you from preparing for this most important moment? Tut tut.");
} else {
    var result = plr.requestPetEvol();
    if (result === 0) {
        plr.takeMesos(10000);
        plr.forceCompleteQuest(8189);
        npc.sendOk("The Dragon has successfully evolved! What do you think? Does it look good or what?");
    } else {
        npc.sendOk("Hmm... something's not right. Check again to make sure you are prepared.");
    }
}
