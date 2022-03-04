package observables

import (
	"context"
	"errors"
	"fmt"
	"github.com/reactivex/rxgo/v2"
	"strconv"
	"time"
)

func RunObservables() {
	justObservable()
	producerObservable()
	obsFromChannel()
	connectableObservable()
}

//create observable using just
func justObservable() {
	fmt.Println("============== Starting just observable ==============")

	//just observable allowing to send multiple items with different type
	obs := rxgo.Just(1, 2, 3, 4, 5, 6, "dor", true, errors.New("Opps, something went wrong"))()
	ch := obs.Observe()

	for item := range ch {
		fmt.Println(item)
	}
	fmt.Println("============== End of just observable ==============")
}

func producerObservable() {
	fmt.Println("============== Starting producer observable ==============")
	producerFunction := func(ctx context.Context, next chan<- rxgo.Item) {
		next <- rxgo.Of(1)
		next <- rxgo.Of(2)
		next <- rxgo.Of(3)
		next <- rxgo.Of("dor")
		next <- rxgo.Error(errors.New("Whoops, you should be handle this error"))
		next <- rxgo.Of(true)
	}
	obs := rxgo.Create([]rxgo.Producer{producerFunction, producerFunction})
	ch := obs.Observe()
	for item := range ch {
		fmt.Println(item)
		if item.Error() {
			fmt.Println("OK, I've already handle it")
		}
	}
	fmt.Println("============== End of producer observable ==============")
}

//obs from channel is hot, it is really hot!
func obsFromChannel() {
	fmt.Println("============== Starting obsFromChannel ==============")
	ch := make(chan rxgo.Item)

	go func() {
		for i := 0; i < 3; i++ {
			ch <- rxgo.Of("item for " + strconv.Itoa(i))
		}
		close(ch)
	}()
	obs := rxgo.FromChannel(ch)

	//we got items for the first consumer
	for item := range obs.Observe() {
		fmt.Print("Observer 1 printing item : ")
		fmt.Print(item.V)
		fmt.Print("\n")
	}

	//for the second consumer we got nothing, but it is ok, because it is hot!
	for item := range obs.Observe() {
		fmt.Print("Observer 2 printing item : ")
		fmt.Print(item.V)
		fmt.Print("\n")
	}

	time.Sleep(time.Second * 3)
	fmt.Println("============== End of obsFromChannel ==============")

}

/*
	connectable observable
	connectable observable will wait until all consumers ready to receive item
*/
func connectableObservable() {
	fmt.Println("============== Starting connectableObservable ==============")
	ch := make(chan rxgo.Item)

	go func() {
		for i := 0; i < 3; i++ {
			ch <- rxgo.Of("item for " + strconv.Itoa(i))
		}
		close(ch)
	}()

	obs := rxgo.FromChannel(ch, rxgo.WithPublishStrategy())

	//do on next will register consumer to publisher
	obs.DoOnNext(func(i interface{}) {
		fmt.Printf("Obs 1 receive item : %v\n", i)
	})

	time.Sleep(1 * time.Second)
	obs.DoOnNext(func(i interface{}) {
		fmt.Printf("Obs 2 receive item : %v\n", i)
	})

	obs.Connect(context.Background())
	time.Sleep(3 * time.Second)
	fmt.Println("============== End of connectableObservable ==============")

}
