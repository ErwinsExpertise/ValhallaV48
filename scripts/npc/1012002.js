function menu(prompt, options) {
    var text = prompt + "#b";
    for (var i = 0; i < options.length; i++) {
        text += "\r\n#L" + i + "# " + options[i] + "#l";
    }
    npc.sendSelection(text);
    return npc.selection();
}

function recipePrompt(recipe, qty) {
    var text = "You want me to make " + (qty === 1 ? "a #t" + recipe.item + "#?" : qty + " #t" + recipe.item + "#?") + " In that case, I'm going to need specific items from you in order to make it. Make sure you have room in your inventory, though!#b";
    for (var i = 0; i < recipe.mats.length; i++) {
        text += "\r\n#i" + recipe.mats[i] + "# " + (recipe.qty[i] * qty) + " #t" + recipe.mats[i] + "#";
    }
    if (recipe.cost > 0) {
        text += "\r\n#i4031138# " + (recipe.cost * qty) + " meso";
    }
    return text;
}

function outputQty(recipe, qty) {
    return recipe.outputQty ? recipe.outputQty * qty : qty;
}

function hasMaterials(recipe, qty) {
    for (var i = 0; i < recipe.mats.length; i++) {
        if (!plr.haveItem(recipe.mats[i], recipe.qty[i] * qty)) {
            return false;
        }
    }
    return true;
}

function takeMaterials(recipe, qty) {
    for (var i = 0; i < recipe.mats.length; i++) {
        plr.gainItem(recipe.mats[i], -(recipe.qty[i] * qty));
    }
}

var type = menu("Hello. I am Vicious, retired Sniper. However, I used to be the top student of Athena Pierce. Though I no longer hunt, I can make some archer items that will be useful for you...", [
    "Create a bow",
    "Create a crossbow",
    "Make a glove",
    "Upgrade a glove",
    "Create materials",
    "Create Arrows"
]);

var recipe;
var qty = 1;

