package pipeline

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		in = executor(done, stageExecutor(done, in, stage))
	}
	return in
}

func executor(done, stageOut In) Out {
	out := make(Bi)

	go func() {
		defer func() {
			close(out)
		}()
		for {
			select {
			case <-done:
				return
			case v, ok := <-stageOut:
				if !ok {
					return
				}
				select {
				case <-done:
					return
				case out <- v:
				}
			}
		}
	}()

	return out
}

func stageExecutor(done, in In, stage Stage) Out {
	out := make(Bi)
	go func() {
		defer func() {
			close(out)
		}()
		for v := range stage(in) {
			select {
			case <-done:
				return
			case out <- v:
			}
		}
	}()

	return out
}
