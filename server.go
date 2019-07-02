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
    "time"
    "chipotle/dialogflow"

	"github.com/gorilla/websocket"
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
	"burrito bowl": []string{"bowl", "burrito bowl"},
	"tacos": []string{"taco"},
	"salad": []string{"salad"},
	"kid's meal": []string{"kid"},
	"sides & drinks": []string{"side", "drink"},
}

var address = map[string][]string{
	"recent": []string{"recent", "recently"},
	"favorites": []string{"favorites", "favorite", "fav"},
	"nearby": []string{"nearby", "nearest", "closest"},
}

var fillings = map[string][]string{
	"steak" : []string{"steak"},
	"carnitas": []string{"carnitas", "pork"},
	"chicken" : []string{"chicken"},
	"barbacoa": []string{"barbacoa", "beef", "barbeque", "bbq"},
	"veggie": []string{"veggie", "vegetable", "guac"},
	"sofritas": []string{"sofritas", "tofu"},
}

var beans = map[string][]string{
	"black beans": []string{"black"},
	"pinto beans": []string{"pinto"},
	"no beans": []string{"no beans"},
}

var rice = map[string][]string{
	"brown rice" : []string{"brown"},
	"white rice" : []string{"white"},
	"no rice" : []string{"no rice"},
}

var salsa = map[string][]string{
    "tomatillo-green chili salsa": []string{"green chili", "medium"},
    "tomatillo-red chili salsa": []string{"red", "hot"},
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
	"22 fl oz soda/iced tea": []string{"small soda", "small fountain soda", "small", "fountain soda"},
	"32 fl oz soda/iced tea": []string{"large soda", "large fountain soda", "large"},
    "pressed apple juice": []string{"apple juice"},
    "blackberry izze": []string{"blackberry izze"},
    "grapefruit izze": []string{"grapefruit izze"},
}

var numbers = map[string][]string{
	"one taco": []string{"one", "1"},
	"two tacos": []string{"two", "2"},
	"three tacos": []string{"three", "3"},
}

var tacotype = map[string][]string{
	"soft flour tortilla": []string{"soft"},
	"crispy corn tortilla": []string{"crispy"},
}

var user = map[float64]Output{}
var tacoflag = map[float64]bool{}

var upgrader = websocket.Upgrader{} // use default options

func SingleMatch (query string, keyword map[string][]string) (wordMatch string){
	for k, v := range keyword {
        for _, item := range v {
            if strings.Contains(query, item) {
            	wordMatch = k
            }
        }
    }
    return wordMatch	
}