if (type === 0) {
    var items0 = [1452002, 1452003, 1452001, 1452000, 1452005, 1452006, 1452007];
    var sel0 = menu("I may have been a Sniper, but bows and crossbows aren't too much different. Anyway, which would you like to make?", [
        "#t1452002##k - Bowman Lv. 10",
        "#t1452003##k - Bowman Lv. 15",
        "#t1452001##k - Bowman Lv. 20",
        "#t1452000##k - Bowman Lv. 25",
        "#t1452005##k - Bowman Lv. 30",
        "#t1452006##k - Bowman Lv. 35",
        "#t1452007##k - Bowman Lv. 40"
    ]);
    recipe = {
        item: items0[sel0],
        mats: [[4003001, 4000000], [4011001, 4003000], [4003001, 4000016], [4011001, 4021006, 4003000], [4011001, 4011006, 4021003, 4021006, 4003000], [4011004, 4021000, 4021004, 4003000], [4021008, 4011001, 4011006, 4003000, 4000014]][sel0],
        qty: [[5, 30], [1, 3], [30, 50], [2, 2, 8], [5, 5, 3, 3, 30], [7, 6, 3, 35], [1, 10, 3, 40, 50]][sel0],
        cost: [800, 2000, 3000, 5000, 30000, 40000, 80000][sel0]
    };
} else if (type === 1) {
    var items1 = [1462001, 1462002, 1462003, 1462000, 1462004, 1462005, 1462006, 1462007];
    var sel1 = menu("I was a Sniper. Crossbows are my specialty. Which would you like me to make for you?", [
        "#t1462001##k - Bowman Lv. 10",
        "#t1462002##k - Bowman Lv. 15",
        "#t1462003##k - Bowman Lv. 20",
        "#t1462000##k - Bowman Lv. 25",
        "#t1462004##k - Bowman Lv. 30",
        "#t1462005##k - Bowman Lv. 35",
        "#t1462006##k - Bowman Lv. 40",
        "#t1462007##k - Bowman Lv. 45"
    ]);
    recipe = {
        item: items1[sel1],
        mats: [[4003001, 4003000], [4011001, 4003001, 4003000], [4011001, 4003001, 4003000], [4011001, 4021006, 4021002, 4003000], [4011001, 4011005, 4021006, 4003001, 4003000], [4021008, 4011001, 4011006, 4021006, 4003000], [4021008, 4011004, 4003001, 4003000], [4021008, 4011006, 4021006, 4003001, 4003000]][sel1],
        qty: [[7, 2], [1, 20, 5], [1, 50, 8], [2, 1, 1, 10], [5, 5, 3, 50, 15], [1, 8, 4, 2, 30], [2, 6, 30, 30], [2, 5, 3, 40, 40]][sel1],
        cost: [1000, 2000, 3000, 10000, 30000, 50000, 80000, 200000][sel1]
    };
} else if (type === 2) {
    var items2 = [1082012, 1082013, 1082016, 1082048, 1082068, 1082071, 1082084, 1082089];
    var sel2 = menu("Okay, so which glove do you want me to make?", [
        "#t1082012#",
        "#t1082013#",
        "#t1082016#",
        "#t1082048#",
        "#t1082068#",
        "#t1082071#",
        "#t1082084#",
        "#t1082089#"
    ]);
    recipe = {
        item: items2[sel2],
        mats: [[4000021, 4000009], [4000021, 4000009, 4011001], [4000021, 4000009, 4011006], [4000021, 4011006, 4021001], [4011000, 4011001, 4000021, 4003000], [4011001, 4021000, 4021002, 4000021, 4003000], [4011004, 4011006, 4021002, 4000030, 4003000], [4011006, 4011007, 4021006, 4000030, 4003000]][sel2],
        qty: [[15, 20], [20, 20, 2], [40, 50, 2], [50, 2, 1], [1, 3, 60, 15], [3, 1, 3, 80, 25], [3, 1, 2, 40, 35], [2, 1, 8, 50, 50]][sel2],
        cost: [5000, 10000, 15000, 20000, 30000, 40000, 50000, 70000][sel2]
    };
} else if (type === 3) {
    var items3 = [1082015, 1082014, 1082017, 1082018, 1082049, 1082050, 1082069, 1082070, 1082072, 1082073, 1082085, 1082083, 1082090, 1082091];
    var sel3 = menu("Upgrade a glove? That shouldn't be too difficult. Which did you have in mind?", [
        "#t1082015#",
        "#t1082014#",
        "#t1082017#",
        "#t1082018#",
        "#t1082049#",
        "#t1082050#",
        "#t1082069#",
        "#t1082070#",
        "#t1082072#",
        "#t1082073#",
        "#t1082085#",
        "#t1082083#",
        "#t1082090#",
        "#t1082091#"
    ]);
    recipe = {
        item: items3[sel3],
        mats: [[1082013, 4021003], [1082013, 4021000], [1082016, 4021000], [1082016, 4021008], [1082048, 4021003], [1082048, 4021008], [1082068, 4011002], [1082068, 4011006], [1082071, 4011006], [1082071, 4021008], [1082084, 4011000, 4021000], [1082084, 4011006, 4021008], [1082089, 4021000, 4021007], [1082089, 4021007, 4021008]][sel3],
        qty: [[1, 2], [1, 1], [1, 3], [1, 1], [1, 3], [1, 1], [1, 4], [1, 2], [1, 4], [1, 2], [1, 1, 5], [1, 2, 2], [1, 5, 1], [1, 2, 2]][sel3],
        cost: [7000, 7000, 10000, 12000, 15000, 20000, 22000, 25000, 30000, 40000, 55000, 60000, 70000, 80000][sel3]
    };
} else if (type === 4) {
    var sel4 = menu("Materials? I know of a few materials that I can make for you...", [
        "Make Processed Wood with Tree Branch",
        "Make Processed Wood with Firewood",
        "Make Screws (packs of 15)"
    ]);
    recipe = [
        { item: 4003001, mats: [4000003], qty: [10], cost: 0 },
        { item: 4003001, mats: [4000018], qty: [5], cost: 0 },
        { item: 4003000, mats: [4011000, 4011001], qty: [1, 1], cost: 0, outputQty: 15 }
    ][sel4];
    qty = npc.sendNumber("So, you want me to make some #t" + recipe.item + "#s? In that case, how many do you want me to make?", 1, 1, 100);
} else {
    var sel5 = menu("Arrows? Not a problem at all.", [
        "#t2060000#",
        "#t2061000#",
        "#t2060001#",
        "#t2061001#",
        "#t2060002#",
        "#t2061002#"
    ]);
    recipe = [
        { item: 2060000, mats: [4003001, 4003004], qty: [1, 1], cost: 0, outputQty: 1000 },
        { item: 2061000, mats: [4003001, 4003004], qty: [1, 1], cost: 0, outputQty: 1000 },
        { item: 2060001, mats: [4011000, 4003001, 4003004], qty: [1, 3, 10], cost: 0, outputQty: 900 },
        { item: 2061001, mats: [4011000, 4003001, 4003004], qty: [1, 3, 10], cost: 0, outputQty: 900 },
        { item: 2060002, mats: [4011001, 4003001, 4003005], qty: [1, 5, 15], cost: 0, outputQty: 800 },
        { item: 2061002, mats: [4011001, 4003001, 4003005], qty: [1, 5, 15], cost: 0, outputQty: 800 }
    ][sel5];
}

var finalQty = outputQty(recipe, qty);

if (!npc.sendYesNo(recipePrompt(recipe, qty))) {
    npc.sendOk("All right. Come back when you're ready.");
} else if (plr.getMesos() < recipe.cost * qty) {
    npc.sendOk("Sorry, but this is how I make my living. No meso, no item.");
} else if (!hasMaterials(recipe, qty)) {
    npc.sendOk("Surely you, of all people, would understand the value of having quality items? I can't do that without the items I require.");
} else if (!plr.canHold(recipe.item, finalQty)) {
    npc.sendOk("Please make sure you have room in your inventory, and talk to me again.");
} else {
    takeMaterials(recipe, qty);
    if (recipe.cost > 0) {
        plr.gainMesos(-(recipe.cost * qty));
    }
    plr.gainItem(recipe.item, finalQty);
    npc.sendOk("A perfect item, as usual. Come and see me if you need anything else.");
}
