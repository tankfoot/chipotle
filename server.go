// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"flag"
	"log"
    "fmt"
	"net/http"
    "encoding/json"
    "os/exec"
    "errors"
    "io/ioutil"
    "bytes"
    "time"

	"github.com/gorilla/websocket"
    sj "github.com/bitly/go-simplejson"
)

//Incoming Json struct
type Data struct {
    Query string
}

type Message struct {
    Header [6]float64
    Data   Data
}

//Dialogflow Query struct
type Text struct {
    Text string `json:"text"`
    LanguageCode string `json:"languageCode"`
}

type TextInput struct {
    TextInput Text `json:"text"`
}

type QueryInput struct {
    QueryInput TextInput `json:"queryInput"`
}

//Output json struct
type DataOutput struct {
    Speech string `json:"speech"`
    Entity map[string]interface{} `json:"entity"`
}

type Output struct {
    Header [7]float64 `json:"header"`
    Data DataOutput `json:"data"`
}

//const letterBytes = "abcdefghijklmnopqrstuvwxyz1234567890"

// func RandStringBytes(n int) string {
// 	b := make([]byte, n)
// 	l := len(letterBytes)
// 	for i := range b {
// 		b[i] = letterBytes[rand.Intn(l)]
// 	}
// 	return string(b)
// }

var user = map[float64]Output{}

var upgrader = websocket.Upgrader{} // use default options

func DetectIntentText(projectID, sessionID, text, languageCode string) (string, string, map[string]interface{}, error) {
    if projectID == "" || sessionID == "" {
        return "", "", nil, errors.New(fmt.Sprintf("Received empty project (%s) or session (%s)", projectID, sessionID))
    }
    basePath := "https://dialogflow.googleapis.com/v2/"
    sessionPath := fmt.Sprintf("projects/%s/agent/sessions/%s", projectID, sessionID)

    client := &http.Client{}
    var jsonData QueryInput
    jsonData.QueryInput = TextInput{Text{Text: text, LanguageCode: languageCode}}
    jsonValue, _ := json.Marshal(jsonData)
    detectIntentUrl := basePath + sessionPath + ":detectIntent"
    r, err := http.NewRequest("POST", detectIntentUrl, bytes.NewBuffer(jsonValue))

    token, _ := GetGcloudToken()
    var bearer = "Bearer " + token
    r.Header.Add("Authorization", bearer)
    r.Header.Add("Content-Type", "application/json; charset=utf-8")

    resp, err := client.Do(r)

    if err != nil {
        fmt.Printf("The HTTP request failed with error %s\n", err)
    } else {
        data, _ := ioutil.ReadAll(resp.Body)
        //fmt.Println(string(data))
        js, _ := sj.NewJson(data)
        speechText := js.Get("queryResult").Get("fulfillmentText").MustString()
        intentName := js.Get("queryResult").Get("intent").Get("displayName").MustString()
        entities := js.Get("queryResult").Get("parameters").MustMap()
        return speechText, intentName, entities, nil
    }

    return "", "", nil, nil
}

func GetGcloudToken() (string, error) {
    cmd := exec.Command("gcloud", 
        "auth",
        "application-default",
        "print-access-token",)

    out, err := cmd.Output()
    if err != nil {
        log.Fatal(err)
        return "", err
    }

    token := string(out)[:len(string(out))-1] // line ending subtract
    return token, nil
}

