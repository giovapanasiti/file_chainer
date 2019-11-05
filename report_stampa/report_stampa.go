package report_stampa

import (
	"bufio"
	"context"
	"fmt"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/giovapanasiti/file_chainer/docx2md"
	"io/ioutil"
	"log"
	"mvdan.cc/xurls"
	"os"

	"math"


	"net/url"
	"github.com/mafredri/cdp/devtool"
	"github.com/mafredri/cdp/rpcc"
	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/target"
	"github.com/mafredri/cdp/protocol/network"
	newPage "github.com/mafredri/cdp/protocol/page"
	"time"
	"github.com/manifoldco/promptui"
	"code.sajari.com/docconv"
)

/*
ReadDocx ad ashdas sadas
*/
func ReadDocx() {

	prompt := promptui.Prompt{
		Label: "Path del file docx",
	}
	c, _ := prompt.Run()
	//c:="/Users/panasiti_g/Downloads/report.docx"
	prompt2 := promptui.Prompt{
		Label: "Path di output",
	}
	f, _ := prompt2.Run()
	//f:="/Users/panasiti_g/Downloads/output"
	_, err := docconv.ConvertPath(c)
	fmt.Printf("Converto il file %v \n", c)

	err, content := docx2md.Convert(c, true)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("File convertito")
	outputFile := WriteFileMd([]byte(content), f)

	urls := ParseFileForUrls(outputFile)

	for i, url := range urls {
		fmt.Println(url)
		VisitPage(url, f, i)
	}

}

func WriteFileMd(body []byte, outFolder string) string {

	fmt.Printf("Salvo il file in %v \n", outFolder)
	name := "output.txt"
	fn := fmt.Sprintf("%v/%v", outFolder, name)
	output, err := os.Create(fn)
	if err != nil {
		fmt.Printf("Non posso creare il file %v \n", fn)
	}
	output.Write(body)
	defer output.Close()
	fmt.Printf("File creato %v/%v \n", outFolder, name)

	return fmt.Sprintf("%v/%v", outFolder, name)
}

func ParseFileForUrls(fp string) []string {
	// (https?):\/\/(www\.)?[a-z0-9\.:].*?(?=\s)

	var result []string

	file, err := os.Open(fp)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		//regex := regexp.MustCompile(`^(?:https?:\/\/)?(?:[^@\/\n]+@)?(?:www\.)?([^:\/\n]+)`)
		//regex := regexp.MustCompile(`\((.*?)\)`)
		//regex2 := regexp.MustCompile(`\([.*?]\)`)

		s := scanner.Text()
		//DownloadFileFromUrl(scanner.Text(), output)
		//match := regex.FindStringSubmatch(scanner.Text())
		//match2 := regex.FindString(s)  // FindStringSubmatch(scanner.Text())
		match := xurls.Relaxed().FindAllString(s, -1)

		if len(match) > 0 {
			result = append(result, match[0])
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return result
}

func VisitPage(url string, fp string, i int) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// capture screenshot of an element
	var buf []byte
	//if err := chromedp.Run(ctx, elementScreenshot(`https://www.google.com/`, `#main`, &buf)); err != nil {
	//	log.Fatal(err)
	//}
	//if err := ioutil.WriteFile("elementScreenshot.png", buf, 0644); err != nil {
	//	log.Fatal(err)
	//}

	// capture entire browser viewport, returning png with quality=90
	if err := chromedp.Run(ctx, fullScreenshotPdf(url, 90, &buf)); err != nil {
		log.Fatal(err)
	}
	path := fmt.Sprintf("%v/%v.pdf", fp, i)
	if err := ioutil.WriteFile(path, buf, 0644); err != nil {
		log.Fatal(err)
	}
}

// fullScreenshot takes a screenshot of the entire browser viewport.
//
// Liberally copied from puppeteer's source.
//
// Note: this will override the viewport emulation settings.
func fullScreenshot(urlstr string, quality int64, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// get layout metrics
			_, _, contentSize, err := page.GetLayoutMetrics().Do(ctx)
			if err != nil {
				return err
			}

			width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))

			// force viewport emulation
			err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).
				Do(ctx)
			if err != nil {
				return err
			}

			// capture screenshot
			*res, err = page.CaptureScreenshot().
				WithQuality(quality).
				WithClip(&page.Viewport{
					X:      contentSize.X,
					Y:      contentSize.Y,
					Width:  contentSize.Width,
					Height: contentSize.Height,
					Scale:  1,
				}).Do(ctx)
			if err != nil {
				return err
			}
			return nil
		}),
	}
}


