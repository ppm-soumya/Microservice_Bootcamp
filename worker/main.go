//https://www.thepolyglotdeveloper.com/2017/07/consume-restful-api-endpoints-golang-application/
// https://redis.io/commands/hmget

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	//"github.com/mediocregopher/radix.v2/redis"
	"github.com/go-redis/redis"
	//"log"
	"strings"
	//"strconv"
	//"encoding/hex"
	//"encoding/json"
)

type test_struct struct {
	Test string
}

func oneWay(w http.ResponseWriter, r *http.Request) {
	var f func()
	var t *time.Timer
	fmt.Println("test worker")
	conn := redis.NewClient(&redis.Options{
		Addr:         "172.25.16.126:6379",
		Password:     "", // no password set
		DB:           0,  // use default DB
		MaxRetries:   3,
		IdleTimeout:  5 * time.Minute,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
	})

	pong, err := conn.Ping().Result()
	fmt.Println("pong : ", pong)
	if err != nil {
		fmt.Println("error : ", err)
		//log.Fatal(err)
	}

	/*conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
			fmt.Println("redis error localip:",err)
			log.Fatal(err)
	}*/
	// Importantly, use defer to ensure the connection is always properly
	// closed before exiting the main() function.
	defer conn.Close()

	f = func() {
		fmt.Println("Calling STG...")
		response, err := http.Get("http://172.25.16.126:10000/stg/tokens/1")

		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		} else {
			data, _ := ioutil.ReadAll(response.Body)
			fmt.Println(string(data))
			/*decoder := json.NewDecoder(data)
			var t test_struct
			err := decoder.Decode(&t)
			if err != nil {
				fmt.Println("ERR:", err)
				panic(err)
			}
			fmt.Println("STG ", t)
			fmt.Println(t.Test)*/
			//jsonValue, err := json.Marshal(t.Test)
			// if err!=nil {
			// 		fmt.Println("error in JSON conversion")
			// } else {
			//fmt.Println("Calling HASHER with...", jsonValue)

			// response, err = http.Post("http://localhost:10001/hasher/", "application/json", bytes.NewBuffer(data))
			fmt.Println("Calling Hasher... Data : ", data)
			url := "http://172.25.16.126:10001/hasher"
			//fmt.Println("URL:>", url)
			// var jsonStr = data//[]byte(`{"title":"Buy cheese and bread for breakfast."}`)
			req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			defer resp.Body.Close()
			if err != nil {
				fmt.Printf("The HTTP request failed with error %s\n", err)
			} else {
				data, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("DATA: ", string(data))

				//! Unmarshaling the data obtained from hasher into a map,
				//! doing HMSET into redis with the dat["hash"] values
				//! doing HMGET to check whether entered into redis database correctly or not
				dat := make(map[string]string)
				err = json.Unmarshal(data, &dat)
				if err != nil {
					fmt.Println("error_46 : ", err)
					panic(err)
				}
				fmt.Println("Unmarshal done", dat["hash"])
				/*resp := conn.Cmd("HMSET", "hashtable", "hash", dat["hash"])
				// Check the Err field of the *Resp object for any errors.
				if resp.Err != nil {
					  fmt.Println("HMSET error", resp.Err)
						//log.Fatal(resp.Err)
				}*/

				val := strings.HasPrefix(dat["hash"], "0") // true
				fmt.Println("Lucky Hash : ", val)
				fmt.Println("Hash : ", dat["hash"])

				pubsub := conn.Subscribe("mychannel1")
				defer pubsub.Close()
				//resp := conn.Publish("hashChannel", dat["hash"]).Err()
				resp := conn.Publish("hashChannel", val).Err()
				if resp != nil {
					fmt.Println(resp)
					panic(err)
				}
			}
		}

		t = time.AfterFunc(time.Duration(1)*time.Second, f)
	}

	t = time.AfterFunc(time.Duration(1)*time.Second, f)

	defer t.Stop()

	//simulate doing stuff
	time.Sleep(time.Minute)
}

func main() {
	http.HandleFunc("/", oneWay)
	http.ListenAndServe(":8080", nil)
}
