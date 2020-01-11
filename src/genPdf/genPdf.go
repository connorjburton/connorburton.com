package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

const PATH = "./../../"

func main() {
	source := fmt.Sprintf("%s%s", PATH, "public/index.html")
	dest := fmt.Sprintf("%s%s", PATH, "public/cv.pdf")
	
	htmlBuf, err := readHtml(source)
	if err != nil {
		panic(err)
	}

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ts := createTestServer(htmlBuf)
	defer ts.Close()

	var buf []byte

	if err := chromedp.Run(ctx, htmlToPdf(ts.URL, &buf)); err != nil {
		panic(err)
	}

	if err := writePdf(dest, &buf); err != nil {
		panic(err)
	}
}

func createTestServer(buf []byte) (*httptest.Server) {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, string(buf))
	}))
}

func writePdf(dest string, buf *[]byte) error {
	return ioutil.WriteFile(dest, *buf, 0644)
}

func readHtml(source string) ([]byte, error) {
	return ioutil.ReadFile(source)
}

func htmlToPdf(url string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().Do(ctx)
			if err != nil {
				return err
			}

			*res = buf
			return nil
		}),
	}
}