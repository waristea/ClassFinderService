

// Updated to support concurrency by adding a notification channel
//package main
/*
import (
	"fmt"
	"golang.org/x/net/html" // html tokenizer
	//"io/ioutil"
	"net/http"
	"os"
	"strings"
)
*/
/* ===getHref===
 * get href value from a token, then set 'ok' to true
 */
// Helper funciton to pull the href attribute from a Token
/*
func getHref(t html.Token) (ok bool, href string) {
	// Iterative over all of the Token's attributes until we find href
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}
	// "bare" return will return the variables (ok, href) as
	// defined in the function definition
	return
}
*/
/* ===crawl===
 * Extract all http** links from a given webpage
 * Flow :
 *  1. Get the html body 'b' with the url provided (L97~L113)
 *      Also, prepare for closure if failure happen when '1' failed
 *  2. Tokenize html body (split html onto iterable tokens 'z')(L115)
 *  3. Use for to iterate over tokens and assign the iterated
 *     item in 'tt' (L117~L118)
 *  4. Switch 'tt' based on token types (L122)
 *  5. Get z's current token 't' (L129)
 *  6. If t=='<a>' then get the url 'url'
 *  7. Pass url to provided channel 'ch'
 */
 /*
func crawl(url string, ch chan string, chFinished chan bool) {
	resp, err := http.Get(url)

	defer func() {
		// notify that we're done after this function
		chFinished <- true
	}()

	if err != nil {
		fmt.Println("ERROR: Failed to crawl \"" + url + "\"")
		return
	}

	b := resp.Body
	defer b.Close() // close Body when the function returns

	z := html.NewTokenizer(b)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return

		// e.g.: if tt == (<a>, <p>, <h1>, etc)
		case tt == html.StartTagToken:
			t := z.Token()

			// Check if the token is an <a> tag
			isAnchor := t.Data == "a"
			if !isAnchor {
				continue
			}

			// extract the href value, if there is one
			ok, url := getHref(t)
			if !ok {
				continue
			}
			//ch <- url
			// Make sure the url begins in "http"
			hasProto := strings.Index(url, "http") == 0
			if hasProto {
				ch <- url
			}

		}
	}
}
*/
/* ===main===
 * Flow :
 *  1. Make a map of found urls,
	    so there won't be loopback in the future
    2. Make channels
	    for finished Urls 'chUrls' and Notification 'chFinished'
    3. Crawl concurrently while running (4) after a channel is done
    4. Put Urls that comes from (3) and list it as done.
	    Then after all is finished, and we get notification from (3)
	    through 'chFinished', do c++ to wait for another input from
	    another 'seedUrls' member
    5. Print everthing in 'foundUrls'
*/
// How to use : (on cmd)
// go run main.go http://golang.org
/*
func main() {
	foundUrls := make(map[string]bool)
	seedUrls := os.Args[1:] // Accept a variable (url) from the terminal

	// Channels
	chUrls := make(chan string)
	chFinished := make(chan bool)

	// Kick off the crawl process (concurrently)
	for _, url := range seedUrls {
		go crawl(url, chUrls, chFinished)
	}

	// Subscribe to both channels
	for c := 0; c < len(seedUrls); {
		fmt.Println(len(seedUrls))
		fmt.Println(c)

		select {
		case url := <-chUrls:
			foundUrls[url] = true
		case <-chFinished:
			c++ // for what?
		}
	}

	// Print the results
	fmt.Println("\nFound", len(foundUrls), "unique urls:\n")

	for url, _ := range foundUrls {
		fmt.Println(" - " + url)
	}

	close(chUrls)
}
*/
