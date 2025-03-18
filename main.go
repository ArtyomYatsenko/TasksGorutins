package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
)

func main() {
	//-------------

}

// Использование wg #1
func exercise1() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Hello from goroutine!")
	}()
	wg.Wait()
}

// Использование wg #2
func exercise2() {
	wg := sync.WaitGroup{}

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Printf("горутина - <%d>\n", i)
		}()
	}

	wg.Wait()
}

// Итерация по каналу, использование close()
func exercise3() { //
	incCh := func(c chan int) {
		for i := 0; i < 5; i++ {
			c <- i
		}
		close(c)
	}
	ch := make(chan int)
	go incCh(ch)
	sum := 0
	for e := range ch {
		sum += e
	}
	fmt.Println(sum)
}

// Использование mutex
func exercise4() {
	wg := sync.WaitGroup{}
	my := sync.Mutex{}
	count := 0

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			my.Lock()
			count++
			my.Unlock()
		}()
	}
	wg.Wait()
	fmt.Println(count)
}

// Использование atomic
func exercise5() {
	var p int64
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt64(&p, 1)

		}()
	}
	wg.Wait()
	fmt.Println(p)
}

// Запросы к url в горутинах
func exercise6(urls []string) map[string]string {
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	urlMap := make(map[string]string, len(urls))

	for _, url := range urls {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := http.Get(url)

			if err != nil {
				mu.Lock()
				urlMap[url] = "error"
				mu.Unlock()
				return
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			mu.Lock()
			urlMap[url] = string(body)[:100]
			mu.Unlock()
		}()
	}
	wg.Wait()
	return urlMap
}

// паттерн funIn
func merge(chls ...<-chan int) <-chan int {
	wg := sync.WaitGroup{}
	wg.Add(len(chls))

	chOut := make(chan int)

	send := func(ch <-chan int) {
		defer wg.Done()
		for i := range ch {
			chOut <- i
		}
	}

	for _, ch := range chls {
		go send(ch)
	}

	go func() {
		wg.Wait()
		close(chOut)
	}()

	return chOut
}

func exercise7() {

	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		ch1 <- 1
		ch1 <- 2
		close(ch1)
	}()

	go func() {
		ch2 <- 1
		ch2 <- 2
		close(ch2)
	}()

	ch := merge(ch1, ch2)

	for el := range ch {
		fmt.Println(el)
	}

}

// паттерн funOut

func split(ch <-chan int, n int) []<-chan int {
	chls := make([]chan int, 0, n)
	for i := 0; i < n; i++ {
		chls = append(chls, make(chan int))
	}

	toCannels := func(ch <-chan int, cnls []chan int) {

		defer func(cnls []chan int) {
			for _, c := range cnls {
				close(c)
			}
		}(chls)

		for {
			for _, c := range cnls {
				select {
				case e, open := <-ch:
					if !open {
						return
					}
					c <- e
				}
			}
		}
	}

	go toCannels(ch, chls)

	result := make([]<-chan int, len(chls))
	for i, c := range chls {
		result[i] = c
	}
	return result
}

func exercise8() {
	wg := sync.WaitGroup{}
	ch := make(chan int)
	go func() {
		for i := 0; i < 5; i++ {
			ch <- i
		}
		close(ch)
	}()

	c := split(ch, 5)
	wg.Add(len(c))
	for _, chl := range c {
		go func() {
			defer wg.Done()
			for i := range chl {
				fmt.Println(i)
			}
		}()
	}
}
