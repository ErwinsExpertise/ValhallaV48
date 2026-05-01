package channel

import (
	"log"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/dop251/goja"
)

func ensureFieldInstance(f *field, id int, rates *rates, server *Server) (*fieldInstance, error) {
	for len(f.instances) <= id {
		f.createInstance(rates, server)
	}
	return f.getInstance(id)
}

type event struct {
	id         int32
	duration   time.Duration
	endTime    time.Time
	finished   chan struct{}
	instanceID int
	playerIDs  []int32
	server     *Server

	startCallback            func()
	beforePortalCallback     func(plr scriptPlayerWrapper, src scriptMapWrapper, dst scriptMapWrapper) bool
	afterPortalCallback      func(plr scriptPlayerWrapper, dst scriptMapWrapper)
	onMapChangeCallback      func(plr scriptPlayerWrapper, dst scriptMapWrapper)
	timeoutCallback          func(plr scriptPlayerWrapper)
	playerLeaveEventCallback func(plr scriptPlayerWrapper)
	scheduledCallbacks       map[string]func()

	program *goja.Program
	vm      *goja.Runtime

	closeFinish func()
	timerReset  chan struct{}
	properties  map[string]interface{}
}

func createEvent(id int32, instID int, players []int32, server *Server, program *goja.Program) (*event, error) {
	ctrl := &event{
		id:                 id,
		finished:           make(chan struct{}),
		instanceID:         instID,
		playerIDs:          players,
		server:             server,
		program:            program,
		vm:                 goja.New(),
		timerReset:         make(chan struct{}, 1),
		scheduledCallbacks: make(map[string]func()),
		properties:         make(map[string]interface{}),
	}

	ctrl.closeFinish = sync.OnceFunc(func() {
		close(ctrl.finished)
	})

	ctrl.vm.SetFieldNameMapper(goja.UncapFieldNameMapper())
	_ = ctrl.vm.Set("ctrl", ctrl)

	_, err := ctrl.vm.RunProgram(ctrl.program)

	if err != nil {
		return nil, err
	}

	if fn := ctrl.vm.Get("start"); fn != nil && !goja.IsUndefined(fn) && !goja.IsNull(fn) {
		err = ctrl.vm.ExportTo(fn, &ctrl.startCallback)
		if err != nil {
			return nil, err
		}
	}

	if fn := ctrl.vm.Get("beforePortal"); fn != nil && !goja.IsUndefined(fn) && !goja.IsNull(fn) {
		err = ctrl.vm.ExportTo(fn, &ctrl.beforePortalCallback)
		if err != nil {
			return nil, err
		}
	}

	if fn := ctrl.vm.Get("afterPortal"); fn != nil && !goja.IsUndefined(fn) && !goja.IsNull(fn) {
		err = ctrl.vm.ExportTo(fn, &ctrl.afterPortalCallback)
		if err != nil {
			return nil, err
		}
	}

	if fn := ctrl.vm.Get("playerLeaveEvent"); fn != nil && !goja.IsUndefined(fn) && !goja.IsNull(fn) {
		err = ctrl.vm.ExportTo(fn, &ctrl.playerLeaveEventCallback)
		if err != nil {
			return nil, err
		}
	}

	if fn := ctrl.vm.Get("onMapChange"); fn != nil && !goja.IsUndefined(fn) {
		_ = ctrl.vm.ExportTo(fn, &ctrl.onMapChangeCallback)
	}

	if fn := ctrl.vm.Get("timeout"); fn != nil && !goja.IsUndefined(fn) && !goja.IsNull(fn) {
		err = ctrl.vm.ExportTo(fn, &ctrl.timeoutCallback)
		if err != nil {
			return nil, err
		}
	}

	for _, name := range []string{"begin", "finish", "earringcheck", "broadcastClock"} {
		if fn := ctrl.vm.Get(name); fn != nil && !goja.IsUndefined(fn) {
			var cb func()
			if err := ctrl.vm.ExportTo(fn, &cb); err == nil && cb != nil {
				ctrl.scheduledCallbacks[name] = cb
			}
		}
	}

	return ctrl, nil
}

func (e *event) start() {
	for _, id := range e.playerIDs {
		if plr, err := e.server.players.GetFromID(id); err == nil {
			plr.event = e
		}
	}

	e.startCallback()

	go func() {
		timeout := time.NewTimer(e.duration)
		defer timeout.Stop()

		for {
			select {
			case <-timeout.C:
				keepRunning := make(chan bool, 1)
				e.server.dispatch <- func() {
					for _, id := range e.playerIDs {
						if plr, err := e.server.players.GetFromID(id); err == nil && e.timeoutCallback != nil {
							e.timeoutCallback(scriptPlayerWrapper{plr: plr, server: e.server})
						}
					}

					keep := false
					if current, ok := e.server.events[e.id]; ok && current == e {
						keep = time.Until(e.endTime) > 0
					}

					if !keep {
						for _, id := range e.playerIDs {
							if plr, err := e.server.players.GetFromID(id); err == nil {
								plr.event = nil
							}
						}

						delete(e.server.events, e.id)
					}

					keepRunning <- keep
				}

				if <-keepRunning {
					remaining := time.Until(e.endTime)
					if remaining <= 0 {
						remaining = time.Second
					}
					timeout.Reset(remaining)
					continue
				}

				return

			case <-e.timerReset:
				if !timeout.Stop() {
					select {
					case <-timeout.C:
					default:
					}
				}
				timeout.Reset(e.duration)

			case <-e.finished:
				e.server.dispatch <- func() {
					for _, id := range e.playerIDs {
						if plr, err := e.server.players.GetFromID(id); err == nil {
							plr.event = nil
						}
					}

					delete(e.server.events, e.id)
				}
				return
			}
		}
	}()
}

