var SUZY_DECISION_KEY = 4994;
var MAPLEMAS_CHOICE_KEY = 4997;
var VERSALMAS_CHOICE_KEY = 4998;
var SUZY_NOTE_ITEM = 4031877;

var state = plr.questData(SUZY_DECISION_KEY);

if (state === "2") {
    npc.sendOk("I already made up my mind! I'm excited for the holiday I picked now.");
} else {
    var choice = npc.askMenu(
        "Excuse me! I'm trying to make a very important decision. Which holiday sounds better to you?#b",
        "Maplemas with Maple Claws",
        "Versalmas with O-Pongo",
        "I'm not sure yet"
    );

    if (choice === 0) {
        plr.setQuestData(SUZY_DECISION_KEY, "2");
        plr.setQuestData(MAPLEMAS_CHOICE_KEY, "2");
        plr.setQuestData(VERSALMAS_CHOICE_KEY, "");
        plr.gainItem(SUZY_NOTE_ITEM, -plr.itemCount(SUZY_NOTE_ITEM));
        npc.sendOk("Maplemas it is! I'll go tell everyone I chose Maplemas. If you want bigger Maplemas rewards, Maple Claws should help you now.");
    } else if (choice === 1) {
        plr.setQuestData(SUZY_DECISION_KEY, "2");
        plr.setQuestData(VERSALMAS_CHOICE_KEY, "2");
        plr.setQuestData(MAPLEMAS_CHOICE_KEY, "");
        plr.gainItem(SUZY_NOTE_ITEM, -plr.itemCount(SUZY_NOTE_ITEM));
        npc.sendOk("Versalmas it is! I'll celebrate with O-Pongo. If you want the bigger Versalmas rewards, go speak with him now.");
    } else {
        if (!plr.haveItem(SUZY_NOTE_ITEM, 1) && !plr.canHold(SUZY_NOTE_ITEM, 1)) {
            npc.sendOk("Come back when you have room in your Etc inventory. I want to give you my little holiday note.");
        } else {
            if (!plr.haveItem(SUZY_NOTE_ITEM, 1)) {
                plr.gainItem(SUZY_NOTE_ITEM, 1);
            }
            npc.sendOk("I'm still thinking... Here, keep my note in case you need to remind me later.");
        }
    }
}
