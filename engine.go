package tmgang

import (
	"time"

	"github.com/gdamore/tcell/v2"
)

type StateKey string

// State definition of the game state machine.
//
// Any computing can be summarized by a state machine. A game, similarly so.
// This is the core definition of the states in the state machine. A state
// will be drawn (and interacted with), based on camera position theoretically
// infinite area the state spans (i.e., position of its entities).
type State interface {
	GetDrawables() []Entity
	GetInteractables() []InteractiveEntity
	GetCamera() Coordinates
	NextState() StateKey
}

// The engine will Run your game.
//
// Yes.
type Engine interface {
	Configure(*EngineOpts)
	ScreenSize() (int, int)
	Run() error
}

type EngineOpts struct {
	Fps          uint32
	InitialState StateKey
	StateMachine map[StateKey]State
	Overlays     []Overlay
}

func NewEngine() (Engine, error) {
	s, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	if err := s.Init(); err != nil {
		return nil, err
	}

	resetStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	s.SetStyle(resetStyle)
	s.Clear()

	eng := &engine{
		screen: s,
	}

	return eng, nil
}

type engine struct {
	screen tcell.Screen

	fps          uint32
	currentState State
	stateMachine map[StateKey]State
	overlays     []Overlay
}

func (eng *engine) Configure(opts *EngineOpts) {
	if opts.Fps == 0 {
		opts.Fps = 30
	}
	eng.fps = opts.Fps
	eng.currentState = opts.StateMachine[opts.InitialState]
	eng.stateMachine = opts.StateMachine
	eng.overlays = opts.Overlays
}

func (eng *engine) ScreenSize() (int, int) {
	return eng.screen.Size()
}

// Run until the end of time!
//
// Any errors that Run returns might be helpful in debugging issues in your game.
// The Run function, given the initialState and the stateMachine, work through the
// entities held in each state, calculate interactions, draw the screen, and eventually
// attempt to calculate state transitions. If a state transitions, the new state
// entities will go through the same process.
func (eng *engine) Run() error {
	return eng.run()
}

func (eng *engine) run() error {
	defer func() {
		// per tcell documentation: recover, finalize and rethrow panic to have the diagnostic trace
		// just make sure to cleanup the terminal
		panicked := recover()
		eng.screen.Clear()
		eng.screen.Fini()
		if panicked != nil {
			panic(panicked)
		}
	}()

	var (
		camera = eng.currentState.GetCamera()

		ratioX, ratioY     float32 = 1, 1
		changedX, changedY int
		initialX, initialY = eng.screen.Size()

		eventChannel = make(chan tcell.Event)
		ev           tcell.Event
	)

	fpsTicker := time.NewTicker(time.Duration(1000000000 / eng.fps))
	defer fpsTicker.Stop()

	go func() {
		quit := make(chan struct{})
		eng.screen.ChannelEvents(eventChannel, quit)
	}()

	for {
		eng.screen.Show()

		select {
		case ev = <-eventChannel:
		case <-fpsTicker.C:
		}

		switch ev := ev.(type) {
		case *tcell.EventResize:
			changedX, changedY = eng.screen.Size()
			ratioX, ratioY = (float32(changedX) / float32(initialX)), (float32(changedY) / float32(initialY))
			eng.screen.Sync()
		case *tcell.EventKey:
			// escape route (usually SIGINT)
			if ev.Key() == tcell.KeyCtrlC {
				return nil
			}

			for _, interactiveEntity := range eng.currentState.GetInteractables() {
				interactiveEntity.ProcessKeyEvent(ev, camera)
			}
		}
		ev = nil

		for _, entity := range eng.currentState.GetDrawables() {
			entity.Draw(eng.screen, camera, ratioX, ratioY)
		}

		for i := 0; i < len(eng.overlays); i++ {
			eng.overlays[i].Draw(eng.screen)
		}

		eng.currentState = eng.stateMachine[eng.currentState.NextState()]
		camera = eng.currentState.GetCamera()
	}

}