func (e *event) Log(msg string) {
	log.Println(msg)
}

func (e *event) GetProperty(key string) interface{} {
	if value, ok := e.properties[key]; ok {
		return value
	}
	return nil
}

func (e *event) SetProperty(key string, value interface{}) interface{} {
	prev := e.GetProperty(key)
	e.properties[key] = value
	return prev
}

func (e *event) RemainingTime() int32 {
	return int32(time.Until(e.endTime).Seconds())
}

func (e *event) PlayerCount() int {
	return len(e.playerIDs)
}

func (e *event) Finished() {
	e.closeFinish()
}

func (e *event) Players() []scriptPlayerWrapper {
	r := make([]scriptPlayerWrapper, 0, len(e.playerIDs))

	for _, id := range e.playerIDs {
		if plr, err := e.server.players.GetFromID(id); err == nil {
			r = append(r, scriptPlayerWrapper{plr: plr, server: e.server})
		}
	}

	return r
}

func (e *event) RemovePlayer(plr scriptPlayerWrapper) {
	for i, v := range e.playerIDs {
		if v == plr.plr.ID {
			e.playerIDs = slices.Delete(e.playerIDs, i, i+1)
			break
		}
	}
	plr.plr.event = nil
}

func (e *event) AddPlayer(plr scriptPlayerWrapper) {
	for _, id := range e.playerIDs {
		if id == plr.plr.ID {
			return
		}
	}
	e.playerIDs = append(e.playerIDs, plr.plr.ID)
	plr.plr.event = e
}

func (e *event) SetDuration(duration string) {
	countdown, err := time.ParseDuration(duration)

	if err != nil {
		countdown = time.Second * 10
	}

	e.duration = countdown
	e.endTime = time.Now().Add(countdown)

	select {
	case e.timerReset <- struct{}{}:
	default:
	}
}

func (e *event) Schedule(name string, duration string) {
	name = strings.TrimSpace(name)
	cb, ok := e.scheduledCallbacks[name]
	if !ok || cb == nil {
		if fn := e.vm.Get(name); fn != nil && !goja.IsUndefined(fn) && !goja.IsNull(fn) {
			var exported func()
			if err := e.vm.ExportTo(fn, &exported); err == nil && exported != nil {
				cb = exported
				e.scheduledCallbacks[name] = exported
			} else {
				return
			}
		} else {
			return
		}
	}
	d, err := time.ParseDuration(duration)
	if err != nil {
		return
	}
	time.AfterFunc(d, func() {
		e.server.dispatch <- func() {
			if _, exists := e.server.events[e.id]; exists {
				cb()
			}
		}
	})
}

func (e *event) GetMap(id int32) scriptMapWrapper {
	if field, ok := e.server.fields[id]; ok {
		inst, err := ensureFieldInstance(field, e.instanceID, &e.server.rates, e.server)

		if err != nil {
			return scriptMapWrapper{}
		}

		return scriptMapWrapper{inst: inst, server: e.server}
	}

	return scriptMapWrapper{}
}

func (e *event) WarpPlayers(dst int32) {
	field := e.server.fields[dst]
	dstInst, err := ensureFieldInstance(field, e.instanceID, &e.server.rates, e.server)
	if err != nil {
		return
	}

	dstPortal, err := dstInst.getRandomSpawnPortal()
	if err != nil {
		return
	}

	for _, id := range e.playerIDs {
		if plr, err := e.server.players.GetFromID(id); err == nil {
			e.server.warpPlayer(plr, field, dstPortal, false)
		}
	}
}

func (e *event) WarpPlayersToPortal(dst int32, portalName string) {
	field := e.server.fields[dst]
	dstInst, err := ensureFieldInstance(field, e.instanceID, &e.server.rates, e.server)
	if err != nil {
		return
	}

	dstPortal, err := dstInst.getPortalFromName(portalName)
	if err != nil {
		return
	}

	for _, id := range e.playerIDs {
		if plr, err := e.server.players.GetFromID(id); err == nil {
			e.server.warpPlayer(plr, field, dstPortal, false)
		}
	}
}

func (e *event) IsParticipantsOnMap(mapID int32) bool {
	for _, id := range e.playerIDs {
		if plr, err := e.server.players.GetFromID(id); err == nil {
			if plr.mapID != mapID {
				return false
			}
		}
	}
	return true
}
