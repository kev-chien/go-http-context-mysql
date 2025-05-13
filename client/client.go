package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("No arguments provided")
		return
	}

	endpoint := args[0]

	c := &http.Client{
		Timeout: 1 * time.Second,
	}
	resp, err := c.Get("http://localhost:8090/" + endpoint)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error making request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading response:", err)
	}
}
