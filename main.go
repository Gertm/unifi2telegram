package main

import (
	"bytes"
	"log"
	"strings"

	// "encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"regexp"

	"github.com/hpcloud/tail"
	"gopkg.in/ini.v1"
)

var cameraRegexp *regexp.Regexp
var cameras map[string]string
var botToken string
var chatId string

func main() {
	log.Println("Loading configuration file...")
	cfg, err := ini.Load("/etc/unifi2telegram.conf")
	if err != nil {
		log.Printf("Fail to read configuration file /etc/unifi2telegram.conf: %v", err)
		os.Exit(1)
	}
	cameras = cfg.Section("cameras").KeysHash()
	botToken = cfg.Section("telegram").Key("bottoken").String()
	chatId = cfg.Section("telegram").Key("chatid").String()

	cameraRegexp, _ = regexp.Compile("Camera\\[[a-zA-Z0-9]*\\|(\\w*)\\]")
	if os.Args != nil && os.Args[1] != "" {
		followLogFile(os.Args[1])
		return
	}
	fmt.Println("You need to specify the log file to follow.")
}

func sendPhoto(photo string) {
	url := fmt.Sprintf("https://api.telegram.org/%s/sendPhoto?chat_id=%s", botToken, chatId)

	b, w := createMultipartFormData("photo", photo)

	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	//defer resp.Body.Close()
	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)
	// body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println("response Body:", string(body))
}

func createMultipartFormData(fieldName, fileName string) (bytes.Buffer, *multipart.Writer) {
	var b bytes.Buffer
	var err error
	w := multipart.NewWriter(&b)
	var fw io.Writer
	file := mustOpen(fileName)
	if fw, err = w.CreateFormFile(fieldName, file.Name()); err != nil {
		fmt.Printf("Error creating writer: %v\n", err)
	}
	if _, err = io.Copy(fw, file); err != nil {
		fmt.Printf("Error with io.Copy: %v\n", err)
	}
	w.Close()
	return b, w
}

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		pwd, _ := os.Getwd()
		fmt.Println("PWD: ", pwd)
		panic(err)
	}
	return r
}

// DownloadFile ...
func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func followLogFile(filename string) {
	fi, err := os.Stat(filename)
	fileSize := fi.Size()
	seekInfo := tail.SeekInfo{Offset: fileSize, Whence: 0}
	// log.Printf("File size is: %d\n", fileSize)
	if err != nil {
		log.Fatal(err)
	}

	t, err := tail.TailFile(filename, tail.Config{Follow: true, ReOpen: true, Location: &seekInfo})
	if err != nil {
		panic(err)
	}
	for line := range t.Lines {
		if strings.Contains(line.Text, "MOTION STARTED") {
			m := cameraRegexp.FindStringSubmatch(line.Text)[1]
			log.Printf("Got motion for: %s\n", m)
			url := fmt.Sprintf("http://%s/snap.jpeg", cameras[m])
			err := DownloadFile("snap.jpeg", url)
			if err != nil {
				panic(err)
			}
			sendPhoto("snap.jpeg")
		}
	}
}
