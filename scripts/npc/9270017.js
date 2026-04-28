if (!npc.sendYesNo("The plane will be taking off soon. Will you leave now? You will have to buy the plane ticket again to come in here.")) {
    npc.sendOk("Please hold on for a sec, and the plane will be taking off. Thanks for your patience.");
} else {
    npc.sendOk("I already told you the ticket is not refundable. Hope to see you again~");
    plr.warp(103000000);
}
