// print Hello/goodbye without order (just WaitGroup )
// package main

// import (
// 	"fmt"
// 	"sync"
// )

// func main() {
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
