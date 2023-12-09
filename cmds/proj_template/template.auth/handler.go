package auth

import (
	"database/sql"
	"net/http"
)

type Handle struct {
	dbs map[string]*sql.DB
}

//func (h *Handle) RegisterGet(w http.ResponseWriter, r *http.Request) {
//	err := Register.Get(w, r)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//}

//func (h *Handle) RegisterPost(w http.ResponseWriter, r *http.Request) {
//	err := Register.Post(w, r)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	http.Redirect(w, r, "/market", http.StatusSeeOther)
//	return
//}

func (h *Handle) LoginGet(w http.ResponseWriter, r *http.Request) {
	err := Auth.LoginGet(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handle) LoginPost(w http.ResponseWriter, r *http.Request) {
	err := Auth.LoginPost(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/market", http.StatusSeeOther)
	return
}

func (h *Handle) Logout(w http.ResponseWriter, r *http.Request) {
	err := Logout.Logout(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/market", http.StatusSeeOther)
	return
}
