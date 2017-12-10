package main

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"gopkg.in/mgo.v2/bson"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

// for saving room and time for each subject

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
func fetch(user User) []schedule {
	formurl := "https://login.itb.ac.id"

	// post variables
	// will be parsed from html form:
	// lt, execution, _eventId, appendix of formurl
	m := make(map[string]string)
	m["postUrl"] = formurl
	m["username"] = user.username // masukkan username
	m["password"] = user.password // masukkan password
	m["nim"] = user.nim           // masukkan nim

	passingUrl := "https://login.itb.ac.id/cas/login?service=https%3A%2F%2Fakademik.itb.ac.id%2Flogin%2FINA"
	targetUrl := "https://akademik.itb.ac.id/app/mahasiswa:" + m["nim"] + "+2017-1/kelas/jadwal/kuliah/list"

	client := NewJarClient()
	fakultas := []string{}
	prodi := make(map[string][]string)

	// get form

	var header http.Header = make(map[string][]string)
	header.Add("Content-Type", "application/x-www-form-urlencoded")

	doc := htmlParser("GET", formurl, client, header)
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
	client.PostForm(m["postUrl"], urlValues)

	doc = htmlParser("GET", passingUrl, client, nil)
	doc = htmlParser("GET", targetUrl, client, nil)

	parseFakultas(doc, &fakultas)

	// ambil prodi untuk semua fakultas
	for i := 0; i < len(fakultas); i++ {
		prodiUrl := targetUrl + "?fakultas=" + fakultas[i]

		doc := htmlParser("GET", prodiUrl, client, nil)

		p := []string{}
		parseProdi(doc, &p)
		prodi[fakultas[i]] = p
	}

	// ambil data jadwal (akhirnya)
	mainSchedule := []schedule{}
	for i := 0; i < len(fakultas); i++ {
		for j := 0; j < len(prodi[fakultas[i]]); j++ {
			jadwalUrl := targetUrl + "?fakultas=" + fakultas[i] + "&prodi=" + prodi[fakultas[i]][j]

			doc := htmlParser("GET", jadwalUrl, client, nil)

			p := []schedule{}
			parseJadwal(doc, &p)

			mainSchedule = append(mainSchedule, p...)
			//fmt.Println(p)
		}
	}

	return mainSchedule
}

func htmlParser(method, url string, client *http.Client, header http.Header) *html.Node {
	req, _ := http.NewRequest(method, url, nil)

	if header != nil {
		req.Header = header
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	// Target html
	b := resp.Body

	defer b.Close() // close Body when the function returns

	doc, err := html.Parse(b)
	if err != nil {
		panic("error")
	}

	if resp.StatusCode >= 400 {
		fmt.Println(resp.StatusCode)
		fmt.Println("Error occurred, please check your credentials or contact admin for further information.")
		fmt.Println("*Please re-fetch after a moment. Possibly you are blocked for a while from akademik.itb.ac.id*")
	}

	return doc
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
				lm := mapNodeAttr(d)
				*fakultas = append(*fakultas, lm["value"])
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
	if n.Type == html.ElementNode && n.Data == "tbody" {
		// empty, numbering, class code, class name, SKS,
		// class number,lecturer, atendee, schedule
		for d := n.FirstChild; d != nil; d = d.NextSibling {
			for e := d.FirstChild; e != nil; e = e.NextSibling {
				if e.Data == "td" {
					f := e.FirstChild
					var kode, matkul, sks, kelas string
					var peserta int
					var lecturerArray []string
					var daytimeArray []daytime

					if f != nil {
						kode = renderNode(f)

						if e.NextSibling != nil {
							e = e.NextSibling
							matkul = renderNode(e.FirstChild)
							matkul = strings.Replace(matkul, "\n", "", -1)
							matkul = strings.Replace(matkul, "  ", "", -1)

						}
						if e.NextSibling != nil {
							e = e.NextSibling
							sks = renderNode(e.FirstChild)
						}

						if e.NextSibling != nil {
							e = e.NextSibling
							kelas = renderNode(e.FirstChild)
						}

						if e.NextSibling != nil {
							e = e.NextSibling
							for f := e.FirstChild; f != nil; f = f.NextSibling {
								for g := f.FirstChild; g != nil; g = g.NextSibling {
									h := g.FirstChild
									if h != nil {
										lct := renderNode(h)
										lct = strings.Replace(lct, "\n", "", -1)
										lct = strings.Replace(lct, "  ", "", -1)
										lecturerArray = append(lecturerArray, lct)
									}
								}
							}
						}

						if e.NextSibling != nil {
							e = e.NextSibling
							peserta, _ = strconv.Atoi(renderNode(e.FirstChild))
						}

						if e.NextSibling != nil {
							e = e.NextSibling
							for f := e.FirstChild; f != nil; f = f.NextSibling {
								for g := f.FirstChild; g != nil; g = g.NextSibling {
									h := g.FirstChild
									if h != nil {
										jdwl := renderNode(h)
										jdwl = strings.Replace(jdwl, "\n", "", -1)
										jdwl = strings.Replace(jdwl, "  ", "", -1)
										stringSlice := strings.Split(jdwl, "/")
										// Parse jadwal
										// Day
										day := stringSlice[0]
										// time
										time := strings.Split(stringSlice[1], "-")
										// TimeStart
										timestart, _ := strconv.Atoi(time[0])
										// TimeEnd
										timeend, _ := strconv.Atoi(time[1])
										// Room
										room := stringSlice[2]
										// Type
										classType := stringSlice[3]

										daytimeIn := daytime{}
										daytimeIn.Day = day
										daytimeIn.Room = room
										daytimeIn.TimeStart = timestart
										daytimeIn.TimeEnd = timeend
										daytimeIn.Type = classType

										daytimeArray = append(daytimeArray, daytimeIn)
									}
								}
							}
						}

					}
					s := schedule{}
					s.ID = bson.NewObjectId()
					s.Code = kode
					s.Subject = matkul
					s.SKS = sks
					s.ClassNum = kelas
					s.Lecturer = lecturerArray
					s.StudentAmt = peserta
					s.Daytime = daytimeArray

					*jadwal = append(*jadwal, s)
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
