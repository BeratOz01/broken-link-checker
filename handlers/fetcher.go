package handlers

import (
	"fmt"
	"net/url"
	"regexp"

	"github.com/gocolly/colly"
)

var (
	UrlRegex = regexp.MustCompile(`^(http(s):\/\/.)[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)$`)
)

func Start(site string) {
	// Checking if the URL is valid or not if not then return an error
	if !UrlRegex.MatchString(site) {
		Error.Println("The URL is not valid")
		return
	}

	fmt.Println("Starting to investigate: " + site)

	// Creating a new slice to store the urls
	urls := make([]string, 0)

	// Creating a new collector
	c := colly.NewCollector()

	// On every a element which has href attribute add it to the urls slice
	c.OnHTML("a", func(e *colly.HTMLElement) {
		urls = append(urls, e.Attr("href"))
	})

	c.OnHTML("picture", func(e *colly.HTMLElement) {
		urls = append(urls, e.Attr("srcset"))
	})

	// On every img element which has src attribute add it to the urls slice
	c.OnHTML("img", func(e *colly.HTMLElement) {
		if e.Attr("src") == "" && e.Attr("data-src") == "" {
			return
		} else if e.Attr("src") == "" {
			urls = append(urls, e.Attr("data-src"))
			return
		}

		urls = append(urls, e.Attr("src"))
	})

	c.OnHTML("source", func(e *colly.HTMLElement) {
		urls = append(urls, e.Attr("data-srcset"))
	})

	c.Visit(site)

	hostname, err := getHostname(site)
	if err != nil {
		Error.Println(err)
		return
	}

	process(&urls, hostname)

}

// It takes a string, and returns a string and an error
func getHostname(site string) (string, error) {
	u, err := url.Parse(site)
	if err != nil {
		return "", err
	}
	return u.Hostname(), nil
}

// This is the function that formats all the URLs to the correct format (`https://xxx/xxx...`) and removes the duplicates
func process(urls *[]string, hostname string) {

}
