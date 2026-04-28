package channel

import "log"

type transportDestination struct {
	mapID      int32
	portalName string
}

func moveTransportPlayers(server *Server, warps map[int32]transportDestination) {
	for src, dst := range warps {
		srcField, ok := server.fields[src]
		if !ok {
			log.Println("Could not move transport players from", src)
			continue
		}

		dstField, ok := server.fields[dst.mapID]
		if !ok {
			log.Println("Could not move transport players to", dst.mapID)
			continue
		}

		for _, srcInst := range srcField.instances {
			dstInst, err := dstField.getInstance(srcInst.id)
			if err != nil {
				dstInst, err = dstField.getInstance(0)
				if err != nil {
					log.Println("Could not find destination instance for", dst.mapID)
					continue
				}
			}

			var portal portal

			if dst.portalName != "" {
				portal, err = dstInst.getPortalFromName(dst.portalName)
			} else {
				portal, err = dstInst.getPortalFromID(0, true)
			}

			if err != nil {
				log.Println("Could not find transport portal for", dst.mapID)
				continue
			}

			for _, plr := range srcInst.players {
				server.warpPlayer(plr, dstField, portal, true)
			}
		}
	}
}
