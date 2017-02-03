package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

const urlPrefix = "http://www.dead.net/show"

func locationByDate(date string,verbose bool,cont bool) {
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
		if cont {
			fmt.Println(date)
			return
		}
		if verbose {
			log.Fatal("no such show")
		}
		return
	}
	city_state := strings.TrimSuffix(string(location[2])," US")

	fmt.Printf("%s - %s, %s\n",date,string(location[1]),city_state)
	return
}

func locationByLines(filename string,verbose bool) {
	file, err := os.Open(filename)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

	re := regexp.MustCompile("(.*)\\s*(\\d{2}-\\d{2}-\\d{2})")
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
		line := scanner.Text()
		if verbose {
			log.Print(line)
		}
		parts := re.FindSubmatch([]byte(line))
		fmt.Print(string(parts[1]))
		if len(parts)>1 {
			locationByDate(string(parts[2]),verbose,true)
		}
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
}

func main() {
	date := flag.String("d","","The show's date YY-MM-DD")
	file := flag.String("f","","File which contains date lines of format YY-MM-DD")
	verbose := flag.Bool("v",false,"Verbose")
	flag.Parse()

	if *date!="" {
		locationByDate(*date,*verbose,false)
	} else if file!=nil {
		locationByLines(*file,*verbose)
	}

}
