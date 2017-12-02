package main

import (
	"fmt"
	//"github.com/PuerkitoBio/goquery"
	//"bufio"
	//"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/net/html"
	"io"
	//"reflect"
	"bytes"
	//"io/ioutil"
	"log"
	"net/http"
	"net/url"
	//"strings"
	//"os"
	"sync"
	//"syscall"
	//"container/list"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// for saving room and time for each subject

type schedule struct {
	//varname type			struct_tag
	ID         bson.ObjectId `json:"id" bson:"_id"`
	Code       string        `json:"code" bson:"code"`
	Subject    string        `json:"subject" bson:"subject"`
	SKS        string        `json:"sks" bson:"sks"`
	ClassNum   string        `json:"class-num" bson:"class-number"`
	Lecturer   []string      `json:"lecturer" bson:"lecturer"`
	StudentAmt int           `json:"student-amount" bson:"student-amount"`
	Daytime    []daytime     `json:"daytime" bson:"daytime"`
}

//
type daytime struct {
	//varname type			struct_tag
	Day       int       `json:"day" bson:"day"` // 1~7 : Mon~Sun
	Room      string    `json:"room" bson:"room"`
	TimeStart time.Time `json:"timestart" bson:"timestart"`
	TimeEnd   time.Time `json:"timeend" bson:"timeend"`
}

type Jar struct {
	sync.Mutex
	cookies map[string][]*http.Cookie
}

type User struct {
	username string
	nim      string
	password string
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

	//fmt.Println(m["username"])
	//fmt.Println(m["nim"])
	//fmt.Println(m["password"])

	passingUrl := "https://login.itb.ac.id/cas/login?service=https%3A%2F%2Fakademik.itb.ac.id%2Flogin%2FINA"
	targetUrl := "https://akademik.itb.ac.id/app/mahasiswa:" + m["nim"] + "+2017-1/kelas/jadwal/kuliah/list"

	// mis:
	//https://akademik.itb.ac.id/app/mahasiswa:18215011+2017-1/kelas/jadwal/kuliah/list?fakultas=SF&prodi=116

	client := NewJarClient()
	fakultas := []string{}
	prodi := make(map[string][]string)
	//jadwal := make(map[string][]schedule)

	//prodiFakultas := make(map[string][]string)

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
	//fmt.Println(resp)

	// see results
	// response
	/*
		fmt.Printf("%v\n\n", m["postUrl"])
		fmt.Printf("====Response====\n%+v\n\n", resp)
		fmt.Printf("====Request of response above====\n%+v\n\n", resp.Request)
	*/
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

	req, _ = http.NewRequest("GET", passingUrl, nil)
	resp, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(resp)

	req, _ = http.NewRequest("GET", targetUrl, nil)
	resp, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	// Target html
	b3 := resp.Body
	defer b3.Close() // close Body when the function returns

	// Write to url
	/*
		bodyBytes, err2 := ioutil.ReadAll(b3)
		html3 := string(bodyBytes)
		ioutil.WriteFile("aim1.html", bodyBytes, 0644)
		if err2 != nil {
			fmt.Println("error")
		}

		io.WriteString(w, html3)
	*/
	// Parse 'fakultas' data
	doc3, err3 := html.Parse(b3)
	if err3 != nil {
		fmt.Println("error")
	}

	//fmt.Println(doc3)
	//fmt.Println(b3)
	parseFakultas(doc3, &fakultas)
	//fmt.Println(fakultas)

	//targetUrl+'?fakultas=FMIPA&prodi=101'

	// ambil prodi untuk semua fakultas
	for i := 0; i < len(fakultas); i++ {
		prodiUrl := targetUrl + "?fakultas=" + fakultas[i]

		req4, _ := http.NewRequest("GET", prodiUrl, nil)
		resp, err = client.Do(req4)
		if err != nil {
			log.Fatal(err)
		}
		// Target html
		b4 := resp.Body
		//fmt.Println(resp.StatusCode)
		defer b4.Close() // close Body when the function returns

		doc4, err4 := html.Parse(b4)
		if err4 != nil {
			fmt.Println("error")
		}

		p := []string{}
		parseProdi(doc4, &p)

		prodi[fakultas[i]] = p

		fmt.Println(prodi[fakultas[i]])

	}

	// ambil data jadwal (akhirnya)
	for i := 0; i < len(fakultas); i++ {
		for j := 0; j < len(prodi[fakultas[i]]); j++ {
			jadwalUrl := targetUrl + "?fakultas=" + fakultas[i] + "&prodi=" + prodi[fakultas[i]][j]

			req5, _ := http.NewRequest("GET", jadwalUrl, nil)
			resp, err = client.Do(req5)
			if err != nil {
				log.Fatal(err)
			}
			// Target html
			b5 := resp.Body
			//fmt.Println(resp.StatusCode)
			defer b5.Close() // close Body when the function returns

			doc5, err5 := html.Parse(b5)
			if err5 != nil {
				fmt.Println("error")
			}

			p5 := []schedule{}
			parseJadwal(doc5, &p5)

			//prodi[fakultas[i]] = p

			//fmt.Println(prodi[fakultas[i]])
		}
	}

	return

}

func mapNodeAttr(n *html.Node) map[string]string {
	m := make(map[string]string)

	for _, element := range n.Attr {
		m[element.Key] = element.Val
	}
	return m
}

func parseForm(n *html.Node, m map[string]string) {
	// getting input variables
	if n.Type == html.ElementNode && n.Data == "input" {
		lm := mapNodeAttr(n)
		// fetch value
		names := [...]string{"lt", "_eventId", "execution"}

		for _, name := range names {
			if lm["name"] == name {
				m[name] = lm["value"]
			}
		}
	}

	if n.Type == html.ElementNode && n.Data == "form" {
		lm := mapNodeAttr(n)

		// fetch form url
		m["postUrl"] += lm["action"]
		//fmt.Println(lm["action"])
	}

	// getting post url

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parseForm(c, m)
	}
}

