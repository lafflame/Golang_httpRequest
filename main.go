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

	client := &http.Client{
		Timeout: 10 * time.Second, //если сервер не отвечает 10 секунд, то ловим таймаут
	}
	start := time.Now()          // Начинаем отсчёт
	resp, err := client.Get(url) // Открываем ссылку
	if err != nil {              //проверка на не открытие ссылки
		results <- fmt.Sprintf("[ERROR] %s: %s", url, err)
		return // Выход из функции
	}

	defer resp.Body.Close()                                                      //Отложенное закрытие тела ответа
	endTime := time.Since(start)                                                 // Останавливаем секундомер
	results <- fmt.Sprintf("[%d] %s (время: %s)", resp.StatusCode, url, endTime) // Запись результата запроса информации из ссылки в канал
}

func standartCheck(choice int) {
	var urls = []string{ // Срез сайтов для проверки
		"http://www.youtube.com",
		"http://www.google.com",
		"http://www.idmb.com", // Кривая ссылка
		"http://www.imdb.com", // Нормальная ссылка
		"http://www.github.com",
		"http://www.amazon.com",
	}
	results := make(chan string) // Канал для результатов
	var wg sync.WaitGroup        // Вейтгрупп для синхронизации горутин

	if choice == 1 {
		for _, url := range urls { //Цикл, который проходит по срезу ЮРЛов
			wg.Add(1)                      // Инкрементируем вейтгрупп
			go checkURL(url, &wg, results) //Запускаем горутину
		}
	} else {
		fmt.Println("Введите ссылку на сайт в формате: https://ссылка.com")
		var link string
		fmt.Scan(&link)

		wg.Add(1)
		go checkURL(link, &wg, results)
	}

	go func() { // Анонимная функция для закрытия канала
		wg.Wait()
		close(results)
	}()

	for res := range results { //Цикл для вывода всех данных из канала
		fmt.Println(res)
	}
}

func main() {
	var choice int
	fmt.Println("Выберите пункт:\n1.Проверка готовых ссылок\n2.Проверка другой ссылки")
	fmt.Scan(&choice)

	if choice == 1 || choice == 2 {
		standartCheck(choice)
	} else {
		fmt.Println("Buenos nochas, senior")
	}
}
