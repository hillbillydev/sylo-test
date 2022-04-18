package main

import (
	"context"
	"errors"
	"fmt"
	"os"
)

// operation defines all of the different operations that can occur on the blockchain.
type operation int

const (
	del    operation = -15
	read   operation = 1
	modify operation = 5
	write  operation = 20
)

// response is the final result that contains the data and the amount you are suppose to pay in gas.
type response struct {
	data []int
	gas  int
	free bool
}

// requestDetails is the details about all different operations that you have interacted with, during a request.
// It gets stored in the request context.
type requestDetails struct {
	gas  int
	free bool
}

// contract represents a smart contract.
type contract struct {
	data map[string][]int
}

// key follows the best practice from the Golang docs see: https://cs.opensource.google/go/go/+/refs/tags/go1.18.1:src/context/context.go;l=136
type key int

// see key description above
var requestKey key

// newContract creates a new contract with some data as default.
func newContract() *contract {
	return &contract{
		data: map[string][]int{
			"unsorted_list": {8, 4, 3, 0, 1, 3, 6, 4},
		},
	}
}

// SetList sets a new unsorted_list in memory while also removing the old sorted list.
func (c *contract) SetList(ctx context.Context, slice []int) *response {
	c.data["unsorted_list"] = slice
	delete(c.data, "sorted_list")

	return toResponse(addOperationsToContext(ctx, write, del), nil)
}

// ReadSortedList attempts to read the sorted list in memory, if it is not sorted it will sort it and save it for later reads.
// This method is free as long as the data has already been sorted, in case of having to sort it it will cost gas.
func (c *contract) ReadSortedList(ctx context.Context) (*response, error) {
	res, ok := c.data["sorted_list"]
	ctx = addOperationsToContext(ctx, read)
	if !ok {
		unsorted, ok := c.data["unsorted_list"]
		if !ok {
			return nil, errors.New("no unsorted list") // This should never happen.
		}
		res = c.sortList(ctx, unsorted)
		c.data["sorted_list"] = res
		ctx = addOperationsToContext(ctx, write)
	}

	return toResponse(ctx, res), nil
}

// sortList simply sorts a slice of integers, it uses bubblesort which is not the best kind of sorting.
// But it gets the job done for now.
func (c *contract) sortList(ctx context.Context, slice []int) []int {
	var (
		res = make([]int, len(slice))
		n   = len(res)
	)
	copy(res, slice)

	// Time complexity of O(n^2), pretty poor choice.
	for i := 0; i < n; i++ {
		for j := 0; j < n-1; j++ {
			if res[j] > res[j+1] {
				lesser := res[j+1]
				res[j+1] = res[j]
				res[j] = lesser
			}
		}
	}

	return res
}

func main() {
	var (
		c   = newContract()
		ctx = context.Background()
	)

	res, err := c.ReadSortedList(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("SortedList: %+v\n, Gas: %d\n Free: %t", res.data, res.gas, res.free)
}

// addOperationsToContext adds a operation to our request context.
func addOperationsToContext(ctx context.Context, ops ...operation) context.Context {
	details, ok := ctx.Value(requestKey).(*requestDetails)
	if !ok {
		details = &requestDetails{
			free: true, // Default free to true.
		}
	}

	for _, op := range ops {
		switch op {
		case modify, write:
			details.free = false
		}

		details.gas += int(op)
	}

	return context.WithValue(ctx, requestKey, details)
}

// toResponse converts our request context and data into a finalized response struct.
func toResponse(ctx context.Context, data []int) *response {
	details, ok := ctx.Value(requestKey).(*requestDetails)
	if !ok {
		return &response{} // Let's assume that everything is fine and that we have not done any work.
	}

	return &response{
		data: data,
		gas:  details.gas,
		free: details.free,
	}
}
