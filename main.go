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
	out, err := exec.Command("ls", "-1", todo.Filepath).Output()
	if err != nil {
		log.Fatal(err)
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
func ValidationOfUserInputInfo(path string) string {
	// check dir
	fInfo, _ := os.Stat(path)
	if !fInfo.IsDir() {
		log.Fatal("not directory: ", path)
	}
	// check suffix
	if path[len(path)-1:] != "/" {
		path += "/"
	}
	return path
}

func main() {
	var path string
	var outputPath string
	fmt.Println("####TODOリストをファイルに書き出します####")
	fmt.Println("####書き出し先を入力してください####")
	// user output file path
	fmt.Scan(&outputPath)
	// user input file path
	fmt.Println("####読み込み先ディレクトリを入力してください####")
	fmt.Scan(&path)

	// validation
	outputPath = ValidationOfUserInputInfo(outputPath)
	path = ValidationOfUserInputInfo(path)

	v := &Todo{Filepath: path}

	fmt.Println("####", path, "以下のTODOを書き出します####")

	//clear or create to_do_list file
	ClearTodoList(outputPath)

	// get fileNameList
	todo := v.getFiles()

	// find to do string and write file
	for _, v := range todo.Files {
		currentFileName := path + v
		todoMessages := BufioScanner(currentFileName)
		WriteTodoList(outputPath+TodoFileName, currentFileName, todoMessages)
	}

	fmt.Println(outputPath+TodoFileName, " へ書き出しました")
}
