package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/emersion/go-vcard"
	"github.com/skip2/go-qrcode"
)

var version = "dev"

//go:embed templates/*
var content embed.FS

type Page struct {
	Title   string
	Version string
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/contactCardGeneratorHandler/", contactCardGeneratorHandler)
	mux.HandleFunc("/urlGeneratorHandler/", urlGeneratorHandler)

	log.Printf("Server (version: %v) started, listening on :8080", version)
	http.ListenAndServe(":8080", mux)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request")
	p := Page{Title: "QR Code Generator", Version: version}

	t, _ := template.ParseFS(content, "templates/generator.html")
	t.Execute(w, p)
}

func contactCardGeneratorHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalln(err)
	}

	vc := vcard.Card{}
	vc.AddValue(vcard.FieldVersion, "3.0")
	vc.AddName(&vcard.Name{
		FamilyName: r.FormValue("lastname"),
		GivenName:  r.FormValue("firstname"),
	})
	vc.AddAddress(&vcard.Address{
		Locality:      r.FormValue("locality"),
		StreetAddress: r.FormValue("street"),
		PostalCode:    r.FormValue("postcode"),
	})
	vc.AddValue(vcard.FieldOrganization, r.FormValue("organization"))
	vc.AddValue(vcard.FieldTitle, r.FormValue("title"))
	vc.AddValue(vcard.FieldEmail, r.FormValue("email"))
	vc.AddValue(vcard.FieldTelephone, r.FormValue("phone"))

	encodedVCard := new(strings.Builder)
	enc := vcard.NewEncoder(encodedVCard)
	err = enc.Encode(vc)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Generating QR-Code for:", encodedVCard.String())

	qr, err := qrcode.Encode(encodedVCard.String(), qrcode.Medium, 256)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "image/png")
	_, err = w.Write(qr)
	if err != nil {
		log.Fatal(err)
	}
}

func urlGeneratorHandler(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	log.Println("Generating QR-Code for:", url)
	err := r.ParseForm()
	if err != nil {
		log.Fatalln(err)
	}
	qr, err := qrcode.Encode(r.FormValue("url"), qrcode.Medium, 256)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "image/png")
	_, err = w.Write(qr)
	if err != nil {
		log.Fatal(err)
	}
}
