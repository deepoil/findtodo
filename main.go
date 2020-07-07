package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const TodoFileName = "todoList.txt"

type Todo struct {
	Filepath string
	Files    []string
}

/**
clear todoListFile
*/
func ClearTodoList(outputPath string) {
	file, err := os.OpenFile(outputPath+TodoFileName, os.O_TRUNC|os.O_CREATE, 0666)
	defer file.Close()
	if err != nil {
		log.Fatal("init file error: ", err)
	}
}

/**
get fileNameList
*/
func (todo *Todo) getFiles() *Todo {

	// Execute the ls command in the target directory.
	out, err := exec.Command("find", "-f", todo.Filepath).Output()
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(err)
		log.Fatal(err.Error())
	}

	// Transfer the acquired file list from byte slice to string slice
	var fileList []string
	for _, v := range regexp.MustCompile("\r\n|\n\r|\n|\r").Split(string(out), -1) {
		fileList = append(fileList, v)
	}

	return &Todo{Filepath: todo.Filepath, Files: fileList}

}

/**
Read the file line by line and find the "to do string"
*/
func BufioScanner(fileName string) (todoList []string) {
	f, err := os.Open(fileName)
	defer f.Close()
	if err != nil {
		log.Fatal("bufio error: ", err)
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(strings.ToUpper(strings.Replace(scanner.Text(), " ", "", -1)), "//TODO") {
			todoList = append(todoList, scanner.Text())
		}
	}
	return todoList
}

/**
write file
*/
func WriteTodoList(todoListFile, currentFileName string, todoMessages []string) {
	file, err := os.OpenFile(todoListFile, os.O_WRONLY|os.O_APPEND, 0666)
	defer file.Close()
	if err != nil {
		log.Fatal("write file error: ", err)
	}

	for _, todo := range todoMessages {
		fmt.Fprintln(file, "["+currentFileName+"]\n"+todo)
	}
}

/**
user input validation
*/
func ValidationOfUserInputInfo(path, outputDirFlag string) string {
	// check dir
	fInfo, _ := os.Stat(path)
	if !fInfo.IsDir() {
		log.Fatal("not directory: ", path)
	}
	// check suffix
	newPath := path
	// output file
	if outputDirFlag == "1" {
		if path[len(path)-1:] != "/" {
			newPath += "/"
		}
		// input file
	} else {
		if path[len(path)-1:] == "/" {
			newPath = path[:len(path)-1]
		}
	}
	return newPath
}

func main() {
	var inputPath string
	var outputPath string
	fmt.Println("####TODOリストをファイルに書き出します####")

	// user output file path
	fmt.Println("####書き出し先を入力してください####")
	fmt.Scan(&outputPath)

	// user input file path
	fmt.Println("####読み込み先ディレクトリを入力してください####")
	fmt.Scan(&inputPath)

	// validation
	outputPath = ValidationOfUserInputInfo(outputPath, "1")
	inputPath = ValidationOfUserInputInfo(inputPath, "0")
	v := &Todo{Filepath: inputPath}

	fmt.Println("####", inputPath, "以下のTODOを書き出します####")
	//clear or create to_do_list file
	ClearTodoList(outputPath)

	// get fileNameList
	todo := v.getFiles()

	// find to do string and write file
	for _, file := range todo.Files {
		fInfo, _ := os.Stat(file)

		// For directories,continue
		if fInfo == nil || fInfo.IsDir() {
			continue
		}

		todoMessages := BufioScanner(file)
		WriteTodoList(outputPath+TodoFileName, file, todoMessages)
	}

	fmt.Println(outputPath+TodoFileName, " へ書き出しました")
}
