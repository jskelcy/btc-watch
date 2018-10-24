package aggregation

import (
	"container/heap"
	"log"
	"os"
	"sync"
	"testing"
)

func TestAggregatorMovingAvg(t *testing.T) {
	dps := &dataPoints{}
	heap.Init(dps)
	a := &aggregator{
		aggWindow:        10,
		collectionWindow: 1,
		dps:              dps,
		log:              log.New(os.Stderr, "", 1),
	}

	wg := sync.WaitGroup{}

	values := []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	for _, value := range values {
		wg.Add(1)
		go func(v float64) {
			defer wg.Done()
			a.Ingest(v)
		}(value)
	}

	wg.Wait()
	avg, err := a.CurrentAggregatedPrice()
	if err != nil {
		t.Error(err)
	}
	if avg != 4.5 {
		t.Errorf("%v does not equal 4.5", avg)
	}

	a.Ingest(20)
	avg, err = a.CurrentAggregatedPrice()
	if err != nil {
		t.Error(err)
	}
	if avg != 6.5 {
		t.Errorf("%v does not equal 6.5", avg)
	}
}
