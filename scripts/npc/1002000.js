var travelMaps = [102000000, 101000000, 100000000, 103000000];
var travelCosts = [1200, 1200, 800, 1000];
var townDescriptions = [
    [
        "The town you are in is Lith Harbor! Alright, I'll tell you more about #bLith Harbor#k. This is where you arrive on Victoria Island after riding The Victoria. Many beginners who leave Maple Island start their journey here.",
        "It's a quiet town with a wide body of water behind it, since the harbor is located on the west side of the island. Most of the people here are, or used to be, fishermen, so they may look intimidating, but if you strike up a conversation with them, they will be friendly.",
        "Around town lies a beautiful prairie. Most of the monsters there are small and gentle, perfect for beginners. If you haven't chosen your job yet, this is a good place to boost up your level."
    ],
    [
        "Alright I'll explain to you more about #bPerion#k. It's a warrior-town located at the northern-most part of Victoria Island, surrounded by rocky mountains. With an unfriendly atmosphere, only the strong survives there.",
        "Around the highland you'll find a really skinny tree, a wild hog running around the place, and monkeys that live all over the island. There's also a deep valley, and when you go deep into it, you'll find a humongous dragon with the power to match his size. Better go in there very carefully, or don't go at all.",
        "If you want to be a #bWarrior#k then find #rDances with Balrog#k, the chief of Perion. If you're level 10 or higher, along with a good STR level, he may make you a warrior after all. If not, better keep training yourself until you reach that level."
    ],
    [
        "Alright I'll explain to you more about #bEllinia#k. It's a magician-town located at the far east of Victoria Island, and covered in tall, mystic trees. You'll find some fairies there, too. They don't like humans in general so it'll be best for you to be on their good side and stay quiet.",
        "Near the forest you'll find green slimes, walking mushrooms, monkeys and zombie monkeys all residing there. Walk deeper into the forest and you'll find witches with the flying broomstick navigating the skies. A word of warning: Unless you are really strong, I recommend you don't go near them.",
        "If you want to be a #bMagician#k, search for #rGrendel the Really Old#k, the head wizard of Ellinia. He may make you a wizard if you're at or above level 8 with a decent amount of INT. If that's not the case, you may have to hunt more and train yourself to get there."
    ],
    [
        "Alright I'll explain to you more about #bHenesys#k. It's a bowman-town located at the southernmost part of the island, made on a flatland in the midst of a deep forest and prairies. The weather's just right, and everything is plentiful around that town, perfect for living. Go check it out.",
        "Around the prairie you'll find weak monsters such as snails, mushrooms, and pigs. According to what I hear, though, in the deepest part of the Pig Park, which is connected to the town somewhere, you'll find a humongous, powerful mushroom called Mushmom every now and then.",
        "If you want to be a #bBowman#k, you need to go see #rAthena Pierce#k at Henesys. With a level of 10 or higher and a decent amount of DEX, she may make you one after all. If not, go train, make yourself stronger, and then try again."
    ],
    [
        "Alright I'll explain to you more about #bKerning City#k. It's a thief-town located at the northwest part of Victoria Island, and there are buildings up there that have just this strange feeling around them. It's mostly covered in black clouds, but if you can go up to a really high place, you'll be able to see a very beautiful sunset there.",
        "From Kerning City, you can go into several dungeons. You can go to a swamp where alligators and snakes are abound, or hit the subway full of ghosts and bats. At the deepest part of the underground, you'll find Lace, who is just as big and dangerous as a dragon.",
        "If you want to be a #bThief#k, seek #rDark Lord#k, the heart of darkness of Kerning City. He may well make you a thief if you're at or above level 10 with a good amount of DEX. If not, go hunt and train yourself to reach there."
    ],
    [
        "Alright I'll explain to you more about #bSleepywood#k. It's a forest town located at the southeast side of Victoria Island. It's pretty much in between Henesys and the ant-tunnel dungeon. There's a hotel there, so you can rest up after a long day at the dungeon ... it's a quiet town in general.",
        "In front of the hotel there's an old buddhist monk by the name of #rChrishrama#k. Nobody knows a thing about that monk. Apparently he collects materials from the travelers and create something, but I am not too sure about the details. If you have any business going around that area, please check that out for me.",
        "From Sleepywood, head east and you'll find the ant tunnel connected to the deepest part of the Victoria Island. Lots of nasty, powerful monsters abound so if you walk in thinking it's a walk in the park, you'll be coming out as a corpse. You need to fully prepare yourself for a rough ride before going in.",
        "And this is what I hear ... apparently, at Sleepywood there's a secret entrance leading you to an unknown place. Apparently, once you move in deep, you'll find a stack of black rocks that actually move around. I want to see that for myself in the near future ..."
    ]
];

npc.sendNext("Do you wanna head over to some other town? With a little money involved, I can make it happen. It's a tad expensive, but I run a special 90% discount for beginners.");

var firstChoice = npc.sendMenu("It's understandable that you may be confused about this place if this is your first time around. If you got any questions about this place, fire away.\r\n#L0##bWhat kind of towns are here in Victoria Island?#l\r\n#L1#Please take me somewhere else.#k#l");

if (firstChoice === 0) {
    var townChoice = npc.sendMenu("There are 6 big towns here in Victoria Island. Which of those do you want to know more of?\r\n#b#L0##m104000000##l\r\n#L1##m102000000##l\r\n#L2##m101000000##l\r\n#L3##m100000000##l\r\n#L4##m103000000##l\r\n#L5##m105040300##l");
    if (townChoice >= 0 && townChoice < townDescriptions.length) {
        var pages = townDescriptions[townChoice];
        for (var i = 0; i < pages.length; i++) {
            npc.sendNext(pages[i]);
        }
    }
} else if (firstChoice === 1) {
    var beginner = plr.job() === 0;
    var travelText = beginner
        ? "There's a special 90% discount for all beginners. Alright, where would you want to go?#b"
        : "Oh you aren't a beginner, huh? Then I'm afraid I may have to charge you full price. Where would you like to go?#b";

    for (var j = 0; j < travelMaps.length; j++) {
        var price = beginner ? Math.floor(travelCosts[j] * 0.10) : travelCosts[j];
        travelText += "\r\n#L" + j + "##m" + travelMaps[j] + "# (" + price + " mesos)#l";
    }

    var destination = npc.sendMenu(travelText);
    if (destination >= 0 && destination < travelMaps.length) {
        var finalCost = beginner ? Math.floor(travelCosts[destination] * 0.10) : travelCosts[destination];
        if (!npc.sendYesNo("I guess you don't need to be here. Do you really want to move to #b#m" + travelMaps[destination] + "##k? Well it'll cost you #b" + finalCost + " mesos#k. What do you think?")) {
            npc.sendOk("There's a lot to see in this town, too. Let me know if you want to go somewhere else.");
        } else if (plr.getMesos() < finalCost) {
            npc.sendOk("You don't have enough mesos. With your abilities, you should have more than that!");
        } else {
            plr.gainMesos(-finalCost);
            plr.warp(travelMaps[destination]);
        }
    }
}
