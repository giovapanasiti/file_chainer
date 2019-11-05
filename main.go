package main

import (
	// "bufio"
	"fmt"
	// "io/ioutil"
	// "net/http"
	// "time"

	//"io/ioutil"
	. "github.com/logrusorgru/aurora"
	"github.com/manifoldco/promptui"
	// "log"
	"os"
	// "path/filepath"
	// "strings"
	"github.com/giovapanasiti/file_chainer/files"
	"github.com/giovapanasiti/file_chainer/report_stampa"
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
		files.ConcatenaFiles()
	case "Scarica Lista Link":
		files.ScaricaHtmlDaTxt()
	case "Report Monitoraggio >> PDF":
		report_stampa.ReadDocx()
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
		Items: []string{"Concatena Files", "Scarica Lista Link", "Report Monitoraggio >> PDF", "Esci"},
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
