package channel

import (
	"math/rand"
	"time"

	"github.com/Hucaru/Valhalla/constant"
)

type sharedTransportRoute struct {
	name             string
	cycleMinutes     int
	departureOffset  int
	transitMinutes   int
	boardingOffset   int
	boardingDuration int
	boardingMaps     []int32
	waitingMap       int32
	transitMaps      []int32
	arrival          transportDestination
	boatVisual       bool
	runArrival       func(*Server)
	runDeparture     func(*Server, *rand.Rand)
	boardingOpen     bool
}

type transportScheduler struct {
	server                 *Server
	routes                 []*sharedTransportRoute
	muLungTimersByPlayerID map[int32]*time.Timer
	rng                    *rand.Rand
}

func newTransportScheduler(server *Server) *transportScheduler {
	return &transportScheduler{
		server: server,
		routes: []*sharedTransportRoute{
			// BMS DataSvr Continent.img
			// 0: Ellinia -> Orbis, 1: Orbis -> Ellinia
			{
				name:             "Ellinia -> Orbis",
				cycleMinutes:     15,
				departureOffset:  0,
				transitMinutes:   10,
				boardingOffset:   10,
				boardingDuration: 4,
				boardingMaps:     []int32{constant.MapStationEllinia},
				waitingMap:       constant.MapBoatElliniaDeparture,
				transitMaps:      []int32{constant.MapBoatElliniaFlight, constant.MapBoatElliniaFlightCabin},
				arrival:          transportDestination{mapID: constant.MapStationOrbis},
				boatVisual:       true,
				runArrival: func(server *Server) {
					checkInvasion(server, true)
				},
			},
			{
				name:             "Orbis -> Ellinia",
				cycleMinutes:     15,
				departureOffset:  0,
				transitMinutes:   10,
				boardingOffset:   10,
				boardingDuration: 4,
				boardingMaps:     []int32{constant.MapStationOrbisEllinaPlatform},
				waitingMap:       constant.MapBoatOrbisElliniaDeparture,
				transitMaps:      []int32{constant.MapBoatOrbisElliniaFlight, constant.MapBoatOrbisElliniaFlightCabin},
				arrival:          transportDestination{mapID: constant.MapStationEllinia},
				boatVisual:       true,
				runArrival: func(server *Server) {
					checkInvasion(server, true)
				},
				runDeparture: maybeStartElliniaBoatInvasion,
			},
			// 2: Orbis -> Ludibrium, 3: Ludibrium -> Orbis
			{
				name:             "Orbis -> Ludibrium",
				cycleMinutes:     10,
				departureOffset:  0,
				transitMinutes:   5,
				boardingOffset:   5,
				boardingDuration: 4,
				boardingMaps:     []int32{constant.MapStationOrbisLudiPlatform},
				waitingMap:       constant.MapBoatOrbisLudiDeparture,
				transitMaps:      []int32{constant.MapBoatOrbisLudiFlight},
				arrival:          transportDestination{mapID: constant.MapStationLudi},
				boatVisual:       true,
			},
			{
				name:             "Ludibrium -> Orbis",
				cycleMinutes:     10,
				departureOffset:  0,
				transitMinutes:   5,
				boardingOffset:   5,
				boardingDuration: 4,
				boardingMaps:     []int32{constant.MapStationLudiOrbisPlatform},
				waitingMap:       constant.MapBoatLudiDeparture,
				transitMaps:      []int32{constant.MapBoatLudiFlight},
				arrival:          transportDestination{mapID: constant.MapStationOrbis},
				boatVisual:       true,
			},
			// 4: Helios elevator 99F -> 2F, 5: Helios elevator 2F -> 99F
			{
				name:             "Helios 99F -> 2F",
				cycleMinutes:     4,
				departureOffset:  0,
				transitMinutes:   1,
				boardingOffset:   3,
				boardingDuration: 1,
				boardingMaps:     []int32{constant.MapHeliosTower99thFloor},
				waitingMap:       constant.MapHeliosTowerKFTWaitingRoom,
				transitMaps:      []int32{constant.MapHeliosTowerKFTElevator},
				arrival:          transportDestination{mapID: constant.MapHeliosTower2ndFloor, portalName: "in00"},
			},
			{
				name:             "Helios 2F -> 99F",
				cycleMinutes:     4,
				departureOffset:  2,
				transitMinutes:   1,
				boardingOffset:   1,
				boardingDuration: 1,
				boardingMaps:     []int32{constant.MapHeliosTower2ndFloor},
				waitingMap:       constant.MapHeliosTowerLudiWaitingRoom,
				transitMaps:      []int32{constant.MapHeliosTowerLudiElevator},
				arrival:          transportDestination{mapID: constant.MapHeliosTower99thFloor, portalName: "in00"},
			},
			// 8: NLC -> Kerning, 9: Kerning -> NLC
			{
				name:             "NLC -> Kerning",
				cycleMinutes:     5,
				departureOffset:  4,
				transitMinutes:   1,
				boardingOffset:   1,
				boardingDuration: 3,
				boardingMaps:     []int32{constant.MapNLCSubwayStation},
				waitingMap:       constant.MapNLCToKerningWaitingRoom,
				transitMaps:      []int32{constant.MapNLCToKerningTrain},
				arrival:          transportDestination{mapID: constant.MapKerningSubwayStation},
			},
			{
				name:             "Kerning -> NLC",
				cycleMinutes:     5,
				departureOffset:  4,
				transitMinutes:   1,
				boardingOffset:   1,
				boardingDuration: 3,
				boardingMaps:     []int32{constant.MapKerningSubwayStation},
				waitingMap:       constant.MapKerningToNLCWaitingRoom,
				transitMaps:      []int32{constant.MapKerningToNLCTrain},
				arrival:          transportDestination{mapID: constant.MapNLCSubwayStation},
			},
		},
		muLungTimersByPlayerID: make(map[int32]*time.Timer),
		rng:                    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (ts *transportScheduler) run() {
	ts.server.post(func() {
		ts.sync(time.Now(), false)
	})

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	lastMinuteStamp := time.Now().Unix() / 60

	for now := range ticker.C {
		minuteStamp := now.Unix() / 60
		if minuteStamp == lastMinuteStamp {
			continue
		}
		lastMinuteStamp = minuteStamp

		ts.server.post(func() {
			ts.sync(now, true)
		})
	}
}

func (ts *transportScheduler) sync(now time.Time, processTransitions bool) {
	minuteOfHour := now.Minute()

	for _, route := range ts.routes {
		boardingOpen := route.isBoardingOpen(minuteOfHour)
		if route.boardingOpen != boardingOpen {
			route.boardingOpen = boardingOpen
			ts.setBoardingState(route, boardingOpen)
		}

		if !processTransitions {
			if !route.isInTransit(minuteOfHour) {
				ts.arriveRoute(route)
			}
			continue
		}

		if route.isArrivalMinute(minuteOfHour) {
			ts.arriveRoute(route)
			if route.runArrival != nil {
				route.runArrival(ts.server)
			}
		}

		if route.isDepartureMinute(minuteOfHour) {
			ts.departRoute(route)
			if route.runDeparture != nil {
				route.runDeparture(ts.server, ts.rng)
			}
		}
	}
}

func (ts *transportScheduler) setBoardingState(route *sharedTransportRoute, open bool) {
	for _, mapID := range route.boardingMaps {
		field, ok := ts.server.fields[mapID]
		if !ok {
			continue
		}

		for _, inst := range field.instances {
			inst.properties["canBoard"] = open
			if route.boatVisual {
				inst.showBoats(open, 0x00)
			}
		}
	}
}

func (ts *transportScheduler) departRoute(route *sharedTransportRoute) {
	moveTransportPlayers(ts.server, map[int32]transportDestination{
		route.waitingMap: {mapID: route.transitMaps[0]},
	})
	if route.boardingOpen {
		route.boardingOpen = false
		ts.setBoardingState(route, false)
	}
}

func (ts *transportScheduler) arriveRoute(route *sharedTransportRoute) {
	warps := make(map[int32]transportDestination, len(route.transitMaps))
	for _, transitMapID := range route.transitMaps {
		warps[transitMapID] = route.arrival
	}
	moveTransportPlayers(ts.server, warps)
}

func (ts *transportScheduler) canBoardFromMap(mapID int32) bool {
	for _, route := range ts.routes {
		for _, boardingMapID := range route.boardingMaps {
			if boardingMapID == mapID {
				return route.boardingOpen
			}
		}
	}
	return false
}

func (ts *transportScheduler) scheduleMuLungArrival(plr *Player) {
	if plr == nil {
		return
	}

	var destinationMapID int32
	switch plr.mapID {
	case constant.MapTransportToMuLung:
		destinationMapID = constant.MapMuLungArrival
	case constant.MapTransportToOrbis:
		destinationMapID = constant.MapOrbisArrival
	default:
		if timer, ok := ts.muLungTimersByPlayerID[plr.ID]; ok {
			timer.Stop()
			delete(ts.muLungTimersByPlayerID, plr.ID)
		}
		return
	}

	if timer, ok := ts.muLungTimersByPlayerID[plr.ID]; ok {
		timer.Stop()
	}

	expectedTransitMapID := plr.mapID
	timer := time.AfterFunc(muLungTransportRideDuration, func() {
		ts.server.post(func() {
			delete(ts.muLungTimersByPlayerID, plr.ID)

			current, err := ts.server.players.GetFromID(plr.ID)
			if err != nil || current == nil || current.mapID != expectedTransitMapID {
				return
			}

			dstField, ok := ts.server.fields[destinationMapID]
			if !ok {
				return
			}

			dstInst, err := dstField.getInstance(current.inst.id)
			if err != nil {
				dstInst, err = dstField.getInstance(0)
				if err != nil {
					return
				}
			}

			portal, err := dstInst.getPortalFromID(0, true)
			if err != nil {
				return
			}

			_ = ts.server.warpPlayer(current, dstField, portal, true)
		})
	})

	ts.muLungTimersByPlayerID[plr.ID] = timer
}

func (route *sharedTransportRoute) isBoardingOpen(minuteOfHour int) bool {
	return route.containsMinute(minuteOfHour%route.cycleMinutes, route.boardingOffset, route.boardingDuration)
}

func (route *sharedTransportRoute) isDepartureMinute(minuteOfHour int) bool {
	return minuteOfHour%route.cycleMinutes == route.departureOffset
}

func (route *sharedTransportRoute) isArrivalMinute(minuteOfHour int) bool {
	return minuteOfHour%route.cycleMinutes == route.arrivalOffset()
}

func (route *sharedTransportRoute) isInTransit(minuteOfHour int) bool {
	return route.containsMinute(minuteOfHour%route.cycleMinutes, route.departureOffset, route.transitMinutes)
}

func (route *sharedTransportRoute) arrivalOffset() int {
	return (route.departureOffset + route.transitMinutes) % route.cycleMinutes
}

func (route *sharedTransportRoute) containsMinute(minute, start, length int) bool {
	if length <= 0 {
		return false
	}
	end := start + length
	if end <= route.cycleMinutes {
		return minute >= start && minute < end
	}
	return minute >= start || minute < end%route.cycleMinutes
}

func maybeStartElliniaBoatInvasion(server *Server, rng *rand.Rand) {
	if rng == nil || rng.Float64() >= 0.3 {
		return
	}

	go func() {
		timer := time.NewTimer(5 * time.Minute)
		defer timer.Stop()
		<-timer.C

		server.post(func() {
			invasion(server)
		})

		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		finish := time.NewTimer(5 * time.Minute)
		defer finish.Stop()

		for {
			select {
			case <-ticker.C:
				server.post(func() {
					checkInvasion(server, false)
				})
			case <-finish.C:
				return
			}
		}
	}()
}
