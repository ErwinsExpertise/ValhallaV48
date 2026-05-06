var QUEST = 7500;
var DARK_CRYSTAL = 4005004;
var WISDOM_NECKLACE = 4031058;

var quizzes = [
    [
        {q: "Question 1. Which NPC would you NOT see in Henesys on Victoria Island?#b", o: ["Maya", "Teo", "Lucas", "Athena Pierce", "Chief Stan"], a: 1},
        {q: "Question 2. Which monster would you NOT see on Maple Island?#b", o: ["Snail", "Orange Mushroom", "Blue Snail", "Pig", "Slime"], a: 3},
        {q: "Question 3. Which item did Maya ask for to cure her illness?#b", o: ["Amazing Remedy", "Bad Remedy", "All-Cure", "Chinese Remedy", "#t4031006#"], a: 4},
        {q: "Question 4. Which of these places is NOT part of Victoria Island?#b", o: ["Sleepywood", "Amherst", "Perion", "Kerning City", "Ellinia"], a: 1},
        {q: "Question 5. Which monster would you NOT find in the Ant Tunnel?#b", o: ["Zombie Mushroom", "Stump", "Wild Boar", "Horny Mushroom", "Balrog"], a: 4}
    ],
    [
        {q: "Question 1. Which monster-drop pair is correct?#b", o: ["Fire Boar - Fire Boar Nose", "Cold Eye - Cold Eye Tail", "Pig - Pig Ear", "Jr. Necki - #t4000042#", "Zombie Mushroom - Zombie Mushroom Cap"], a: 3},
        {q: "Question 2. Which NPC would you NOT see in Perion?#b", o: ["Ayan", "Sophia", "Manji", "The Rememberer", "Dances with Balrog"], a: 1},
        {q: "Question 3. Which NPC in Kerning City is Alex's father?#b", o: ["Nella", "Irene", "Cody", "Pason", "John"], a: 4},
        {q: "Question 4. Which item do you receive after collecting 30 Dark Marbles for the 2nd job test?#b", o: ["#t4031012#", "Hero Necklace", "Hero Pendant", "Hero Medal", "Hero Mark"], a: 0},
        {q: "Question 5. Which first-job requirement pair is correct?#b", o: ["Warrior - STR 30+", "Magician - INT 25+", "Bowman - DEX 25+", "Thief - DEX 20+", "Thief - LUK 20+"], a: 2}
    ],
    [
        {q: "Question 1. Which NPC do you meet first in MapleStory?#b", o: ["Rain", "Biggs", "Sera", "Lith Harbor Chief", "Shanks"], a: 4},
        {q: "Question 2. How much EXP is needed to go from level 1 to 2?#b", o: ["10", "15", "20", "25", "30"], a: 1},
        {q: "Question 3. Which NPC would you NOT see in El Nath?#b", o: ["Chief's Residence Guide", "El Nath Chief", "Alcaster", "El Nath Resident", "El Nath Shopkeeper"], a: 2},
        {q: "Question 4. Which job can you NOT obtain as a 2nd job?#b", o: ["Page", "Assassin", "Bandit", "Cleric", "Mage"], a: 4},
        {q: "Question 5. Which quest can be repeated?#b", o: ["Maya and the Weird Medicine", "Alex the Runaway Kid", "Pia and the Blue Mushroom", "Arwen and the Glass Shoe", "Alpha Platoon Network"], a: 3}
    ],
    [
        {q: "Question 1. Which NPC is NOT part of Alpha Platoon?#b", o: ["Sergeant Peter", "Corporal An", "Private Allie", "Rolf", "Moppie"], a: 0},
        {q: "Question 2. Which item is NOT needed for Manji's Old Gladius quest?#b", o: ["#t4003002#", "#t4021009#", "#t4001006#", "#t4003003#", "#t4001005#"], a: 3},
        {q: "Question 3. Which NPC would you NOT see in Kerning City?#b", o: ["Jazz", "Vicky", "Maria", "Icarus", "Pia"], a: 4},
        {q: "Question 4. Which monster-drop pair is NOT correct?#b", o: ["Drake - #t4000059#", "Nependeath - Nependeath's Honey", "Jr. Pepe - #t4000040#", "Lunar Pixie - #t4000050#", "Lupin - #t4000029#"], a: 1},
        {q: "Question 5. Which of these monsters flies?#b", o: ["Jr. Cellion", "Nependeath", "Iron Hog", "Malady", "Stirge"], a: 4}
    ],
    [
        {q: "Question 1. Which status effect and result is NOT correct?#b", o: ["Darkness - lowers accuracy", "Curse - lowers EXP gain", "Weakness - lowers speed", "Seal - disables skills", "Poison - drains HP over time"], a: 2},
        {q: "Question 2. Which NPC would you NOT see in Orbis?#b", o: ["Spinel", "Dances with Balrog", "Nuri", "Dr. Kim", "Moppie"], a: 1},
        {q: "Question 3. Which quest requires the highest level to start?#b", o: ["Manji's Old Gladius", "Luke the Security Guy", "Searching for the Ancient Book", "Alcaster and the Dark Crystal", "Alpha Platoon Network"], a: 3},
        {q: "Question 4. Which NPC has nothing to do with refining, upgrading, or making items?#b", o: ["JM From the Streetz", "Mr. Thunder", "Ronnie", "Vogen", "Chrishrama"], a: 2},
        {q: "Question 5. Which potion-effect pair is correct?#b", o: ["#t2000001# - restores 200 HP", "#t2001001# - restores 2000 MP", "#t2010004# - restores 100 MP", "#t2020001# - restores 300 HP", "#t2020003# - restores 400 HP"], a: 4}
    ],
    [
        {q: "Question 1. Which NPC would you NOT see in Ellinia?#b", o: ["Vicious", "Rowen the Fairy", "Grendel the Really Old", "Arwen", "Maya"], a: 4},
        {q: "Question 2. Which monster would you NOT fight in Ossyria?#b", o: ["Jr. Balrog", "Dark Stone Golem", "Jr. Yeti", "Dark Yeti", "Stone Golem"], a: 1},
        {q: "Question 3. Which monster has the highest level?#b", o: ["Orange Mushroom", "Ribbon Pig", "Green Mushroom", "Stump", "Blue Mushroom"], a: 3},
        {q: "Question 4. Which potion-effect pair is NOT correct?#b", o: ["#t2050003# - cures curse or seal", "#t2020014# - restores 3000 MP", "#t2020004# - restores 400 HP", "#t2020000# - restores 200 MP", "#t2000003# - restores 100 MP"], a: 1},
        {q: "Question 5. Which NPC has nothing to do with pets?#b", o: ["John", "Pia", "Cloy", "Nella", "Maya"], a: 3}
    ]
];

