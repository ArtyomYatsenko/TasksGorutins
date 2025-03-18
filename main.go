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
	//-------
}

func FetchURLs(urls []string) map[string]string {
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

func incCh(c chan int) {
	for i := 1; i <= 5; i++ {
		c <- i
	}
	close(c)
}

func ex1() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Hello from goroutine!")
	}()
	wg.Wait()
}

func ex2() {
	wg := sync.WaitGroup{}

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Printf("горутина - <%d>\n", i)
		}()
	}

	wg.Wait()
}

func ex3() {
	ch := make(chan int)
	s := 0
	go incCh(ch)
	for e := range ch {
		s += e
	}
	fmt.Println(s)
}

func ex4() {
	wg := sync.WaitGroup{}
	my := sync.Mutex{}
	count := 0

	for i := 1; i <= 5; i++ {
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

func ex5() {
	var p int64
	wg := sync.WaitGroup{}
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt64(&p, 1)

		}()
	}
	wg.Wait()
	fmt.Println(p)
}

func ex6() { //funIn
	// Функция generate отправляет числа в канал
	generate := func(id int, ch chan<- int) {
		defer close(ch)
		for i := 1; i < 5; i++ {
			fmt.Println("Горутина", id, "отправила", i)
			ch <- i
		}
		// Закрываем канал после отправки всех данных
	}

	// Функция fanIn объединяет данные из нескольких каналов в один
	fanIn := func(c chan int, cs ...chan int) {
		wg := sync.WaitGroup{}
		wg.Add(len(cs))

		for _, ch := range cs {
			go func(ch chan int) {
				defer wg.Done()
				for el := range ch {
					c <- el
				}
			}(ch) // Передаем ch как аргумент, чтобы избежать захвата переменной
		}
		go func() {
			wg.Wait()
			close(c)
		}()
	}

	// Создаем выходной канал
	chOut := make(chan int)

	// Создаем срез каналов
	chls := make([]chan int, 3)
	for i := 0; i < 3; i++ {
		chls[i] = make(chan int)
		go generate(i, chls[i]) // Запускаем горутины для генерации данных
	}

	// Объединяем данные из всех каналов в один
	fanIn(chOut, chls...)

	// Читаем данные из выходного канала
	for e := range chOut {
		fmt.Println(e)
	}
}

func ex7() { //funOut
	wg := sync.WaitGroup{}

	worker := func(id int, chIn chan int, chOut chan int) {
		defer wg.Done()
		for ch := range chIn {
			v := ch * 2
			fmt.Printf("Горутина %d выполнила работу %d\n", id, v)
			chOut <- v
		}
	}

	chIn := make(chan int)
	chOut := make(chan int)

	countWorker := 5
	wg.Add(countWorker)

	//Пишем в канал
	go func() {
		for i := 0; i < 15; i++ {
			chIn <- i
		}
		close(chIn)
	}()

	for i := 0; i < countWorker; i++ {
		go worker(i, chIn, chOut)
	}

	go func() {
		wg.Wait()
		close(chOut)
	}()

	for el := range chOut {
		fmt.Println(el)
	}
}
