package main

import "fmt"

// // starter is
// func starter(wg *sync.WaitGroup) {
// 	fmt.Println("This is starter on call")
// 	defer wg.Done()
// }

// // follow is
// func follow() {
// 	fmt.Println("This is the follower call")
// }

func starter(entry chan<- string, message string) {
	entry <- message
}

func follower(sender <-chan string, receiver chan<- string) {
	message := <-sender
	receiver <- message
}

func main() {
	send := make(chan string, 1)
	receive := make(chan string, 1)
	starter(send, "Successfully sent message")
	follower(send, receive)
	fmt.Println(<-receive)
}
