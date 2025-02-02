package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Проверка ссылки из массива
func checkURL(url string, wg *sync.WaitGroup, results chan string) {
	defer wg.Done() // Отложенное закрытие функции

	start := time.Now()        // Начинаем отсчёт
	resp, err := http.Get(url) // Открываем ссылку
	if err != nil {            //проверка на не открытие ссылки
		results <- fmt.Sprintf("[ERROR] %s: %s", url, err)
		return // Выход из функции
	}

	defer resp.Body.Close()                                                      //Отложенное закрытие тела ответа
	endTime := time.Since(start)                                                 // Останавливаем секундомер
	results <- fmt.Sprintf("[%d] %s (время: %s)", resp.StatusCode, url, endTime) // Запись результата запроса информации из ссылки в канал
}

func main() {
	urls := []string{ // Срез сайтов для проверки
		"http://www.youtube.com",
		"http://www.google.com",
		"http://www.idmb.com", // Кривая ссылка
		"http://www.imdb.com", // Нормальная ссылка
		"http://www.github.com",
		"http://www.amazon.com",
	}
	results := make(chan string, 5) // Канал для результатов
	var wg sync.WaitGroup           // Вейтгрупп для синхронизации горутин

	for _, url := range urls { //Цикл, который проходит по срезу ЮРЛов
		wg.Add(1)                      // Инкрементируем вейтгрупп
		go checkURL(url, &wg, results) //Запускаем горутину
	}

	go func() { // Анонимная функция для закрытия канала
		wg.Wait()
		close(results)
	}()

	for res := range results { //Цикл для вывода всех данных из канала
		fmt.Println(res)
	}
}
