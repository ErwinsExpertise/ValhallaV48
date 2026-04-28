npc.sendSelection("It seems that you managed to get to the top of the mission. Congratulations, #h #...\r\n\r\n\t#b#L0#Yes, I did it!#l");

if (npc.selection() === 0) {
    plr.warp(105040300);
}
