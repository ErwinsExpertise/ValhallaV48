package channel

import (
	"time"

	"github.com/Hucaru/Valhalla/constant"
)

const muLungTransportRideDuration = time.Minute

func scheduleMuLungTransport(server *Server) {
	wait := func(duration time.Duration) {
		timer := time.NewTimer(duration)
		<-timer.C
		timer.Stop()
	}

	arrivals := map[int32]transportDestination{
		constant.MapTransportToMuLung: {mapID: constant.MapMuLungArrival},
		constant.MapTransportToOrbis:  {mapID: constant.MapOrbisArrival},
	}

	for {
		wait(muLungTransportRideDuration)

		server.dispatch <- func() {
			moveTransportPlayers(server, arrivals)
		}
	}
}
