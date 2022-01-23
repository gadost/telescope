package alert

import "sync"

var wg sync.WaitGroup

func New() {
	wg.Add(1)
	go Alert()
}

func Alert() {
	defer wg.Done()
	Send()
}

func Send() {

}