func fullScreenshotPdf(urlstr string, quality int64, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// get layout metrics
			_, _, contentSize, err := page.GetLayoutMetrics().Do(ctx)
			if err != nil {
				return err
			}

			width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))

			// force viewport emulation
			err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
				WithScreenOrientation(&emulation.ScreenOrientation{
				Type:  emulation.OrientationTypePortraitPrimary,
				Angle: 0,
			}).
				Do(ctx)
			if err != nil {
				return err
			}

			//capture screenshot
			*res, _, err = page.PrintToPDF().
				WithLandscape(false).
				WithPaperHeight(11.7).
				WithPaperWidth(8.27).Do(ctx)
			if err != nil {
				fmt.Println(err)
				fmt.Println("Errore cazzo")
				return err
			}
			return nil
		}),
	}
}



func CreatePdf(urlRequest string) ([]byte, error) {

	width := 21.5  //float64,  float64
	height := 21.5

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()


	// Use the DevTools API to manage targets
	devt := devtool.New("http://127.0.0.1:9222")
	pt, err := devt.Get(ctx, devtool.Page)
	if err != nil {
		pt, err = devt.Create(ctx)
		if err != nil {
			return nil, err
		}
	}

	// Initiate a new RPC connection to the Chrome DevTools Protocol target.
	conn, err := rpcc.DialContext(ctx, pt.WebSocketDebuggerURL)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Create new browser context
	baseBrowser := cdp.NewClient(conn)
	newContextTarget, _ := baseBrowser.Target.CreateBrowserContext(ctx)

	// Create a new blank target with the new browser context
	newTargetArgs := target.NewCreateTargetArgs("about:blank").
		SetBrowserContextID(newContextTarget.BrowserContextID)
	newTarget, _ := baseBrowser.Target.CreateTarget(ctx, newTargetArgs)

	// Connect the client to the new target
	newTargetWsURL := fmt.Sprintf("ws://127.0.0.1:9222/devtools/page/%s", newTarget.TargetID)
	newContextConn, _ := rpcc.DialContext(ctx, newTargetWsURL)
	defer newContextConn.Close()
	c := cdp.NewClient(newContextConn)

	// Close the target when done
	// (In development, skip this step to leave tabs open!)
	closeTargetArgs := target.NewCloseTargetArgs(newTarget.TargetID)
	defer baseBrowser.Target.CloseTarget(ctx, closeTargetArgs)

	// Enable the runtime
	c.Runtime.Enable(ctx)

	// Enable the network
	c.Network.Enable(ctx, network.NewEnableArgs())

	// Enable events on the Page domain
	c.Page.Enable(ctx)

	// Create a client to listen for the load event to be fired
	loadEventFiredClient, _ := c.Page.LoadEventFired(ctx)
	defer loadEventFiredClient.Close()

	// Tell the page to navigate to the URL
	url.ParseRequestURI(urlRequest)
	navArgs := newPage.NewNavigateArgs(urlRequest)
	c.Page.Navigate(ctx, navArgs)

	// Wait for the page to finish loading
	loadEventFiredClient.Recv()

	// Print to PDF
	printToPDFArgs := newPage.NewPrintToPDFArgs().
		SetLandscape(true).
		SetPrintBackground(true).
		SetMarginTop(0).
		SetMarginBottom(0).
		SetMarginLeft(0).
		SetMarginRight(0).
		SetPrintBackground(true).
		SetPaperWidth(width).
		SetPaperHeight(height)
	print, _ := c.Page.PrintToPDF(ctx, printToPDFArgs)
	return print.Data, nil
}


func WriteFileFromBytes(i int, body []byte, outFolder string) {

	fn := fmt.Sprintf("%v/%v.%v", outFolder, i, "pdf")
	output, err := os.Create(fn)

	if err != nil {
		fmt.Printf("Non posso creare il file %v", fn)
	}

	output.Write(body)
	defer output.Close()
}