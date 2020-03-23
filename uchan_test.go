package uchan

import (
	"fmt"
	"sync"
	"time"
	"testing"
)

func TestUnboundedChannel(t *testing.T) {

	uc := NewUnboundedChannel()


	var wg sync.WaitGroup
	lastVal := -1

	wg.Add(1)
	go func() {
		for v := range uc.Out {
			vi := v.(int)
			if lastVal + 1 != vi {
				t.Errorf("Unexpected value, expected %d got %d", lastVal + 1, vi)
			}
			lastVal = vi
			fmt.Printf("Read: %d\n", vi)
			time.Sleep(time.Millisecond * 50)
		}
		wg.Done()
	}()

	for i := 0; i < 100; i++ {
		fmt.Printf("Writing: %d\n", i)
		uc.In <- i
	}

	fmt.Println("Closing")
	uc.Close()
	fmt.Println("Waiting for reader to finish")
	wg.Wait()

	if lastVal != 99 {
		t.Errorf("Didn't get all values: last value %d", lastVal)
	}

}