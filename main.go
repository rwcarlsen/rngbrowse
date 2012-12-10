
package main

import (
	"encoding/json"
	"io/ioutil"
	"flag"
	"log"
	"net/http"
	"html/template"
)

var (
	stales = []*StaleUse{}
	freshes = []*FreshVal{}
	tmpl = template.Must(template.New("rng-results").Parse(myhtml))
)
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

    http.HandleFunc("/", handler)
	log.Print("listening on localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	if err := tmpl.Execute(w, &TmplData{stales, freshes}); err != nil {
		log.Fatal(err)
	}
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