func HeaderProcess(headerIn [6]float64, intent string, speech string, entity map[string]interface{}) (
        [7]float64, string, map[string]interface{}, error) {
    var headerOut [7]float64
    var talkback string
    entityback := make(map[string]interface{})

    headerOut[0] = headerIn[0]
    headerOut[1] = headerIn[1]
    headerOut[2] = headerIn[2]

    switch intent {
    case "chipotle.burrito":
        switch speech {
        case "address":
            headerOut[3] = 2000
            talkback = "please select address, you can say recent, favorite, or nearby"
            entityback["ordertype"] = "burrito"
            entityback["address"] = entity["address"]
            entity = entityback
        case "fillings":
            headerOut[3] = 1100
            entityback["ordertype"] = "burrito"
            entityback["address"] = entity["address"]
            entity = entityback
            talkback = "which fillings do you want?"
        case "rice":
            headerOut[3] = 1110
            talkback = "fillings added, Any rice?"
        case "beans":
            headerOut[3] = 1120
            talkback = "Any beans?"
        case "toppings":
            headerOut[3] = 1130
            talkback = "Any toppings?"
        case "sides":
            headerOut[3] = 1140
            talkback = "Any sides?"
        case "drinks":
            headerOut[3] = 1150
            talkback = "Any drinks?"
        case "Done":
            headerOut[3] = 1160
            talkback = "Okay, Do you want to add item to cart"
        default:
            talkback = speech
        }
    case "chipotle.burrito - yes": 
        headerOut[3] = 1900
        talkback = speech
    case "chipotle.bowl":
        switch speech {
        case "address":
            headerOut[3] = 2000
            talkback = "please select address, you can say recent, favorite, or nearby"
            entityback["ordertype"] = "bowl"
            entityback["address"] = entity["address"]
            entity = entityback
        case "fillings":
            headerOut[3] = 1200
            entityback["ordertype"] = "bowl" 
            entity = entityback
            talkback = "which fillings do you want?"
        case "rice":
            headerOut[3] = 1210
            talkback = "Any rice?"
        case "beans":
            headerOut[3] = 1220
            talkback = "Any beans?"
        case "toppings":
            headerOut[3] = 1230
            talkback = "Any toppings?"
        case "sides":
            headerOut[3] = 1240
            talkback = "Any sides?"
        case "drinks":
            headerOut[3] = 1250
            talkback = "Any drinks?"
        case "Done":
            headerOut[3] = 1260
            talkback = "Okay, Do you want to add item to cart"
        default:
            talkback = speech
        }
    case "chipotle.bowl - yes": 
        headerOut[3] = 1900
        talkback = speech
    case "chipotle.salad":
        switch speech {
        case "address":
            headerOut[3] = 2000
            talkback = "please select address, you can say recent, favorite, or nearby"
            entityback["ordertype"] = "bowl"
            entityback["address"] = entity["address"]
            entity = entityback
        case "fillings":
            headerOut[3] = 1300
            entityback["ordertype"] = "salad" 
            entity = entityback
            talkback = "which fillings do you want?"
        case "rice":
            headerOut[3] = 1310
            talkback = "Any rice?"
        case "beans":
            headerOut[3] = 1320
            talkback = "Any beans?"
        case "toppings":
            headerOut[3] = 1330
            talkback = "Any toppings?"
        case "sides":
            headerOut[3] = 1340
            talkback = "Any sides?"
        case "drinks":
            headerOut[3] = 1350
            talkback = "Any drinks?"
        case "Done":
            headerOut[3] = 1360
            talkback = "Okay, Do you want to add item to cart"
        default:
            talkback = speech
        }
    case "chipotle.salad - yes": 
        headerOut[3] = 1900
        talkback = speech
    case "chipotle.tacos":
        switch speech {
        case "address":
            headerOut[3] = 2000
            talkback = "please select address, you can say recent, favorite, or nearby"
        case "number":
            headerOut[3] = 1400
            entityback["ordertype"] = "tacos"
            entity = entityback
            talkback = "how many tacos do you want?"
        case "tortilla":
            headerOut[3] = 1410
            talkback = "soft or crispy tortilla"
        case "fillings":
            headerOut[3] = 1400
            talkback = "which fillings do you want?"
        case "rice":
            headerOut[3] = 1410
            talkback = "Any rice?"
        case "beans":
            headerOut[3] = 1420
            talkback = "Any beans?"
        case "toppings":
            headerOut[3] = 1430
            talkback = "Any toppings?"
        case "sides":
            headerOut[3] = 1440
            talkback = "Any sides?"
        case "drinks":
            headerOut[3] = 1450
            talkback = "Any drinks?"
        case "Done":
            headerOut[3] = 1460
            talkback = "Okay, Do you want to add item to cart"
        default:
            talkback = speech
        }
    case "chipotle.tacos - yes": 
        headerOut[3] = 1900
        talkback = speech
    case "chipotle.kids": 
        switch speech {
        case "address":
            headerOut[3] = 2000
            talkback = "please select address, you can say recent, favorite, or nearby"
        case "choose":
            headerOut[3] = 2100
            talkback = "build your own or quesadilla?"
        default:
            talkback = speech
        }
    case "chipotle.kids - buildyourown": 
        switch speech {
        case "tortilla":
            headerOut[3] = 1500
            talkback = "soft or crispy tortilla?"
        case "fillings":
            headerOut[3] = 1510
            talkback = "which fillings do you want?"
        case "beans":
            headerOut[3] = 1520
            talkback = "Any beans?"
        default:
            talkback = speech
        }
    case "chipotle.kids - quesadilla": 
        switch speech {
        case "fillings":
            headerOut[3] = 1600
            talkback = "which fillings do you want?"
        case "rice":
            headerOut[3] = 1610
            talkback = "Any rice?"
        case "beans":
            headerOut[3] = 1620
            talkback = "Any beans?"
        case "kidsides":
            headerOut[3] = 1630 
            talkback = "Any sides for kids?"
        case "kidsdrinks":
            headerOut[3] = 1640 
            talkback = "Any drinks for kids?"
        case "Done":
            headerOut[3] = 1720
            talkback = "Okay, Do you want to add item to cart"
        default:
            talkback = speech
        }
    case "chipotle.sides&drinks":
        switch speech {
        case "address":
            headerOut[3] = 100
            talkback = "please select address, you can say recent, favorite, or nearby"
        case "sides":
            headerOut[3] = 1700
            talkback = "Any sides?"
        case "drinks":
            headerOut[3] = 1710
            talkback = "Any drinks?"
        case "Done":
            headerOut[3] = 1720
            talkback = "Okay, Do you want to add item to cart"
        default:
            talkback = speech
        }
    case "chipotle.sides&drinks - yes": 
        headerOut[3] = 1900
        talkback = speech
    case "chipotle.addtobag":
        headerOut[3] = 1900
        talkback = speech
    case "chipotle.cart":
        headerOut[3] = 5000
        talkback = speech
    case "chipotle.recents":
        headerOut[3] = 3000
        talkback = speech 
    case "chipotle.recents - select.number":
        headerOut[3] = 5000
        talkback = speech
    case "chipotle.confirm":
        switch speech {
        case "time":
            headerOut[3] = 6000
            talkback = "please tell me the pickup time"
        case "payment":
            headerOut[3] = 6100
            str := fmt.Sprintf("%v", entity["time"])
            str1, _ := time.Parse(time.RFC3339, str)
            entityback["time"] = str1.Format("3:04 PM")
            entityback["payment"] = entity["payment"]
            entity = entityback
            talkback = "please tell me payment type, you can say google pay or credit card"
        case "Done":
            headerOut[3] = 6200
            talkback = "Okay, Do you want to submit order?"
        default:
            talkback = speech
        }
    case "chipotle.confirm - yes":
        headerOut[3] = 7000
        talkback = speech
    default:
    	headerOut[3] = headerOut[2]
        talkback = speech
    }

    headerOut[4] = float64(time.Now().UnixNano() / 1000000)
    headerOut[5] = 3
    return headerOut, talkback, entity, nil
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("\nrecv: %s", message)

        var m Message
        err1 := json.Unmarshal(message, &m)

        if err1 != nil {
            log.Fatalln("error:", err1)
        }

        // check user status
        if val, ok := user[m.Header[0]]; ok {
        	var r Output
        	r = val
        	if m.Data.Query == "actionFalse" {
        		r.Data.Speech = "Sorry I can't perform the action right Now."
        	} else if m.Data.Query == "actionTrue" {
        		fmt.Println("Task performed")
        	} else {
        		r.Data.Speech = "Someting is wrong."
        	}
        	delete(user, m.Header[0])
        	b, _ := json.Marshal(r)
    		err = c.WriteMessage(mt, b)

			if err != nil {
				log.Println("write:", err)
				break
			}
        } else {
        	s, i, e, _ := DetectIntentText("chipotle-flat", "123", m.Data.Query, "en")
        	var p Output
        	p.Header, p.Data.Speech, p.Data.Entity, _ = HeaderProcess(m.Header, i, s, e)
        	if p.Header[2] != p.Header[3] {
        		user[m.Header[0]] = p
        		p.Data.Speech = "Performing task now."
        	}
        	fmt.Println(p)
        	b, _ := json.Marshal(p)
    		err = c.WriteMessage(mt, b)

			if err != nil {
				log.Println("write:", err)
				break
			}
        }
        // s, i, e, _ := DetectIntentText("chipotle-flat", "123", m.Data.Query, "en")
        // user[m.Header[0]] = m.Header[2]
        //fmt.Println(user)
        // var p Output
        // p.Header, p.Data.Speech, p.Data.Entity, _ = HeaderProcess(m.Header, i, s, e)
        // b, _ := json.Marshal(p)
        //fmt.Println(string(b))
		// err = c.WriteMessage(mt, b)

		// if err != nil {
		// 	log.Println("write:", err)
		// 	break
		// }
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/chipotle", echo)
	//log.Fatal(http.ListenAndServe(*addr, nil))
    log.Fatal(http.ListenAndServe(":8080", nil))
}
