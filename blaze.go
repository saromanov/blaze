package blaze

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var errNoSteps = errors.New("steps is not defined")

// ExecuteFunc defines type for the execution func
// on test
type ExecuteFunc func() (interface{}, error)

// Blaze implemenths the basic method
type Blaze struct {
	mainExec  ExecuteFunc
	timeout   time.Time
	duration  time.Duration
	tickEvery time.Duration
	steps     []step
}

// Config provides configuration for the Blaze
type Config struct {
	MainExec  ExecuteFunc
	TickEvery time.Duration
	Steps     []Step
	Duration  time.Duration
}

// Step implements step
type Step struct {
	Name     string
	Duration time.Duration
	Execute  ExecuteFunc
}

func (s *Step) makeStep() step {
	return step{
		Name:     s.Name,
		Duration: s.Duration,
		Execute:  s.Execute,
		m:        &sync.RWMutex{},
	}
}

// Step implements step
type step struct {
	Name      string
	Duration  time.Duration
	Execute   ExecuteFunc
	executed  bool
	m         *sync.RWMutex
	order     int
	started   bool
	startTime time.Time
	endTime   time.Time
}

func (s *step) updateExecuted() {
	s.m.RLock()
	defer s.m.RUnlock()
	s.executed = true
}

/*
Example of using

type Method struct {
	Url string
}

func (m *Method) Execute() {
	fmt.Println("URL: ", m.Url)
	_, err := http.Get(m.Url)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

url := "https://www.google.ru"
b := NewBlaze(&BlazeConfig {
	Duration: 1 * time.Minute,
	Steps: []Step {
		Name: "First Step"
		Duration: 10 * time.Second,
		Execuite: func()(interface{}, error ) {
			url = "https://www.yandex.ru"
		},
	},
})

err := b.Do()
if err != nil {
	fmt.Fatal("unable to pass test")
}
*/
func New(conf *Config) *Blaze {
	steps := make([]step, len(conf.Steps))
	startTime := time.Now()
	for i, s := range conf.Steps {
		steps[i] = s.makeStep()
		steps[i].order = i
		steps[i].startTime = startTime
		startTime = startTime.Add(steps[i].Duration)
		steps[i].endTime = startTime
	}
	return &Blaze{
		mainExec:  conf.MainExec,
		steps:     steps,
		duration:  conf.Duration,
		tickEvery: conf.TickEvery,
	}
}

// Do is a main method for executing
func (b *Blaze) Do() error {
	err := b.checkConfig()
	if err != nil {
		return err
	}
	fmt.Println(b.steps)
	//startTime := time.Now()
	ticker := time.NewTicker(b.tickEvery)
	go func() {
		for t := range ticker.C {
			step, i := b.getStep()
			if step.executed {
				continue
			}
			if b.steps[i].executed {
				continue
			}
			fmt.Println("STEP:", t, b.steps[i].started)
			seconds := time.Since(b.steps[i].startTime).Seconds()
			if seconds > 0 {
				b.steps[i].started = true
			}
			end := time.Since(b.steps[i].endTime).Seconds()
			if end > 0 {
				b.steps[i].executed = true
			}
			if !b.steps[i].executed {
				b.steps[i].Execute()
			}
			b.mainExec()
		}
	}()
	time.Sleep(b.duration)
	ticker.Stop()
	return nil
}

// getStep returns first not executed step
func (b *Blaze) getStep() (step, int) {
	for i, s := range b.steps {
		if !s.executed {
			return s, i
		}
	}
	return step{}, 0
}

// checkConfig provides checking of required
// arguments on config
func (b *Blaze) checkConfig() error {

	if len(b.steps) == 0 {
		return errNoSteps
	}
	return nil
}
