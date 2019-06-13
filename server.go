// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"flag"
	"log"
    "fmt"
    "os"
    "strings"
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
    Result string
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

var ordertype = map[string][]string{
	"burrito": []string{"burrito"},
	"burrito bowl": []string{"bowl"},
	"tacos": []string{"taco"},
	"salad": []string{"salad"},
	"kid's Meal": []string{"kid's meal ", "kid", "kids"},
}

var address = map[string][]string{
	"recent": []string{"recent", "recently"},
	"favorites": []string{"favorites", "favorite", "fav"},
	"nearby": []string{"nearby", "nearest", "closest"},
}

var fillings = map[string][]string{
	"steak" : []string{"steak", "beef"},
	"carnitas": []string{"carnitas", "pork"},
	"chicken" : []string{"chicken"},
	"barbacoa": []string{"barbacoa"},
	"veggie": []string{"veggie", "vegetable"},
	"sofritas": []string{"sofritas", "sofrito"},
}

var beans = map[string][]string{
	"black beans": []string{"black"},
	"pinto beans": []string{"pinto"},
	"no beans": []string{"no beans", "no"},
}

var rice = map[string][]string{
	"brown rice" : []string{"brown"},
	"white rice" : []string{"white"},
	"no rice" : []string{"no rice", "no"},
}

var salsa = map[string][]string{
    "tomatillo-green chili salsa": []string{"green chili", "medium"},
    "tomatillo-red chili salsa": []string{"red chili", "hot"},
    "fresh tomato salsa": []string{"salsa", "mild"},
}

var tops = map[string][]string{
	"roasted chili-corn salsa": []string{"corn"},
	"queso": []string{"queso"},
	"romaine lettuce": []string{"lettuce"},
	"tomatillo-green chili salsa": []string{"green chili"},
	"fajita veggies": []string{"fajita"},
	"fresh tomato salsa": []string{"salsa"},
	"guacamole": []string{"guac"},
	"cheese": []string{"cheese"},
	"tomatillo-red chili salsa": []string{"red chili"},
	"double wrap with tortilla": []string{"tortilla"},
	"sour cream": []string{"sour cream"},
}

var sides = map[string][]string{
	"chips": []string{"chips"},
    "tortilla on the side": []string{"tortilla"},
}

