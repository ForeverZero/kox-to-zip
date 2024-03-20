package main

import (
	"archive/zip"
	"encoding/xml"
	"io"
)

// 定义结构体映射到XML
type HTML struct {
	XMLName xml.Name `xml:"html"`
	Head    Head     `xml:"head"`
	Body    Body     `xml:"body"`
}

type Head struct {
	Title string `xml:"title"`
	Link  Link   `xml:"link"`
	Meta  Meta   `xml:"meta"`
}

type Link struct {
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
	Href string `xml:"href,attr"`
}

type Meta struct {
	HTTPequiv string `xml:"http-equiv,attr"`
	Content   string `xml:"content,attr"`
}

type Body struct {
	Div ImgDiv `xml:"div"`
}

type ImgDiv struct {
	Div []Div `xml:"div"`
}

type Div struct {
	Class string `xml:"class,attr"`
	Img   []Img  `xml:"img"`
}

type Img struct {
	Src   string `xml:"src,attr"`
	Alt   string `xml:"alt,attr"`
	Class string `xml:"class,attr"`
}

func readImgFromHtml(z *zip.ReadCloser, htmlPath string) ([]string, error) {
	page, err := z.Open(htmlPath)
	if err != nil {
		return nil, err
	}
	defer page.Close()

	bytes, err := io.ReadAll(page)
	if err != nil {
		return nil, err
	}

	var p HTML
	if err := xml.Unmarshal(bytes, &p); err != nil {
		return nil, err
	}
	println(p.Head.Title)

	var imgPaths []string

	divs := p.Body.Div.Div
	for i := range divs {
		imgDiv := divs[i]
		for j := range imgDiv.Img {
			img := imgDiv.Img[j]
			imgPaths = append(imgPaths, img.Src)
		}
	}

	return imgPaths, nil
}
