npc.sendSelection("What would you like to do?\r\n#L0#Exit the Guild Quest#l");
if (npc.selection() === 0) {
    if (!npc.sendYesNo("Are you sure you want to leave? You will not be able to return.")) {
        npc.sendOk("Good luck on finishing the Guild Quest!");
    } else {
        plr.leaveEvent();
    }
}