var drinks = map[string][]string{
	"bottled water": []string{"water"},
	"22 fl oz soda/iced tea": []string{"small soda"},
	"32 fl oz soda/iced tea": []string{"large soda"},
    "pressed apple juice": []string{"apple juice"},
    "blackberry izze": []string{"blackberry izze"},
    "grapefruit izze": []string{"grapefruit izze"},
}

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
        talkback = "please select address, you can say recent, favorite, or nearby"
        headerOut[3] = 2000
        entityback["ordertype"] = "burrito"
        entity = entityback
    case "chipotle.bowl":
        talkback = "please select address, you can say recent, favorite, or nearby"
        headerOut[3] = 2000
        entityback["ordertype"] = "bowl"
        entity = entityback
    case "chipotle.salad":
        talkback = "please select address, you can say recent, favorite, or nearby"
        headerOut[3] = 2000
        entityback["ordertype"] = "salad"
        entity = entityback
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
    case "chipotle.cancel":
        headerOut[3] = 0
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
		log.Printf("recv: %s", message)

        var m Message
        err1 := json.Unmarshal(message, &m)

        if err1 != nil {
            log.Fatalln("error:", err1)
        }

        // check user status
        if val, ok := user[m.Header[0]]; ok {
        	var r Output
        	r = val
	        switch m.Data.Result {
	        case "actionFalse":
	            r.Data.Speech = "Sorry I can't perform the action right Now."
	        case "actionTrue":
				r.Header[2] = r.Header[3]
			case "pageWrong":
				r.Data.Speech = fmt.Sprintf("Wrong page, go to level %d to perform action", int(r.Header[2]))
	        case "itemNotFound: ":
	            r.Data.Speech = ""
	        default:
	        	r.Data.Speech = "Error type not found"
	        }
        	r.Header[3] = 9999
        	delete(user, m.Header[0])
        	b, _ := json.Marshal(r)
        	fmt.Println(string(b))
        	log.Printf("sent: %s", string(b))
    		err = c.WriteMessage(mt, b)

			if err != nil {
				log.Println("write:", err)
				break
			}
        } else {
        	var p Output
        	p.Header[0] = m.Header[0]
        	p.Header[1] = m.Header[1]
        	p.Header[2] = m.Header[2]
   			p.Header[4] = float64(time.Now().UnixNano() / 1000000)
    		p.Header[5] = 3
    		entityback := make(map[string]interface{})
        	switch m.Header[2] {
            case 100:
                if strings.Contains(m.Data.Query, "pick up") {
                    p.Header[3] = 2000
                    entityback["servicetype"] = "pick up"
                    p.Data.Entity = entityback
                    p.Data.Speech = "Which store do you want to pick up? you can say recent, favorite, or nearby"
                    user[m.Header[0]] = p
                    p.Data.Speech = "selecting"
                } else if strings.Contains(m.Data.Query, "deliver") {
                    p.Header[3] = 2000
                    entityback["servicetype"] = "deliver"
                    p.Data.Entity = entityback
                    p.Data.Speech = "What is your deliver address? you can say recent."
                    user[m.Header[0]] = p
                    p.Data.Speech = "selecting"
                } else if strings.Contains(m.Data.Query, "menu") {
                    p.Header[3] = 1000
                    p.Data.Entity = entityback
                    p.Data.Speech = "What item do you want, burrito, bowl or tacos?"
                    user[m.Header[0]] = p
                    p.Data.Speech = "selecting"
                } else {
                    s, i, e, _ := DetectIntentText("chipotle-flat", "123", m.Data.Query, "en")
                    p.Header, p.Data.Speech, p.Data.Entity, _ = HeaderProcess(m.Header, i, s, e)
                    if strings.Contains(p.Data.Speech, "cancel"){
                        p.Header[3] = 0
                    } else if p.Header[2] != p.Header[3] {
                        user[m.Header[0]] = p
                        p.Data.Speech = "Performing task now."
                    } else {
                        p.Header[3] = 9999
                        p.Data.Speech = "Hello, this is chipotle, Do you want to pick up in store or deliver to an address?"
                    }
                }
            case 1000:
                p.Data.Speech = "what items do you want, burrito, bowl, or tacos?"
                p.Header[3] = 9999
                for k, v := range ordertype {
                    for _, item := range v {
                        if strings.Contains(m.Data.Query, item) {   
                            entityback["ordertype"] = k
                            p.Data.Entity = entityback
                            p.Data.Speech = "Choose your meat or veggie"
                            user[m.Header[0]] = p
                            p.Data.Speech = "selecting"
                            p.Header[3] = 1100
                        }
                    }
                }
                if strings.Contains(m.Data.Query, "cancel") {
                    p.Header[3] = 0
                    p.Data.Speech = "Okay, Cancel ordering"
                }
        	case 2000:
        		p.Data.Speech = "please select address, you can say recent, favorite, or nearby"
        		p.Header[3] = 9999
	        	for k, v := range address {
	        		for _, item := range v {
	        			if strings.Contains(m.Data.Query, item) {   
	        				entityback["address"] = k
	        				p.Data.Entity = entityback
	        				p.Data.Speech = "Choose your meat or veggie"
                            p.Header[3] = 1100
	        				user[m.Header[0]] = p
	        				p.Data.Speech = "selecting"
	        			}
	        		}
	        	}
	        	if strings.Contains(m.Data.Query, "cancel") {
	        		p.Header[3] = 0
	        		p.Data.Speech = "Okay, Cancel ordering"
	        	}
	        case 1100:
	        	p.Data.Speech = "Choose your meat or veggie"
	        	p.Header[3] = 9999
	        	var s []string
	        	for k, v := range fillings {
	        		for _, item := range v {
	        			if strings.Contains(m.Data.Query, item) {
	        				s = append(s, k)
	        				entityback["fillings"] = s
	        			    p.Data.Entity = entityback
	        			    p.Data.Speech = "Now add your rice and beans"
                            p.Header[3] = 1110
	        				user[m.Header[0]] = p
	        				p.Data.Speech = "selecting"
	        			} 
	        		}
	        	}
	        	if strings.Contains(m.Data.Query, "cancel") {
	        		p.Header[3] = 0
	        		p.Data.Speech = "Okay, Cancel ordering"
	        	}
	        case 1110:
	        	p.Data.Speech = "Now add your rice and beans"
	        	p.Header[3] = 9999
	        	var s_rice []string
                var s_beans []string
	        	for k, v := range rice {
	        		for _, item := range v {
	        			if strings.Contains(m.Data.Query, item) {
	        				s_rice = append(s_rice, k)
	        				entityback["rice"] = s_rice
	        			    p.Data.Entity = entityback
	        			    p.Data.Speech = "Do you want to add salsa? Mild, medium, or hot?"
	        				user[m.Header[0]] = p
	        				p.Data.Speech = "selecting"
                            p.Header[3] = 1120
	        			} 
	        		}
	        	}
                for k, v := range beans {
                    for _, item := range v {
                        if strings.Contains(m.Data.Query, item) {
                            s_beans = append(s_beans, k)
                            entityback["beans"] = s_beans
                            p.Data.Entity = entityback
                            p.Data.Speech = "Do you want to add salsa? Mild, medium, or hot?"
                            user[m.Header[0]] = p
                            p.Data.Speech = "selecting"
                            p.Header[3] = 1120
                        } 
                    }
                }
	        	if strings.Contains(m.Data.Query, "cancel") {
	        		p.Header[3] = 0
	        		p.Data.Speech = "Okay, Cancel ordering"
	        	}
	        case 1120:
	        	p.Data.Speech = "Do you want to add salsa? Mild, medium, or hot?"
	        	p.Header[3] = 9999
	        	var s []string
	        	for k, v := range salsa {
	        		for _, item := range v {
	        			if strings.Contains(m.Data.Query, item) {
	        				s = append(s, k)
	        				entityback["tops"] = s
	        			    p.Data.Entity = entityback
	        			    p.Data.Speech = "Do you want queso, guac, or corn?"
                            p.Header[3] = 1130
	        				user[m.Header[0]] = p
	        				p.Data.Speech = "selecting"
	        			} 
	        		}
	        	}
	        	if strings.Contains(m.Data.Query, "cancel") {
	        		p.Header[3] = 100
	        		p.Data.Speech = "Okay, Cancel ordering"
	        	}
	        case 1130:
	        	p.Data.Speech = "Do you want queso, guac, or corn?"
	        	p.Header[3] = 9999
	        	var s []string
	        	for k, v := range tops {
	        		for _, item := range v {
	        			if strings.Contains(m.Data.Query, item) {
	        				s = append(s, k)
	        				entityback["tops"] = s
	        			    p.Data.Entity = entityback
	        			    p.Data.Speech = "how about sour cream, fajita veggies, cheese, and lettuce?"
	        				p.Header[3] = 1140
                            user[m.Header[0]] = p
	        				p.Data.Speech = "selecting"
	        			} 
	        		}
	        	}
	        	if strings.Contains(m.Data.Query, "cancel") {
	        		p.Header[3] = 0
	        		p.Data.Speech = "Okay, Cancel ordering"
	        	}
	        	if strings.Contains(m.Data.Query, "no") {
	        		p.Header[3] = 9999
	        		p.Data.Entity = entityback
	        		p.Data.Speech = "how about sour cream, fajita veggies, cheese, and lettuce?"
	        		user[m.Header[0]] = p
	        		p.Data.Speech = "okay no sides"
	        	}
	        case 1140:
	        	p.Data.Speech = "how about sour cream, fajita veggies, cheese, and lettuce?"
	        	p.Header[3] = 9999
	        	var s []string
	        	for k, v := range tops {
	        		for _, item := range v {
	        			if strings.Contains(m.Data.Query, item) {
	        				s = append(s, k)
	        				entityback["top"] = s
	        			    p.Data.Entity = entityback
	        			    p.Data.Speech = "Do you want chips or tortilla as sides?"
                            p.Header[3] = 1150
	        				user[m.Header[0]] = p
	        				p.Data.Speech = "selecting"
	        			} 
	        		}
	        	}
	        	if strings.Contains(m.Data.Query, "no") {
	        		p.Header[3] = 9999
	        		p.Data.Speech = "Do you want chips or tortilla as sides?"
	        	}
	        	if strings.Contains(m.Data.Query, "cancel") {
	        		p.Header[3] = 0
	        		p.Data.Speech = "Okay, Cancel ordering"
	        	}
	        case 1150:
	        	p.Data.Speech = "Do you want chips or tortilla as sides?"
	        	p.Header[3] = 9999
	        	var s []string
	        	for k, v := range sides {
	        		for _, item := range v {
	        			if strings.Contains(m.Data.Query, item) {
	        				s = append(s, k)
	        				entityback["sides"] = s
	        			    p.Data.Entity = entityback
	        			    p.Data.Speech = "Do you want fountain soda, or bottled juice?"
                            p.Header[3] = 1160
	        				user[m.Header[0]] = p
	        				p.Data.Speech = "selecting"
	        			} 
	        		}
	        	}
	        	if strings.Contains(m.Data.Query, "no") {
	        		p.Header[3] = 9999
	        		p.Data.Speech = "Do you want fountain soda, or bottled juice?"
	        	}
	        	if strings.Contains(m.Data.Query, "cancel") {
	        		p.Header[3] = 0
	        		p.Data.Speech = "Okay, Cancelled"
	        	}
            case 1160:
                p.Data.Speech = "Do you want fountain soda, or bottled juice?"
                p.Header[3] = 9999
                var s []string
                for k, v := range drinks {
                    for _, item := range v {
                        if strings.Contains(m.Data.Query, item) {
                            s = append(s, k)
                            entityback["drinks"] = s
                            p.Data.Entity = entityback
                            p.Data.Speech = "Do you want to add item to cart?"
                            p.Header[3] = 1170
                            user[m.Header[0]] = p
                            p.Data.Speech = "selecting"
                        } 
                    }
                }
                if strings.Contains(m.Data.Query, "no") {
                    p.Header[3] = 9999
                    p.Data.Speech = "Do you want to add item to cart?"
                }
                if strings.Contains(m.Data.Query, "cancel") {
                    p.Header[3] = 0
                    p.Data.Speech = "Okay, Cancelled"
                }
	        case 1170:
	        	p.Data.Speech = "Do you want to add item to cart?"
	        	p.Header[3] = 9999
	        	if strings.Contains(m.Data.Query, "yes") {
	        		p.Data.Speech = "Okay, item add to cart"
                    p.Header[3] = 1900
	        		user[m.Header[0]] = p
	        		p.Data.Speech = "selecting"
	        	}
	        	if strings.Contains(m.Data.Query, "no") {
	        		p.Header[3] = 9999
	        		p.Data.Speech = "Sure please tell me when you ready."
	        	}
	        	if strings.Contains(m.Data.Query, "cancel") {
	        		p.Header[3] = 0
	        		p.Data.Speech = "Okay, Cancelled"
	        	}        
        	default:
	        	s, i, e, _ := DetectIntentText("chipotle-flat", "123", m.Data.Query, "en")
	        	p.Header, p.Data.Speech, p.Data.Entity, _ = HeaderProcess(m.Header, i, s, e)
	        	if strings.Contains(p.Data.Speech, "cancel"){
	        		p.Header[3] = 0
	        	} else if p.Header[2] != p.Header[3] {
	        		user[m.Header[0]] = p
	        		p.Data.Speech = "Performing task now."
	        	} else {
	        		p.Header[3] = 9999
	        	}
        	}
        	b, _ := json.Marshal(p)
	        fmt.Println(string(b))
	        log.Printf("sent: %s", string(b))
	    	err = c.WriteMessage(mt, b)

			if err != nil {
				log.Println("write:", err)
				break
			}
        }
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/chipotle", echo)
	f, err := os.OpenFile("log/logfile.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
	    log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	//log.Fatal(http.ListenAndServe(*addr, nil))
    log.Fatal(http.ListenAndServe(":8080", nil))
}
