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
	ela, err := elastic.New([]string{"http://0.0.0.0:9200"})
	if err != nil {
		log.Fatalln(err)
	}
	if err := ela.CreateIndex("post"); err != nil {
		log.Fatalln(err)
	}

	// Bootstrap storage.
	storage, err := elastic.NewPostStorage(*ela)
	if err != nil {
		log.Fatalln(err)
	}

	// Bootstrap API.
	postAPI := post.NewHandler(storage)

	// Bootstrap HTTP router.
	router := httprouter.New()
	router.HandlerFunc("POST", "/api/v1/posts", postAPI.Create)

	// Start HTTP server.
	log.Fatalln(http.ListenAndServe(":3000", router))
}
