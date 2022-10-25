package main

import ( 
	"os"						// access operating system functionality
	"fmt"
	"log"						// log output
	"time"
	"net/url"
	"net/http"					// provides http client and server implementations 
	"html/template"				// allows templating for html
	"strconv"
	"math"
	"bytes"

	// "news/news.go"
	// news "news/news"
	// "main/news"
	"github.com/freshman-tech/news-demo-starter-files/news"
	"github.com/joho/godotenv" 	// uses github.com/joho/godotenv to get environement
)

// points to template definition from provided files
// template.Must will throw error if  fails
var tpl = template.Must(template.ParseFiles("index.html"))

type Search struct {
	Query 		string 
	NextPage	int 
	TotalPages 	int 
	Results 	*news.Results
}

// navigating to next page after 20 page limit
func (s *Search) IsLastPage() bool {
	return s.NextPage >= s.TotalPages
}

// navigatingto check current page 
func (s *Search) CurrentPage() int {
	if s.NextPage == 1 {
		return s.NextPage
	}

	return s.NextPage - 1
}

func (s *Search) PreviousPage() int {
	return s.CurrentPage() - 1
}

/*
 * handle index
 * :input w http.ResponseWriter used to send responses to a https requiest
 * :input r *http.Request is the request received from the client 
 * :return
 */
func indexHandler(w http.ResponseWriter, r *http.Request) {
	// // execute with where we want to write output and data passed to template
	// tpl.Execute(w, nil)
	// w.Write([]byte("<h1>Hellow World!</h1>"))

	// instead we no longer execute directly to responsewriter
	buf := &bytes.Buffer{} 
	err := tpl.Execute(buf, nil) 
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 
	}

	buf.WriteTo(w)

}


// /*
//  * extracts the parameters q and page from the requests url and prints
//  * them both to the standard output
//  * :input w http.ResponseWriter, 
//  * :input r http.Requst
//  */
// func searchHandler(w http.ResponseWriter, r *http.Request) {
// 	u, err := url.Parse(r.URL.String())
	
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}

// 	params := u.Query() 
// 	searchQuery :=  params.Get("q")
// 	page := params.Get("page")
// 	if page == "" {
// 		page = "1"
// 	}

// 	fmt.Println("Search Query is: ", searchQuery) 
// 	fmt.Println("Page is: ", page)
// }


/*
 * extracts the parameters q and page from the requests url and prints
 * them both to the standard output
 * :input w http.ResponseWriter, 
 * :input r http.Requst
 */
 func searchHandler(newsapi *news.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := url.Parse(r.URL.String())
		
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		params := u.Query() 
		searchQuery :=  params.Get("q")
		page := params.Get("page")
		if page == "" {
			page = "1"
		}

		fmt.Println("Search Query is: ", searchQuery) 
		fmt.Println("Page is: ", page)

		results, err := newsapi.FetchEverything(searchQuery, page) 
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// convert next page from string to int
		nextPage, err := strconv.Atoi(page) 
		if err != nil { 
			http.Error(w,err.Error(), http.StatusInternalServerError) 
			return
		}

		// create search from struct
		search := &Search {
			Query: 		searchQuery, 
			NextPage: 	nextPage, 
			TotalPages: int(math.Ceil(
				float64(results.TotalResults) / float64(newsapi.PageSize))),
			Results: results, 
		}

		// ok is last page, if ok is true, increment nextpage
		if ok := !search.IsLastPage(); ok {
			search.NextPage++
		}

		// template tpl executed into empty buffer tocheck for errors, before being 
		// written to responsewriter
		buf := &bytes.Buffer{} 
		err = tpl.Execute(buf, search) 
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		buf.WriteTo(w)
		fmt.Printf("%+v", results)
	}
}




/*  
 *
 *
 */
func main() {
	// read env file and laods into environment such that they can be accessed
	// by os.Getenv()
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	// if port not found, set to localhost:3000
	port:= os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// declare api key 
	apiKey := os.Getenv("NEWS_API_KEY")
	if apiKey == "" {
		log.Fatal("Env: apiKey must be set")
	}

	// start client and newsapi
	myClient := &http.Client{Timeout: 10 * time.Second} 
	newsapi := news.NewClient(myClient, apiKey, 20)

	// create file server for css assets
	fs := http.FileServer(http.Dir("assets"))

	// create http request multiplexer
	mux := http.NewServeMux() 
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
	mux.HandleFunc("/", indexHandler)

	// create search handler
	mux.HandleFunc("/search", searchHandler(newsapi))
	
	// starts server on the port and runs mux
	http.ListenAndServe(":"+port, mux)
}