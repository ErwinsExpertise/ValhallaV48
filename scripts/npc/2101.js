if (!npc.sendYesNo("Are you done with your training? If you wish, I will send you out from this training camp.")) {
    npc.sendOk("Haven't you finished the training program yet? If you want to leave this place, please do not hesitate to tell me.");
} else {
    plr.warp(3);
}
