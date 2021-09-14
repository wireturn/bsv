package main

import (
	"fmt"
	"image/png"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

// Page comment.
type Page struct {
	Title string
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/generator/", viewCodeHandler)
	log.Fatal(http.ListenAndServe(":8111", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	p := Page{Title: "QR Code Generator"}
	wd, _ := os.Getwd()
	fmt.Println(wd)
	t, _ := template.ParseFiles("./cmd/qrcode/generator.html")
	_ = t.Execute(w, p)
}

func viewCodeHandler(w http.ResponseWriter, r *http.Request) {
	dataString := r.FormValue("dataString")

	qrCode, _ := qr.Encode(dataString, qr.L, qr.Auto)
	qrCode, _ = barcode.Scale(qrCode, 512, 512)

	_ = png.Encode(w, qrCode)
}
