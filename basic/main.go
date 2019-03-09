package main

import (
	"fmt"
	"html/template"
	"math"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type Data struct {
	Name string
}

func main() {

	//fmt
	fmt.Println("print")
	fmt.Printf("%d, %f, %s, %q, %q \n", 1, 1.0, "string", 2, "number")

	//time
	now := time.Now()
	fmt.Println("time ", now)
	after := now.Add(time.Hour * 2)
	before := now.Add(time.Hour * 2)
	fmt.Println("after 2hours ", after)
	fmt.Println("before 2hours ", before)
	fmt.Println("After ", now.After(after))
	fmt.Println("Before ", now.Before(after))
	unix := time.Unix(0, 0)
	fmt.Println("Unix ", unix)
	fmt.Println("Unix nano ", unix.UnixNano())
	nanoUnix := now.UnixNano()
	fmt.Println("Now Unix nano ", nanoUnix)

	//math
	val := -123.12
	fmt.Println("Abs ", math.Abs(val))
	compare1 := 1
	compare2 := 2
	fmt.Println("Min ", math.Min(float64(compare1), float64(compare2)))
	fmt.Println("Max ", math.Max(float64(compare1), float64(compare2)))

	//math/rand
	fmt.Println("Rand number is", rand.Int())
	fmt.Println("Rand number is", rand.Intn(100))

	//strings
	text := "Hi, I am a tester."
	fmt.Printf("result: [%q]\n", text)
	fmt.Printf("Split %q\n", strings.Split(text, " "))
	fmt.Printf("Trim [%q]\n", strings.Join(strings.Split(text, " "), ""))

	http.HandleFunc("/", router)
	http.ListenAndServe(":8080", nil)
}

func router(res http.ResponseWriter, req *http.Request) {
	name := req.URL.Path[1:]
	fmt.Println(name)
	data := &Data{Name: name}
	t, err := template.ParseFiles("view/index.html")
	if err == nil {
		t.Execute(res, data)
	}
}
