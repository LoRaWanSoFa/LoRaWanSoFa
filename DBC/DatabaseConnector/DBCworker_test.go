package DatabaseConnector

import (
	"testing"
	"time"
)

func TestNewWorker(t *testing.T) {
	WorkerQueueTest := make(chan chan WorkRequest, 1)
	worker := NewWorker(1, WorkerQueueTest)
	if 1 != worker.ID {
		t.Errorf("Expected %d, was %d", 1, worker.ID)
	}

}

func TestStart(t *testing.T) {
	//By stesting the StartDispatcher we also test the start function.
	startDispatcher(1)

	result := make(chan WorkResult)
	args := make([]interface{}, 1)
	go func(args []interface{}, result chan (WorkResult)) {
		WorkQueue <- WorkRequest{Query: "", Arguments: args, ResultChannel: result, F: func(w *WorkRequest) {
			w.ResultChannel <- WorkResult{Result: true, err: nil}
		}}
	}(args, result)
	time.Sleep(200 * time.Millisecond)
	fail := <-result
	if fail.Result != true {
		t.Errorf("Expected %t, was %+v", true, fail)
	}
	if stopWorker() != true {
		t.Errorf("Could not stop the worker!")
	}
	//For 100% coverage we will add a fucntion to stop the worker
	defer close(result)

}

func TestStop(t *testing.T) {
	WorkerQueue := make(chan chan WorkRequest, 1)
	worker := NewWorker(1, WorkerQueue)
	if 1 != worker.ID {
		t.Errorf("Expected %d, was %d", 1, worker.ID)
	}
	worker.stop()
	if 1 != worker.ID {
		t.Errorf("Expected %d, was %d", 1, worker.ID)
	}
}
