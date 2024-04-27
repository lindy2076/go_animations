package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var frames = []rune{'-', '\\', '|', '/'}
var frames2 = []rune{'0', '1', '2', '3', '4', '5'}

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
func (l *LoadingAnimation) PrintNextFrame(ctx context.Context) {
	if !l.reverse {
		fmt.Printf("%c", l.NextFrame())
	} else {
		fmt.Printf("%c", l.NextFrameReverse())
	}
	fmt.Printf("\b")
	select {
	case <-ctx.Done():
	case <-l.ticker.C:
	}
}

// Plays animation until the context is cancelled.
func (l *LoadingAnimation) Play(ctx context.Context) {
	for {
		if ctx.Err() != nil {
			return
		}
		l.PrintNextFrame(ctx)
	}
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

	mainCtx := context.Background()

	ctx, cancel := context.WithCancelCause(mainCtx)
	// some time consuming job
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(2 * time.Second)
		cancel(fmt.Errorf("job done"))
	}()
	// one way 5fps animation
	loadingAnimation := NewLoadingAnimation(frames, 5, false)
	fmt.Printf("LOADING:")
	go loadingAnimation.Play(ctx)

	wg.Wait()
	fmt.Printf("%c\n", '✓')

	fmt.Println("AND ANOTHER ONE")

	ctx, cancel = context.WithCancelCause(mainCtx)
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(3 * time.Second)
		cancel(fmt.Errorf("job done"))
	}()
	// reverse 8fps animation
	loadingAnimation2 := NewLoadingAnimation(frames2, 8, true)
	fmt.Printf("PROCESSING:")
	go loadingAnimation2.Play(ctx)

	wg.Wait()
	fmt.Printf("%c\n", '✓')
}
