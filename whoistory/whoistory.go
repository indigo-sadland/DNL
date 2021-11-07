package whoistory

import (
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

// GetWhoistory makes request to whoishistory.com and parses request.
func GetWhoistory(keyword, date string) {

	baseURL := "https://whoistory.com/"

	split := strings.Split(date, ".")
	var correctMD []string
	for _, d := range split {
		if len (d) == 1 {
			d = "0" + d
		}
		correctMD = append(correctMD, d)
	}

	trail := strings.Join(correctMD, "/")
	targetURL := baseURL + trail

	// Set up http client and make request.
	httpClient := &http.Client{}
	req, _ := http.NewRequest("GET", targetURL, nil)
	req.Header.Set("Content-Type", "text/html; charset=UTF-8")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:92.0) Gecko/20100101 Firefox/92.0")
	time.Sleep(time.Millisecond * 300) // Sleep for the sake of the website.
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	if !strings.Contains(string(html), "Домены, зарегистрированные") {
		fmt.Printf("There is no records for the specified date :c\n")
		os.Exit(1)

	}

	re := regexp.MustCompile(`<div class="left">((.|\n)*)<div class="right">`)
	submatchall := re.FindAllStringSubmatch(string(html), -1)

	// We use StripTagsPolicy to cut off all HTML tags.
	stripPolicy := bluemonday.StripTagsPolicy()
	var matched bool
	for n, element := range submatchall {
		s := strings.Split(element[n], "<br />")
		for _, r := range s[:] {
			if strings.Contains(r,"</h2>") {
				split := strings.Split(r, "</h2>")
				r = split[1]
			}
			sanitizedString := stripPolicy.Sanitize(r)
			sanitizedString = strings.TrimPrefix(sanitizedString, "\n")

			// Print all domains.
			if keyword == "" {
				fmt.Println(sanitizedString)
			}

			// Print domains that contain specified keyword.
			if strings.Contains(sanitizedString, keyword ) {
				matched = true
				fmt.Println(sanitizedString)
			}

		}
	}

	if matched == false {
		fmt.Println("No domains matching the given keyword were found")
	}

}
