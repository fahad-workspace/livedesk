// Command matchfab generates synthetic, football-realistic match events and
// writes them as NDJSON, one event per line. It is LiveDesk's data innovation
// deliverable and its ingestion source.
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/fahad-workspace/livedesk/internal/matchfab"
)

func main() {
	matchID := flag.String("match", "match-001", "match identifier")
	seed := flag.Int64("seed", 42, "random seed for reproducible matches")
	outPath := flag.String("out", "", "output file path (default stdout)")
	speed := flag.Float64("speed", 0, "real-time pacing multiplier (0 = as fast as possible, 1 = real time, 60 = one minute per second)")
	flag.Parse()

	eng := matchfab.NewEngine(*matchID, rand.New(rand.NewSource(*seed)))

	var w *bufio.Writer
	if *outPath == "" {
		w = bufio.NewWriter(os.Stdout)
	} else {
		f, err := os.Create(*outPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer f.Close()
		w = bufio.NewWriter(f)
	}
	defer w.Flush()

	enc := json.NewEncoder(w)
	var lastClock int64
	eng.Run(func(ev matchfab.Event) {
		if *speed > 0 && lastClock > 0 {
			delta := time.Duration(float64(ev.ClockMs-lastClock) / *speed * float64(time.Millisecond))
			time.Sleep(delta)
		}
		lastClock = ev.ClockMs
		if err := enc.Encode(ev); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	})
}
