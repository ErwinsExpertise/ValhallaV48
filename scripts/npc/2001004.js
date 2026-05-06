if (npc.sendYesNo("Finished decorating the tree? If you leave now, I will send you back to Happyville.")) {
    plr.warp(209000000);
} else {
    npc.sendOk("Take your time. Speak to me again when you are ready to go back.");
}
