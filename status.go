package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
)

const (
	ColorRed      = "\033[91m"
	ColorOrange   = "\033[38;5;208m"
	ColorLime     = "\033[92m"
	ColorSkyBlue  = "\033[38;5;39m" // Define sky blue color
	ColorReset    = "\033[0m"
)

func getStatusCodeColor(statusCode int) string {
	switch {
	case statusCode >= 200 && statusCode <= 290:
		return ColorLime
	case statusCode >= 400 && statusCode <= 490:
		return ColorRed
	case statusCode >= 300 && statusCode <= 390:
		return ColorOrange
	case statusCode >= 500 && statusCode <= 590:
		return ColorSkyBlue // Set sky blue color for 500-590
	default:
		return ""
	}
}

func checkSubdomainsStatus(subdomains []string) map[string]int {
	statusCodes := make(map[string]int)
	var wg sync.WaitGroup

	for _, subdomain := range subdomains {
		wg.Add(1)
		go func(sd string) {
			defer wg.Done()
			url := "http://" + sd // You can use "https://" as well if needed

			resp, err := http.Head(url)
			if err != nil {
				statusCodes[url] = http.StatusNotFound
				return
			}
			defer resp.Body.Close()

			statusCodes[url] = resp.StatusCode
		}(subdomain)
	}

	wg.Wait()
	return statusCodes
}

func readSubdomainsFromFile(filename string) ([]string, error) {
	var subdomains []string

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		subdomains = append(subdomains, strings.TrimSpace(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return subdomains, nil
}

func main() {
	// Replace with your subdomains file path
	subdomainsFile := "subdomains.txt"

	subdomains, err := readSubdomainsFromFile(subdomainsFile)
	if err != nil {
		fmt.Println("Error reading subdomains file:", err)
		return
	}

	statusCodes := checkSubdomainsStatus(subdomains)

	// Separate subdomains by status code
	status200 := make(map[string]int)
	status300 := make(map[string]int)
	status400 := make(map[string]int)
	status500 := make(map[string]int) // Map for 500-590 status codes

	for subdomain, status := range statusCodes {
		switch {
		case status >= 200 && status <= 290:
			status200[subdomain] = status
		case status >= 300 && status <= 390:
			status300[subdomain] = status
		case status >= 400 && status <= 490:
			status400[subdomain] = status
		case status >= 500 && status <= 590:
			status500[subdomain] = status
		}
	}

	// Print status codes for each subdomain with colors
	printStatusCodes := func(statusMap map[string]int, color string) {
		for subdomain, status := range statusMap {
			resetColor := ColorReset

			switch color {
			case ColorRed, ColorOrange, ColorLime, ColorSkyBlue:
				fmt.Printf("%s%s : %d%s\n", color, subdomain, status, resetColor)
			}
		}
	}

	fmt.Println("Status 200:")
	printStatusCodes(status200, ColorLime)

	fmt.Println("\nStatus 300:")
	printStatusCodes(status300, ColorOrange)

	fmt.Println("\nStatus 400:")
	printStatusCodes(status400, ColorRed)

	fmt.Println("\nStatus 500:") // Print 500-590 status codes
	printStatusCodes(status500, ColorSkyBlue)
}


