if (plr.hasStoreBankItems()) {
    npc.sendStoreBank(9030000);
} else {
    npc.sendOk("I'm Fredrick, the Store Banker. If your hired merchant closes, I keep your unsold items and mesos safe until you come reclaim them.");
}
