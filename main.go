package main

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"./internal/list"
)

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

	fmt.Printf("alg;exp;threads;add;remove;contains;duration;ops;miss;total\n")
	for y := 0; y <= 4; y++ {

		for x := 1; x <= 3; x = x + 1 {
			threads := int(math.Pow(2.0, float64(x)))
			//lC := list.NewCoarseList()
			//fmt.Printf("Coarse;%d;%d", y, threads)
			//exec(lC, 100000, 100000, threads)

			//lF := list.NewFineList()
			//fmt.Printf("Fine;%d;%d", y, threads)
			//exec(lF, 100000, 100000, threads)

			//lO := list.NewOptimisticList()
			//fmt.Printf("Optimistic;%d;%d", y, threads)
			//exec(lO, 100000, 100000, threads)

			//lL := list.NewLazyList()
			//fmt.Printf("Lazy;%d;%d", y, threads)
			//exec(lL, 100000, 100000, threads)

			lN := list.NewNonBlockingList()
			fmt.Printf("LockFree;%d;%d", y, threads)
			exec(lN, 100000, 100000, threads)
		}
	}
}
