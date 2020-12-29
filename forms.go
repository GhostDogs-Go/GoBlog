package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type PostDetails struct {
	title         string
	draft         string
	tags          string
	author        string
	categories    string
	images        string
	altImages     string
	stretchImages string
	story         string
	toml          string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "templates/forms.html")
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		details := PostDetails{
			title:         r.FormValue("title"),
			tags:          r.FormValue("tags"),
			author:        r.FormValue("author"),
			story:         r.FormValue("story"),
			categories:    r.FormValue("categories"),
			images:        r.FormValue("images"),
			altImages:     r.FormValue("altImages"),
			stretchImages: r.FormValue("stretchImages"),
			draft:         r.FormValue("draft"),
			toml:          r.FormValue("toml"),
		}

		// do something with details
		_ = details

		fmt.Println(details.title)
		fmt.Println(details.tags)
		fmt.Println(details.categories)
		fmt.Println(details.images)
		fmt.Println(details.altImages)
		fmt.Println(details.tags)
		fmt.Println(details.story)
		fmt.Println(details.author)
		fmt.Println(details.toml)

		if details.draft == "on" {
			details.draft = "false"
		}
		if details.draft == "" {
			details.draft = "true"
		}

		if details.stretchImages == "on" {
			details.stretchImages = "true"
		}
		if details.stretchImages == "" {
			details.stretchImages = "false"
		}

		fmt.Println("Stretch " + details.stretchImages)
		fmt.Println("Draft " + details.draft)

		f, error := os.Create(details.title + ".md")
		check(error)
		defer f.Close()

		markToml := "+++"

		tm := time.Now()
		postDate := tm.Format(time.ANSIC)

		f.WriteString(markToml + "\n")
		f.WriteString("title = '" + details.title + "'\n")
		f.WriteString("date = '" + postDate + "'\n")
		f.WriteString("draft = " + details.draft + "\n")
		f.WriteString("tags = [\"" + details.tags + "\"]\n")
		f.WriteString("categories = [\"" + details.categories + "\"]\n")
		f.WriteString("author = \"" + details.author + "\"\n")
		f.WriteString("[[images]]\n")
		f.WriteString("\tsrc = \"img/" + details.images + "\"\n")
		f.WriteString("\talt = \"" + details.altImages + "\"\n")
		f.WriteString("\tstretch = " + details.stretchImages + "\n")
		f.WriteString(markToml + "\n")
		f.WriteString("\n\n\n")
		f.WriteString(details.story + "\n")

		f.Sync()
		f.Close()

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func main() {
	http.HandleFunc("/", hello)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	http.Handle("/jpg/", http.StripPrefix("/jpg/", http.FileServer(http.Dir("jpg"))))

	fmt.Printf("Starting server for testing HTTP POST...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
