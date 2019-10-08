package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"./internal/list"
	"github.com/tkanos/gonfig"
)

// Configuration is ..
type Configuration struct {
	InitialSize int
	Operations  int
	Threads     int
	Jump        int
	Experiments int
}

func exec(l list.List, initialSize int, operations int, threads int) {

	rand.Seed(time.Now().UnixNano())
	seedSize := 200000

	//runtime.GOMAXPROCS(8)

	for x := 0; x < initialSize; x++ {
		var earlyAdd func()
		earlyAdd = func() {
			r := l.Add(rand.Intn(seedSize) + 1)
			if !r {
				earlyAdd()
			}
		}
		earlyAdd()
	}
	start := time.Now()
	var wg sync.WaitGroup

	var add, contains, remove, total, miss int64
	add, contains, remove, total, miss = 0, 0, 0, 0, 0

	for y := 0; y < threads; y++ {
		wg.Add(1)
		go func() {
			for x := 0; x < operations; x++ {
				if atomic.LoadInt64(&total) >= int64(operations) {
					break
				}
				op := rand.Intn(3)
				//op := 2
				if op == 0 {
					l.Add(rand.Intn(seedSize) + 1)
					atomic.AddInt64(&add, 1)
				} else if op == 1 {
					a := l.Remove(rand.Intn(seedSize) + 1)
					if a == false {
						atomic.AddInt64(&miss, 1)
					}
					atomic.AddInt64(&remove, 1)
				} else if op == 2 {
					l.Contains(rand.Intn(seedSize) + 1)
					atomic.AddInt64(&contains, 1)
				}
				atomic.AddInt64(&total, 1)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	stop := time.Now()
	elapsed := stop.Sub(start)
	ops := float64(operations) / elapsed.Seconds()
	fmt.Printf(";%d;%d;%d;%f;%.2f;%d;%d\n", add, remove, contains, elapsed.Seconds(), ops, miss, total)

}

func main() {

	configuration := Configuration{}
	err := gonfig.GetConf("config.json", &configuration)
	if err != nil {
		fmt.Println("Config file not found |", err)
		return
	}

	fmt.Printf("alg;exp;threads;add;remove;contains;duration;ops;miss;total\n")
	for y := 0; y < configuration.Experiments; y++ {

		for x := 2; x <= configuration.Threads; x = x + configuration.Jump {
			lC := list.NewCoarseList()
			fmt.Printf("Coarse;%d;%d", y, x)
			exec(lC, configuration.InitialSize, configuration.Operations, x)

			lF := list.NewFineList()
			fmt.Printf("Fine;%d;%d", y, x)
			exec(lF, configuration.InitialSize, configuration.Operations, x)

			lO := list.NewOptimisticList()
			fmt.Printf("Optimistic;%d;%d", y, x)
			exec(lO, configuration.InitialSize, configuration.Operations, x)

			lL := list.NewLazyList()
			fmt.Printf("Lazy;%d;%d", y, x)
			exec(lL, configuration.InitialSize, configuration.Operations, x)

			lN := list.NewNonBlockingList()
			fmt.Printf("LockFree;%d;%d", y, x)
			exec(lN, configuration.InitialSize, configuration.Operations, x)
		}
	}
}
