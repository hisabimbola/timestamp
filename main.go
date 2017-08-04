package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("path", r.URL.Path)
	cwd, _ := os.Getwd()
	t, err := template.ParseFiles(filepath.Join(cwd, "./static/home.html"))
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(w, nil)
}

// I feel there is a more elegant way to parse date in Go.
func convertMonthToInt(month string) int {
	switch month {
	case "January":
		return 1
	case "February":
		return 2
	case "March":
		return 3
	case "April":
		return 4
	case "May":
		return 5
	case "June":
		return 6
	case "July":
		return 7
	case "August":
		return 8
	case "September":
		return 9
	case "October":
		return 10
	case "November":
		return 11
	case "December":
		return 12
	default:
		return 0
	}
}

type TimeRes struct {
	Unix    int    `json:"unix"`
	Natural string `json:"natural"`
}

func dateToStr(t time.Time) string {
	return t.Month().String() + " " + strconv.Itoa(t.Day()) + ", " + strconv.Itoa(t.Year())
}

func writeRes(w http.ResponseWriter, res TimeRes) {
	j, err := json.Marshal(res)
	if err != nil {
		fmt.Println("err when marshalling json", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func TimeHandler(w http.ResponseWriter, r *http.Request) {
	arg := r.URL.Path[1:]
	r1 := regexp.MustCompile(`^\d{10}$`)
	r2 := regexp.MustCompile(`^[[:alpha:]]{4,9}[[:blank:]][[:digit:]]{2},[[:blank:]][[:digit:]]{4}$`)
	isUnixTimestamp := r1.MatchString(arg)
	isNaturalDate := r2.MatchString(arg)

	if isNaturalDate {
		timeArr := strings.Split(arg, " ")
		year, _ := strconv.Atoi(timeArr[2])
		day, _ := strconv.Atoi(strings.TrimSuffix(timeArr[1], ","))
		month := convertMonthToInt(timeArr[0])

		t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
		res := TimeRes{int(t.Unix()), dateToStr(t)}
		writeRes(w, res)
	} else if isUnixTimestamp {
		i, _ := strconv.ParseInt(arg, 10, 64)
		t := time.Unix(i, 0)
		res := TimeRes{int(t.Unix()), dateToStr(t)}
		writeRes(w, res)
	} else {
		res := TimeRes{}
		writeRes(w, res)

	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/{timestamp}", TimeHandler)
	http.Handle("/", r)
	port := os.Getenv("PORT")
	log.Println("Port to listen to", port)
	if port == "" {
		port = "3000"
	}
	log.Fatal(http.ListenAndServe(port, nil))
}
