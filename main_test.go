package main

import (
	"context"
	"reflect"
	"testing"
)

func TestContract(t *testing.T) {
	// ok  	sylo	0.398s	coverage: 80.5% of statements

	t.Run("SetList should set a new list in memory", func(t *testing.T) {
		con := newContract()

		res := con.SetList(context.Background(), []int{2, 1, 3})
		if res.free == true {
			t.Fatal("Should not have been a free transaction.")
		}

		if res.gas != 5 {
			t.Fatalf("Expected gas to be %d but was %d.", 5, res.gas)
		}

		res, err := con.ReadSortedList(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		// Checking that the new list has 3 elements instead of default.
		if len(res.data) != 3 {
			t.Fatalf("Expected length to be %d not %d", 3, len(res.data))
		}
	})

	t.Run("ReadSortedList", func(t *testing.T) {
		con := newContract()
		res := con.SetList(context.Background(), []int{2, 1, 3})
		if res.free == true {
			t.Fatal("Should not have been a free transaction.")
		}

		if res.gas != 5 {
			t.Fatalf("Expected gas to be %d but was %d.", 5, res.gas)
		}

		res, err := con.ReadSortedList(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		if res.gas != 21 {
			t.Fatalf("Expected gas to be %d but was %d.", 21, res.gas)
		}

		if !reflect.DeepEqual(res.data, []int{1, 2, 3}) {
			t.Fatalf("%+v is not sorted", res.data)
		}

		// This time it should be free..
		res, err = con.ReadSortedList(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		if res.free != true {
			t.Fatal("Should not have been a free transaction.")
		}
	})
}
