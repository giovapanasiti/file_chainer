package files

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	//"io/ioutil"
	// . "github.com/logrusorgru/aurora"
	"github.com/manifoldco/promptui"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func ConcatenaFiles() {
	prompt := promptui.Prompt{
		Label: "Cartella dei file",
	}
	c, _ := prompt.Run()
	promptBis := promptui.Prompt{
		Label: "Estensione dei file da concatenare",
	}
	e, _ := promptBis.Run()
	files := GetRightFiles(c, e)
	ChainFiles(c, e, files)
}

func GetRightFiles(folder string, ext string) []string {
	files, err := FilePathWalkDir(folder)
	if err != nil {
		fmt.Println(err)
	}
	results := FileMatchExt(files, ext)
	fmt.Println(results)
	return results
}

func FilePathWalkDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func FileMatchExt(files []string, ext string) []string {
	var results []string
	for _, res := range files {
		s := strings.Split(res, ".")
		if s[len(s)-1] == ext {
			results = append(results, res)
		}
	}
	return results
}

func ChainFiles(dir string, ext string, files []string) string {
	fn := fmt.Sprintf("%v/output.%v", dir, ext)
	output, err := os.Create(fn)

	if err != nil {
		fmt.Printf("Non posso creare il file %v", fn)
	}
	for _, file := range files {
		fi, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}
		defer fi.Close()

		scanner := bufio.NewScanner(fi)

		for scanner.Scan() {
			output.Write(scanner.Bytes())
			output.Write([]byte("\n"))
		}
	}
	defer output.Close()

	return "fatto"
}

func ScaricaHtmlDaTxt() {
	prompt := promptui.Prompt{
		Label: "Copia il path completo del file txt con le URL (C://......../file.txt)",
	}
	f, _ := prompt.Run()
	promptBis := promptui.Prompt{
		Label: "Il path completo della cartella dove salvare i file",
	}
	o, _ := promptBis.Run()
	ReadFileLineByLine(f, o)
}

func ReadFileLineByLine(fp string, output string) {
	//file, err := os.Open(fp)
	file, err := os.Open(fp)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		DownloadFileFromUrl(scanner.Text(), output)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func DownloadFileFromUrl(url string, outFolder string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("NON POSSO SCARICARE LA URL")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	//fmt.Println(body)
	WriteFileFromBody(body, outFolder)
}

func WriteFileFromBody(body []byte, outFolder string) {
	now := time.Now()
	name := fmt.Sprintf("%v", now.Unix())
	fn := fmt.Sprintf("%v/%v.%v", outFolder, name, "html")
	output, err := os.Create(fn)

	if err != nil {
		fmt.Printf("Non posso creare il file %v", fn)
	}

	output.Write(body)
	defer output.Close()
}
