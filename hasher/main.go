// main.go Hasher
//https://stackoverflow.com/questions/10701874/generating-the-sha-hash-of-a-string-using-golang
//https://gist.github.com/andreagrandi/97263aaf7f9344d3ffe6

package main

import (
    "fmt"
    "log"
    "net/http"
    //"crypto/sha1"
    "encoding/hex"
    "github.com/gorilla/mux"
    "encoding/json"
    "io/ioutil"
    "crypto/sha256"
)

type test_struct struct {
	Test string
}

func hasher(token_value string) string {

    //s := "sha1 this string"
    //h := sha1.New()
    //h.Write([]byte(token_value))
    //sha1_hash := hex.EncodeToString(h.Sum(nil))

    h2 := sha256.New()
	  h2.Write([]byte(token_value))
	  fmt.Println("hashed value ==", h2.Sum(nil))
    sha1_hash := hex.EncodeToString(h2.Sum(nil))

    fmt.Println(token_value, sha1_hash)
    return sha1_hash
}

func returnHashedNumber(w http.ResponseWriter, r *http.Request){
  fmt.Println("Method: ", r.Method)
    if r.Method!="POST"{
        fmt.Println("hasher called with Method other than POST")
        return
    }

    fmt.Println("hasher called with Method POST")
    data, _ := ioutil.ReadAll(r.Body)
    fmt.Println(string(data))
    //var dat map[string]interface{}
    dat := make(map[string]string)
    err := json.Unmarshal(data, &dat);
    if err != nil {
      fmt.Println("error_46 : ",err)
      panic(err)
    }
    fmt.Println("Unmarshal done", dat["token"])

    /*decoder := json.NewDecoder(r.Body)
    var t test_struct
    err = decoder.Decode(&t)
    if err != nil {
      fmt.Println("ERR ",err)
      panic(err)

    }*/
    //fmt.Println(t.Test)

    fmt.Println("Hasher Input : ", dat["token"])//t.Test)
    sha1_hash := hasher(dat["token"])//t.Test)
    m := make(map[string]string)
    m["hash"]=sha1_hash
    js, err := json.Marshal(m)
      if err != nil {
          fmt.Println("Error Marshaling")
          http.Error(w, err.Error(), http.StatusInternalServerError)
          return
      }
    fmt.Println("JSON Output ",js)
    w.Header().Set("Content-Type", "application/json")
    w.Write(js)
}

func handleRequests() {

    myRouter := mux.NewRouter().StrictSlash(true)
    myRouter.HandleFunc("/hasher", returnHashedNumber).Methods("POST")
    log.Fatal(http.ListenAndServe(":10001", myRouter))
}

func main() {
    fmt.Println("Rest API v2.0 - Mux Routers")
    handleRequests()
}
