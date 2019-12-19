package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func isGradeGreater(gradeA string, gradeB string) bool {
	grades := map[string]int{
		"":   15,
		"A+": 1,
		"A":  2,
		"B+": 3,
		"B":  4,
		"C+": 5,
		"C":  6,
		"D+": 7,
		"D":  8,
		"E+": 9,
		"E":  10,
		"F+": 11,
		"F":  12,
		"M":  13,
		"T":  14,
	}

	return grades[gradeA] > grades[gradeB]
}

func GetPageBody(domain string) (string, error) {
	response, err := http.Get(fmt.Sprintf("https://%s", domain))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	// Get the response body as a string
	data, err := ioutil.ReadAll(response.Body)
	return string(data), nil
}

func DownloadImage(filepath string, url string) error {
	// Create the path
	os.MkdirAll(filepath, os.ModePerm)

	savePath := fmt.Sprintf("%s/%s", filepath, "favicon.ico")

	// If image has already been downloaded, do nothing
	if _, err := os.Stat(savePath); !os.IsNotExist(err) {
		return nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("status: %d", http.StatusNotFound)
	}

	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
