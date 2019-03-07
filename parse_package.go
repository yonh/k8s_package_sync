package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

//func main() {
//	str := read_file("Packages")
//
//	parse_package(str)
//}

func parse_package(s string) [][2]string {
	r, _ := regexp.Compile("Filename: (.*)[\\s\\S]*?SHA256: (.*)")

	arr := r.FindAllStringSubmatch(s, -1)

	array := make([][2]string, len(arr))

	for i, s := range arr {
		array[i][0] = s[1]
		array[i][1] = s[2]
		//fmt.Printf("index:%d\n",i)
		//fmt.Printf("%v\n", )read file fail
		//fmt.Printf("%v\n\n", s[2])
	}
	//fmt.Printf("%v", array)

	return array
}


func read_file(file string) string {
	data, err := ioutil.ReadFile(file)

	if err != nil{
		fmt.Println("read file fail")
		os.Exit(1)
	}

	str := fmt.Sprintf("%s", data)

	return str
}
