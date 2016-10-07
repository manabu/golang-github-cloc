package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/go-github/github"
)

func getCountLineOfCode(owner string, repo string) int {
	client := github.NewClient(nil)
	var codes []*github.WeeklyStats
	var loopcount = 0
	for loopcount < 5 {
		var err error
		codes, _, err = client.Repositories.ListCodeFrequency(owner, repo)
		fmt.Println(err)
		if _, ok := err.(*github.RateLimitError); ok {
			log.Println("hit rate limit")
		}
		if err != nil {
			fmt.Println("has error")
		} else {
			fmt.Println("no error")
			break
		}
		time.Sleep(1000 * time.Millisecond)
		loopcount = loopcount + 1
	}
	var num int
	num = -1
	for _, item := range codes {
		num = num + *item.Additions
		num = num - *item.Deletions
	}
	fmt.Println(num)
	return num
}
func handler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Hello, World")
}

func main() {
	fmt.Println(getCountLineOfCode("manabu", "golang-github-cloc"))
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8081", nil)
}
