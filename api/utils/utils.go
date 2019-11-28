package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/html"
)

var baseDir, _ = filepath.Abs(filepath.Dir(os.Args[0]))

// Represent the server info result from the SSLabs API
type ServersInfo struct {
	Endpoints []struct {
		IPAddress string `json:"ipAddress"`
		Grade     string `json:"grade"`
	} `json:"endpoints"`
}

var servers ServersInfo

func APIInfo(domain string) (ServersInfo, error) {
	url := fmt.Sprint("https://api.ssllabs.com/api/v3/analyze?host=", domain)
	client := http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if err := json.NewDecoder(res.Body).Decode(&servers); err != nil {
		log.Println(err)
	}

	fmt.Println(WhoIsInfo(domain))
	fmt.Println(TitleMetaInfo(domain))
	fmt.Println(LogoMetaInfo(domain))

	for _, v := range servers.Endpoints {
		fmt.Printf("Ip Address: %s | Grade: %s\n", v.IPAddress, v.Grade)
	}
	return servers, nil
}

func WhoIsInfo(domain string) map[string]string {

	info := make(map[string]string)

	command := fmt.Sprintf("whois %s | grep \"Admin Organization\\|Admin Country\" ", domain)
	output, err := exec.Command("bash", "-c", command).Output()
	if err != nil {
		log.Fatal(err)
	}
	outputList := strings.Split(string(output), "\n")

	organization := outputList[0][20:]
	country := outputList[1][15:]

	info["organization"] = organization
	info["country"] = country

	return info
}

func TitleMetaInfo(domain string) string {

	body := getPageBody(domain)

	re := regexp.MustCompile("(?:<title.*?>)(.*)(?:</title>)")
	title := re.FindAllStringSubmatch(body, -1)

	if title == nil {
		log.Fatal("Page has no title tag")
	}

	return title[0][1]
}

/*
 Works 8/10 times as some pages like amazon.com) block their access to
 crawlers
*/
func LogoMetaInfo(domain string) error {
	body := getPageBody(domain)
	doc, _ := html.Parse(strings.NewReader(body))
	var link string

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "link" {
			for _, l := range n.Attr {

				if l.Key == "rel" {
					if strings.Contains(l.Val, "icon") {
						// Have to loop again the attr slice as the href key is not always
						// on the same position
						for _, l := range n.Attr {
							if l.Key == "href" {
								link = l.Val
							}
						}
						return
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	logoPath := fmt.Sprintf("%s%s%s", baseDir, "/static/", domain)
	if link != "" {
		// Fetch image location depending if its a relative path or not
		if link[0] == '/' {
			faviconPath := fmt.Sprintf("https://%s%s", domain, link)
			fmt.Println(faviconPath)
			err := downloadImage(logoPath, faviconPath)
			if err != nil {
				return err
			}
		} else {
			err := downloadImage(logoPath, link)
			if err != nil {
				return err
			}
		}
	} else {
		return errors.New("No image found")
	}

	return nil
}

func getPageBody(domain string) string {
	response, err := http.Get(fmt.Sprintf("https://%s", domain))
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Get the response body as a string
	data, err := ioutil.ReadAll(response.Body)
	return string(data)
}

func downloadImage(filepath string, url string) error {
	// Create the path
	fmt.Println(filepath)
	os.MkdirAll(filepath, os.ModePerm)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(fmt.Sprintf("%s/%s", filepath, "/favicon.ico"))
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
