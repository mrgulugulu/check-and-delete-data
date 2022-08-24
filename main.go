package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
)

var (
	// 正则匹配规则
	pattern1 = `.*.csv`
	// 文件目录
	filePath = "/data/tkeclusters-cpu-and-mem/"
)

func main() {
	csvFileList := []string{}
	err := check(filePath, &csvFileList)
	if err != nil {
		fmt.Printf("check files errors: %v", err)
	}
	// data文件大于9个，就说明有过期的数据。因为已经再check里面排序了，所以后9个csv文件肯定是新的。
	// 需要删除前n-9个csv文件
	if len(csvFileList) > 9 {
		toBeDeletedFiles := csvFileList[:len(csvFileList)-9]
		deleteOldFiles(toBeDeletedFiles)
	}
}

// check 用一个数据将csv文件都存起来，然后只保留时间最大的6个文件（也就是最新的文件）
func check(dirname string, csvFileList *[]string) error {
	reg := regexp.MustCompile(pattern1)
	fileInfos, err := ioutil.ReadDir(dirname)
	if err != nil {
		return fmt.Errorf("read file errors: %v", err)
	}
	for _, fi := range fileInfos {
		filename := dirname + "/" + fi.Name()
		if fi.IsDir() {
			//继续遍历fi这个目录
			check(filename, csvFileList)
		} else {
			if reg.MatchString(fi.Name()) {
				*csvFileList = append(*csvFileList, fi.Name())
			}
		}
	}
	sort.Strings(*csvFileList)
	return nil
}

func deleteOldFiles(fileList []string) {
	var err error
	for i := range fileList {
		fileName := filePath + fileList[i]
		err := os.Remove(fileName)
		if err != nil {
			fmt.Printf("remove error: %v, filename: %s", err, fileList[i])
		}
	}
	if err == nil {
		fmt.Println("remove all old files")
	}
}
