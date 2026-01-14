package main

import (
	"bytes"
	"encoding/binary"
	"math"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

const sampleRate = 44100

var audioContext = audio.NewContext(sampleRate)

func playBounce(kind string) {
	freq := 600.0 // wall
	if kind == "paddle" {
		freq = 1000.0
	}
	if kind == "start" {
		freq = 500.0
	}
	if kind == "out" {
		freq = 250.0
	}

	const (
		duration = 0.04
		volume   = 0.25
	)

	n := int(sampleRate * duration)
	buf := make([]byte, n*2) // 16-bit mono

	for i := 0; i < n; i++ {
		t := math.Sin(2 * math.Pi * freq * float64(i) / sampleRate)
		var v float64
		if t >= 0 {
			v = 1
		} else {
			v = -1
		}

		s := int16(v * volume * math.MaxInt16)
		binary.LittleEndian.PutUint16(buf[i*2:], uint16(s))
	}

	p, _ := audioContext.NewPlayer(bytes.NewReader(buf))
	p.Play()
}
