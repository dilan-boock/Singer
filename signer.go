package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// сюда писать код
func ExecutePipeline(jobs ...job) {
	in := make(chan interface{})

	for _, jobItem := range jobs {
		out := make(chan interface{})

		// Запускаем каждую job в отдельной горутине
		go func(j job, in, out chan interface{}) {
			defer close(out) // Закрываем выходной канал после завершения работы
			j(in, out)
		}(jobItem, in, out)

		in = out // Выход текущей job становится входом следующей
	}

	// Дожидаемся завершения последней job
	for range in {
	}
}

func SingleHash(in chan interface{}, out chan interface{}) {
	wg := &sync.WaitGroup{}
	md5Mutex := &sync.Mutex{} // Для защиты DataSignerMd5

	for data := range in {
		wg.Add(1)
		go func(d interface{}) {
			defer wg.Done()

			str := fmt.Sprintf("%v", d)

			// Вычисляем crc32(data) параллельно с md5
			crc32Chan := make(chan string, 1)
			go func() {
				crc32Chan <- DataSignerCrc32(str)
			}()

			// Вычисляем md5 с защитой от перегрева
			md5Chan := make(chan string, 1)
			go func() {
				md5Mutex.Lock()
				defer md5Mutex.Unlock()
				md5Chan <- DataSignerMd5(str)
			}()

			md5Hash := <-md5Chan
			crc32md5 := DataSignerCrc32(md5Hash)
			crc32 := <-crc32Chan

			out <- crc32 + "~" + crc32md5
		}(data)
	}

	wg.Wait()
}

func MultiHash(in chan interface{}, out chan interface{}) {
	wg := &sync.WaitGroup{}

	for data := range in {
		wg.Add(1)
		go func(d interface{}) {
			defer wg.Done()

			str := fmt.Sprintf("%v", d)
			results := make([]string, 6)

			innerWg := &sync.WaitGroup{}
			for th := 0; th < 6; th++ {
				innerWg.Add(1)
				go func(i int) {
					defer innerWg.Done()
					results[i] = DataSignerCrc32(strconv.Itoa(i) + str)
				}(th)
			}
			innerWg.Wait()

			out <- strings.Join(results, "")
		}(data)
	}

	wg.Wait()
}

func CombineResults(in chan interface{}, out chan interface{}) {
	var results []string

	for data := range in {
		results = append(results, fmt.Sprintf("%v", data))
	}

	sort.Strings(results)
	out <- strings.Join(results, "_")
}
