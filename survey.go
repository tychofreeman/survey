package main

import (
    "encoding/json"
    "net/http"
    "io"
    "fmt"
    "log"
    "text/template"
)

type Fields struct {
    Path string
}

func copyIndexFile(out io.Writer, path string) {
    tmpl, err := template.ParseFiles("survey.html")
    if err != nil {
        fmt.Printf("Could not load template survey.html: %v\n", err)
    } else {
        tmpl.Execute(out, Fields{Path: path})
    }
}

func copyToTotals(totals map[string]int, questions map[string]interface{}) map[string]int {
    for k, v := range questions {
        totals[k] = totals[k]
        if v == "y" {
            totals[k] = totals[k] + 1
        }
    }
    fmt.Printf("Totals:\n")
    for k, v := range totals {
        fmt.Printf("\t%v = %v\n", k, v)
    }
    return totals
}

func main() {
    totals := make(map[string]int)
    records := make(map[string][]map[string]interface{})
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        copyIndexFile(w, r.URL.Path)
    })
    http.Handle("/a/", http.StripPrefix("/a/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        copyIndexFile(w, r.URL.Path)
    })))

    http.HandleFunc("/add-result", func(w http.ResponseWriter, r *http.Request) {
        r.ParseForm()

        questions := make(map[string]interface{})
        id := r.Form["id"][0]
        err := json.Unmarshal([]byte(r.Form["questions"][0]), &questions)
        if err != nil {
            fmt.Printf("Error: %v\n", err)
        } else {
            totals = copyToTotals(totals, questions)
            records[id] = append(records[id], questions)
            fmt.Printf("Appending record for %v to records: %v\n", id, records)
            for id, recs := range records {
                fmt.Printf("-- %v -- \n", id)
                for _,r := range recs {
                    fmt.Printf("\t-------\n")
                    for q, a := range r {
                        fmt.Printf("\t%v -- %v\n", q, a)
                    }
                }
            }
        }
    })

    log.Fatal(http.ListenAndServe(":8080", nil))
}
