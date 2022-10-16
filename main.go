package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

const defaultLimit = 100000

type Parser interface {
	Split(r rune) bool
	ParseValues(r []string) (int64, int64, error)
}

type JsonParser struct{}

func (p JsonParser) Split(r rune) bool {
	return r == ':' || r == ',' || r == ' ' || r == '}'
}

func (p JsonParser) ParseValues(r []string) (int64, int64, error) {
	if len(r) < 5 {
		return 0, 0, errors.New("short line for splitting")
	}
	price, err := strconv.Atoi(r[3])
	if err != nil {
		return 0, 0, err
	}
	rating, err := strconv.Atoi(r[5])
	if err != nil {
		return 0, 0, err
	}
	return int64(price), int64(rating), nil
}

type CsvParser struct{}

func (p CsvParser) Split(r rune) bool {
	return r == ';'
}

func (p CsvParser) ParseValues(r []string) (int64, int64, error) {
	if len(r) < 3 {
		return 0, 0, errors.New("short line for splitting")
	}
	price, err := strconv.Atoi(r[1])
	if err != nil {
		return 0, 0, err
	}
	rating, err := strconv.Atoi(r[2])
	if err != nil {
		return 0, 0, err
	}
	return int64(price), int64(rating), nil
}

func initParser(filename string) Parser {
	y := strings.Split(filename, ".")
	expansion := y[len(y)-1]
	switch expansion {
	case "json":
		return JsonParser{}
	case "csv":
		return CsvParser{}
	default:
		log.Fatal("Unsupported file expansion")
	}
	return nil
}

func initFlags() (string, int) {
	filename := flag.String("filename", "file", "flag \"filename\" - file name")
	countLimit := flag.Int("limit", defaultLimit, "flag \"limit\" - the number of rows being processed at the same time")
	flag.Parse()
	return *filename, *countLimit
}

func main() {
	filename, countLimit := initFlags()
	var parser = initParser(filename)

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	limit := make(chan struct{}, countLimit)
	scanner := bufio.NewScanner(file)
	var maxPrice int64
	var maxRating int64

	for scanner.Scan() {
		values := strings.FieldsFunc(scanner.Text(), parser.Split)
		price, rating, err := parser.ParseValues(values)
		if err == nil {
			wg.Add(1)
			limit <- struct{}{}
			go func() {
				defer func() {
					wg.Done()
					<-limit
				}()
				maxPriceValue := atomic.LoadInt64(&maxPrice)
				if maxPriceValue < price {
					atomic.StoreInt64(&maxPrice, price)
				}
				maxRatingValue := atomic.LoadInt64(&maxRating)
				if maxRatingValue < rating {
					atomic.StoreInt64(&maxRating, rating)
				}
			}()
		}
	}
	wg.Wait()
	fmt.Println("Max price:", maxPrice, "Max rating:", maxRating)
}
