
package main

import (
	"encoding/json"
	"io/ioutil"
	"flag"
	"log"
	"net/http"
	"html/template"
	"fmt"
)

var (
	stales = []*StaleUse{}
	freshes = []*FreshVal{}
	tmpl = template.Must(template.New("rng-results").Parse(myhtml))
)

var web = flag.Bool("web", false, "serve results via webserver in browser")

func main() {
	flag.Parse()
	name := flag.Arg(0)

	data, err := ioutil.ReadFile(name)
	if err != nil {
		log.Fatal(err)
	}

	raw := []map[string]string{}

	if err := json.Unmarshal(data, &raw); err != nil {
		log.Fatal(err)
	}

	for _, entry := range raw {
		if v, ok := entry["FreshRand"]; ok {
			fresh := &FreshVal{
				FreshRand: v,
				File: entry["File"],
				RLine: entry["RLine"],
			}
			freshes = append(freshes, fresh)
		} else {
			stale := &StaleUse{
				File: entry["File"],
				RandomValueAssignmentLine: entry["RandomValueAssignmentLine"],
				RandomValueLink: entry["RandomValueLink"],
				BlockingFunction: entry["BlockingFunction"],
				BlockingLine: entry["BlockingLine"],
				BlockingLink: entry["BlockingLink"],
				StaleRandomVariable: entry["StaleRandomVariable"],
				StaleUseLine: entry["StaleUseLine"],
				StaleLink: entry["StaleLink"],
			}
			stales = append(stales, stale)

		}
	}

	fmt.Print("[")
	for key, stale := range aggregateStale() {
		v := &Printout{
			File: stale[0].File,
			RandType: "Stale",
			FuncKey: key,
		}
		data, _ := json.Marshal(v)
		fmt.Println(string(data), ",")
	}
	count := 0
	aggFresh := aggregateFresh()
	for key, fresh := range aggFresh {
		v := &Printout{
			File: fresh[0].File,
			RandType: "Fresh",
			FuncKey: key,
		}
		data, _ := json.Marshal(v)
		count++
		if count == len(aggFresh) {
			fmt.Println(string(data), "]")
		} else {
			fmt.Println(string(data), ",")
		}
	}

	if *web {
		http.HandleFunc("/", handler)
		log.Print("listening on localhost:7777")
		if err := http.ListenAndServe(":7777", nil); err != nil {
			log.Fatal(err)
		}
	}
}

type Printout struct {
	File string
	RandType string
	FuncKey string
}

func handler(w http.ResponseWriter, r *http.Request) {
	if err := tmpl.Execute(w, &TmplData{stales, freshes}); err != nil {
		log.Fatal(err)
	}
}

func aggregateStale() map[string][]*StaleUse {
	result := map[string][]*StaleUse{}
	for _, v := range stales {
		key := fmt.Sprintf("%v:%v:%v", v.File, v.StaleRandomVariable, v.RandomValueAssignmentLine)
		result[key] = append(result[key], v)
	}
	return result
}

func aggregateFresh() map[string][]*FreshVal {
	result := map[string][]*FreshVal{}
	for _, v := range freshes {
		key := fmt.Sprintf("%v:%v", v.File, v.RLine)
		result[key] = append(result[key], v)
	}
	return result
}

type TmplData struct {
	Stales []*StaleUse
	Freshes []*FreshVal
}

type StaleUse struct {
	File string
	RandomValueAssignmentLine string
	RandomValueLink string
	BlockingFunction string
	BlockingLine string
	BlockingLink string
	StaleRandomVariable string
	StaleUseLine string
	StaleLink string
}

type FreshVal struct {
	FreshRand string
	File string
	RLine string
}

var myhtml =
`
<!DOCTYPE html>
<html>
  <head>
    <title>Banana Kingdom</title>
	<link href="http://twitter.github.com/bootstrap/assets/css/bootstrap.css" rel="stylesheet">
  </head>

  <body>
  	<table class="table">  
        <thead>  
          <tr>  
			<th>File</th>  
            <th>Rand Val Assignment Line</th>  
            <th>Blocking function</th>  
            <th>Blocking Line</th>  
            <th>Stale Rand Var</th>  
            <th>Stale Use Line</th>  
          </tr>  
        </thead>  

        <tbody>  
{{range $index, $entry := .Stales}}

          <tr>  
			<th>{{$entry.File}}</th>  
			<th><a href="{{$entry.RandomValueLink}}">{{$entry.RandomValueAssignmentLine}}</a></th>  
			<th>{{$entry.BlockingFunction}}</th>  
			<th><a href="{{$entry.BlockingLink}}">{{$entry.BlockingLine}}</a></th>  
			<th>{{$entry.StaleRandomVariable}}</th>  
			<th><a href="{{$entry.StaleLink}}">{{$entry.StaleUseLine}}</a></th>  
          </tr>  
{{end}}
        </tbody>  
      </table>  

  </body>
</html>

`

