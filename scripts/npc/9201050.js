var minLevel = 10;

if (plr.questCompleted(4911)) {
    npc.sendOk("Good job! You've solved all of my questions about NLC. Enjoy your trip!");
} else if (plr.questCompleted(4900) || plr.questStarted(4900)) {
    npc.sendOk("Hey, pay attention, I'm trying to quiz you on another question, fam!");
} else {
    var selection = npc.sendMenu(
        "What up! Name's Icebyrd Slimm, mayor of New Leaf City! Happy to see you accepted my invite. So, what can I do for you?",
        "What is this place?",
        "Who is Professor Foxwit?",
        "What's a Foxwit Door?",
        "Where are the MesoGears?",
        "What is the Krakian Jungle?",
        "What's a Gear Portal?",
        "What do the street signs mean?",
        "What's the deal with Jack Masque?",
        "Lita Lawless looks like a tough cookie, what's her story?",
        "When will new boroughs open up in the city?",
        "I want to take the quiz!"
    );

    if (selection === 0) {
        npc.sendOk("I've always dreamed of building a city. Not just any city, but one where everyone was welcome. I used to live in Kerning City, so I decided to see if I could create a city. As I went along in finding the means to do so, I encountered many people, some of whom I've come to regard as friends.");
    } else if (selection === 1) {
        npc.sendOk("Professor Foxwit is a pretty spry guy for being 97. He's a time-traveller I ran into outside the city one day, and he agreed to help build our museum.");
    } else if (selection === 2) {
        npc.sendOk("Foxwit Doors are warp points. Pressing Up will warp you to another location. I recommend getting the hang of them, they're our transport system.");
    } else if (selection === 3) {
        npc.sendOk("The MesoGears are beneath Bigger Ben. It's a monster-infested section of Bigger Ben that Barricade discovered. Be careful though, the Wolf Spiders in there are no joke.");
    } else if (selection === 4) {
        npc.sendOk("The Krakian Jungle is located on the outskirts of New Leaf City. Many new and powerful creatures roam there, so you'd better be prepared to fight if you head out.");
    } else if (selection === 5) {
        npc.sendOk("Gear Portals are old warp devices. They don't cycle through destinations like the Foxwit Door, but they can still get you around.");
    } else if (selection === 6) {
        npc.sendOk("You'll see street signs just about everywhere. Red lights mean an area isn't finished yet, but green lights mean it's open. Check back often, we're always building!");
    } else if (selection === 7) {
        npc.sendOk("Jack Masque is from Amoria. He's almost too smooth of a talker for his own good, and he hides behind that mask for reasons only he can fully explain.");
    } else if (selection === 8) {
        npc.sendOk("Lita is an old friend from Kerning City and one of the toughest thieves you'll ever meet. When it was time to pick a sheriff, it was a no-brainer.");
    } else if (selection === 9) {
        npc.sendOk("Soon, my friend. Even though you can't see them, the city developers are hard at work. When they're ready, we'll open them.");
    } else if (selection === 10) {
        if (plr.level() >= minLevel) {
            plr.startQuest(4900);
            npc.sendOk("No problem. I'll give you something nice if you answer the quiz correctly!");
        } else {
            npc.sendOk("Eager, are we? How about you explore a bit more before I let you take the quiz?");
        }
    }
}
