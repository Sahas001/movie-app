package main

import (
	"log"
	"net/http"

	"github.com/Sahas001/movieapp/movie/internal/controller/movie"
	metadatagateway "github.com/Sahas001/movieapp/movie/internal/gateway/metadata/http"
	ratinggateway "github.com/Sahas001/movieapp/movie/internal/gateway/rating/http"
	httphandler "github.com/Sahas001/movieapp/movie/internal/handler/http"
)

func main() {
	log.Println("Starting the movie service")
	metadataGateway := metadatagateway.New("localhost:8081")
	ratingGateway := ratinggateway.New("localhost:8082")
	ctrl := movie.New(ratingGateway, metadataGateway)
	h := httphandler.New(ctrl)
	http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))
	if err := http.ListenAndServe(":8083", nil); err != nil {
		panic(err)
	}
}
