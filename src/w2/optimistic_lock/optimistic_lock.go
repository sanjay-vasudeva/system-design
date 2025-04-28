package optimisticlock

import (
	"fmt"
	"sync/atomic"
)

var count atomic.Int32

func IncrementCounter() {
	oldValue := count.Load()
	newValue := oldValue + 1
	if !count.CompareAndSwap(oldValue, newValue) {
		fmt.Println("Increment failed")
	}
}
