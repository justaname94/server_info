package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"../models"

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

func GetWebsiteData(domain string) models.Site {
	servers, _ := aPIInfo(domain)
	whois := whoIsInfo(domain)
	title, _ := titleMetaInfo(domain)
	logo := logoMetaInfo(domain)
	var lowestGrade string

	site := models.Site{
		Title:  title,
		Domain: domain,
	}

	for i, v := range servers.Endpoints {
		server := models.Server{
			Address: v.IPAddress,
			Grade:   v.Grade,
			Country: whois["country"],
			Owner:   whois["owner"],
		}
		site.Servers = append(site.Servers, server)
		if i > 1 {
			if isGradeGreater(lowestGrade, v.Grade) {
				lowestGrade = v.Grade
			}
		} else {
			lowestGrade = v.Grade
		}
	}

	site.Grade = lowestGrade

	if logo != nil {
		site.Logo = "favicon"
	}

	return site
}

func aPIInfo(domain string) (ServersInfo, error) {
	url := fmt.Sprint("https://api.ssllabs.com/api/v3/analyze?host=", domain)
	client := http.Client{
		Timeout: time.Second * 2,
	}
	var servers ServersInfo

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return ServersInfo{}, err
	}

	res, err := client.Do(req)
	if err != nil {
		return ServersInfo{}, err
	}

	if err := json.NewDecoder(res.Body).Decode(&servers); err != nil {
		return ServersInfo{}, err
	}

	return servers, nil
}

func whoIsInfo(domain string) map[string]string {

	info := make(map[string]string)

	command := fmt.Sprintf("whois %s | grep \"Admin Organization\\|Admin Country\" ", domain)
	output, err := exec.Command("bash", "-c", command).Output()
	if err != nil {
		log.Fatal(err)
	}
	outputList := strings.Split(string(output), "\n")

	organization := outputList[0][20:]
	country := outputList[1][15:]

	info["owner"] = organization
	info["country"] = country

	return info
}

func titleMetaInfo(domain string) (string, error) {

	body := GetPageBody(domain)

	re := regexp.MustCompile("(?:<title.*?>)(.*)(?:</title>)")
	title := re.FindAllStringSubmatch(body, -1)

	if title == nil {
		return "", errors.New("Page has no title")
	}

	return title[0][1], nil
}

/*
 Works 8/10 times as some pages like amazon.com) block their access to
 crawlers
*/
func logoMetaInfo(domain string) error {
	body := GetPageBody(domain)
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
			err := DownloadImage(logoPath, faviconPath)
			if err != nil {
				return err
			}
		} else {
			err := DownloadImage(logoPath, link)
			if err != nil {
				return err
			}
		}
	} else {
		return errors.New("No image found")
	}

	return nil
}
