package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

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
	fmt.Println(MetaInfo(domain))

	for _, v := range servers.Endpoints {
		fmt.Println(v)
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

func MetaInfo(domain string) string {
	response, err := http.Get(fmt.Sprintf("https://%s", domain))
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Get the response body as a string
	data, err := ioutil.ReadAll(response.Body)
	body := string(data)

	re := regexp.MustCompile("(?:<title.*?>)(.*)(?:</title>)")
	title := re.FindAllStringSubmatch(body, -1)

	if title == nil {
		log.Fatal("Page has no title tag")
	}

	return title[0][1]
}
