package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Мінімальна вхідна задача
type Job struct {
	ID    int
	Value int
}

// Мінімальна структура результату[cite: 2]
type Result struct {
	JobID int
	Value int
}

// Worker отримує задачу з jobs, виконує її та передає результат у results[cite: 2]
func worker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan Job, results chan<- Result) {
	defer wg.Done()

	for {
		// Використання select для можливості сигналу скасування[cite: 2]
		select {
		case <-ctx.Done():
			// Stop processing[cite: 2]
			return
		case job, ok := <-jobs:
			if !ok {
				return // канал закритий, завершуємо роботу[cite: 2]
			}

			// Імітація обчислення (щоб помітити різницю в мілісекундах)
			time.Sleep(100 * time.Microsecond)

			// Наприклад, result = value * value[cite: 2]
			res := Result{
				JobID: job.ID,
				Value: job.Value * job.Value,
			}
			results <- res
		}
	}
}

// Функція для запуску одного тесту конвеєра
func runPipeline(tasks, workers, bufferSize int) int64 {
	// Для синхронізації застосовуємо context.Context[cite: 2]
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Channels можуть бути unbuffered або buffered залежно від опцій[cite: 2]
	jobs := make(chan Job, bufferSize)
	results := make(chan Result, bufferSize)

	var wg sync.WaitGroup

	// 1. Producer: створює задану кількість задач[cite: 2]
	go func() {
		for i := 1; i <= tasks; i++ {
			jobs <- Job{ID: i, Value: i}
		}
		// Після завершення генерації задач channel jobs повинен бути коректно закритий[cite: 2]
		close(jobs)
	}()

	// 2. Worker Pool: запуск необхідної кількості воркерів
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker(ctx, &wg, jobs, results)
	}

	// Горутина, яка закриває канал результатів, коли всі воркери відпрацюють
	go func() {
		wg.Wait() // Для синхронізації застосовуємо WaitGroup[cite: 2]
		close(results)
	}()

	// 3. Aggregator: отримує результати, підраховує кількість, накопичує загальний результат[cite: 2]
	start := time.Now()

	completedTasks := 0
	totalSum := 0

	for res := range results {
		completedTasks++
		totalSum += res.Value
	}

	elapsed := time.Since(start).Milliseconds() // визначає час завершення обробки[cite: 2]
	return elapsed
}

func main() {
	tasks := 10000
	// Передбачити можливість запускати 1, 2, 4, 8 workers з буфером 0 та 100[cite: 2]
	configs := []struct {
		workers int
		buffer  int
	}{
		{1, 0},
		{2, 0},
		{4, 0},
		{8, 0},
		{4, 100},
		{8, 100},
	}

	// Формування підсумкової статистики у вигляді таблиці[cite: 2]
	fmt.Printf("%-10s %-10s %-10s %-10s\n", "Tasks", "Workers", "Buffer", "Time, ms")
	fmt.Println("--------------------------------------------")

	for _, cfg := range configs {
		timeMs := runPipeline(tasks, cfg.workers, cfg.buffer)
		fmt.Printf("%-10d %-10d %-10d %-10d\n", tasks, cfg.workers, cfg.buffer, timeMs)
	}
}
