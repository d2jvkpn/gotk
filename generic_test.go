package gotk

import (
	"fmt"
	"testing"
)

func TestExtendJoin(t *testing.T) {
	type Item struct {
		Id    int64
		Value string
	}

	type Data struct {
		Id    int64
		XX    string
		Value string
	}

	var (
		dataList []Data
		items    []Item
	)

	dataList = []Data{
		{Id: 1, XX: "one"},
		{Id: 2, XX: "two"},
		{Id: 3, XX: "three"},
		{Id: 4, XX: "four"},
		{Id: 5, XX: "five"},
	}

	items = []Item{
		{Id: 1, Value: "I"},
		{Id: 3, Value: "II"},
		{Id: 4, Value: "IV"},
	}

	var mp map[int64]string = Slice2Map(
		items,
		func(item *Item) int64 { return item.Id },
		func(item *Item) string { return item.Value },
	)
	fmt.Printf("==> map: %+v\n", mp)

	JoinSlices(dataList, mp, func(d *Data) int64 { return d.Id }, func(d *Data, v string) {
		d.Value = v
	})

	for i := range dataList {
		fmt.Printf("==> %d: Id=%d, XX=%s, Value: %v\n", i, dataList[i].Id, dataList[i].XX, dataList[i].Value)
	}
}
