package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	ipsArg := flag.String("ips", "", "IPs (string or filename)")
	hostsArg := flag.String("hosts", "", "Hosts (string or filename)")
	curlArg := flag.Bool("curl", false, "Output as Curl command")
	silent := flag.Bool("silent", false, "Silent mode")
	flag.Parse()

	ips, err := readLines(*ipsArg)
	if err != nil {
		fmt.Println("Error reading IP addresses:", err)
		return
	}

	hosts, err := readLines(*hostsArg)
	if err != nil {
		fmt.Println("Error reading hosts:", err)
		return
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 5) // Limit concurrent calls to 5

	for _, ip := range ips {
		for _, host := range hosts {
			wg.Add(1)
			go func(ip, host string) {
				defer wg.Done()

				semaphore <- struct{}{} // Acquire a slot from the semaphore
				defer func() {
					<-semaphore // Release the slot back to the semaphore
				}()

				// Do HTTP request
				respStatusCode, contentLength, error := doHttpRequest(ip, host, *silent)

				if *silent && (respStatusCode != http.StatusOK || error != nil) {
					return
				}

				if respStatusCode == http.StatusOK {
					if *curlArg {
						fmt.Printf("curl -ik http://%s -H \"Host: %s\"\t(Content-Length: %d)\n", ip, host, contentLength)
					} else {
						fmt.Printf("[%d]\t[%d]\t%s\t%s\n", respStatusCode, contentLength, ip, host)
					}
				}

				// Do HTTPS request
				respStatusCode, contentLength, error = doHttpsRequest(ip, host, *silent)

				if *silent && (respStatusCode != http.StatusOK || error != nil) {
					return
				}

				if respStatusCode == http.StatusOK {
					if *curlArg {
						fmt.Printf("curl -ik http://%s -H \"Host: %s\"\t(Content-Length: %d)\n", ip, host, contentLength)
					} else {
						fmt.Printf("[%d]\t[%d]\t%s\t%s\thttps\n", respStatusCode, contentLength, ip, host)
					}
				}
			}(ip, host)
		}
	}

	wg.Wait()
}

func doHttpRequest(ip string, host string, silent bool) (int, int64, error) {
	return doRequest(ip, host, silent, false)
}

func doHttpsRequest(ip string, host string, silent bool) (int, int64, error) {
	return doRequest(ip, host, silent, true)
}

func doRequest(ip string, host string, silent bool, https bool) (int, int64, error) {
	var scheme string
	if https {
		scheme = "https"
	} else {
		scheme = "http"
	}

	client := http.Client{
		Timeout: 1 * time.Second,
	}

	url := fmt.Sprintf("%s://%s", scheme, ip)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		if !silent {
			fmt.Printf("Error creating request for IP %s and host %s: %s\n", ip, host, err)
		}
		return 0, 0, err
	}
	req.Host = cleanupHost(host)

	resp, err := client.Do(req)
	if err != nil {
		if !silent {
			fmt.Printf("Error making request for IP %s and host %s: %s\n", ip, host, err)
		}
		return 0, 0, err
	}
	defer resp.Body.Close()

	contentLength := resp.ContentLength

	return resp.StatusCode, contentLength, nil
}

func cleanupHost(host string) string {
	host = strings.TrimSpace(host)
	host = strings.TrimSuffix(host, "/")

	u, err := url.Parse(host)
	if err == nil {
		host = u.Host
	}

	return host
}

func readLines(arg string) ([]string, error) {
	lines := []string{}

	if fileExists(arg) {
		file, err := os.Open(arg)
		if err != nil {
			return lines, err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
	} else {
		lines = append(lines, arg)
	}

	return lines, nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}
