package main

import (
	"net/http"
	"log"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"github.com/gorilla/mux"
	"encoding/json"
	"os"
)

type Food struct {
	Date       string `json:"date"`
	FirstFood  string `json:"first_food"`
	SecondFood string `json:"second_food"`
	ThirdFood  string `json:"third_food"`
	FourthFood string `json:"fourth_food"`
}

type Foods []Food

const BASE_URL string = "http://ihale.manas.edu.kg/kki.php"

func parseFoods() (foods Foods, err error) {
	resp, err := http.Get(BASE_URL)
	if (err != nil) {
		return foods, err
	}

	root, err := html.Parse(resp.Body)
	if (err != nil) {
		return foods, err
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
		food := Food{}
		for id, rawFood := range rawFoods {
			if id == 0 {
				food.Date = scrape.Text(rawFood)
			} else if id == 1 {
				food.FirstFood = scrape.Text(rawFood)
			} else if id == 3 {
				food.SecondFood = scrape.Text(rawFood)
			} else if id == 5 {
				food.ThirdFood = scrape.Text(rawFood)
			} else if id == 7 {
				food.FourthFood = scrape.Text(rawFood)
			}
		}
		foods = append(foods, food)
	}

	return foods, nil
}

func FoodsHandler(w http.ResponseWriter, r *http.Request) {
	foods, _ := parseFoods()

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(foods)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/foods", FoodsHandler)
	log.Fatal(http.ListenAndServe(":" + port, router))
}