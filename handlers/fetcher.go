package handlers

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

var (
	UrlRegex         = regexp.MustCompile(`^(http(s):\/\/.)[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)$`)
	ExcludedPrefixes = []string{
		"mailto",
		"tel",
		"whatsapp",
		"javascript",
		"skype",
		"viber",
		"tg",
		"file",
		"data",
		"chrome",
		"about",
	}
	TagsToCheck = []string{
		"a[href]",
		"iframe[src]",
		"source[src]",
		"source[srcset]",
		"source[data-srcset]",
		"img[src]",
		"video[src]",
		"audio[src]",
		"embed[src]",
		"link[href]",
		"script[src]",
	}
)

func Start(site string) {
	// Checking if the URL is valid or not if not then return an error
	if !UrlRegex.MatchString(site) {
		Error.Println("The URL is not valid")
		return
	}

	fmt.Println("Starting to investigate: " + site)

	// Creating a new slice to store the ch
	ch := make([]string, 0)

	// Creating a new collector
	coll := colly.NewCollector(
		colly.MaxDepth(1),
	)

	/* Implementation changed. It's much easier to add new elements this way   */
	coll.OnHTML(strings.Join(TagsToCheck, ","), func(e *colly.HTMLElement) {
		switch {
		case e.Attr("href") != "":
			ch = append(ch, e.Attr("href"))
		case e.Attr("src") != "":
			ch = append(ch, e.Attr("src"))
		case e.Attr("srcset") != "":
			ch = append(ch, e.Attr("srcset"))
		case e.Attr("data-srcset") != "":
			ch = append(ch, e.Attr("data-srcset"))
		}
	})

	// Visiting the site
	coll.Visit(site)

	hostname, scheme, err := parseHost(site)
	if err != nil {
		Error.Println(err)
		return
	}

	urls := process(&ch, hostname, scheme)

	fmt.Println("Found", len(urls), "URLs in"+site)
}

/*
- Function takes the URL and returns the hostname and scheme of the website using `url` package
*/
func parseHost(site string) (string, string, error) {
	u, err := url.Parse(site)
	if err != nil {
		return "", "", err
	}

	return u.Hostname(), u.Scheme + "://", nil
}

/*
- Function takes slice of all urls, hostname and scheme of the website
- returns a slice of formatted URLs after processing them
*/
func process(urls *[]string, hostname string, scheme string) []string {
	/*
		- Since we are using `colly` to get the URLs, it will return the URLS but we neded to format them to the correct format :
			- `https://xxx/xxx` - Leave that way
			- `/#xxx/xxx` - Add `https://hostname` to the beginning
			- `//xxx/xxx` - Add `https:` to the beginning
			- `/xxx/xxx` - Add `https://hostname` to the beginning
			- `xxx/xxx` - Add `https://hostname/` to the beginning

		- Also, it will return the same URL multiple times, so we need to remove the duplicates.
	*/

	/* Creating a new slice to store the URLs without duplicates */
	formattedUrls := make([]string, 0)

	/* Creating a new map for the duplicates check (key: URL, value: true) */
	duplicates := make(map[string]bool)

	for _, url := range *urls {
		if len(url) == 0 {
			continue
		}

		// If URL starts with `https://` or `http://` then add it to the slice and duplicate check map
		if strings.HasPrefix(url, scheme) || strings.HasPrefix(url, "http://") {
			if _, ok := duplicates[url]; !ok {
				duplicates[url] = true
				formattedUrls = append(formattedUrls, url)
			}
			continue
		}

		/* If URL starts with excluded prefixes, then skip it */
		for _, prefix := range ExcludedPrefixes {
			if strings.HasPrefix(url, prefix) {
				continue
			}
		}

		switch {
		case len(url) == 1 && (url[0] == '#' || url[0] == '/' || url[0] == '.'):
			// If URL starts with `#` or `/` or `.` then add `https://hostname` to the beginning
			url = fmt.Sprintf("%s%s%s", scheme, hostname, url)
		case url[0] == '/' && url[1] != '/':
			// If URL starts with `/` then add `hostname` to the beginning
			url = fmt.Sprintf("%s%s%s", scheme, hostname, url)
		case strings.HasPrefix(url, "//"):
			// If URL starts with `//` then add `https:` to the beginning
			url = fmt.Sprintf("%s%s", "https:", url)
		case url[0] != '/' && url[0] != '#':
			// If URL starts with anything else then add `https://hostname/` to the beginning
			url = strings.Join([]string{scheme, hostname, "/", url}, "")
		default:
			continue
		}

		if _, ok := duplicates[url]; !ok {
			duplicates[url] = true
			formattedUrls = append(formattedUrls, url)
		}
	}

	return formattedUrls
}