func parseFakultas(n *html.Node, fakultas *[]string) {
	if n.Type == html.ElementNode && n.Data == "select" {
		//fmt.Println("=====" + n.Data + "====")

		lm := mapNodeAttr(n)

		if lm["id"] == "fakultas" {
			for d := n.FirstChild; d != nil; d = d.NextSibling {
				//fmt.Println(d.Data)
				lm := mapNodeAttr(d)
				//fmt.Println(lm["value"])
				*fakultas = append(*fakultas, lm["value"])
				//fmt.Println(*fakultas)
			}
			return
		}
	}

	// getting traverse recursively
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parseFakultas(c, fakultas)
	}
}

func parseProdi(n *html.Node, prodi *[]string) {
	//fmt.Println("=====" + n.Data + "====")

	if n.Type == html.ElementNode && n.Data == "select" {
		//fmt.Println("=====" + n.Data + "====")

		lm := mapNodeAttr(n)

		if lm["id"] == "prodi" {
			for d := n.FirstChild; d != nil; d = d.NextSibling {
				//fmt.Println(d.Data)
				lm := mapNodeAttr(d)
				//fmt.Println(lm)
				/*
					fmt.Println("====Label====")
					fmt.Println(lm["label"])
					fmt.Println("====End Label====")
				*/
				if lm["label"] == "Sarjana" {
					//fmt.Println("Entered 1")
					for e := d.FirstChild; e != nil; e = e.NextSibling {
						lm := mapNodeAttr(e)
						//fmt.Println(lm)
						*prodi = append(*prodi, lm["value"])
					}
				}
				//fmt.Println(*prodi)
			}
			return
		}
	}

	// getting traverse recursively
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parseProdi(c, prodi)
	}
}

func parseJadwal(n *html.Node, jadwal *[]schedule) {
	//fmt.Println("=====" + n.Data + "====")

	if n.Type == html.ElementNode && n.Data == "tbody" {
		// empty, numbering, class code, class name, SKS,
		// class number,lecturer, atendee, schedule
		for d := n.FirstChild; d != nil; d = d.NextSibling {
			for e := d.FirstChild; e != nil; e = e.NextSibling {
				if e.Data == "td" {
					for f := e.FirstChild; f != nil; f = f.NextSibling {
						// cari cara untuk memisahkan antar elemennya!
						if f.Data == "ul" {
							for g := f.FirstChild; g != nil; g = g.NextSibling {
								text := renderNode(g.FirstChild)
								fmt.Print("Lecturer : ") // or jadwal
								fmt.Println(text)
							}
						} else {
							text := renderNode(f)
							fmt.Print("Text : ")
							fmt.Println(text)
						}

					}
				}
			}
		}
		return

	}

	// getting traverse recursively
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parseJadwal(c, jadwal)
	}
}

func renderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	return buf.String()
}

var user = User{}

func main() {

	// Masukkan credentials
	/*
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter username: ")
		user.username, _ = reader.ReadString('\n')

		reader = bufio.NewReader(os.Stdin)
		fmt.Print("Enter nim: ")
		user.nim, _ = reader.ReadString('\n')

		fmt.Println("Enter password: ")
		password, _ := terminal.ReadPassword(int(syscall.Stdin))
		user.password = string(password)
	*/
	http.HandleFunc("/", fetch)

	fmt.Println("Website is now served")
	http.ListenAndServe("127.0.0.1:49721", nil)
}
