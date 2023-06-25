package main

import "sync"

var mutex = sync.Mutex{}

func main() {
	mutex.Lock()
	mutex.Unlock()
	mutex.Unlock() // panic 발생함
}
