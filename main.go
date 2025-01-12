package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
)

type JSON struct {
	Dependencies map[string]string `json:"dependencies"`
}

func getMainJsonFileContent() ([]byte, string) {
	fmt.Println("Welcome to our cli app! 🎉")
	var fileUrl string
	fmt.Println("Enter the URL of package.json of your project 🔗(e.g. D:/.../package.json): ")
	fmt.Scanln(&fileUrl)

	for fileUrl == "" || fileUrl == " " || strings.LastIndex(fileUrl, "package.json") == -1 {
		color.Red("Incorrect URL 🙄. Please try again.")
		fmt.Scanln(&fileUrl)
	}

	if strings.HasPrefix(fileUrl, "\"") && strings.HasSuffix(fileUrl, "\"") {
		fileUrl = strings.Trim(fileUrl, "\"")
	}

	content, err := ioutil.ReadFile(fileUrl)
	if err != nil {
		log.Fatal(err)
	}

	return content, fileUrl
}

func getUsedFiles(url string, usedFiles *[]string) {
	if url == "" || strings.HasSuffix(url, "node_modules") {
		return
	}
	files, err := ioutil.ReadDir(url)
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			subDirUrl := url + "/" + file.Name()
			getUsedFiles(subDirUrl, usedFiles)
		} else if strings.HasSuffix(file.Name(), ".js") || strings.HasSuffix(file.Name(), ".jsx") || strings.HasSuffix(file.Name(), ".ts") || strings.HasSuffix(file.Name(), ".tsx") {
			(*usedFiles) = append((*usedFiles), url+"/"+file.Name())
		}

	}
}

func main() {
	content, fileUrl := getMainJsonFileContent()

	var jsonData JSON
	err := json.Unmarshal(content, &jsonData)
	if err != nil {
		fmt.Println(err)
	}

	if jsonData.Dependencies == nil {
		color.Yellow("No dependencies found in package.json 🙅‍♂️")
		return
	}

	dirPath := strings.Split(fileUrl, "package.json")[0]

	var usedFiles []string
	getUsedFiles(dirPath, &usedFiles)

	fileInDepondency := make(map[string][]string)

	setUsedDep := make(map[string]bool)

	for _, file := range usedFiles {
		fileContent, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}

		var dependenceList []string
		for dependence := range jsonData.Dependencies {
			if strings.Contains(string(fileContent), dependence) {
				dependenceList = append(dependenceList, dependence)
				setUsedDep[dependence] = true
			}
		}

		fileInDepondency[file] = dependenceList

	}

	fmt.Println("Used dependencies: ")
	for key := range setUsedDep {
		color.Green("%s", key)
	}

	fmt.Println("\nUnused dependencies: ")
	for key := range jsonData.Dependencies {
		if !setUsedDep[key] {
			color.Red("%s", key)
		}
	}

	fmt.Println("Is more result write in file 🤔 ? (y/n)")
	var answer string
	fmt.Scanln(&answer)
	if answer == "y" {
		f, err := os.Create(dirPath + "result.txt")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		for file, depondencies := range fileInDepondency {
			f.WriteString(file + ":\n")
			for _, depondence := range depondencies {
				f.WriteString("\t-> " + depondence + "\n")
			}
		}
		fmt.Println("Result saved in file 💾: ", dirPath+"result.txt")
	}
	fmt.Println("Thanks for using! :)")
}