function ask(question) {
    return npc.sendMenu.apply(npc, [question.q].concat(question.o));
}

function runQuiz(quiz) {
    for (var i = 0; i < quiz.length; i++) {
        if (ask(quiz[i]) !== quiz[i].a) {
            npc.sendOk("Wrong... Start over from the beginning.");
            return false;
        }
    }
    return true;
}

if (plr.questData(QUEST) !== "end1") {
    npc.sendOk("#b(A mysterious energy surrounds this stone. It feels bitterly cold.)");
} else if (!npc.sendYesNo("... ... ...\r\nIf you want to test your wisdom, you must offer #b#t4005004##k as a sacrifice.\r\nAre you ready to offer one and answer my questions?")) {
    npc.sendOk("Come back when you're prepared.");
} else if (plr.getEtcInventoryFreeSlot() < 1) {
    npc.sendOk("Your Etc inventory is full. Free up at least one slot first.");
} else if (plr.haveItem(WISDOM_NECKLACE, 1)) {
    npc.sendOk("You already have #b#t4031058##k. Bring it back to your chief.");
} else if (!plr.haveItem(DARK_CRYSTAL, 1)) {
    npc.sendOk("If you want to test your wisdom, you must offer #b#t4005004##k as a sacrifice.");
} else {
    plr.gainItem(DARK_CRYSTAL, -1);
    npc.sendOk("Very well. I will test your wisdom now. Answer every question correctly. If you fail even once, you must start again from the beginning.");
    var quiz = quizzes[Math.floor(Math.random() * quizzes.length)];
    if (runQuiz(quiz)) {
        if (!plr.gainItem(WISDOM_NECKLACE, 1)) {
            npc.sendOk("Your Etc inventory is full.");
        } else {
            npc.sendOk("Excellent. Your wisdom has been proven. Take this necklace back to your chief.");
        }
    }
}
