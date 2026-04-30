if (plr.gender() === 0) {
    if (portal.id() === 15) portal.warp(670010200, "ma01");
    else portal.block("This portal is only available for ladies.");
} else {
    if (portal.id() === 16) portal.warp(670010200, "fe01");
    else portal.block("This portal is only available for men.");
}
