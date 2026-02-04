package main

import (
	"os"
	"fmt"
	"io"
	"errors"
	"path/filepath"
)

func main(){
	if len(os.Args) <= 1{
		fmt.Println("bcp file1 Dir/")
		fmt.Println("bcp file1 file2 ... fileN Dir/")
		os.Exit(1)
	}

	lenArgs := len(os.Args)
	files := os.Args[1:lenArgs]
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
			fmt.Printf("%s ", file)
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

	var copied int64
	var copiedCounter int64

	for {
		blockSize := calcBlock(fileInfo.Size(), copied)
		copied, copyStatus := io.CopyN(dstFd, srcFd, blockSize)
		copiedCounter += copied
		progress := (float64(copiedCounter) / float64(fileInfo.Size())) * 100
		fmt.Printf("%.2f%% %s %s\r", progress, file, path)

		if copyStatus == io.EOF{
			fmt.Println()
			return nil
		}	
	}
}
