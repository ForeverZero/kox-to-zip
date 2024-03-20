package main

import (
	"encoding/xml"
	"fmt"
	"io/fs"
	"io/ioutil"
)

type Package struct {
	XMLName  xml.Name `xml:"package"`
	Version  string   `xml:"version,attr"`
	Metadata Metadata `xml:"metadata"`
	Manifest Manifest `xml:"manifest"`
	Spine    Spine    `xml:"spine"`
	Guide    Guide    `xml:"guide"`
}

type Metadata struct {
	BookType    string `xml:"meta>book-type"`
	Orientation string `xml:"meta>orientation-lock"`
	Resolution  string `xml:"meta>original-resolution"`
	Title       string `xml:"title"`
	Language    string `xml:"language"`
	Creator     string `xml:"creator"`
	Publisher   string `xml:"publisher"`
	Date        string `xml:"date"`
	Rights      string `xml:"rights"`
	Series      string `xml:"series"`
	Cover       string `xml:"meta>cover"`
}

type Manifest struct {
	Items []Item `xml:"item"`
}

type Item struct {
	ID        string `xml:"id,attr"`
	Href      string `xml:"href,attr"`
	MediaType string `xml:"media-type,attr"`
}

type Spine struct {
	Toc      string    `xml:"toc,attr"`
	ItemRefs []ItemRef `xml:"itemref"`
}

type ItemRef struct {
	IDRef string `xml:"idref,attr"`
}

type Guide struct {
	References []Reference `xml:"reference"`
}

type Reference struct {
	Type  string `xml:"type,attr"`
	Href  string `xml:"href,attr"`
	Title string `xml:"title,attr"`
}

func readOpf(opfFile fs.File) (*Package, error) {
	// 读取文件内容
	xmlData, _ := ioutil.ReadAll(opfFile)

	// 解析XML
	var pkg Package
	if err := xml.Unmarshal(xmlData, &pkg); err != nil {
		fmt.Println(err)
		return nil, err
	}

	// 打印解析结果
	return &pkg, nil
}
