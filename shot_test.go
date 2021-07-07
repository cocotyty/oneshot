package oneshot

import (
	"runtime"
	"sync"
	"testing"
	"time"
)
const TestSize = 2000000

// BenchmarkChannel-8   	 8432158	      1020 ns/op
func BenchmarkChannel(b *testing.B) {
	wg := sync.WaitGroup{}
	chs := make([]chan struct{}, b.N)
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		chs[i] = make(chan struct{})
		go func(ch chan struct{}) {
			<-ch
			wg.Done()
		}(chs[i])
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		chs[i] <- struct{}{}
	}
	wg.Wait()
}

// BenchmarkShot-8   	10965412	       889.4 ns/op
func BenchmarkShot(b *testing.B) {
	wg := sync.WaitGroup{}
	shots := make([]*Shot, b.N)

	for i := 0; i < b.N; i++ {
		i := i
		wg.Add(1)
		shots[i] = &Shot{}
		go func(s *Shot) {
			s.Wait()
			wg.Done()
		}(shots[i])
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		shots[i].WakeUp()
	}
	wg.Wait()
}

func TestShots(t *testing.T) {
	N := TestSize
	wg := sync.WaitGroup{}
	shots := make([]Shot, N)
	for i := 0; i < N; i++ {
		i := i
		wg.Add(1)
		go func(s *Shot) {
			s.Wait()
			wg.Done()
		}(&shots[i])
	}
	startTime := time.Now()
	for x := 0; x < N; x += 1 {
		shots[x].WakeUp()
	}
	wg.Wait()
	t.Log("usage:", time.Now().Sub(startTime))
}


func TestChans(t *testing.T) {
	N := TestSize
	wg := sync.WaitGroup{}
	shots := make([]chan struct{}, N)
	for i := 0; i < N; i++ {
		ch:=make(chan struct{})
		shots[i] = ch
		wg.Add(1)
		go func(ch chan struct{}) {
			<-ch
			wg.Done()
		}(ch)
	}
	startTime := time.Now()
	for x := 0; x < N; x += 1 {
		shots[x]<- struct{}{}
	}
	wg.Wait()
	t.Log("usage:", time.Now().Sub(startTime))
}


func TestShotParallel(t *testing.T) {
	N := TestSize
	wg := sync.WaitGroup{}
	shots := make([]*Shot, N)
	for i := 0; i < N; i++ {
		i := i
		shots[i] = &Shot{}
		wg.Add(1)
		go func(s *Shot) {
			s.Wait()
			wg.Done()
		}(shots[i])
	}
	startTime := time.Now()
	cpu := runtime.NumCPU()
	batch := N / cpu
	xshots := shots
	for {
		if len(xshots) <= batch {
			go func(shots []*Shot) {
				for _, s := range shots {
					s.WakeUp()
				}
			}(xshots)
			break
		}
		go func(shots []*Shot) {
			for _, s := range shots {
				s.WakeUp()
			}
		}(xshots[:batch])
		xshots = xshots[batch:]
	}
	wg.Wait()
	t.Log("usage:", time.Now().Sub(startTime))
}

func TestChannelParallel(t *testing.T) {
	N := TestSize
	wg := sync.WaitGroup{}
	chs := make([]chan struct{}, N)
	for i := 0; i < N; i++ {
		i := i
		wg.Add(1)
		ch := make(chan struct{})
		chs[i] = ch
		go func(ch chan struct{}) {
			<-ch
			wg.Done()
		}(ch)
	}
	startTime := time.Now()
	batch := N / runtime.NumCPU()
	xchs := chs
	for {
		if len(xchs) <= batch {
			go func(chs []chan struct{}) {
				for _, ch := range chs {
					ch <- struct{}{}
				}
			}(xchs)
			break
		}
		go func(chs []chan struct{}) {
			for _, ch := range chs {
				ch <- struct{}{}
			}
		}(xchs[:batch])
		xchs = xchs[batch:]
	}
	wg.Wait()
	t.Log("usage:", time.Now().Sub(startTime))
}
