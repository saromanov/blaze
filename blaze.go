package blaze

import (
	"fmt"
	"time"
)

type ExecuteFunc func() (interface{}, error)

// Blaze implemenths the basic method
type Blaze struct {
	mainExec ExecuteFunc
	timeout  time.Time
	duration time.Duration
	steps    []BlazeStep
}

// BlazeConfig provides configuration for the Blaze
type BlazeConfig struct {
	MainExec ExecuteFunc
	Steps    []BlazeStep
	Duration time.Duration
}

// BlazeStep implements step
type BlazeStep struct {
	Name     string
	Duration time.Duration
	Execute  ExecuteFunc
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
func (b *Blaze) NewBlaze(conf *BlazeConfig) *Blaze {
	return &Blaze{
		mainExec: conf.MainExec,
		steps:    conf.Steps,
		duration: conf.Duration,
	}
}

// Do is a main method for executing
func (b *Blaze) Do() error {
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for t := range ticker.C {
			fmt.Println(t)
			b.mainExec()
		}
	}()
	time.Sleep(b.duration)
	ticker.Stop()
	return nil
}