func MultipleMatch (query string, keyword map[string][]string) (wordMatch []string){
	for k, v := range keyword {
        for _, item := range v {
            if strings.Contains(query, item) {
            	wordMatch = append(wordMatch, k)
            }
        }
    }
    return wordMatch
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
        talkback = "please select address, you can say recent, favorite, or nearby"
        headerOut[3] = 2000
        entityback["ordertype"] = "tacos"
        entity = entityback
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
            talkback = "Okay, Do you want to add item to bag"
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
            talkback = "Okay, Do you want to add item to bag"
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
        case "Done":
            headerOut[3] = 6100
            str := fmt.Sprintf("%v", entity["time"])
            str1, _ := time.Parse(time.RFC3339, str)
            entityback["time"] = str1.Format("3:04 PM")
            // entityback["payment"] = entity["payment"]
            entity = entityback
            talkback = "Please touch to make your payment and submit order."
        // case "Done":
        //     headerOut[3] = 6200
        //     talkback = "Okay, Do you want to submit order?"
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
                    p.Data.Speech = "OK, which store do you prefer? you can say recent, favorite, or nearby."
                    user[m.Header[0]] = p
                    p.Data.Speech = "selecting"
                } else if strings.Contains(m.Data.Query, "deliver") {
                    p.Header[3] = 2000
                    entityback["servicetype"] = "deliver"
                    p.Data.Entity = entityback
                    p.Data.Speech = "OK, what's your address? You can choose from recent, or add a new one."
                    user[m.Header[0]] = p
                    p.Data.Speech = "selecting"
                } else if strings.Contains(m.Data.Query, "menu") {
                    p.Header[3] = 1000
                    p.Data.Entity = entityback
                    p.Data.Speech = "What item do you want, burrito, bowl or tacos?"
                    user[m.Header[0]] = p
                    p.Data.Speech = "selecting"
                } else {
                    s, i, e, _ := dialogflow.DetectIntentText("chipotle-flat", "123", m.Data.Query, "en")
                    p.Header, p.Data.Speech, p.Data.Entity, _ = HeaderProcess(m.Header, i, s, e)
                    if strings.Contains(p.Data.Speech, "cancel"){
                        p.Header[3] = 0
                    } else if p.Header[2] != p.Header[3] {
                        user[m.Header[0]] = p
                        p.Data.Speech = "Performing task now."
                    } else {
                        p.Header[3] = 9999
                        p.Data.Speech = "Hello, this is Chipotle, Do you want to pick up in store or deliver to an address?"
                    }
                }
        	case 2000:
        		p.Data.Speech = "OK, which store do you prefer? you can say recent, favorite, or nearby."
        		p.Header[3] = 9999
        		if matched := SingleMatch(m.Data.Query, address); len(matched) != 0 {
        			entityback["address"] = matched
        			p.Data.Entity = entityback
        			p.Data.Speech = "what items do you want, burrito, bowl, or tacos?"
        			p.Header[3] = 1000
        			user[m.Header[0]] = p
        			p.Data.Speech = "selecting"
        		}else if strings.Contains(m.Data.Query, "cancel") {
	        		p.Header[3] = 0
	        		p.Data.Speech = "Okay, Cancel ordering"
	        	}
            case 1000:
                p.Data.Speech = "what items do you want, burrito, bowl, or tacos?"
                p.Header[3] = 9999
                if matched := SingleMatch(m.Data.Query, ordertype); len(matched) != 0 {
	              	entityback["ordertype"] = matched
	                p.Data.Entity = entityback
                	if matched == "tacos" {
                		p.Header[3] = 1100
	                	p.Data.Speech = "OK. How many tacos do you want?"
	                	user[m.Header[0]] = p
	                	p.Data.Speech = "selecting"
	                	tacoflag[m.Header[0]] = true
                	} else if matched == "kid's meal" {
                		p.Header[3] = 9999
	                	p.Data.Speech = "Voice not be enabled for kid's meal right now, please touch to progress"
                	} else if matched == "sides & drinks" {
                		p.Header[3] = 9999
	                	p.Data.Speech = "Voice not be enabled for sides & drinks right now, please touch to progress"
                	} else {
                		p.Header[3] = 1100
	                	p.Data.Speech = "Choose your meat or veggie"
	                	user[m.Header[0]] = p
	                	p.Data.Speech = "selecting"
	                }
                }else if strings.Contains(m.Data.Query, "cancel") {
                    p.Header[3] = 0
                    p.Data.Speech = "Okay, Cancel ordering"
                }
	        case 1100:
	        	if tacoflag[m.Header[0]] {
		        	p.Data.Speech = "OK. How many tacos do you want?"
		        	p.Header[3] = 9999
		     		if matched := MultipleMatch(m.Data.Query, numbers); len(matched) != 0 {
		        		fmt.Println(matched)
		        		entityback["quantity"] = matched
		        		p.Data.Entity = entityback
			        	p.Data.Speech = "Do you want soft or crispy taco?"
			        	p.Header[3] = 1101
			        	user[m.Header[0]] = p
			        	p.Data.Speech = "selecting"
			        	tacoflag[m.Header[0]] = false
		        	}else if strings.Contains(m.Data.Query, "cancel") {
		        		p.Header[3] = 0
		        		p.Data.Speech = "Okay, Cancel ordering"
		        	}	        		
	        	} else {
		        	p.Data.Speech = "OK. First choose your meat or veggie."
		        	p.Header[3] = 9999
		        	if matched := MultipleMatch(m.Data.Query, fillings); len(matched) != 0 {
		        		fmt.Println(matched)
		        		entityback["fillings"] = matched
		        		p.Data.Entity = entityback
			        	p.Data.Speech = "Now add your rice and beans"
			        	p.Header[3] = 1110
			        	user[m.Header[0]] = p
			        	p.Data.Speech = "selecting"
		        	}else if strings.Contains(m.Data.Query, "cancel") {
		        		p.Header[3] = 0
		        		p.Data.Speech = "Okay, Cancel ordering"
		        	}
		        }
	        case 1101:
	        	p.Data.Speech = "Do you want soft or crispy taco?"
	        	p.Header[3] = 9999
	        	if matched := MultipleMatch(m.Data.Query, tacotype); len(matched) != 0 {
        			entityback["tacotype"] = matched
        			p.Data.Entity = entityback
        			p.Data.Speech = "OK. choose your meat or veggie."
        			p.Header[3] = 1102
        			user[m.Header[0]] = p
        			p.Data.Speech = "selecting"
        		}else if strings.Contains(m.Data.Query, "cancel") {
	        		p.Header[3] = 0
	        		p.Data.Speech = "Okay, Cancel ordering"
	        	}
	        case 1102:
	        	p.Data.Speech = "OK. choose your meat or veggie."
	        	p.Header[3] = 9999
	        	if matched := MultipleMatch(m.Data.Query, fillings); len(matched) != 0 {
	        		fmt.Println(matched)
	        		entityback["fillings"] = matched
	        		p.Data.Entity = entityback
		        	p.Data.Speech = "Now add your rice and beans"
		        	p.Header[3] = 1110
		        	user[m.Header[0]] = p
		        	p.Data.Speech = "selecting"
	        	}else if strings.Contains(m.Data.Query, "cancel") {
	        		p.Header[3] = 0
	        		p.Data.Speech = "Okay, Cancel ordering"
	        	}
	        case 1110:
	        	p.Data.Speech = "Now add your rice and beans"
	        	p.Header[3] = 9999
	        	s_rice := MultipleMatch(m.Data.Query, rice)
                s_beans := MultipleMatch(m.Data.Query, beans)
                if len(s_rice) != 0 && len(s_beans) != 0 {
                	entityback["rice"] = s_rice
                	entityback["beans"] = s_beans
                	p.Data.Entity = entityback
                	p.Data.Speech = "Do you want to add salsa? Mild, medium, or hot?"
                    p.Header[3] = 1120
                	user[m.Header[0]] = p
                	p.Data.Speech = "selecting"
                }else if len(s_rice) != 0 {
                	entityback["rice"] = s_rice
                	p.Data.Entity = entityback
                	p.Data.Speech = "Now add your beans"
                	p.Header[3] = 1112
                	user[m.Header[0]] = p
                	p.Data.Speech = "selecting"
                }else if len(s_beans) != 0 {
                	entityback["beans"] = s_beans
                	p.Data.Entity = entityback
                	p.Data.Speech = "Now add your rice"
                	p.Header[3] = 1111
                	user[m.Header[0]] = p
                	p.Data.Speech = "selecting"
                }else if strings.Contains(m.Data.Query, "cancel") {
	        		p.Header[3] = 0
	        		p.Data.Speech = "Okay, Cancel ordering"
	        	}
	        case 1111:
	        	p.Data.Speech = "Now add your rice"
	        	p.Header[3] = 9999
	        	if matched := MultipleMatch(m.Data.Query, rice); len(matched) != 0 {
        			entityback["rice"] = matched
        			p.Data.Entity = entityback
        			p.Data.Speech = "Do you want to add salsa? Mild, medium, or hot?"
        			p.Header[3] = 1120
        			user[m.Header[0]] = p
        			p.Data.Speech = "selecting"
        		}else if strings.Contains(m.Data.Query, "cancel") {
	        		p.Header[3] = 0
	        		p.Data.Speech = "Okay, Cancel ordering"
	        	}
	        case 1112:
	        	p.Data.Speech = "Now add your beans"
	        	p.Header[3] = 9999
	        	if matched := MultipleMatch(m.Data.Query, beans); len(matched) != 0 {
        			entityback["beans"] = matched
        			p.Data.Entity = entityback
        			p.Data.Speech = "Do you want to add salsa? Mild, medium, or hot?"
        			p.Header[3] = 1120
        			user[m.Header[0]] = p
        			p.Data.Speech = "selecting"
        		}else if strings.Contains(m.Data.Query, "cancel") {
	        		p.Header[3] = 0
	        		p.Data.Speech = "Okay, Cancel ordering"
	        	}
	        case 1120:
	        	p.Data.Speech = "Do you want to add salsa? Mild, medium, or hot?"
	        	p.Header[3] = 9999

	        	if matched := MultipleMatch(m.Data.Query, salsa); len(matched) != 0 {
	        		entityback["tops"] = matched
	        		p.Data.Entity = entityback
	        		p.Data.Speech = "Do you want queso, guac, or corn?"
	        		p.Header[3] = 1130
	        		user[m.Header[0]] = p
	        		p.Data.Speech = "selecting"
	        	} else if strings.Contains(m.Data.Query, "cancel") {
	        		p.Header[3] = 100
	        		p.Data.Speech = "Okay, Cancel ordering"
	        	} else if strings.Contains(m.Data.Query, "no") {
                    entityback["tops"] = []string{}
	        		p.Data.Entity = entityback
	        	    p.Data.Speech = "Do you want queso, guac, or corn?"
                    p.Header[3] = 1130
      				user[m.Header[0]] = p
	        		p.Data.Speech = "selecting"
                }
	        case 1130:
	        	p.Data.Speech = "Do you want queso, guac, or corn?"
	        	p.Header[3] = 9999
	        	if matched := MultipleMatch(m.Data.Query, tops); len(matched) != 0 {
	        		entityback["tops"] = matched
	        		p.Data.Entity = entityback
	        		p.Data.Speech = "how about sour cream, fajita veggies, cheese, and lettuce?"
	        		p.Header[3] = 1140
	        		user[m.Header[0]] = p
	        		p.Data.Speech = "selecting"
	        	} 
	        	if strings.Contains(m.Data.Query, "cancel") {
	        		p.Header[3] = 0
	        		p.Data.Speech = "Okay, Cancel ordering"
	        	}
	        	if strings.Contains(m.Data.Query, "no") {
                    entityback["tops"] = []string{}
	        		p.Data.Entity = entityback
	        	    p.Data.Speech = "how about sour cream, fajita veggies, cheese, and lettuce?"
                    p.Header[3] = 1140
      				user[m.Header[0]] = p
	        		p.Data.Speech = "selecting"
	        	}
	        case 1140:
	        	p.Data.Speech = "how about sour cream, fajita veggies, cheese, and lettuce?"
	        	p.Header[3] = 9999
	        	var s []string
	        	for k, v := range tops {
	        		for _, item := range v {
	        			if strings.Contains(m.Data.Query, item) {
	        				s = append(s, k)
	        				entityback["tops"] = s
	        			    p.Data.Entity = entityback
	        			    p.Data.Speech = "Any tortilla or chips?"
                            p.Header[3] = 1150
	        				user[m.Header[0]] = p
	        				p.Data.Speech = "selecting"
	        			} 
	        		}
	        	}
	        	if strings.Contains(m.Data.Query, "no") {
                    entityback["tops"] = []string{}
	        		p.Data.Entity = entityback
	        	    p.Data.Speech = "Any tortilla or chips?"
                    p.Header[3] = 1150
      				user[m.Header[0]] = p
	        		p.Data.Speech = "selecting"
	        	}
	        	if strings.Contains(m.Data.Query, "cancel") {
	        		p.Header[3] = 0
	        		p.Data.Speech = "Okay, Cancel ordering"
	        	}
	        case 1150:
	        	p.Data.Speech = "Any tortilla or chips?"
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
                    entityback["sides"] = []string{}
	        		p.Data.Entity = entityback
	        	    p.Data.Speech = "Do you want fountain soda, or bottled juice?"
                    p.Header[3] = 1160
      				user[m.Header[0]] = p
	        		p.Data.Speech = "selecting"
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
                            p.Data.Speech = "Do you want to add item to bag?"
                            p.Header[3] = 1170
                            user[m.Header[0]] = p
                            p.Data.Speech = "selecting"
                        } 
                    }
                }
                if strings.Contains(m.Data.Query, "no") {
                    entityback["drinks"] = []string{}
	        		p.Data.Entity = entityback
	        	    p.Data.Speech = "Do you want to add item to bag?"
                    p.Header[3] = 1170
      				user[m.Header[0]] = p
	        		p.Data.Speech = "selecting"
                }
                if strings.Contains(m.Data.Query, "cancel") {
                    p.Header[3] = 0
                    p.Data.Speech = "Okay, Cancelled"
                }
	        case 1170:
	        	p.Data.Speech = "Do you want to add item to bag?"
	        	p.Header[3] = 9999
	        	if strings.Contains(m.Data.Query, "ye") || strings.Contains(m.Data.Query, "sure"){
	        		p.Data.Speech = "Okay, item add to bag"
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
	        case 5000:
	        	if strings.Contains(m.Data.Query, "edit") || strings.Contains(m.Data.Query, "remove") || strings.Contains(m.Data.Query, "duplicate") {
	        		p.Data.Speech = "Voice is not enabled for this yet, please proceed with touch."
                    p.Header[3] = 9999
	        	} else {
		        	s, i, e, _ := dialogflow.DetectIntentText("chipotle-flat", "123", m.Data.Query, "en")
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
        	default:
	        	s, i, e, _ := dialogflow.DetectIntentText("chipotle-flat", "123", m.Data.Query, "en")
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
