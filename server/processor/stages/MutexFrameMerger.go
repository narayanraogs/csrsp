package stages

import (
	"context"
	"csrspServer/pipeline"
	"sync"
)

type MutexFrameMergerConfig struct{}

func NewMutexFrameMergerStage(config MutexFrameMergerConfig, errChan chan<- error) pipeline.StageManyToOne {
	return func(ctx context.Context, inputs []<-chan pipeline.Frame) <-chan pipeline.Frame {
		output := make(chan pipeline.Frame)

		go func() {
			defer close(output)

			var wg sync.WaitGroup
			wg.Add(len(inputs))

			winnerChan := make(chan int, 1)
			winnerFound := make(chan struct{})

			for i, inputChan := range inputs {
				go func(index int, ch <-chan pipeline.Frame) {
					defer wg.Done()
					select {
					case frame, ok := <-ch:
						if !ok {
							return
						}
						select {
						case winnerChan <- index:
							output <- frame
							for subsequentFrame := range ch {
								output <- subsequentFrame
							}
						case <-winnerFound:
							return
						case <-ctx.Done():
							return
						}
					case <-winnerFound:
						return
					case <-ctx.Done():
						return
					}
				}(i, inputChan)
			}

			select {
			case <-winnerChan:
				close(winnerFound)
			case <-ctx.Done():
				close(winnerFound)
			}
			wg.Wait()
		}()

		return output
	}
}
