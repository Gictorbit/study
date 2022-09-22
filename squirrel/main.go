package main

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/go-co-op/gocron"
	"os"
	"os/signal"
	"strconv"
	"time"
)

func main() {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Select("*").From("table").Join("")
	fmt.Println(query.ToSql())
	//init the loc
	loc, _ := time.LoadLocation("Asia/Tehran")

	//set timezone,
	now := time.Now().In(loc)
	fmt.Printf("%d-%02d-%02d", now.Year(), int(now.Month()), now.Day())
	s := gocron.NewScheduler(time.UTC)

	s.Every(1).Seconds().Do(func() {
		fmt.Println("hello")
	})
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		time.Sleep(time.Second * 10)
		fmt.Println("stop service")
		s.Stop()
		s.Clear()
	}()
	s.StartBlocking()
}

func test() {
	numbers := map[int]string{
		0: "Zero",
		1: "One",
		2: "Tow",
		3: "Three",
		4: "Four",
		5: "Five",
		6: "Six",
		7: "Seven",
		8: "Eight",
		9: "Nine",
	}
	inputs := make([]*int, 0)
	for i := 0; i < 3; i++ {
		var input int
		fmt.Scanln(&input)
		inputs = append(inputs, &input)
	}
	for _, num := range numbers {
		intNum, err := strconv.Atoi(num)
		if err != nil {
			panic(err)
		}
		fmt.Println(numbers[intNum])
	}
}
