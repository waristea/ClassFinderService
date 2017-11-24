// Uncomment main function for direct use
package main

import (
	"fmt"
	//"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"io"
	//"reflect"
	//"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	//"strings"
	//"os"
	"sync"
)

type Jar struct {
	sync.Mutex
	cookies map[string][]*http.Cookie
}

func NewJar() *Jar {
	jar := new(Jar)
	jar.cookies = make(map[string][]*http.Cookie)
	return jar
}

func (jar *Jar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	jar.Lock()
	if _, ok := jar.cookies[u.Host]; ok {
		for _, c := range cookies {
			jar.cookies[u.Host] = append(jar.cookies[u.Host], c)
		}
	} else {
		jar.cookies[u.Host] = cookies
	}
	jar.Unlock()
}

func (jar *Jar) Cookies(u *url.URL) []*http.Cookie {
	return jar.cookies[u.Host]
}

func NewJarClient() *http.Client {
	return &http.Client{
		Jar: NewJar(),
	}
}

// fetch()
// akan di ambil 1 halaman html yang berisi jadwal perkuliahan untuk
// satu jurusan tertentu sesuai nim yang dimasukkan.
// Halaman yang diambil akan ditulis ke 'w'
// Cara memakai :
// - Masukkan username, password, dan nim sebagai elemen map
func fetch(w http.ResponseWriter, r *http.Request) {
	formurl := "https://login.itb.ac.id"

	fmt.Println("Website is now served")

	// post variables
	// will be parsed from html form:
	// lt, execution, _eventId, appendix of formurl
	m := make(map[string]string)
	m["postUrl"] = formurl
	m["username"] = "" // masukkan username
	m["password"] = "" // masukkan password
	m["nim"] = ""      // masukkan nim

	client := NewJarClient()

	// get form
	req, _ := http.NewRequest("GET", formurl, nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	// get posted values from form
	b := resp.Body
	defer b.Close() // close Body when the function returns

	doc, err := html.Parse(b)
	if err != nil {
		fmt.Println("error")
	}

	parseForm(doc, m)

	// fill the form
	// setting url values
	urlValues := url.Values{}

	for k, v := range m {
		if k != "formUrl" && k != "nim" {
			urlValues.Add(k, v)
		}
	}

	// post the form
	resp, _ = client.PostForm(m["postUrl"], urlValues)

	// see results
	// response
	fmt.Printf("%v\n\n", m["postUrl"])
	fmt.Printf("====Response====\n%+v\n\n", resp)
	fmt.Printf("====Request of response above====\n%+v\n\n", resp.Request)

	// html
	/*
		b2 := resp.Body
		defer b2.Close() // close Body when the function returns

		bodyBytes, err2 := ioutil.ReadAll(b2)
		html := string(bodyBytes)
		if err2 != nil {
			fmt.Println("error")
		}

		io.WriteString(w, html)

		return
	*/

	passingUrl := "https://login.itb.ac.id/cas/login?service=https%3A%2F%2Fakademik.itb.ac.id%2Flogin%2FINA"

	targetUrl := "https://akademik.itb.ac.id/app/mahasiswa:" + m["nim"] + "+2017-1/kelas/jadwal/kuliah/list"

	req, _ = http.NewRequest("GET", passingUrl, nil)
	resp, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	req, _ = http.NewRequest("GET", targetUrl, nil)
	resp, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	b3 := resp.Body
	defer b3.Close() // close Body when the function returns

	bodyBytes, err2 := ioutil.ReadAll(b3)
	html := string(bodyBytes)
	if err2 != nil {
		fmt.Println("error")
	}

	io.WriteString(w, html)

	return

}

func parseForm(n *html.Node, m map[string]string) {
	// getting input variables
	if n.Type == html.ElementNode && n.Data == "input" {
		lm := make(map[string]string) // lm : local map

		for _, element := range n.Attr {
			lm[element.Key] = element.Val
		}
		// fetch value
		names := [...]string{"lt", "_eventId", "execution"}

		for _, name := range names {
			if lm["name"] == name {
				m[name] = lm["value"]
			}
		}
	}

	if n.Type == html.ElementNode && n.Data == "form" {
		lm := make(map[string]string)
		for _, element := range n.Attr {
			lm[element.Key] = element.Val
		}

		// fetch form url
		m["postUrl"] += lm["action"]
		//fmt.Println(lm["action"])
	}

	// getting post url

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parseForm(c, m)
	}
}
/*
func main() {
	http.HandleFunc("/", fetch)

	fmt.Println("Website is now served")
	http.ListenAndServe("127.0.0.1:49721", nil)
}
*/
