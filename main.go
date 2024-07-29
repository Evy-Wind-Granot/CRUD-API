package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Movie struct represents a movie with an ID, ISBN, title, and a director
type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

// Director struct represents a director with a firstname and lastname
type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// movies slice to seed initial movie data
var movies []Movie

// getMovies function to return all movies as JSON
func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

// deleteMovie function to delete a movie by its ID
func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // get parameters from the request
	for index, item := range movies {
		if item.ID == params["id"] {
			// remove the movie from the slice
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(movies) // return the updated list
}

// getMovie function to return a single movie by its ID
func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // get parameters from the request
	for _, item := range movies {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item) // return the specific movie
			return
		}
	}
}

// createMovie function to add a new movie
func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)       // decode the request body into movie struct
	movie.ID = strconv.Itoa(rand.Intn(1000000))      // generate a random ID for the new movie
	movies = append(movies, movie)                   // add the new movie to the slice
	json.NewEncoder(w).Encode(movie)                 // return the new movie as JSON
}

// updateMovie function to update an existing movie
func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // get parameters from the request
	for index, item := range movies {
		if item.ID == params["id"] {
			// remove the existing movie
			movies = append(movies[:index], movies[index+1:]...)
			var movie Movie
			_ = json.NewDecoder(r.Body).Decode(&movie) // decode the request body into movie struct
			movie.ID = params["id"]                    // set the ID to the existing movie's ID
			movies = append(movies, movie)             // add the updated movie to the slice
			json.NewEncoder(w).Encode(movie)           // return the updated movie as JSON
			return
		}
	}
}

func main() {
	r := mux.NewRouter() // create a new router

	// Seed initial movie data
	movies = append(movies, Movie{ID: "1", Isbn: "438227", Title: "Movie One", Director: &Director{Firstname: "John", Lastname: "Doe"}})
	movies = append(movies, Movie{ID: "2", Isbn: "45455", Title: "Movie Two", Director: &Director{Firstname: "Steve", Lastname: "Smith"}})

	// Route handlers & endpoints
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	// Start the server
	fmt.Printf("Starting server at port 8000\n")
	log.Fatal(http.ListenAndServe(":8000", r))
}
