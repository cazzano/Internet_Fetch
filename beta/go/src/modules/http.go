package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func downloadFile(urlStr string) error {
        u, err := url.Parse(urlStr)
        if err != nil {
                return fmt.Errorf("invalid URL: %v", err)
        }
        path := u.Path

        filename := filepath.Base(path)
        ext := filepath.Ext(filename)

        if filename == "." || filename == "/" {
                filename = "downloaded_file"
                ext = ""
                i := 1
                for {
                        if _, err := os.Stat(filename + ext); os.IsNotExist(err) {
                                break
                        }
                        filename = fmt.Sprintf("downloaded_file_%d", i)
                        i++
                }

        } else {
                if ext == "" {
                        ext = ".dat"
                        i := 1
                        for {
                                if _, err := os.Stat(filename + ext); os.IsNotExist(err) {
                                        break
                                }
                                filename = fmt.Sprintf("%s_%d", filename, i)
                                i++
                        }
                } else {
                        i := 1
                        baseFilename := strings.TrimSuffix(filename, ext) // Store filename without extension
                        for {
                                newFilename := baseFilename + ext // Reconstruct filename with current extension
                                if i > 1 {
                                        newFilename = fmt.Sprintf("%s_%d%s", baseFilename, i, ext)
                                }
                                if _, err := os.Stat(newFilename); os.IsNotExist(err) {
                                        filename = newFilename // Update filename for the loop
                                        break
                                }
                                i++
                        }
                }

        }

        filepath := filepath.Join(".", filename)
        out, err := os.Create(filepath)
        if err != nil {
                return err
        }
        defer out.Close()

        resp, err := http.Get(urlStr)
        if err != nil {
                return err
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
                return fmt.Errorf("bad status: %s", resp.Status)
        }

        _, err = io.Copy(out, resp.Body)
        if err != nil {
                return err
        }

        fmt.Println("Downloaded", urlStr, "to", filepath)
        return nil
}

func main() {
        flag.Parse()

        if flag.NArg() != 1 {
                fmt.Println("Usage: ./http <url>")
                return
        }

        url := flag.Arg(0)

        if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
                fmt.Println("Invalid URL.  Must start with http:// or https://")
                return
        }

        err := downloadFile(url)
        if err != nil {
                fmt.Println("Error:", err)
        }
}
