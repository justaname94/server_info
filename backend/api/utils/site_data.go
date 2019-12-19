package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"../models"

	"golang.org/x/net/html"
)

var baseDir, _ = filepath.Abs(filepath.Dir(os.Args[0]))

const (
	// InProgressMsg means the site analysis is not completed
	InProgressMsg = "IN_PROGRESS"
	// DNSMsg means that either is still resolving the name or its an invalid URL
	DNSMsg = "DNS"
	// ReadyMsg means that all the checks have been completed
	ReadyMsg = "READY"
)

const (
	whoIsURL = "https://www.whoisxmlapi.com/whoisserver/WhoisService?" +
		"apiKey=at_jUM8u6lTS16Q6TJFJnxwiYw6H2W9l&domainName=%s&outputFormat=json"
	apiInfoURL = "https://api.ssllabs.com/api/v3/analyze?host="
)

// ServersInfo Represent the server info result from the SSLabs API
type ServersInfo struct {
	Status    string `json:"status"`
	Endpoints []struct {
		IPAddress string `json:"ipAddress"`
		Grade     string `json:"grade"`
	} `json:"endpoints"`
}

type whoisData struct {
	WhoisRecord struct {
		RegistryData struct {
			Registrant struct {
				Country      string `json:"country"`
				Organization string `json:"organization"`
			} `json:"registrant"`
		} `json:"registryData"`
	} `json:"WhoisRecord"`
}

func GetWebsiteData(domain string) (models.Site, string) {
	site := models.Site{}
	body, err := GetPageBody(domain)
	if err != nil {
		site.IsDown = true
		return site, ""
	}
	servers, metadata := getServerData(domain)

	/*
	 A DNS response msg on a working domain its just sometimes the first
	 response of the API before the InProgressMsg as is still resolving the ip
	 addresses, but afterwards it just gives the InProgressMsg
	*/
	if metadata["status"] != ReadyMsg {
		return models.Site{}, InProgressMsg
	}

	site.Grade = metadata["lowestGrade"]
	title, _ := titleMetaInfo(body)
	logo, logoErr := logoMetaInfo(body, domain)

	site.Title = title
	site.Domain = domain

	for _, s := range servers {
		site.Servers = append(site.Servers, s)
	}

	// site.Grade = lowestGrade
	if logoErr != nil {
		site.Logo = ""
	}
	site.Logo = logo

	return site, ReadyMsg
}

func getServerData(domain string) (map[string]models.Server, map[string]string) {
	servers := map[string]models.Server{}
	metadata := map[string]string{}
	endpoints, _ := apiInfo(domain)
	metadata["status"] = endpoints.Status
	if endpoints.Status != ReadyMsg {
		return servers, metadata
	}

	var lowestGrade string

	for i, v := range endpoints.Endpoints {
		whois, _ := whoisInfo(v.IPAddress)
		server := models.Server{
			Address: v.IPAddress,
			Grade:   v.Grade,
			Country: whois.WhoisRecord.RegistryData.Registrant.Country,
			Owner:   whois.WhoisRecord.RegistryData.Registrant.Organization,
		}

		servers[server.Address] = server
		if i > 1 {
			if isGradeGreater(lowestGrade, v.Grade) {
				lowestGrade = v.Grade
			}
		} else {
			lowestGrade = v.Grade
		}
	}
	metadata["lowestGrade"] = lowestGrade
	return servers, metadata
}

func HasServersUpdated(domain string, servers []models.Server) bool {
	newServers, _ := getServerData(domain)

	if len(newServers) != len(servers) {
		return true
	}

	for _, s := range servers {
		_, present := newServers[s.Address]
		if !present {
			return true
		}
	}

	return false
}

// TODO: Refactor apiInfo and whoisInfo
func apiInfo(domain string) (ServersInfo, error) {
	url := fmt.Sprint(apiInfoURL, domain)
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

func whoisInfo(address string) (whoisData, error) {
	url := fmt.Sprintf(whoIsURL, address)
	client := http.Client{
		Timeout: time.Second * 2,
	}

	whois := whoisData{}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return whoisData{}, err
	}

	res, err := client.Do(req)
	if err != nil {
		return whoisData{}, err
	}
	if err := json.NewDecoder(res.Body).Decode(&whois); err != nil {
		return whoisData{}, err
	}

	return whois, nil
}

func titleMetaInfo(body string) (string, error) {
	re := regexp.MustCompile("(?:<title.*?>)(.*)(?:</title>)")
	title := re.FindAllStringSubmatch(body, -1)

	if title == nil {
		return "", errors.New("Page has no title")
	}

	return title[0][1], nil
}

//  Works 8/10 times as some pages like amazon.com) block their access to
//  crawlers or save their favicon in non-conventional ways
func logoMetaInfo(body string, domain string) (string, error) {
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

	logoPath := filepath.Join(baseDir, "static", domain)
	if link != "" {
		// Fetch image location depending if its a relative path or not
		if link[0] == '/' {
			faviconPath := fmt.Sprintf("https://%s%s", domain, link)
			err := DownloadImage(logoPath, faviconPath)
			if err != nil {
				return "", err
			}
		} else {
			err := DownloadImage(logoPath, link)
			if err != nil {
				return "", err
			}
		}
	} else {
		// Some pages have the logo at the relative path {{host}}/favicon.ico, so
		// if the page is not found, try to look at this path
		faviconURL := fmt.Sprintf("https://%s/favicon.ico", domain)
		err := DownloadImage(logoPath, faviconURL)
		if err != nil {
			return "", errors.New("No image found")
		}
	}

	return fmt.Sprintf("/static/%s/favicon.ico", domain), nil
}
