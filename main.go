package main

import (
	"fmt"
	"sync"
	"time"
)

var frames = []rune{'-', '\\', '|', '/'}
var frames2 = []rune{'R', 'A', 'C', 'E'}

type LoadingAnimation struct {
	frames  []rune
	reverse bool

	frame  int
	ticker time.Ticker
}

// Returns next frame for the animation. Goes to frame 0 after the last one.
func (l *LoadingAnimation) NextFrame() rune {
	l.frame = l.frame%len(l.frames) + 1
	return l.frames[l.frame-1]
}

// Returns next frame for the animation with reverse (0,1,...,n-1,n,n-1,...,0).
func (l *LoadingAnimation) NextFrameReverse() rune {
	n := len(l.frames)
	maxframes := n*2 - 2
	l.frame = l.frame%maxframes + 1

	var idx int = l.frame
	if l.frame > n {
		idx = n - l.frame%n
	}

	return l.frames[idx-1]
}

// Prints next frame and then blocks a thread for 1000/fps millis.
func (l *LoadingAnimation) PrintNextFrame() {
	if !l.reverse {
		fmt.Printf("%c", l.NextFrame())
	} else {
		fmt.Printf("%c", l.NextFrameReverse())
	}
	<-l.ticker.C
	fmt.Printf("\b")
}

// Creates new LoadingAnimation. Frames is slice of runes that will be displayed in
// corresponding order. fps is self-explanatory and reverse enables playing animation
// backwards after hitting the end frame.
func NewLoadingAnimation(frames []rune, fps int, reverse bool) *LoadingAnimation {
	return &LoadingAnimation{
		frames:  frames,
		reverse: reverse,
		ticker:  *time.NewTicker(time.Millisecond * time.Duration(1000/fps)),
	}
}

func main() {
	fmt.Println("WHAT A WONDERFUL LOADING ANIMATION")

	// some time consuming job
	var wg sync.WaitGroup
	wg.Add(1)
	jobDone := make(chan bool)
	go func(jobDone chan<- bool) {
		defer func() {
			wg.Done()
			jobDone <- true
		}()

		time.Sleep(time.Second * 3)
	}(jobDone)

	// one way 5fps animation
	loadingAnimation := NewLoadingAnimation(frames, 5, false)
	fmt.Printf("LOADING:")
	go func(jobDone <-chan bool) {
		for {
			select {
			case <-jobDone:
				return
			default:
				loadingAnimation.PrintNextFrame()
			}
		}
	}(jobDone)

	wg.Wait()
	fmt.Println("\nDONE")

	fmt.Println("AND ANOTHER ONE")
	// 3fps animation with reverse
	loadingAnimation2 := NewLoadingAnimation(frames2, 3, true)
	fmt.Printf("LOADING:")
	go func() {
		for {
			loadingAnimation2.PrintNextFrame()
		}
	}()
	time.Sleep(5 * time.Second)
	fmt.Println("\nDONE")
}
