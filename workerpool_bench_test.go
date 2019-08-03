package workerpool

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/Jeffail/tunny"
	slaves "github.com/dgrr/GoSlaves"
	"github.com/panjf2000/ants"
	"github.com/stretchr/testify/assert"
)

func BenchmarkWorkerPool(b *testing.B) {
	var wg sync.WaitGroup
	sum := uint64(0)
	wg.Add(b.N)
	count := func(v interface{}) {
		atomic.AddUint64(&sum, v.(uint64))
		wg.Done()
	}
	wp, err := New(HandlerFunc(count), WithWorkerPoolMaxWorkersCount(b.N))
	assert.Nil(b, err)
	defer wp.Stop()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// i escape to the heap
		wp.Serve(uint64(i))
	}
	wg.Wait()
}

func BenchmarkAntsWorkerPool(b *testing.B) {
	var wg sync.WaitGroup
	wg.Add(b.N)
	sum := uint64(0)
	p, _ := ants.NewPoolWithFunc(b.N, func(i interface{}) {
		atomic.AddUint64(&sum, i.(uint64))
		wg.Done()
	})
	defer p.Release()

	b.ReportAllocs()
	b.ResetTimer()

	// Submit tasks one by one.
	for i := 0; i < b.N; i++ {
		p.Invoke(uint64(i))
	}
	wg.Wait()
}

func BenchmarkTunnyWorkerPool(b *testing.B) {
	var wg sync.WaitGroup
	wg.Add(b.N)
	sum := uint64(0)
	pool := tunny.NewFunc(b.N, func(payload interface{}) interface{} {
		atomic.AddUint64(&sum, payload.(uint64))
		wg.Done()
		return nil
	})
	defer pool.Close()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pool.Process(uint64(i))
	}
	wg.Wait()
}

func BenchmarkSlaveWorkerPool(b *testing.B) {
	var wg sync.WaitGroup
	sum := uint64(0)
	wg.Add(b.N)
	pool := slaves.NewPool(b.N, func(obj interface{}) {
		atomic.AddUint64(&sum, obj.(uint64))
		wg.Done()
	})
	defer pool.Close()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pool.Serve(uint64(i))
	}
	wg.Wait()
}
