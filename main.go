package main

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	args := os.Args
	for i := range args {
		if i == 0 {
			continue
		}
		path := args[i]
		println("转换格式: ", path)
		if err := convert(path); err != nil {
			println("发生异常: ", err.Error())
		}
	}
}

func convert(path string) error {
	z, err := zip.OpenReader(path)
	if err != nil {
		return err
	}
	defer z.Close()

	opfFile, err := z.Open("vol.opf")
	if err != nil {
		return err
	}
	defer opfFile.Close()

	opf, err := readOpf(opfFile)
	if err != nil {
		return err
	}

	outputPath := filepath.Dir(path) + "/" + fmt.Sprintf("%s - %s.zip", opf.Metadata.Creator, opf.Metadata.Title)
	println("输出: ", outputPath)
	if _, err := os.Stat(outputPath); err == nil {
		return errors.New("文件已存在")
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	zipWriter := zip.NewWriter(outputFile)
	defer zipWriter.Close()

	itemMap := make(map[string]string)
	for i := range opf.Manifest.Items {
		item := opf.Manifest.Items[i]
		itemMap[item.ID] = item.Href
	}

	imgSeq := 0
	for i := range opf.Spine.ItemRefs {
		ref := opf.Spine.ItemRefs[i]
		id := ref.IDRef
		htmlPath := itemMap[id]
		println(id, htmlPath)
		imgs, err := readImgFromHtml(z, htmlPath)
		if err != nil {
			return err
		}

		for idxImg := range imgs {
			println(imgs[idxImg])
			imgFile, err := getImgFile(z, imgs[idxImg])
			if err != nil {
				return err
			}
			stat, err := imgFile.Stat()
			imgName := stat.Name()
			println(imgName)
			idx := strings.LastIndex(imgName, ".")
			if idx < 0 {
				return errors.New(fmt.Sprintf("处理文件名异常: %s", imgName))
			}
			suffix := imgName[idx:]

			// 写zip
			header, err := zip.FileInfoHeader(stat)
			if err != nil {
				return err
			}

			// 计数
			header.Name = fmt.Sprintf("%d%s", imgSeq, suffix)
			imgSeq += 1
			w, err := zipWriter.CreateHeader(header)
			if err != nil {
				return err
			}

			if _, err := io.Copy(w, imgFile); err != nil {
				return err
			}

			imgFile.Close()
		}
	}
	return nil
}

func getImgFile(z *zip.ReadCloser, imgPath string) (fs.File, error) {
	if strings.Index(imgPath, "/") == 0 || strings.Index(imgPath, "\\") == 0 {
		return z.Open(imgPath)
	}

	idx := strings.LastIndex(imgPath, "/")
	if idx < 0 {
		idx = strings.LastIndex(imgPath, "\\")
	}
	if idx < 0 {
		return nil, errors.New("处理imgPath异常")
	}

	return z.Open("image" + imgPath[idx:])
}
