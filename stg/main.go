// main.go STG
//https://tutorialedge.net/golang/creating-restful-api-with-golang/
// Rendering XML, JSON : https://www.alexedwards.net/blog/golang-response-snippets
// check for server output in terminal : $ curl -i localhost:10000/stg/tokens

package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "crypto/rand"
    "strconv"
    "encoding/json"
)

type token struct {
  token_value string
  //token_string string
}

func tokenGenerator(num int) string {
	b := make([]byte, num)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func returnRandomNumber(w http.ResponseWriter, r *http.Request){
    if r.Method!="GET"{
        fmt.Println("tokenGenerator called with Method other than GET")
        return
    }

    fmt.Println("tokenGenerator called with Method GET")
    vars := mux.Vars(r)
    key := vars["id"]
    i1, err := strconv.Atoi(key)
    if err == nil {
      a := tokenGenerator(i1)
      //fmt.Fprintf(w, a)
      tokenGen := token{}
      //tokenGen.token_string = "token"
      tokenGen.token_value = a
      m := make(map[string]string)
      m["token"]=a
      js, err := json.Marshal(m)
      //var jsonStr = []byte(`{"token":"` + a + `"}`)
      if err != nil {
          http.Error(w, err.Error(), http.StatusInternalServerError)
          return
      }
      fmt.Println("Generated token : ", a)
      fmt.Println("hex encoded value : ", js)
      w.Header().Set("Content-Type", "application/json")
      w.Write(js)
    }
}

func handleRequests() {

    myRouter := mux.NewRouter().StrictSlash(true)
    myRouter.HandleFunc("/stg/tokens/{id}", returnRandomNumber).Methods("GET")
    log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
    fmt.Println("Rest API v2.0 - Mux Routers")
    handleRequests()
}
