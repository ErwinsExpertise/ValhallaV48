function buildPrompt(recipe) {
    var text = "So we are going for a #t" + recipe.item + "#, right? In that case, I'm going to need specific items from you in order to make it. Make sure you have room in your inventory, though!#b";
    for (var i = 0; i < recipe.mats.length; i++) {
        text += "\r\n#i" + recipe.mats[i] + "# " + recipe.qty[i] + " #t" + recipe.mats[i] + "#";
    }
    text += "\r\n#i4031138# " + recipe.cost + " meso";
    return text;
}

function hasMaterials(recipe) {
    for (var i = 0; i < recipe.mats.length; i++) {
        if (!plr.haveItem(recipe.mats[i], recipe.qty[i])) {
            return false;
        }
    }
    return true;
}

var recipes = [
    { item: 4080100, mats: [4030012], qty: [99], cost: 10000 },
    { item: 4080006, mats: [4030009, 4030013, 4030014], qty: [1, 99, 99], cost: 25000 },
    { item: 4080007, mats: [4030009, 4030013, 4030016], qty: [1, 99, 99], cost: 25000 },
    { item: 4080008, mats: [4030009, 4030014, 4030016], qty: [1, 99, 99], cost: 25000 },
    { item: 4080009, mats: [4030009, 4030015, 4030013], qty: [1, 99, 99], cost: 25000 },
    { item: 4080010, mats: [4030009, 4030015, 4030014], qty: [1, 99, 99], cost: 25000 },
    { item: 4080011, mats: [4030009, 4030015, 4030016], qty: [1, 99, 99], cost: 25000 }
];

var sel = npc.sendMenu("Hey there! My name is #p2040014#, and I am a specialist in mini-games. What kind of mini-game you want me to make? #b", "#i4080100# #t4080100#", "#i4080006# #t4080006#", "#i4080007# #t4080007#", "#i4080008# #t4080008#", "#i4080009# #t4080009#", "#i4080010# #t4080010#", "#i4080011# #t4080011#");
var recipe = recipes[sel];

if (!npc.sendYesNo(buildPrompt(recipe))) {
    npc.sendOk("Come back if you want another set.");
} else if (plr.getMesos() < recipe.cost) {
    npc.sendOk("See, I need to specify my wages to support my career, that cannot be bypassed. I will gladly help you once you've got the money.");
} else if (!hasMaterials(recipe)) {
    npc.sendOk("You are lacking some items for the set you want to make. Please provide them so that we can assemble the game set.");
} else if (!plr.canHold(recipe.item, 1)) {
    npc.sendOk("I can't make a set for you if there's no room in your ETC inventory for it. Please free a space first and then talk to me.");
} else {
    for (var i = 0; i < recipe.mats.length; i++) {
        plr.gainItem(recipe.mats[i], -recipe.qty[i]);
    }
    plr.gainMesos(-recipe.cost);
    plr.gainItem(recipe.item, 1);
    npc.sendOk("There is your game set. Have fun!");
}
