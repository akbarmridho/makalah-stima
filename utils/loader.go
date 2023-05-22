package utils

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/yeka/zip"
)

func ReadFile(hash string, password string) []byte {
	resp, err := http.PostForm("https://mb-api.abuse.ch/api/v1/", url.Values{
		"query":       {"get_file"},
		"sha256_hash": {hash}})
	// resp, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Reading file:", zipReader.File[0].Name)
	unzippedFileBytes, err := readZipFile(zipReader.File[0], password)
	if err != nil {
		log.Println(err)
	}

	return unzippedFileBytes
}

func readZipFile(zf *zip.File, password string) ([]byte, error) {
	zf.SetPassword(password)
	f, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}
