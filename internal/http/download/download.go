package download

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

func File(url string, p *mpb.Progress, wg *sync.WaitGroup) {
	defer wg.Done()

	fileName := filepath.Base(url)
	destination := getDestination(fileName)

	if err := os.MkdirAll(filepath.Dir(destination), os.ModePerm); err != nil {
		fmt.Printf("Failed to create directory %s: %s\n", filepath.Dir(destination), err)
		return
	}

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Request error:", url, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Bad status:", url, resp.Status)
		return
	}

	out, err := os.Create(destination)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer out.Close()

	total := resp.ContentLength

	bar := p.AddBar(total,
		mpb.PrependDecorators(
			decor.Name(fileName, decor.WC{W: len(fileName) + 1, C: decor.DindentRight}),
			decor.CountersKibiByte("% .2f / % .2f"),
		),
		mpb.AppendDecorators(
			decor.EwmaETA(decor.ET_STYLE_GO, 60),
			decor.Name(" ]"),
			decor.EwmaSpeed(decor.SizeB1024(0), "% .2f", 60),
		),
	)

	reader := bar.ProxyReader(resp.Body)
	defer reader.Close()

	_, err = io.Copy(out, reader)
	if err != nil {
		fmt.Println("Ошибка копирования:", err)
	}
}

func getDestination(fileName string) string {
	return filepath.Join("downloads", fileName)
}
