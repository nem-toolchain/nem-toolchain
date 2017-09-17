package example

import (
	"errors"
	"fmt"
)

const SearchWorkersNum uint = 100

func search() error {
	channels := make([]chan struct {
		string
		error
	}, SearchWorkersNum)
	for i := uint(0); i < SearchWorkersNum; i++ {
		ch := make(chan struct {
			string
			error
		})
		go func() {
			acc, err := GenerateAccount(TestnetId)
			ch <- struct {
				string
				error
			}{string: acc, error: err}
		}()
		channels[i] = ch
	}

	for i := uint(0); i < SearchWorkersNum; i++ {
		msg, ok := <-channels[i]
		if !ok {
			return errors.New("can't read account")
		}

		if msg.error != nil {
			return msg.error
		}
		fmt.Printf("Account %v\n", msg.string)
	}

	return nil
}
