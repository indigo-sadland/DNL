package main

import (
	"flag"
	"fmt"
	"github.com/indigo-sadland/dnl/whoistory"
	"os"
	"regexp"
)

var (

	keyword 	*string
	date		*string

)


func Intro() {

	fmt.Println(`
                                           
		Yes, I'm a librarian. Why you ask?  
		^..^      /			                           
		/_/\_____/                                     
		   /\   /\
		  /  \ /  \
	`)

}

func init() {

	Intro()

	keyword = flag.String("keyword", "", "Pattern to use in search.\n" +
		"Will print all domains for the given date if not specified.")
	date = flag.String("date", "", "Registration date of domains in format yyyy.mm.dd")
	flag.Parse()

	if *date == "" {
		fmt.Printf("Date must be specified\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Regex to check correct date format
	re := regexp.MustCompile("((19|20)\\d\\d).(0?[1-9]|1[012]).(0?[1-9]|[12][0-9]|3[01])")
	if !re.MatchString(*date) {
		fmt.Println("Date format is incorrect. Example: 2021.05.20")
		os.Exit(1)
	}

}

func main()  {

	whoistory.GetWhoistory(*keyword, *date)
}
