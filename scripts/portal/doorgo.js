var routes = {
    18: { reactor: "gate00", portal: "gt00PIA" },
    19: { reactor: "gate01", portal: "gt01PIA" },
    20: { reactor: "gate02", portal: "gt02PIA" },
    21: { reactor: "gate03", portal: "gt03PIA" },
    22: { reactor: "gate04", portal: "gt04PIA" },
    23: { reactor: "gate05", portal: "gt05PIA" },
    24: { reactor: "gate06", portal: "gt06PIA" }
};

var route = routes[portal.id()];
if (!route) {
    portal.block("This portal is not available now.");
} else if (map.reactorStateByName(route.reactor) === 4) {
    portal.warp(670010600, route.portal);
} else {
    portal.block("This portal is not available now.");
}
