// app.go

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"github.com/go-gl/mathgl/mgl64"
	"net/http"
	"strconv"
	"github.com/twinj/uuid"
	"github.com/gorilla/mux"
	"github.com/alfredxing/calc/compute"
)

import "github.com/leprosus/golang-ttl-map"



type user struct {
	Login string `json:"login"`
	Pass  string    `json:"pass"`
}

type token struct {
	Token string `json:"token"`
}

type result struct{
	Result string `json:"result"`	
}

type computerequest struct {
	Token string `json:"token"`
	Expression string `json:"expression"`	
}


type fault struct {
	Fault string `json:"fault"`
}


type App struct {
	Router *mux.Router
//	ValidToken string
	ttlMap ttl_map.Heap
}

func (a *App) Initialize() {

	a.Router = mux.NewRouter()
	a.initializeRoutes()
	a.ttlMap = ttl_map.New("token.tsv")
	a.ttlMap.Restore()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/login", a.login).Methods("POST")
	a.Router.HandleFunc("/compute", a.compute).Methods("POST")
//	a.Router.HandleFunc("/user/{id:[0-9]+}", a.getUser).Methods("GET")
//	a.Router.HandleFunc("/user/{id:[0-9]+}", a.updateUser).Methods("PUT")
//	a.Router.HandleFunc("/user/{id:[0-9]+}", a.deleteUser).Methods("DELETE")
}

func (a *App) login(w http.ResponseWriter, r *http.Request) {
	var u user
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if (u.Login == "admin" && u.Pass=="admin"){
		uid := uuid.NewV4();
		a.ttlMap.Set(uid.String(), "valid", 15*60)		
		a.ttlMap.Save()		
		var t token = token{uid.String()}
		respondWithJSON(w, http.StatusOK, t)	
	}else{
		respondWithJSON(w, http.StatusBadRequest, fault{"Invalid User"})
	}

}

func (a *App) compute(w http.ResponseWriter, r *http.Request) {
	var req computerequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	//Check valid token
	valuev := a.ttlMap.Get(req.Token)
	if (valuev == "valid"){
		// Split expression
		if strings.Contains(req.Expression,"=="){
			s := strings.Split(req.Expression,"==")
			exp1, exp2 := s[0], s[1]
			exp1 = strings.Replace(exp1, " ", "", -1)
			exp2 = strings.Replace(exp2, " ", "", -1)
			res1, err1 := compute.Evaluate(exp1)
			if err1 != nil {
				respondWithJSON(w, http.StatusBadRequest, fault{fmt.Sprint("Error: " + err1.Error())})
				return;
			}
			res2, err2 := compute.Evaluate(exp2)
			if err2 != nil {
				respondWithJSON(w, http.StatusBadRequest, fault{fmt.Sprint("Error: " + err2.Error())})
				return;
			}
			if (mgl64.FloatEqual(res1,res2)){
				respondWithJSON(w, http.StatusOK, result{"True"})				
			}else{
				respondWithJSON(w, http.StatusOK, result{"False"})				
			}			
		}else{ 
			exp1 := req.Expression
			exp1 = strings.Replace(exp1, " ", "", -1)
	//		fmt.Println(exp1)
			res1, err1 := compute.Evaluate(exp1)
			if err1 != nil {
				respondWithJSON(w, http.StatusBadRequest, fault{fmt.Sprint("Error: " + err1.Error())})
				return;
			}
			respondWithJSON(w, http.StatusOK, result{fmt.Sprint(strconv.FormatFloat(res1, 'G', -1, 64))})							
		}
	}else{
		respondWithJSON(w, http.StatusBadRequest, fault{"User not authentified"})
	}	
}


func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

