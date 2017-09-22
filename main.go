package main

import (
	"net/http"
	"log"
	"encoding/json"
	"os"
	"regexp"
	"errors"
	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"github.com/gorilla/mux"
)

type Menu struct {
	Date          string `json:"date"`
	Foods         []Food `json:"foods"`
	TotalCalories string `json:"total_calories"`
}

type Food struct {
	Name     string `json:"name"`
	Calories string `json:"calories"`
}

type ApiError struct {
	Code        uint16 `json:"code"`
	Description string `json:"description"`
}

type Menus []Menu

const BASE_URL string = "http://manasbis.manas.edu.kg/menu/"

func parseFoods() (menus Menus, err error) {
	resp, err := http.Get(BASE_URL)
	if err != nil {
		return menus, err
	}

	root, err := html.Parse(resp.Body)
	if err != nil {
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

	trs := scrape.FindAll(root, matcher)

	for _, tr := range trs {
		tds := scrape.FindAll(tr, scrape.ByTag(atom.Td))
		var elements []string
		for _, td := range tds {
			elements = append(elements, scrape.Text(td))
		}
		menu := Menu{}

		elementsQty := len(elements)
		if elementsQty > 0 {
			menu.Date = elements[0]
		}

		for i := 1; i <= 10 && elementsQty > i; i += 2 {
			menu.Foods = append(menu.Foods, Food{Name: elements[i], Calories: elements[i+1]})
		}

		if elementsQty == 12 {
			menu.TotalCalories = elements[11]
		}

		menus = append(menus, menu)
	}

	return menus, nil
}

func isValidDate(date string) bool {
	validDate := regexp.MustCompile(`^[0-9]{4}-[0-9]{2}-[0-9]{2}$`)
	return validDate.MatchString(date)
}

func getMenuByDate(menus Menus, date string) (foundMenu Menu, err error) {
	for _, menu := range menus {
		if menu.Date == date {
			return menu, nil
		}
	}
	return foundMenu, errors.New("Menu by date not found")
}

func MenusHandler(w http.ResponseWriter, r *http.Request) {
	menus, err := parseFoods()
	if err != nil {
		err := ApiError{Code: http.StatusBadRequest, Description: "Error occurred during parsing"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(menus)
}

func MenuHandler(w http.ResponseWriter, r *http.Request) {
	date := mux.Vars(r)["date"]

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	if !isValidDate(date) {
		err := ApiError{Code: http.StatusBadRequest, Description: "Invalid date format"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
	} else {
		menus, err := parseFoods()

		if err != nil {
			err := ApiError{Code: http.StatusBadRequest, Description: "Error occurred during parsing"}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err)
			return
		}

		menu, err := getMenuByDate(menus, date)
		if err != nil {
			err := ApiError{Code: http.StatusNotFound, Description: http.StatusText(http.StatusNotFound)}
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(err)
		} else {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(menu)
		}
	}
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	err := ApiError{Code: http.StatusNotFound, Description: http.StatusText(http.StatusNotFound)}
	json.NewEncoder(w).Encode(err)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/menus", MenusHandler)
	router.HandleFunc("/menus/{date:[0-9-]+}", MenuHandler)
	router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
