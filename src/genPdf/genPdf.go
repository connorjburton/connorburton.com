package main

import (
	"context"
	"fmt"
	"path"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

const PATH = "./../../"

func main() {
	source := path.Join(PATH, "public/index.html")
	dest := path.Join(PATH, "public/cv.pdf")
	
	htmlBuf, err := ioutil.ReadFile(source)
	if err != nil {
		panic(err)
	}

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, string(htmlBuf))
	}))
	defer ts.Close()

	var buf []byte

	if err := chromedp.Run(ctx, htmlToPdf(ts.URL, &buf)); err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile(dest, buf, 0644); err != nil {
		panic(err)
	}
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