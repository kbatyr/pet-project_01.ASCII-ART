package server

import (
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"web/ascii_art"
)

type Data struct {
	Font, Input, Button, Art string
	STfs, SHfs, THfs         bool
}
type Error struct {
	Code int
	Msg  string
}

func renderTemplate(w http.ResponseWriter, tmpl string, d *Data) {

	_, err := template.ParseFiles("templates/index.html")
	if err != nil {
		printError(w, &Error{
			Code: http.StatusInternalServerError,
			Msg:  http.StatusText(500),
		})
		return
	}

	var templates = template.Must(template.ParseFiles("templates/index.html"))

	if err := templates.ExecuteTemplate(w, tmpl+".html", d); err != nil {
		printError(w, &Error{
			Code: http.StatusInternalServerError,
			Msg:  http.StatusText(500),
		})
		return
	}
}

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		printError(w, &Error{
			Code: http.StatusMethodNotAllowed,
			Msg:  http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}

	data := &Data{}

	if r.URL.Path != "/" {
		printError(w, &Error{
			Code: http.StatusNotFound,
			Msg:  http.StatusText(http.StatusNotFound),
		})
		return
	}
	renderTemplate(w, "index", data)
}

func AsciiHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		printError(w, &Error{
			Code: http.StatusMethodNotAllowed,
			Msg:  http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}

	if r.URL.Path != "/ascii-art" {
		http.Error(w, http.StatusText(404), http.StatusNotFound)
		return
	}

	data := &Data{}
	defer r.Body.Close()

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		printError(w, &Error{
			Code: http.StatusBadRequest,
			Msg:  http.StatusText(400),
		})
		return
	}
	values, err := url.ParseQuery(string(bytes))
	if err != nil {
		printError(w, &Error{
			Code: http.StatusBadRequest,
			Msg:  http.StatusText(400),
		})
		return
	}

	for i, v := range values {
		switch i {
		case "text":
			data.Input = v[0]
		case "font":
			data.Font = v[0]
		case "button":
			data.Button = v[0]
		default:
			printError(w, &Error{
				Code: http.StatusBadRequest,
				Msg:  http.StatusText(400),
			})
			return
		}
	}
	if data.Input == "" || len(values["text"]) > 1 || len(values["font"]) > 1 || len(values["button"]) > 1 {
		printError(w, &Error{
			Code: http.StatusBadRequest,
			Msg:  http.StatusText(400),
		})
		return
	}

	switch data.Font {
	case "standard.txt":
		data.STfs = true
	case "shadow.txt":
		data.SHfs = true
	case "thinkertoy.txt":
		data.THfs = true
	default:
		printError(w, &Error{
			Code: http.StatusBadRequest,
			Msg:  http.StatusText(400),
		})
		return
	}

	data.Art, err = ascii_art.AsciiArt(data.Input, data.Font)

	if err != nil {
		if err.Error() == "400, Bad Request" {
			printError(w, &Error{
				Code: http.StatusBadRequest,
				Msg:  http.StatusText(400),
			})
		} else {
			printError(w, &Error{
				Code: http.StatusInternalServerError,
				Msg:  http.StatusText(500),
			})
		}
		return
	}

	if data.Button == "Download" {
		sendFileToClient(w, r, data.Art)
		return
	}
	renderTemplate(w, "index", data)
}

func printError(w http.ResponseWriter, errMsg *Error) {
	tmpl, err := template.ParseFiles("templates/err.html")

	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(errMsg.Code)
	tmpl.Execute(w, errMsg)
}

func sendFileToClient(w http.ResponseWriter, r *http.Request, res string) {

	file := strings.NewReader(res)
	fileSize := strconv.Itoa(int(file.Size()))

	w.Header().Set("Content-Disposition", "attachment; filename=output.txt")
	w.Header().Set("Content-Type", "plain/text; charset=utf-8")
	w.Header().Set("Content-Length", fileSize)
	io.Copy(w, file)
}
