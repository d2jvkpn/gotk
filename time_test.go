package gox

import (
	"fmt"
	"testing"
	"time"
)

func TestParseTime_t1(t *testing.T) {
	strs := []string{
		"2022-06-22",
		"05:32:22",
		"2022-06-22T05:32:03",
		"2022-06-22 05:32:03",
		"2022-06-22 05:32:03+08:00",
	}

	for _, str := range strs {
		tm, err := ParseDatetime(str)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(tm)
	}
}

func TestParseTime_t2(t *testing.T) {
	str := "2022-06-22T05:32:03"
	fmt.Println(time.ParseInLocation("2006-01-02T15:04:05", str, time.Local))
}
