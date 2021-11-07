package cubdomain

import (
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (

	baseURL = "https://www.cubdomain.com/domains-registered-by-date/"
	page = 0
	matched bool

)

func GetCubdomain(keyword, date string) {

	split := strings.Split(date, ".")
	var correctMD []string
	for _, d := range split {
		if len (d) == 1 {
			d = "0" + d
		}
		correctMD = append(correctMD, d)
	}

	date = strings.Join(correctMD, "-")

	dateURL := baseURL + date

	for true {

		page ++

		html := request(dateURL)

		re := regexp.MustCompile(`<div class="row">((.|\n)*)<div class="text-center">`)
		submatchall := re.FindAllStringSubmatch(html, -1)

		// Break the loop when last page is reached.
		if !strings.Contains(strings.Join(submatchall[0][:], ""), "<div class=\"col-md-4\">") {
			break
		}

		// We use StripTagsPolicy to cut off all HTML tags.
		stripPolicy := bluemonday.StripTagsPolicy()
		for n, element := range submatchall {
			s := strings.Split(element[n], "<div class=\"col-md-4\">")
			for _, r := range s[:] {
				sanitizedString := stripPolicy.Sanitize(r)
				sanitizedString = strings.Trim(sanitizedString, "\n")

				// Print all domains.
				if keyword == "" {
					fmt.Println(sanitizedString)
				}

				// Print domains that contain specified keyword.
				if strings.Contains(sanitizedString, keyword) {
					matched = true
					fmt.Println(sanitizedString)
				}
			}
		}

	}

	if keyword != "" && matched == false {
		fmt.Println("No domains matching the given keyword were found")
	}

}

func request(url string) string {

	targetURL := url + "/" + strconv.Itoa(page)

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
	if resp.StatusCode != 200 {
		fmt.Printf("There is no records for the specified date :c\n")
		os.Exit(1)
	}

	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return string(html)
}
