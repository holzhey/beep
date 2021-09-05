package main

import (
	"io/ioutil"
	"time"

	"github.com/ariejan/i6502"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
)

type Noise struct{}

var a float64 = 0
var b float64 = 0.05

func (n Noise) Stream(samples [][2]float64) (qtd int, ok bool) {
	for i := range samples {
		a += b
		if (a > 10 && b > 0) || (a < -10 && b < 0) {
			b = -b
		}
		samples[i][0] = a
		samples[i][1] = a
	}
	return len(samples), true
}

func (n Noise) Err() error {
	return nil
}

func main() {
	sr := beep.SampleRate(44100)
	speaker.Init(sr, sr.N(time.Second/10))

	done := make(chan bool)
	speaker.Play(beep.Seq(beep.Take(sr.N(5*time.Second), Noise{}), beep.Callback(func() {
		done <- true
	})))
	speaker.Play(beep.Seq(beep.Take(sr.N(2*time.Second), Noise{}), beep.Callback(func() {
		done <- true
	})))
	<-done

	// Create Ram, 64kB in size
	ram, _ := i6502.NewRam(0x10000)

	// Create the AddressBus
	bus, _ := i6502.NewAddressBus()

	// And attach the Ram at 0x0000
	bus.Attach(ram, 0x0000)

	// Create the Cpu, with the AddressBus
	cpu, _ := i6502.NewCpu(bus)

	program, _ := ioutil.ReadFile("path")

	// This will load the program (if it fits within memory)
	// at 0x0200 and set cpu.PC to 0x0200 as well.
	cpu.LoadProgram(program, 0x0200)

	go cpu.Step()
}
