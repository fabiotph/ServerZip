package main

import (
	"archive/zip"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"runtime"
	"sync"
	"time"
)

func main() {
	PrintMemoryUsage()
	var await sync.WaitGroup
	const count = 20
	fmt.Printf("Fazendo a operação %d vezes\n", count)
	start := time.Now()
	for _, _ = range [count]int{} {
		await.Add(1)
		go WriteZip(&await)
	}

	await.Wait()
	totalTime := time.Since(start)
	fmt.Println("Tempo de execução:" + totalTime.String())
}

func WriteZip(await *sync.WaitGroup) {
	defer await.Done()
	basePathPdf := "pdf/"
	basePathZip := "outFile/"

	if _, err := os.Stat(basePathZip); os.IsNotExist(err) {
		_ = os.Mkdir(basePathZip, 0700)
	}

	outFile, err := os.CreateTemp(basePathZip, "compress*.zip")

	if err != nil {
		fmt.Println(err)
	}

	//defer os.Remove(outFile.Name())

	w := zip.NewWriter(outFile)

	files, err := ioutil.ReadDir(basePathPdf)
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range files {
		AddFile(file, w, basePathPdf, basePathZip)
	}
	PrintMemoryUsage()

	//err = w.Close()
	if err != nil {
		fmt.Println(err)
	}

}

func AddFile(file fs.FileInfo, w *zip.Writer, basePathPdf, basePathZip string) {
	if !file.IsDir() {
		data, _ := ioutil.ReadFile(basePathPdf + file.Name())
		f, _ := w.Create(basePathPdf + file.Name())
		_, _ = f.Write(data)

	}
}

func PrintMemoryUsage() {
	var memory runtime.MemStats
	runtime.ReadMemStats(&memory)
	fmt.Printf("Memory usage: %vMB\n", memory.Alloc/1024/1024)
	fmt.Printf("Total Memory usage: %vMB\n", memory.TotalAlloc/1024/1024)
}
