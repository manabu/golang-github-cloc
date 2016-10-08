package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"strconv"
	"strings"
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
		if _, ok := err.(*github.RateLimitError); ok {
			log.Println("hit rate limit")
		}
		if err != nil {
			log.Println("has error " + "owner=[" + owner + "], repository=[" + repo + "]")
		} else {
			//fmt.Println("no error")
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
	return num
}
func handler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// GET
	//ownerandrepository := r.URL.Query().Get("ownerandrepository")
	// POST
	ownerandrepository := r.Form.Get("ownerandrepository")
	var cloc = ""
	if len(ownerandrepository) != 0 {
		var fields = strings.Split(ownerandrepository, "/")
		var owner = fields[0]
		var repository = fields[1]
		num := getCountLineOfCode(owner, repository)
		// escape ownerandrepository and num
		var escapedownerandrepository = html.EscapeString(ownerandrepository)
		if num != -1 {
			cloc = "<h1>Result</h1>" + escapedownerandrepository + "<br>" + html.EscapeString(strconv.Itoa(num)) + "<br>"
		} else {
			cloc = "<h1>Result</h1>" + escapedownerandrepository + "<br>" + "Repository not found or error happens" + "<br>"
		}
	}
	str := `<html><head><title>Count Line of Code on GitHub Repository</title></head>
<body><form id="inputform" method="POST">
<input type="text" name="ownerandrepository" autofocus />
<input type="submit" value="Count" />
</form>
ex. jquery/jquery<br>
Request: <a href="https://api.github.com/repos/jquery/jquery/stats/code_frequency">https://api.github.com/repos/jquery/jquery/stats/code_frequency</a><br>
`
	str = str + cloc +
		`
</body>
</html>`
	fmt.Fprintf(w, str)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8081", nil)
}
