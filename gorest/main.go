package main

import (
	"log"
	"net/http"

	"github.com/SGTYang/gorest/gorest/elastic"
	"github.com/SGTYang/gorest/gorest/post"

	"github.com/julienschmidt/httprouter"
)

func main() {
	// Bootstrap elasticsearch.
	elastic, err := elastic.New([]string{"http://0.0.0.0:9200"})
	if err != nil {
		log.Fatalln(err)
	}
	if err := elastic.CreateIndex("post"); err != nil {
		log.Fatalln(err)
	}

	// Bootstrap storage.
	storage, err := elastic.NewPostStorage(*elastic)
	if err != nil {
		log.Fatalln(err)
	}

	// Bootstrap API.
	postAPI := post.NewHandler(storage)

	// Bootstrap HTTP router.
	router := httprouter.New()
	router.HandlerFunc("POST", "/api/v1/posts", postAPI.Create)
	router.HandlerFunc("PATCH", "/api/v1/posts/:id", postAPI.Update)
	router.HandlerFunc("DELETE", "/api/v1/posts/:id", postAPI.Delete)
	router.HandlerFunc("GET", "/api/v1/posts/:id", postAPI.Find)

	// Start HTTP server.
	log.Fatalln(http.ListenAndServe(":3000", router))
}
