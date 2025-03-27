package main

// сюда писать код
func ExecutePipeline(jobs ...job) {
	// Создаем каналы для передачи данных между функциями
	in := make(chan interface{})
	out := make(chan interface{})

	// Запускаем каждую функцию в отдельной горутине
	for _, job := range jobs {
		go func(j job, in, out chan interface{}) {
			j(in, out)
			close(out) // Закрываем канал после завершения функции
		}(job, in, out)

		// Переопределяем каналы для следующей функции
		in = out
		out = make(chan interface{})
	}
}

func SingleHash(in chan interface{}, out chan interface{}) {

}

func MultiHash(in chan interface{}, out chan interface{}) {

}

func CombineResults(in chan interface{}, out chan interface{}) {

}
