package server

import (
	"bufio"
	"fe-tool/common"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const TargetDir = "."

func Download() {
	client := http.Client{}

	// 检查是否需要清理缓存
	resp, err := client.Head("https://codeload.github.com/msojocs/fiddler-everywhere-enhance/zip/refs/heads/v8.x")
	if err != nil {
		log.Fatalln("Check server error:" + err.Error())
	}
	defer resp.Body.Close()
	etag := resp.Header.Get("ETag")
	etag = etag[1 : len(etag)-1] // 去掉双引号
	log.Println("ETag:", etag)
	repoFile := "cache/server-" + etag + ".zip"
	if s, err := os.Stat(repoFile); err == nil && !s.IsDir() {
		log.Println(repoFile, "exists.")
		return
	}

	{
		// 清理旧的缓存文件，server*.zip
		files, err := os.ReadDir("cache")
		if err != nil {
			log.Fatalln("Read cache dir error:" + err.Error())
		}
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if matched, _ := filepath.Match("server*.zip", file.Name()); matched {
				os.Remove("cache/" + file.Name())
			}
		}
	}

	file, err := os.Create(repoFile + ".tmp")
	if err != nil {
		log.Fatalln("Create file error:" + err.Error())
	}

	writer := bufio.NewWriter(file)
	resp, err = client.Get("https://github.com/msojocs/fiddler-everywhere-enhance/archive/refs/heads/v8.x.zip")
	if err != nil {
		file.Close()
		log.Fatalln("Download server error:" + err.Error())
	}
	defer resp.Body.Close()

	fileSize, err := io.Copy(writer, resp.Body)
	file.Close()
	if err != nil {
		log.Fatalln("Write file error:" + err.Error())
	}

	err = os.Rename(repoFile+".tmp", repoFile)
	if err != nil {
		log.Fatalln("Rename server.zip.tmp error", err)
	}
	log.Println("Download server.zip end, file size:", fileSize)
}

func Extract() {
	log.Println("Extract server.zip start...")
	client := http.Client{}

	// 检查是否需要清理缓存
	resp, err := client.Head("https://codeload.github.com/msojocs/fiddler-everywhere-enhance/zip/refs/heads/v8.x")
	if err != nil {
		log.Fatalln("Check server error:" + err.Error())
	}
	defer resp.Body.Close()
	etag := resp.Header.Get("ETag")
	etag = etag[1 : len(etag)-1] // 去掉双引号
	log.Println("ETag:", etag)
	repoFile := "cache/server-" + etag + ".zip"

	err = common.ExtractZipArchive(repoFile, TargetDir)
	if err != nil {
		log.Fatalln("Extract error:", err)
	}
}
