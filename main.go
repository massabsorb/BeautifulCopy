package main

import (
	"os"
	"fmt"
	"io"
	"errors"
	"path/filepath"
)

const GREEN = "\033[1;32m"
const BLUE = "\033[1;34m"
const WHITE = "\033[1;37m"
const RESET = "\033[00m"

func main(){
	if len(os.Args) < 3{
		fmt.Printf("%s# %sbcp file1 Directory/\n", BLUE, RESET)
		fmt.Printf("%s# %sbcp file1 file2 ... fileN Directory/\n", BLUE, RESET)
		os.Exit(1)
	}

	lenArgs := len(os.Args)
	files := os.Args[1:lenArgs - 1]
	dstPath := os.Args[lenArgs-1] 

	for counter := 0; counter != len(files); counter++{
		copyFileErr := copyFile(files[counter], dstPath)
		if copyFileErr != nil{
			fmt.Println(copyFileErr)
		}
	}				
}

func calcBlock(filesize int64, copied int64)(int64){
	var defaultBlockSize, diff int64 = 64000, 0 
	diff = filesize - copied	

	switch {
		case diff < defaultBlockSize:
			return diff;
		default:
			return defaultBlockSize
	}
}


func copyFile(file string, path string)(error){

	fileInfo, _ := os.Stat(file)	
	
	switch GetMode := fileInfo.Mode();{
		case GetMode.IsDir():
			fmt.Printf("%s%s%s ", BLUE, file, RESET)
			return errors.New("Directory skipped.")
	}	

	srcFd, srcErr := os.Open(file)
	if srcErr != nil{
		return srcErr
	}

	_, fileName := filepath.Split(file)
	dstFd, dstErr := os.Create(path + fileName)
	if dstErr != nil{
		return dstErr
	}

	var copied int64 = 0
	var copiedCounter int64 = 0

	for {
		blockSize := calcBlock(fileInfo.Size(), copied)
		copied, copyStatus := io.CopyN(dstFd, srcFd, blockSize)
		copiedCounter += copied
		progress := (float64(copiedCounter) / float64(fileInfo.Size())) * 100
		fmt.Printf("%s%.2f%% %s %s%s\r", GREEN, progress, RESET, file, path)

		if copyStatus == io.EOF{
			fmt.Println()
			return nil
		}	
	}
}
