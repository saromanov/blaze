package blaze

import (
	"errors"
	"fmt"
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
	steps     []BlazeStep
}

// Config provides configuration for the Blaze
type Config struct {
	MainExec  ExecuteFunc
	TickEvery time.Duration
	Steps     []BlazeStep
	Duration  time.Duration
}

// BlazeStep implements step
type BlazeStep struct {
	Name     string
	Duration time.Duration
	Execute  ExecuteFunc
	executed bool
}

func (s *BlazeStep) updateExecuted() {
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
	Steps: []BlazeStep {
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
	return &Blaze{
		mainExec:  conf.MainExec,
		steps:     conf.Steps,
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
	startTime := time.Now()
	ticker := time.NewTicker(b.tickEvery)
	go func() {
		for t := range ticker.C {
			for i, s := range b.steps {
				fmt.Println(t, s)
				seconds := time.Since(startTime).Seconds()
				if !s.executed && seconds > s.Duration.Seconds() {
					s.Execute()
					b.steps[i].updateExecuted()
				}
			}
			b.mainExec()
		}
	}()
	time.Sleep(b.duration)
	ticker.Stop()
	return nil
}

// checkConfig provides checking of required
// arguments on config
func (b *Blaze) checkConfig() error {

	if len(b.steps) == 0 {
		return errNoSteps
	}
	return nil
}
