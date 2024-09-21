package pipeline

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if len(stages) == 0 {
		return in
	}

	if in == nil {
		return in
	}

	for _, stage := range stages {
		in = executor(done, stage(in))
	}

	return in
}

func executor(done, stageOut In) Out {
	out := make(Bi)

	go func() {
		defer func() {
			close(out)
			for range stageOut {
				_ = stageOut
			}
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
