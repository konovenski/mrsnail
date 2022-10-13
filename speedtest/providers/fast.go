package providers

import (
	"errors"
	"fmt"
	"github.com/gesquive/fast-cli/fast"
	"golang.org/x/sync/errgroup"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// This implementation relied on what is described in the following article.
// https://netflixtechblog.com/building-fast-com-4857fe0f8adb
// However, this implementation isn't using the 0-2048 chunk due to lack
// of information about how to use it during calculations

const FastDefaultTargetAmount uint = 5
const FastDefaultFileSizeInMBytes uint = 25

// FastProvider is being used for fast.com based speedtest.
type FastProvider struct {
	FileSizeInMBytes uint
	TargetAmount     uint
	targets          []string
	initialized      bool
}

// Name returns human-readable name
func (f *FastProvider) Name() string {
	return "fast"
}

// Init verifies initial config and receives test URL's for fast.com provider.
// In case of unexpected failures here the most probable reason is that netflix
// changed the token/layout format, because API token is being extracted from the
// website files
func (f *FastProvider) Init() error {
	if f.TargetAmount <= 0 {
		f.TargetAmount = FastDefaultTargetAmount
	}

	if f.FileSizeInMBytes <= 0 {
		f.FileSizeInMBytes = FastDefaultFileSizeInMBytes
	}

	f.targets = fast.GetDlUrls(uint64(f.TargetAmount))
	if uint(len(f.targets)) != f.TargetAmount {
		return errors.New("can't fetch any targets")
	}
	f.initialized = true
	return nil
}

// DownloadTest performs download speedtest.
func (f *FastProvider) DownloadTest() (bits uint64, err error) {
	if !f.initialized {
		return 0, errors.New("provider was not initialized")
	}
	return f.runCompositeSpeedTest("download")
}

// UploadTest performs upload speedtest.
func (f *FastProvider) UploadTest() (bits uint64, err error) {
	if !f.initialized {
		return 0, errors.New("provider was not initialized")
	}
	return f.runCompositeSpeedTest("upload")
}

// CompleteTest performs both download and upload speedtest.
func (f *FastProvider) CompleteTest() (dBits uint64, uBits uint64, err error) {
	r, err := f.DownloadTest()
	if err != nil {
		return 0, 0, err
	}
	u, err := f.UploadTest()
	if err != nil {
		return 0, 0, err
	}
	return r, u, nil
}

// runCompositeSpeedTest starting a fast.com methodology-based test.
// For more info check the article. https://netflixtechblog.com/building-fast-com-4857fe0f8adb
func (f *FastProvider) runCompositeSpeedTest(mode string) (uint64, error) {
	rampResult, err := f.runSpeedTest(mode, 1)
	if err != nil {
		return 0, err
	}

	speedResult, err := f.runSpeedTest(mode, 25*1024*1024)
	if err != nil {
		return 0, err
	}

	return f.calculateSpeed(speedResult - rampResult), nil
}

// runSpeedTest starting an upload or download speedtest with the given number of bytes.
func (f *FastProvider) runSpeedTest(mode string, payloadSize int) (time.Duration, error) {
	if payloadSize < 1 {
		return 0, errors.New("payload size should be at least 1")
	}
	eg := errgroup.Group{}
	startTime := time.Now()

	for _, target := range f.targets {
		eg.Go(func() error {
			switch mode {
			case "download":
				return downloadFastSample(target, payloadSize)
			case "upload":
				return uploadFastSample(target, payloadSize)
			default:
				return errors.New("unknown run mode")
			}
		})
	}
	if err := eg.Wait(); err != nil {
		return 0, err
	}
	return time.Now().Sub(startTime), nil
}

// calculateSpeed returns amount of bits per second. Aggregated from all the other config and result variables
// obtained during testing
func (f *FastProvider) calculateSpeed(seconds time.Duration) (bits uint64) {
	if seconds.Seconds() == 0 {
		return 0
	}
	t := seconds.Seconds()
	mBits := float64(f.FileSizeInMBytes*f.TargetAmount*8) / t
	return uint64(mBits * 1024 * 1024)
}

// downloadFastSample downloads a sample of given size from the fast.com CDN.
// Being used for speedtest
func downloadFastSample(target string, payloadSize int) error {
	var client = http.Client{}

	u, _ := url.Parse(target)
	u.Path = fmt.Sprintf("/speedtest/range/0-%d", payloadSize)
	target = u.String()

	resp, err := client.Get(target)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	return err
}

// uploadFastSample downloads a sample of given size from the fast.com CDN.
// Being used for speedtest
func uploadFastSample(target string, payloadSize int) error {
	var client = http.Client{}

	data := url.Values{}
	data.Add("content", strings.Repeat("0", payloadSize))

	resp, err := client.PostForm(target, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	return err
}
