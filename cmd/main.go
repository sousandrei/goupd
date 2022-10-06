package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const ROOT_URL = "https://go.dev"

func main() {
	link, err := getLink()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fileName, err := downloadFile(link)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := removeOldGo(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err = unpackNewGo(fileName); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getLink() (string, error) {
	fmt.Println("Getting link to latest go version")

	res, err := http.Get(ROOT_URL + "/dl")
	if err != nil {
		return "", fmt.Errorf("failed to get %s: %w", ROOT_URL, err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", fmt.Errorf("failed to get %s: status code %d", ROOT_URL, res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to parse html: %w", err)
	}

	div := doc.Find("body > main > article > a:nth-child(10)").First()
	return div.Get(0).Attr[1].Val, nil
}

func downloadFile(path string) (string, error) {
	fileName := strings.Split(path, "/")[2]

	fmt.Println("Downloading", fileName)

	res, err := http.Get(ROOT_URL + path)
	if err != nil {
		return "", fmt.Errorf("failed to get %s: %w", ROOT_URL+path, err)
	}
	defer res.Body.Close()

	out, err := os.Create("/tmp/" + fileName)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to copy: %w", err)
	}

	return fileName, nil
}

func removeOldGo() error {
	fmt.Println("Removing old go version")

	args := []string{"-rf", "/usr/local/go"}
	cmd := exec.Command("rm", args...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to remove old go version: %w", err)
	}

	return nil
}

func unpackNewGo(fileName string) error {
	fmt.Println("Unpacking new go version")

	args := []string{"-C", "/usr/local", "-xzf", fileName}
	cmd := exec.Command("tar", args...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to remove old go version: %w", err)
	}

	return nil
}
