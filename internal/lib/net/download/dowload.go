package download

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"losh/internal/lib/unit"

	n "losh/internal/lib/net"
	"losh/internal/lib/net/request"

	"github.com/aisbergg/go-errors/pkg/errors"
)

// ErrTooLarge indicates that the download content is too large.
type ErrTooLarge struct {
	Size  unit.ByteSize
	Limit unit.ByteSize
}

// Error implements the error interface.
func (e *ErrTooLarge) Error() string {
	if e.Size >= 0 {
		return fmt.Sprintf("content to large (%s > %s limit)", e.Size, e.Limit)
	}
	return fmt.Sprintf("content to large (> %s)", e.Limit)
}

// Context returns the context map.
func (e *ErrTooLarge) Context() map[string]interface{} {
	return map[string]interface{}{
		"size":  e.Size,
		"limit": e.Limit,
	}
}

// A Downloader defines the parameters for running a web crawler.
type Downloader struct {
	// The client used to make requests.
	requester *request.HTTPRequester

	// The user agent to use for regular HTTP requests.
	userAgent string
	// The headers to add to regular HTTP requests.
	headers http.Header
	// The return codes that are considered successful. Defaults to [200].
	okCodes []int
}

// NewDownloaderWithRequester creates a new Downloader with given HTTP client.
func NewDownloaderWithRequester(requester *request.HTTPRequester) *Downloader {
	return &Downloader{
		requester: requester,
		headers:   make(http.Header),
		okCodes:   []int{200},
	}
}

// NewDownloaderWithClient creates a new Downloader with given HTTP client.
func NewDownloaderWithClient(httpClient *http.Client) *Downloader {
	requester := request.NewHTTPRequester(httpClient).
		SetRetryCount(5).
		SetMaxWaitTime(60 * time.Second)
	return &Downloader{
		requester: requester,
		headers:   make(http.Header),
		okCodes:   []int{200},
	}
}

// NewDownloader creates a new Downloader with default HTTP client.
func NewDownloader() *Downloader {
	return NewDownloaderWithClient(http.DefaultClient)
}

// SetUserAgent sets the user agent to use for requests.
func (r *Downloader) SetUserAgent(userAgent string) *Downloader {
	r.userAgent = userAgent
	return r
}

// AddHeader adds a header to the requests.
func (r *Downloader) AddHeader(key, value string) *Downloader {
	r.headers.Add(key, value)
	return r
}

// AddOKCode adds a return code that is considered successful.
func (r *Downloader) AddOKCode(code int) *Downloader {
	r.okCodes = append(r.okCodes, code)
	return r
}

// SetOKCodes sets the return codes that are considered successful.
func (r *Downloader) SetOKCodes(codes []int) *Downloader {
	r.okCodes = codes
	return r
}

// DownloadContent downloads the content of the given request and returns it as
// bytes.
func (r *Downloader) DownloadContent(ctx context.Context, urlStr string) ([]byte, error) {
	return r.DownloadContentWithMaxSize(ctx, urlStr, 0)
}

// DownloadContentWithMaxSize downloads the content of the given request and
// returns it as bytes. If the content is larger than the given size, an error
// is returned.
func (r *Downloader) DownloadContentWithMaxSize(ctx context.Context, urlStr string, maxSize unit.ByteSize) ([]byte, error) {
	var buf bytes.Buffer
	if err := r.download(ctx, urlStr, maxSize, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// DownloadFile downloads a file from the given request and writes it to the given path.
func (r *Downloader) DownloadFile(ctx context.Context, urlStr, filename string) error {
	return r.DownloadFileWithMaxSize(ctx, urlStr, filename, 0)
}

// DownloadFileWithMaxSize downloads a file from the given request and writes it
// to the given path. If the content is larger than the given size, an error is
// returned and any downloaded content is discarded.
func (r *Downloader) DownloadFileWithMaxSize(ctx context.Context, urlStr, filename string, maxSize unit.ByteSize) error {
	// create blank file
	file, err := os.Create(filename)
	if err != nil {
		return errors.Wrap(err, "failed to create download file")
	}
	defer file.Close()

	// download content and write to file
	if err := r.download(ctx, urlStr, maxSize, file); err != nil {
		// remove partially downloaded file
		if err := file.Close(); err != nil {
			return errors.Wrap(err, "failed to delete partially downloaded file")
		}
		if err := os.Remove(filename); err != nil {
			return errors.Wrap(err, "failed to delete partially downloaded file")
		}
		return errors.Wrap(err, "file download failed")
	}

	return nil
}

func (r *Downloader) download(ctx context.Context, urlStr string, maxSize unit.ByteSize, w io.Writer) error {
	var (
		req    *http.Request
		resp   *http.Response
		reader io.Reader
		err    error
	)

	// create request and execute it
	req, err = http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}
	resp, err = r.requester.Do(req)
	if err != nil {
		return errors.Wrap(err, "request failed")
	}
	defer resp.Body.Close()

	// check if response is ok
	successful := false
	for _, code := range r.okCodes {
		if resp.StatusCode == code {
			successful = true
			break
		}
	}
	if !successful {
		return errors.Errorf("request failed with status code %d", resp.StatusCode)
	}

	// continue with downloading the request body
	reader = resp.Body

	// decompress gzip if necessary
	if strings.EqualFold(resp.Header.Get(n.HdrContentEncodingKey), "gzip") && resp.ContentLength != 0 {
		if _, ok := reader.(*gzip.Reader); !ok {
			gz, err := gzip.NewReader(reader)
			if err != nil {
				return errors.Wrap(err, "failed to create gzip reader")
			}
			defer gz.Close()
			reader = gz
		}
	}

	// create custom writer to limit downloaded size
	if maxSize > 0 {
		// quit early if content is too large
		if resp.ContentLength > int64(maxSize) {
			return &ErrTooLarge{Size: unit.ByteSize(resp.ContentLength), Limit: maxSize}
		}

		// create a new reader that will limit the size of the downloaded content
		reader = &limitedReader{R: resp.Body, N: maxSize}
	}

	// copy content to writer
	if _, err := io.Copy(w, reader); err != nil {
		if err == errLimitExceeded {
			return errors.Wrap(&ErrTooLarge{Limit: maxSize}, "failed to download content")
		}
		return errors.Wrap(err, "failed to download content")
	}

	return nil
}

func copyHeaders(headers http.Header) http.Header {
	nh := http.Header{}
	for k, v := range headers {
		nh[k] = v
	}
	return nh
}
