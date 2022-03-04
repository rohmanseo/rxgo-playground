package operators

import (
	"context"
	"fmt"
	"github.com/reactivex/rxgo/v2"
)

func RunOperators() {
	mapOperator()
}

func mapOperator() {
	fmt.Println("============== Starting mapOperator ==============")
	//I will transform this int below to bool, if it is odd then it will be false,
	//otherwise if it is even then it will be true
	obs := rxgo.Just(1, 2, 3, 4, 5)().
		Map(func(ctx context.Context, i interface{}) (interface{}, error) {
			if i.(int)%2 == 0 {
				return true, nil
			}
			return false, nil
		})

	for item := range obs.Observe() {
		fmt.Printf("Receiving item: %v\n", item)
	}
	fmt.Println("============== Starting mapOperator ==============")
}
