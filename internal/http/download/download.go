package download

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

const stallTimeout = 2 * time.Second

type stallReader struct {
	r     io.ReadCloser
	timer *time.Timer
	d     time.Duration
}

func (s *stallReader) Read(p []byte) (int, error) {
	n, err := s.r.Read(p)
	if n > 0 {
		s.timer.Reset(s.d)
	}

	return n, err
}

func (s *stallReader) Close() error {
	s.timer.Stop()
	return s.r.Close()
}

func File(url string, p *mpb.Progress) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		fmt.Println("Request error:", url, err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Request error:", url, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Bad status:", url, resp.Status)
		return
	}

	fileName := fetchFileName(resp)
	destination := getDestination(fileName)

	if err := os.MkdirAll(filepath.Dir(destination), os.ModePerm); err != nil {
		fmt.Printf("Failed to create directory %s: %s\n", filepath.Dir(destination), err)
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
			decor.OnAbort(decor.CountersKibiByte("% .2f / % .2f"), "canceled: no data"),
		),
		mpb.AppendDecorators(
			decor.OnAbort(decor.EwmaETA(decor.ET_STYLE_GO, 60), ""),
			decor.OnAbort(decor.Name(" ]"), ""),
			decor.OnAbort(decor.EwmaSpeed(decor.SizeB1024(0), "% .2f", 60), ""),
		),
	)

	timer := time.AfterFunc(stallTimeout, cancel)
	src := &stallReader{r: resp.Body, timer: timer, d: stallTimeout}

	reader := bar.ProxyReader(src)
	defer reader.Close()

	if _, err := io.Copy(out, reader); err != nil {
		bar.Abort(false)
		out.Close()
		os.Remove(destination)
	}
}

func fetchFileName(resp *http.Response) string {
	if cd := resp.Header.Get("Content-Disposition"); cd != "" {
		_, params, err := mime.ParseMediaType(cd)
		if err == nil {
			if fn, ok := params["filename"]; ok {
				return fn
			}
		}
	}

	return path.Base(resp.Request.URL.Path)
}

func getDestination(fileName string) string {
	return filepath.Join("downloads", fileName)
}
