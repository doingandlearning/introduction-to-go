package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func makeCounter() func() int {
	Count := 0
	return func() int {
		Count++
		return Count
	}
}

func invokeAndLog(f func()) {
	start := time.Now()
	f()
	log.Printf("Function took %v", time.Since(start))
}

func main() {
	invokeAndLog(func() { makeCounter() })
	till1 := makeCounter()
	till2 := makeCounter()

	fmt.Println(till1())
	fmt.Println(till1())
	fmt.Println(till1())
	fmt.Println(till1())

	fmt.Println(till2())
	fmt.Println(till2())
	fmt.Println(till1())
	fmt.Println(till1())
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}
