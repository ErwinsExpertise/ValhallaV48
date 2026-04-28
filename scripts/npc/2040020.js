function menu(prompt, options) {
    var text = prompt + "#b";
    for (var i = 0; i < options.length; i++) {
        text += "\r\n#L" + i + "# " + options[i] + "#l";
    }
    npc.sendSelection(text);
    return npc.selection();
}

function buildPrompt(recipe, stimId, stimulator) {
    var text = "You want me to make a #t" + recipe.item + "#? In that case, I'm going to need specific items from you in order to make it. Make sure you have room in your inventory, though!#b";
    if (stimulator) {
        text += "\r\n\r\n#rStimulator forging here keeps the original 10% destruction risk, but if it succeeds you'll receive the normal item rather than a random-stat reroll.#k";
        text += "\r\n#i" + stimId + "# 1 #t" + stimId + "#";
    }
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

function takeMaterials(recipe) {
    for (var i = 0; i < recipe.mats.length; i++) {
        plr.gainItem(recipe.mats[i], -recipe.qty[i]);
    }
}

var stimId = 4130000;
var type = menu("Hello, and welcome to the Ludibrium Glove Store. How can I help you today?", [
    "What's a stimulator?",
    "Create a Warrior glove",
    "Create a Bowman glove",
    "Create a Magician glove",
    "Create a Thief glove",
    "Create a Warrior glove with a Stimulator",
    "Create a Bowman glove with a Stimulator",
    "Create a Magician glove with a Stimulator",
    "Create a Thief glove with a Stimulator"
]);

if (type === 0) {
    npc.sendOk("A stimulator is a special potion that I can add into the process of creating certain items. It gives it stats as though it had dropped from a monster. However, it is possible to have no change, and it is also possible for the item to be below average. There's also a 10% chance of not getting any item when using a stimulator, so please choose wisely.");
} else {
    var stimulator = type > 4;
    var recipe;
    var selection;

    if (type === 1) {
        selection = menu("Warrior glove? Sure thing, which kind?", ["#t1082007#", "#t1082008#", "#t1082023#", "#t1082009#"]);
        recipe = [
            { item: 1082007, mats: [4011000, 4011001, 4003000], qty: [3, 2, 15], cost: 18000 },
            { item: 1082008, mats: [4000021, 4011001, 4003000], qty: [30, 4, 15], cost: 27000 },
            { item: 1082023, mats: [4000021, 4011001, 4003000], qty: [50, 5, 40], cost: 36000 },
            { item: 1082009, mats: [4011001, 4021007, 4000030, 4003000], qty: [3, 2, 30, 45], cost: 45000 }
        ][selection];
    } else if (type === 2) {
        selection = menu("Bowman glove? Sure thing, which kind?", ["#t1082048#", "#t1082068#", "#t1082071#", "#t1082084#"]);
        recipe = [
            { item: 1082048, mats: [4000021, 4011006, 4021001], qty: [50, 2, 1], cost: 18000 },
            { item: 1082068, mats: [4011000, 4011001, 4000021, 4003000], qty: [1, 3, 60, 15], cost: 27000 },
            { item: 1082071, mats: [4011001, 4021000, 4021002, 4000021, 4003000], qty: [3, 1, 3, 80, 25], cost: 36000 },
            { item: 1082084, mats: [4011004, 4011006, 4021002, 4000030, 4003000], qty: [3, 1, 2, 40, 35], cost: 45000 }
        ][selection];
    } else if (type === 3) {
        selection = menu("Magician glove? Sure thing, which kind?", ["#t1082051#", "#t1082054#", "#t1082062#", "#t1082081#"]);
        recipe = [
            { item: 1082051, mats: [4000021, 4021006, 4021000], qty: [60, 1, 2], cost: 22500 },
            { item: 1082054, mats: [4000021, 4011006, 4011001, 4021000], qty: [70, 1, 3, 2], cost: 27000 },
            { item: 1082062, mats: [4000021, 4021000, 4021006, 4003000], qty: [80, 3, 3, 30], cost: 36000 },
            { item: 1082081, mats: [4021000, 4011006, 4000030, 4003000], qty: [3, 2, 35, 40], cost: 45000 }
        ][selection];
    } else if (type === 4) {
        selection = menu("Thief glove? Sure thing, which kind?", ["#t1082042#", "#t1082046#", "#t1082075#", "#t1082065#"]);
        recipe = [
            { item: 1082042, mats: [4011001, 4000021, 4003000], qty: [2, 50, 10], cost: 22500 },
            { item: 1082046, mats: [4011001, 4011000, 4000021, 4003000], qty: [3, 1, 60, 15], cost: 27000 },
            { item: 1082075, mats: [4021000, 4000101, 4000021, 4003000], qty: [3, 100, 80, 30], cost: 36000 },
            { item: 1082065, mats: [4021005, 4021008, 4000030, 4003000], qty: [3, 1, 40, 30], cost: 45000 }
        ][selection];
    } else if (type === 5) {
        selection = menu("Warrior glove with a stimulator? Sure thing, which kind?", ["#t1082005#", "#t1082006#", "#t1082035#", "#t1082036#", "#t1082024#", "#t1082025#", "#t1082010#", "#t1082011#"]);
        recipe = [
            { item: 1082005, mats: [1082007, 4011001], qty: [1, 1], cost: 18000 },
            { item: 1082006, mats: [1082007, 4011005], qty: [1, 2], cost: 22500 },
            { item: 1082035, mats: [1082008, 4021006], qty: [1, 3], cost: 27000 },
            { item: 1082036, mats: [1082008, 4021008], qty: [1, 1], cost: 36000 },
            { item: 1082024, mats: [1082023, 4011003], qty: [1, 4], cost: 40500 },
            { item: 1082025, mats: [1082023, 4021008], qty: [1, 2], cost: 45000 },
            { item: 1082010, mats: [1082009, 4011002], qty: [1, 5], cost: 49500 },
            { item: 1082011, mats: [1082009, 4011006], qty: [1, 4], cost: 54000 }
        ][selection];
    } else if (type === 6) {
        selection = menu("Bowman glove with a stimulator? Sure thing, which kind?", ["#t1082049#", "#t1082050#", "#t1082069#", "#t1082070#", "#t1082072#", "#t1082073#", "#t1082085#", "#t1082083#"]);
        recipe = [
            { item: 1082049, mats: [1082048, 4021003], qty: [1, 3], cost: 13500 },
            { item: 1082050, mats: [1082048, 4021008], qty: [1, 1], cost: 18000 },
            { item: 1082069, mats: [1082068, 4011002], qty: [1, 4], cost: 19800 },
            { item: 1082070, mats: [1082068, 4011006], qty: [1, 2], cost: 22500 },
            { item: 1082072, mats: [1082071, 4011006], qty: [1, 4], cost: 27000 },
            { item: 1082073, mats: [1082071, 4021008], qty: [1, 2], cost: 36000 },
            { item: 1082085, mats: [1082084, 4011000, 4021000], qty: [1, 1, 5], cost: 49500 },
            { item: 1082083, mats: [1082084, 4011006, 4021008], qty: [1, 2, 2], cost: 54000 }
        ][selection];
    } else if (type === 7) {
        selection = menu("Magician glove with a stimulator? Sure thing, which kind?", ["#t1082052#", "#t1082053#", "#t1082055#", "#t1082056#", "#t1082063#", "#t1082064#", "#t1082082#", "#t1082080#"]);
        recipe = [
            { item: 1082052, mats: [1082051, 4021005], qty: [1, 3], cost: 31500 },
            { item: 1082053, mats: [1082051, 4021008], qty: [1, 1], cost: 36000 },
            { item: 1082055, mats: [1082054, 4021005], qty: [1, 3], cost: 36000 },
            { item: 1082056, mats: [1082054, 4021008], qty: [1, 1], cost: 40500 },
            { item: 1082063, mats: [1082062, 4021002], qty: [1, 4], cost: 40500 },
            { item: 1082064, mats: [1082062, 4021008], qty: [1, 2], cost: 45000 },
            { item: 1082082, mats: [1082081, 4021002], qty: [1, 5], cost: 49500 },
            { item: 1082080, mats: [1082081, 4021008], qty: [1, 3], cost: 54000 }
        ][selection];
    } else {
        selection = menu("Thief glove with a stimulator? Sure thing, which kind?", ["#t1082043#", "#t1082044#", "#t1082047#", "#t1082045#", "#t1082076#", "#t1082074#", "#t1082067#", "#t1082066#"]);
        recipe = [
            { item: 1082043, mats: [1082042, 4011004], qty: [1, 2], cost: 13500 },
            { item: 1082044, mats: [1082042, 4011006], qty: [1, 1], cost: 18000 },
            { item: 1082047, mats: [1082046, 4011005], qty: [1, 3], cost: 19800 },
            { item: 1082045, mats: [1082046, 4011006], qty: [1, 2], cost: 22500 },
            { item: 1082076, mats: [1082075, 4011006], qty: [1, 4], cost: 36000 },
            { item: 1082074, mats: [1082075, 4021008], qty: [1, 2], cost: 45000 },
            { item: 1082067, mats: [1082065, 4021000], qty: [1, 5], cost: 49500 },
            { item: 1082066, mats: [1082065, 4011006, 4021008], qty: [1, 2, 1], cost: 54000 }
        ][selection];
    }

    if (!npc.sendYesNo(buildPrompt(recipe, stimId, stimulator))) {
        npc.sendOk("All right. Come back when you're ready.");
    } else if (!plr.canHold(recipe.item, 1)) {
        npc.sendOk("Check your inventory for a free slot first.");
    } else if (plr.getMesos() < recipe.cost) {
        npc.sendOk("Sorry, we only accept meso.");
    } else if (!hasMaterials(recipe) || (stimulator && !plr.haveItem(stimId, 1))) {
        npc.sendOk("Sorry, but I have to have those items to get this exactly right. Perhaps next time.");
    } else {
        takeMaterials(recipe);
        plr.gainMesos(-recipe.cost);
        if (stimulator) {
            plr.gainItem(stimId, -1);
            if (Math.floor(Math.random() * 10) === 0) {
                npc.sendOk("Eek! I think I accidently added too much stimulator and, well, the whole thing is unusable now... Sorry, but I can't offer a refund.");
            } else {
                plr.gainItem(recipe.item, 1);
                npc.sendOk("There, the gloves are ready. Be careful, they're still hot. The stimulator kept its 10% destruction risk here, but successful crafts use the base item stats.");
            }
        } else {
            plr.gainItem(recipe.item, 1);
            npc.sendOk("There, the gloves are ready. Be careful, they're still hot.");
        }
    }
}
