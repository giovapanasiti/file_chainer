package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	//"io/ioutil"
	. "github.com/logrusorgru/aurora"
	"github.com/manifoldco/promptui"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	//reader := bufio.NewReader(os.Stdin)
	//fmt.Print("Incollare qui la cartella dove si trovano i file: ")
	//root, _ := reader.ReadString('\n')
	//fmt.Println(root)

	wm := `
                                                                      
                                               dddddddd                    
    MMMMMMMM               MMMMMMMM            d::::::d       GGGGGGGGGGGGG
    M:::::::M             M:::::::M            d::::::d    GGG::::::::::::G
    M::::::::M           M::::::::M            d::::::d  GG:::::::::::::::G
    M:::::::::M         M:::::::::M            d:::::d  G:::::GGGGGGGG::::G
    M::::::::::M       M::::::::::M    ddddddddd:::::d G:::::G       GGGGGG
    M:::::::::::M     M:::::::::::M  dd::::::::::::::dG:::::G              
    M:::::::M::::M   M::::M:::::::M d::::::::::::::::dG:::::G              
    M::::::M M::::M M::::M M::::::Md:::::::ddddd:::::dG:::::G    GGGGGGGGGG
    M::::::M  M::::M::::M  M::::::Md::::::d    d:::::dG:::::G    G::::::::G
    M::::::M   M:::::::M   M::::::Md:::::d     d:::::dG:::::G    GGGGG::::G
    M::::::M    M:::::M    M::::::Md:::::d     d:::::dG:::::G        G::::G
    M::::::M     MMMMM     M::::::Md:::::d     d:::::d G:::::G       G::::G
    M::::::M               M::::::Md::::::ddddd::::::dd G:::::GGGGGGGG::::G
    M::::::M               M::::::M d:::::::::::::::::d  GG:::::::::::::::G
    M::::::M               M::::::M  d:::::::::ddd::::d    GGG::::::GGG:::G
    MMMMMMMM               MMMMMMMM   ddddddddd   ddddd       GGGGGG   GGGG

_____________________________________________________________________________
_____________________________________________________________________________

 - Applicazione:	Noesi Multi Utility tool
 - Versione:  		0.1
 - Autore:        	Giovanni Panasiti
 - Email:			giovanni@montedelgallo.com
_____________________________________________________________________________
`

	fmt.Println(Yellow(wm))

	doAnAction(askWhatToDo)
}

func doAnAction(ask func() string) {
	scelta := ask()
	switch scelta {
	case "Concatena Files":
		ConcatenaFiles()
	case "Scarica Lista Link":
		ScaricaHtmlDaTxt()
	case "Esci":
		os.Exit(1)
	}
	ancora := againAskWhatToDo()
	switch ancora {
	case "SI":
		doAnAction(ask)
	case "NO":
		fmt.Println(Red(`Ok! Alla Prossima!`))
		os.Exit(1)
	}
	doAnAction(ask)
}

func askWhatToDo() (result string) {
	prompt := promptui.Select{
		Label: "Cosa vuoi fare?",
		Items: []string{"Concatena Files", "Scarica Lista Link", "Esci"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
	return result
}

func againAskWhatToDo() (result string) {
	prompt := promptui.Select{
		Label: "Vuoi fare altro?",
		Items: []string{"SI", "NO"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
	return result
}

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
	file, err := os.Open("./pippo.txt")
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
	fn := fmt.Sprintf("%v/%v.%v", outFolder, name,"html")
	output, err := os.Create(fn)

	if err != nil {
		fmt.Printf("Non posso creare il file %v", fn)
	}

	output.Write(body)
	defer output.Close()
}
