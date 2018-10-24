package aggregation

import (
	"container/heap"
	"errors"
	"log"
	"os"
	"sync"
	"time"
)

var (
	errNotPrimedError = errors.New("Data is still being aggregated please wait")
)

// Aggregator ingests data points and uses an aggregation function on new data points.
// All operations on an Aggregator are thread safe.
type Aggregator interface {
	// CurrentAggregatedPrice returns current price after aggregation.
	CurrentAggregatedPrice() (float64, error)
	// Ingest takes a new value to be added to he aggregation.
	Ingest(float64)
}

// Config contains the configuration for an Aggregator.
type Config struct {
	AggWindow          int
	CollectionInterval int
}

// NewAggregator returns a the default aggregator from config.
// The default aggregator uses a moving 1 minute moving average
// as the aggregation function.
func NewAggregator(cfg Config) Aggregator {
	dps := &dataPoints{}
	heap.Init(dps)
	return &aggregator{
		aggWindow:          cfg.AggWindow,
		collectionInterval: cfg.CollectionInterval,
		dps:                dps,
		log:                log.New(os.Stderr, "", 1),
	}
}

// movingAvgAlerter takes in new values and keeps a moving average. Struct is
// not thread safe.
type aggregator struct {
	sync.Mutex
	dps                  *dataPoints
	primed               bool
	total                float64
	aggWindow            int
	collectionInterval   int
	lastMinMovingAverage float64
	log                  *log.Logger
}

func (a *aggregator) CurrentAggregatedPrice() (float64, error) {
	a.Lock()
	defer a.Unlock()

	if !a.primed {
		return 0, errNotPrimedError
	}

	return a.lastMinMovingAverage, nil
}

func (a *aggregator) Ingest(value float64) {
	a.Lock()
	defer a.Unlock()
	// if aggWindow is full pop off the oldest values
	if len(*a.dps) == (a.aggWindow / a.collectionInterval) {
		oldValue := heap.Pop(a.dps).(dataPoint)
		a.total = a.total - oldValue.value
	}

	newDataPoint := dataPoint{
		timestamp: time.Now().UnixNano(),
		value:     value,
	}

	heap.Push(a.dps, newDataPoint)
	a.total = a.total + newDataPoint.value

	a.lastMinMovingAverage = a.total / float64(len(*a.dps))
	a.primed = true
	a.log.Printf("average: %.2f", a.lastMinMovingAverage)
}

type dataPoint struct {
	timestamp int64
	value     float64
}

// Heap interface implementation to calculate moving average --
type dataPoints []dataPoint

func (dps dataPoints) Len() int           { return len(dps) }
func (dps dataPoints) Less(i, j int) bool { return dps[i].timestamp < dps[j].timestamp }
func (dps dataPoints) Swap(i, j int)      { dps[i], dps[j] = dps[j], dps[i] }

func (dps *dataPoints) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*dps = append(*dps, x.(dataPoint))
}

func (dps *dataPoints) Pop() interface{} {
	old := *dps
	n := len(old)
	x := old[n-1]
	*dps = old[0 : n-1]
	return x
}

// -----------------------------------------------------------
