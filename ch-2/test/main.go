package main

import (
	"fmt"
	"net"
	"sort"
	"sync"
)

func worker(ports <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for p := range ports {
		address := fmt.Sprintf("storestero.com:%d", p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- p
	}
}

func main() {
	const maxPorts = 1024

	ports := make(chan int, 65000)
	results := make(chan int, maxPorts)
	var openports []int
	var wg sync.WaitGroup

	for i := 0; i < cap(ports); i++ {
		wg.Add(1)
		go worker(ports, results, &wg)
	}

	// Отправляем порты в канал
	go func() {
		for i := 1; i <= maxPorts; i++ {
			ports <- i
		}
		close(ports) // Закрываем канал после отправки всех портов
	}()

	// Читаем результаты из канала
	go func() {
		wg.Wait()
		close(results) // Закрываем канал после завершения всех горутин
	}()

	for port := range results {
		if port != 0 {
			fmt.Printf("%d\n", port)
			openports = append(openports, port)
		}
	}

	sort.Ints(openports)
	for _, port := range openports {
		fmt.Printf("%d open\n", port)
	}
}
