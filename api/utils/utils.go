package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func isGradeGreater(gradeA string, gradeB string) bool {
	grades := map[string]int{
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

func GetPageBody(domain string) string {
	response, err := http.Get(fmt.Sprintf("https://%s", domain))
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Get the response body as a string
	data, err := ioutil.ReadAll(response.Body)
	return string(data)
}

func DownloadImage(filepath string, url string) error {
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
