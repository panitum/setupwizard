package download

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

const stallTimeout = 15 * time.Second

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

	if alreadyDownloaded(destination, resp) {
		resp.Body.Close()

		bar := p.AddBar(resp.ContentLength,
			mpb.BarFillerClearOnComplete(),
			mpb.PrependDecorators(
				decor.Name(fileName, decor.WC{W: len(fileName) + 1, C: decor.DindentRight}),
				decor.OnComplete(decor.Name(""), "уже скачано"),
			),
		)
		bar.SetCurrent(resp.ContentLength)
		bar.Wait()

		return
	}

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

	tmp := destination + ".part"

	out, err = os.Create(tmp)
	if err != nil {
		fmt.Println("Error creating temp file:", err)
		return
	}
	defer out.Close()

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

	if _, err := io.Copy(out, reader); err != nil {
		bar.Abort(false)
		out.Close()
		os.Remove(tmp)
		return
	}

	out.Close()
	if err := os.Rename(tmp, destination); err != nil {
		bar.Abort(false)
		return
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

func alreadyDownloaded(destination string, resp *http.Response) bool {
	fi, err := os.Stat(destination)
	if err != nil {
		return false
	}

	if resp.ContentLength < 0 {
		return false
	}
	if fi.Size() != resp.ContentLength {
		return false
	}

	if sum, ok := etagMD5(resp.Header.Get("ETag")); ok {
		local, err := fileMD5(destination)
		if err != nil || local != sum {
			return false
		}
	}
	return true
}

func etagMD5(etag string) (string, bool) {
	etag = strings.Trim(etag, `"`)
	if len(etag) != 32 || strings.Contains(etag, "-") {
		return "", false
	}
	if _, err := hex.DecodeString(etag); err != nil {
		return "", false
	}
	return strings.ToLower(etag), true
}

func fileMD5(p string) (string, error) {
	f, err := os.Open(p)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
