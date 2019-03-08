package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

func main() {
	download("https://packages.cloud.google.com/apt/doc/apt-key.gpg", "mirror/apt/doc/apt-key.gpg")
	download("https://packages.cloud.google.com/apt/dists/kubernetes-xenial/InRelease", "mirror/apt/dists/kubernetes-xenial/InRelease")
	download("https://packages.cloud.google.com/apt/dists/kubernetes-xenial/Release", "mirror/apt/dists/kubernetes-xenial/Release")
	download("https://packages.cloud.google.com/apt/dists/kubernetes-xenial/Release.gpg", "mirror/apt/dists/kubernetes-xenial/Release.gpg")

	read_and_download_packages("mirror/apt/dists/kubernetes-xenial/Release")
	read_and_download_packages("mirror/apt/dists/kubernetes-xenial/InRelease")

	fmt.Printf("[done]\n\n")
}

func read_and_download_packages(file string) {
	base := "https://packages.cloud.google.com/apt/dists/kubernetes-xenial/"

	//resp, err := http.Get("https://packages.cloud.google.com/apt/dists/kubernetes-xenial/InRelease")
	//if err != nil{
	//	fmt.Println("download error")
	//	os.Exit(1)
	//}

	data, err := ioutil.ReadFile(file)

	if err != nil {
		fmt.Println("read file fail")
		os.Exit(1)
	}

	str := fmt.Sprintf("%s", data)

	r, _ := regexp.Compile("(.{64}) (\\d+) (.*)")

	//fmt.Println(r.FindAllString(str, -1))
	arr := r.FindAllString(str, -1)

	for _, s := range arr {
		fmt.Printf("%v\n", s)
	}

	fmt.Println()

	for _, s := range r.FindAllStringSubmatch(str, -1) {
		//fmt.Printf("%s, %s, %s\n", s[1],s[2],s[3])
		url := base + s[3]

		dir := "mirror/apt/dists/kubernetes-xenial/" + filepath.Dir(s[3])
		filename := filepath.Base(s[3])

		if _, err := os.Stat(dir); os.IsNotExist(err) {
			os.MkdirAll(dir, 0755)
		}

		if (!file_exists(dir + filename) || !verify_hash(dir + "/" + filename, s[1])) {
			download(url, dir+"/"+filename)
		}

		// download packages
		if (filepath.Base(filename) == "Packages") {
			packages_content := read_file(dir+"/"+filename)
			packages_arr := parse_package(packages_content)

			for _, p_arr := range packages_arr {
				p_filename := "mirror/apt/" + p_arr[0]

				if (!file_exists(p_filename) || !verify_hash(p_filename,  p_arr[1])) {
					download("https://packages.cloud.google.com/apt/" + p_arr[0], p_filename)
				}
			}
		}
	}
}

func file_exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func verify_hash(file string, hash string) bool {

	file_hash, err := file_hash_sha256(file)
	if (err != nil) {
		fmt.Printf("%v", err)
	}

	return hash == file_hash
}

func file_hash_sha256(filePath string) (string, error) {
	var hashValue string
	file, err := os.Open(filePath)
	if err != nil {
		return hashValue, err
	}
	defer file.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return hashValue, err
	}
	hashInBytes := hash.Sum(nil)
	hashValue = hex.EncodeToString(hashInBytes)
	return hashValue, nil
}

func download(url string, file string) {
	fmt.Printf("download from: %s\n", url)

	resp, err := http.Get(url)

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error : %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// if folder not exists, create it.
	dir := filepath.Dir(file)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}

	write_err := ioutil.WriteFile(file, data, 0644)

	if write_err != nil {
		fmt.Println(write_err)
		//f.Close()
		return
	}
}
