package main

import (
	"net/http"
	"log"
	"encoding/json"
	"os"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"github.com/gorilla/mux"
	"github.com/Machiel/slugify"
)

type Menu struct {
	Date  string `json:"date"`
	Foods []Food `json:"foods"`
}

type Food struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type Menus []Menu

const BASE_URL string = "http://ihale.manas.edu.kg/kki.php"

func parseFoods() (menus Menus, err error) {
	resp, err := http.Get(BASE_URL)
	if (err != nil) {
		return menus, err
	}

	root, err := html.Parse(resp.Body)
	if (err != nil) {
		return menus, err
	}

	matcher := func(node *html.Node) bool {
		if node.FirstChild == nil {
			return false
		}
		if node.DataAtom == atom.Tr && node.Parent != nil && node.Parent.Parent != nil {
			return node.FirstChild.DataAtom != atom.Th && node.Parent.DataAtom == atom.Tbody && node.Parent.Parent.DataAtom == atom.Table
		}
		return false
	}

	lines := scrape.FindAll(root, matcher)

	for _, line := range lines {
		rawFoods := scrape.FindAll(line, scrape.ByTag(atom.Td))
		menu := Menu{}
		for id, rawFood := range rawFoods {
			if id == 0 {
				menu.Date = scrape.Text(rawFood)
			} else if (id == 1 || id == 3 || id == 5 || id == 7) {
				name := scrape.Text(rawFood)
				menu.Foods = append(menu.Foods, Food{Name:name, Slug:slugify.Slugify(name)})
			}
		}
		menus = append(menus, menu)
	}

	return menus, nil
}

func FoodsHandler(w http.ResponseWriter, r *http.Request) {
	menus, _ := parseFoods()

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(menus)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/menus", FoodsHandler)
	router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
	log.Fatal(http.ListenAndServe(":" + port, router))
}