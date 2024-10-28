// print Hello/goodbye without order (just WaitGroup )
// package main

// import (
// 	"fmt"
// 	"sync"
// )

// func main2() {
// 	var wg sync.WaitGroup

// 	wg.Add(2)
// 	go func() {
// 		fmt.Println("Hello")
// 		wg.Done()
// 	}()

// 	go func() {
// 		fmt.Println("Goodbye")
// 		wg.Done()
// 	}()
// 	wg.Wait()
// }

// ***************************************************
// print Hello/goodbye by order (just WaitGroup )

// package main

// import (
// 	"fmt"
// 	"sync"
// )

// func main() {
// 	var wg sync.WaitGroup

// 	wg.Add(1)
// 	go func() {
// 		fmt.Println("Hello")
// 		wg.Done()
// 	}()

// 	wg.Wait()

// 	wg.Add(1)
// 	go func() {
// 		fmt.Println("Goodbye")
// 		wg.Done()
// 	}()
// 	wg.Wait()
// }

// ***************************************************
// print Hello/goodbye by order (using two channel)

// package main

// import (
// 	"fmt"
// )

// func main() {
// 	ch := make(chan string)
// 	done := make(chan bool)

// 	go func() {
// 		ch <- "Hello"
// 		done <- true
// 	}()

// 	go func() {
// 		<-done
// 		ch <- "Goodbye"
// 	}()
// 	msg1 := <-ch
// 	fmt.Println(msg1)

// 	msg2 := <-ch
// 	fmt.Println(msg2)
// }

// ***************************************************
// print Hello/goodbye by order (using two channel)

// package main

// import (
// 	"fmt"
// )

// func main() {
// 	ch1 := make(chan bool)
// 	ch2 := make(chan bool)

// 	go func() {
// 		fmt.Println("Hello")
// 		ch1 <- true
// 	}()

// 	go func() {
// 		<-ch1
// 		fmt.Println("Goodbye")
// 		ch2 <- true

// 	}()
// 	<-ch2
// }

// ***************************************************
// increment by atomic

// package main

// import (
// 	"fmt"
// 	"sync"
// 	"sync/atomic"
// )

// // var x = 0
// var x atomic.Int32

// func increment(wg *sync.WaitGroup) {
// 	// x = x + 1
// 	x.Add(1)
// 	wg.Done()
// }

// func main() {
// 	for i := 0; i < 5; i++ {
// 		var w sync.WaitGroup
// 		for i := 0; i < 1000; i++ {
// 			w.Add(1)
// 			go increment(&w)
// 		}
// 		w.Wait()
// 		// fmt.Println("final value of x", x)
// 		fmt.Println("final value of x", x.Load())
// 	}
// }

// ***************************************************
// checking speed

// package main

// import (
// 	"fmt"
// 	"time"
// )

// func main(){
// 	start:=time.Now()
// 	defer func(){
// 		fmt.Println(time.Since(start))
// 	}()

// 	list:=[]string{"a","b","c","d"}
// 	for _,v:=range list{
// 		go output(v)
// 	}
// 	time.Sleep(time.Second)
// }

// func output(v string){
// 	fmt.Println(v)
// 	time.Sleep(time.Second)
// }
