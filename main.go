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

	var add, contains, remove, addTime, containsTime, removeTime, miss int64
	add, contains, remove, addTime, containsTime, removeTime, miss = 0, 0, 0, 0, 0, 0, 0

	for y := 0; y < threads; y++ {
		wg.Add(1)
		go func() {
			for x := 0; x < operations/threads; x++ {

				op := rand.Intn(3)
				funcStart := time.Now()
				if op == 0 {
					l.Add(rand.Intn(seedSize) + 1)
					atomic.AddInt64(&add, 1)
					stopFunc := time.Now()
					atomic.AddInt64(&addTime, stopFunc.Sub(funcStart).Nanoseconds()/1000)
				} else if op == 1 {
					a := l.Remove(rand.Intn(seedSize) + 1)
					if a == false {
						atomic.AddInt64(&miss, 1)
					}
					atomic.AddInt64(&remove, 1)
					stopFunc := time.Now()
					atomic.AddInt64(&removeTime, stopFunc.Sub(funcStart).Nanoseconds()/1000)

				} else if op == 2 {
					l.Contains(rand.Intn(seedSize) + 1)
					atomic.AddInt64(&contains, 1)
					stopFunc := time.Now()
					atomic.AddInt64(&containsTime, stopFunc.Sub(funcStart).Nanoseconds()/1000)
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	stop := time.Now()
	elapsed := stop.Sub(start)
	ops := float64(operations) / elapsed.Seconds()
	fmt.Printf(";%d;%d;%d;%f;%.2f;%d\n", add, remove, contains, elapsed.Seconds(), ops, miss)

}

func main() {

	for x := 0; x <= 3; x = x + 1 {
		//lC := list.NewCoarseList()
		//fmt.Printf("CG;%d", x)
		//exec(lC, 100000, 100000, int(math.Pow(2.0, float64(x))))

		//lF := list.NewFineList()
		//fmt.Printf("FG;%d", x)
		//exec(lF, 30000, 30000, int(math.Pow(2.0, float64(x))))

		lO := list.NewOptimisticList()
		fmt.Printf("OP;%d", x)
		exec(lO, 30000, 30000, int(math.Pow(2.0, float64(x))))
	}
}
