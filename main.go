package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var ErrIDNotFound = fmt.Errorf("ID not found")

type Wish struct {
	ID          uint64 `json:"id"`
	Description string `json:"description"`
}

type WishList struct {
	identity uint64
	list     []Wish
}

func (wl *WishList) Add(w Wish) []Wish {
	w.ID = uint64(wl.identity)
	wl.list = append(wl.list, w)
	wl.identity++
	return wl.list
}

func (wl *WishList) Remove(id uint64) ([]Wish, error) {
	found := false
	for i, wish := range wl.list {
		if wish.ID == id {
			wl.list = append(wl.list[:i], wl.list[i+1:]...)
			found = true
		}
	}

	if found {
		return wl.list, nil
	}

	return wl.list, ErrIDNotFound
}

func (wl *WishList) Update(uw Wish) ([]Wish, error) {
	found := false
	for i := range wl.list {
		if wl.list[i].ID == uw.ID {
			wl.list[i].Description = uw.Description
			found = true
		}
	}

	if found {
		return wl.list, nil
	}

	return wl.list, ErrIDNotFound
}

var wishList = &WishList{
	identity: 0,
	list:     []Wish{},
}

func getWishes(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(wishList.list)
}

func addWish(w http.ResponseWriter, r *http.Request) {
	var wish Wish
	err := json.NewDecoder(r.Body).Decode(&wish)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(wishList.Add(wish))
}

func updateWish(w http.ResponseWriter, r *http.Request) {
	var wish Wish
	err := json.NewDecoder(r.Body).Decode(&wish)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	list, err := wishList.Update(wish)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(list)
}

func deleteWish(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	list, err := wishList.Remove(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(list)
}

func main() {
	wishList.Add(Wish{Description: "Two Front Teeth"})

	r := mux.NewRouter()
	r.HandleFunc("/", getWishes).Methods("GET")
	r.HandleFunc("/", addWish).Methods("POST")
	r.HandleFunc("/", updateWish).Methods("PUT")
	r.HandleFunc("/{id}", deleteWish).Methods("DELETE")
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	srv.ListenAndServe()
}
