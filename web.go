package web

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type Web struct {
	port string
	main *Router
}

func New(port string, main *Router) *Web {
	primary := &Router{
		matcher: nil,
		hops:    []*Router{main},
	}
	return &Web{port: port, main: primary}
}

func ViewHTML(path string, data interface{}, w http.ResponseWriter, r *http.Request) error {
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		return err
	}

	return tmpl.Execute(w, data)
}

func ViewFile(path string, data interface{}, w http.ResponseWriter, r *http.Request) error {
	fmt.Println(path)
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	fileHeader := make([]byte, 512)
	_, err = file.Read(fileHeader)
	if err != nil {
		return err
	}
	contentType := http.DetectContentType(fileHeader)

	//fileStat, err := file.Stat()
	if err != nil {
		return err
	}
	//fileSize := strconv.FormatInt(fileStat.Size(), 10)

	//w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(path))
	w.Header().Set("Content-Type", contentType)
	//w.Header().Set("Content-Length", fileSize)

	fmt.Println("heer")
	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}

	fmt.Println("heer")
	_, err = io.Copy(w, file) //'Copy' the file to the client
	return err
}

func (web *Web) handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Requested URL path", r.URL.Path)
	route := strings.Split(r.URL.Path, "/")
	if route[len(route)-1] == "" {
		route = route[:len(route)-1]
	}
	log.Println("Requested route", route)

	router := web.main
	for indie, segment := range route {
		log.Println("Segment", indie, len(route), segment)
		var jndie int
		var hop *Router

		// Use manual for over range.  Range only goes through specific indices.
		for jndie = 0; jndie < len(router.hops); jndie++ {
			hop = router.hops[jndie]
			match := hop.matcher.Match(segment)
			log.Printf("Current segment \"%s\" vs hop \"%s\" %t", segment, hop.matcher.String(), match)
			if match {
				log.Println("Hop found", hop, "at", jndie, len(router.hops))
				break
			}
		}
		log.Println("Router jndie", jndie)
		if jndie == len(router.hops) {
			log.Println("Failed to find router for", route[:indie+1])
			http.NotFound(w, r)
			return
		}
		router = hop
	}

	resolved := router.resolver.resolve(route)
	log.Println("Hop path:", resolved)
	var data interface{}
	var err error
	if router.handler != nil {
		data, err = router.handler(resolved, w, r)
		if err != nil {
			log.Println("Data Handler Error:", err.Error())
			http.Error(w, "", http.StatusInternalServerError)
			return
		} /*else if data == nil {
			log.Println("No data returned to process?")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}*/
	}

	if router.writer != nil {
		err = router.writer(resolved, data, w, r)
		if err != nil {
			log.Println("Data Writer Error:", err.Error())
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}

func (web *Web) Run() error {
	http.HandleFunc("/", web.handler)
	return http.ListenAndServe(":"+web.port, nil)
}
