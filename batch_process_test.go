package gox

import (
	"flag"
	"fmt"
	"testing"
	"time"
)

// go test -run  TestBatchProcess_t1 -- 5m
func TestBatchProcess_t1(t *testing.T) {
	var (
		delay time.Duration
		err   error
		args  []string
	)

	flag.Parse()

	delay = 180 * time.Millisecond
	if args = flag.Args(); len(args) > 0 {
		if delay, err = time.ParseDuration(args[0]); err != nil {
			t.Fatal(err)
		}
		if delay < 0 {
			t.Fatal("invalid delay argument")
		}
	}

	bp, _ := NewBatchProcess[string](10, time.Second, func(items []string) {
		fmt.Printf(
			">>> %s process %d items: %#v\n",
			time.Now().Format("2006-01-02T15:05:05.000-07:00"), len(items), items,
		)
		// time.Now().Format(time.RFC3339)
	})

	go func() {
		for i := 0; i < 1000; i++ {
			if err = bp.Recv(fmt.Sprintf("X-%d", i)); err != nil {
				fmt.Printf("!!! Recv error %d: %v\n", i, err)
				break
			}
			fmt.Println("-->", i)
			time.Sleep(delay)
		}
	}()

	time.Sleep(10 * time.Second)
	bp.Down()
	time.Sleep(2 * time.Second)
	fmt.Println("exit")
}

// go test -run  TestBatchProcess_t2
func TestBatchProcess_t2(t *testing.T) {
	var (
		delay time.Duration
		rps   int64
		err   error
		args  []string
	)

	flag.Parse()

	delay = 100 * time.Millisecond

	if args = flag.Args(); len(args) > 0 {
		if delay, err = time.ParseDuration(args[0]); err != nil {
			t.Fatal(err)
		}
		if delay < 0 {
			t.Fatal("invalid delay argument")
		}
	}

	delay = 20 * time.Millisecond

	bp, _ := NewBatchProcess[time.Time](100, time.Second, func(items []time.Time) {
		length := len(items)
		if length <= 1 {
			return
		}

		rps = int64(length) * 1e9 / int64(items[length-1].Sub(items[0]))

		fmt.Printf(">>> %s %d\n", time.Now().Format("2006-01-02T15:05:05.000-07:00"), rps)
	})

	go func() {
		for i := 0; i < 1000; i++ {
			if err = bp.Recv(time.Now()); err != nil {
				break
			}
			// fmt.Println("-->", i)
			time.Sleep(delay)
		}
	}()

	time.Sleep(10 * time.Second)
	bp.Down()
	time.Sleep(2 * time.Second)
	fmt.Println("exit")

}
