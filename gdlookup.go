package main

import (
	"fmt"
	"flag"
	"time"
	"log"
	"strings"
	"net/http"
	"io/ioutil"
	"regexp"
)

const urlPrefix = "http://www.dead.net/show"

func locationByDate(date string,verbose bool) {
	re := regexp.MustCompile(">([^>]+)</a></h3>.+>([^<]+)</a></h4>")
	t, err := time.Parse("06-01-02",date)
	if err!=nil {
		log.Fatal("Illegal date",date)
	}

	searchString := fmt.Sprintf("%s/%s-%d-%d",urlPrefix,
		strings.ToLower(t.Month().String()),int(t.Day()),t.Year());
	if verbose {
		fmt.Println("Searching for ",searchString)
	}

	resp, err := http.Get(searchString)
	if err!=nil {
		log.Fatal("Failed getting",searchString,err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err!=nil {
		log.Fatal("not found",err)
	}
	location := re.FindSubmatch([]byte(body))
	if location==nil {
		if verbose {
			log.Fatal("no such show")
		}
		return
	}
	fmt.Printf("%s, %s",string(location[1]),string(location[2]))
	return
}

func main() {
	date := flag.String("d","","The show's date YY-MM-DD")
	verbose := flag.Bool("v",false,"Verbose")
	flag.Parse()

	locationByDate(*date,*verbose)
}
