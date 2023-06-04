package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {
	ipsArg := flag.String("ips", "", "IPs (string or filename)")
	hostsArg := flag.String("hosts", "", "Hosts (string or filename)")
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

	client := http.Client{
		Timeout: 1 * time.Second,
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

				url := fmt.Sprintf("http://%s", ip)
				req, err := http.NewRequest("GET", url, nil)
				if err != nil {
					if !*silent {
						fmt.Printf("Error creating request for IP %s and host %s: %s\n", ip, host, err)
					}
					return
				}
				req.Host = host

				resp, err := client.Do(req)
				if err != nil {
					if !*silent {
						fmt.Printf("Error making request for IP %s and host %s: %s\n", ip, host, err)
					}
					return
				}
				defer resp.Body.Close()

				if *silent && resp.StatusCode != http.StatusOK {
					return
				}

				if resp.StatusCode == http.StatusOK {
					fmt.Printf("%s\t%s\n", ip, host)
				}
			}(ip, host)
		}
	}

	wg.Wait()
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
