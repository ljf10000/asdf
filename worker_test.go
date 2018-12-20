package asdf

import (
	"testing"
)

func TestWorkerProgress(t *testing.T) {
	newWorker := func(begin, end, current int) WorkerProgress {
		return WorkerProgress{
			Begin:   begin,
			End:     end,
			Current: current,
		}
	}

	workers := []WorkerProgress{
		newWorker(1, 2, 1),
		newWorker(1, 3, 1),
		newWorker(1, 3, 2),
		newWorker(1, 10000, 1),
		newWorker(1, 10000, 10),
		newWorker(1, 10000, 100),
		newWorker(1, 10000, 1000),
		newWorker(1, 10000, 9999),
	}

	for _, w := range workers {
		w.Calc()

		t.Logf("worker: %s", &w)
	}
}
