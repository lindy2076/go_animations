# Some animations in terminal

### loading animation

The animation is built from slice of runes.  

Can be played in two modes: 

1. Direct: `0 -> 1 -> ... -> n-1 -> n -> 0 -> 1 -> ...`
2. Reverse (back and forth): `0 -> 1 -> ... -> n-1 -> n -> n-1 -> ... -> 1 -> 0 -> ...`

```go
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

```

![qw](https://github.com/lindy2076/go_animations/assets/67479681/1337c243-e230-4b8b-a0c6-1307a3d4559b)


