package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"howmg.p/utils"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func main() {
	tempFile := "oui.txt"
	jsonFile := "oui.js"
	//读取mac地址表的文件
	resp, err := http.Get("http://standards-oui.ieee.org/oui.txt")
	if err != nil {
		fmt.Println("http get error:", err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("http read error:", err.Error())
		return
	}

	src := string(body)
	//fmt.Println(src)

	//写入临时文件中
	myFile, err := os.Create(tempFile)
	if err != nil {
		fmt.Println("temp file create error:", err.Error())
		return
	}
	myFile.WriteString(src)
	myFile.Close()

	//逐行读取文件
	//\s*[0-9A-F]{2}\-[0-9A-F]{2}\-[0-9A-F]{2}\s*\(hex\)\s*[\w\s]*
	//\s*[0-9A-F]{6}\s*\(base 16\)\s*[\w\s]*
	var items = make([]utils.MacItem, 0)

	f, _ := os.Open(tempFile)
	scanner := bufio.NewScanner(f)
	i := 1
	num := 0
	tempItem := new(utils.MacItem)
	for scanner.Scan() {
		str := scanner.Text()
		if strings.TrimSpace(str) == "" {
			//fmt.Println("make new item...", num)
			if i > 0 && tempItem.Hex != "" { //说明后面读取了操作或者是初始化
				items = append(items, *tempItem)
				//fmt.Printf("mac hex:%s , base 16:%s , company:%s\n", tempItem.Hex, tempItem.Base16, tempItem.Company)
				//fmt.Println("make new item...", num)
				tempItem = new(utils.MacItem)
				num++
				i = 0 //把标记归零
			}
		} else if m, _ := regexp.MatchString("\\s*[0-9A-F]{2}\\-[0-9A-F]{2}\\-[0-9A-F]{2}\\s*\\(hex\\)\\s*[\\w\\s]*", str); m {
			//读取  00-00-00   (hex)		XEROX CORPORATION
			//fmt.Println("read hex[", num, "]->", str)
			re, _ := regexp.Compile("[0-9A-F]{2}\\-[0-9A-F]{2}\\-[0-9A-F]{2}")
			//mac hex
			macHex := re.FindString(str)
			tempItem.Hex = macHex
			// company
			re, _ = regexp.Compile("[0-9A-F]{2}\\-[0-9A-F]{2}\\-[0-9A-F]{2}\\s*\\(hex\\)\\s*")
			company := strings.TrimSpace(re.ReplaceAllString(str, ""))
			tempItem.Company = company
			i++
		} else if m, _ := regexp.MatchString("\\s*[0-9A-F]{6}\\s*\\(base 16\\)\\s*[\\w\\s]*", str); m {
			//读取  000000     (base 16)		XEROX CORPORATION
			//fmt.Println("read base16[", num, "]->", str)
			re, _ := regexp.Compile("[0-9A-F]{6}")
			//mac base 16
			baseHex := re.FindString(str)
			tempItem.Base16 = baseHex
			i++
		}
	}
	//fmt.Println("items length->", len(items))
	fmt.Printf("items[0]:mac hex:%s , base 16:%s , company:%s\n", items[0].Hex, items[0].Base16, items[0].Company)
	fmt.Printf("items[1]:mac hex:%s , base 16:%s , company:%s\n", items[1].Hex, items[1].Base16, items[1].Company)
	fmt.Printf("items[len]:mac hex:%s , base 16:%s , company:%s\n", items[len(items)-1].Hex, items[len(items)-1].Base16, items[len(items)-1].Company)

	myJson, err := json.Marshal(items)
	if err != nil {
		fmt.Println("marshal json error:", err.Error())
		return
	}
	//存json到文件中
	myFile, err = os.Create(jsonFile)
	if err != nil {
		fmt.Println("temp file create error:", err.Error())
		return
	}
	myFile.WriteString(string(myJson))
	myFile.Close()
}
