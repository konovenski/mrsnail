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

//https://ipv6-c001-lis006-meo-isp.1.oca.nflxvideo.net/speedtest?c=pt&n=3243&v=95&e=1665638321&t=G_5bchBzYW92mDTv2oq4EQc1xCa3P9qdTsqmhA
//https://ipv6-c001-lis006-meo-isp.1.oca.nflxvideo.net/speedtest/range/0-0?c=pt&n=3243&v=95&e=1665638321&t=G_5bchBzYW92mDTv2oq4EQc1xCa3P9qdTsqmhA

// This implementation is rely on what is described in the following article.
// However, this implementation isn't using the 0-2048 chunk due to lack
// of information about how to use it during calculations

const FastDefaultTargetAmount int = 5
const FastDefaultFileSizeInMBytes int = 25

type FastProvider struct {
	fileSizeInMBytes int
	targetAmount     int
	targets          []string
}

func (f *FastProvider) Name() string {
	return "fast"
}

var client = http.Client{}

func (f *FastProvider) Init() error {
	if f.targetAmount == 0 {
		f.targetAmount = FastDefaultTargetAmount
	}

	if f.fileSizeInMBytes == 0 {
		f.fileSizeInMBytes = FastDefaultFileSizeInMBytes
	}

	f.targets = fast.GetDlUrls(uint64(f.targetAmount))
	if len(f.targets) != f.targetAmount {
		return errors.New("can't fetch any targets")
	}
	return nil
}

func (f *FastProvider) DownloadTest() (bits uint64, err error) {
	return f.runCompositeSpeedTest("download")
}

func (f *FastProvider) UploadTest() (bits uint64, err error) {
	return f.runCompositeSpeedTest("upload")
}

func (f *FastProvider) CompleteTest() (dBits uint64, uBits uint64, err error) {
	r, _ := f.DownloadTest()
	u, _ := f.UploadTest()
	return r, u, nil
}

func dw(target string, payloadSize int) error {
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

func up(target string, payloadSize int) error {
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

func (f *FastProvider) runSpeedTest(mode string, payloadSize int) (time.Duration, error) {
	eg := errgroup.Group{}
	startTime := time.Now()

	for _, target := range f.targets {
		eg.Go(func() error {
			switch mode {
			case "download":
				return dw(target, payloadSize)
			case "upload":
				return up(target, payloadSize)
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

func (f *FastProvider) calculateSpeed(seconds time.Duration) uint64 {
	if seconds.Seconds() == 0 {
		return 0
	}
	t := seconds.Seconds()
	mbits := float64(f.fileSizeInMBytes*f.targetAmount*8) / t
	return uint64(mbits * 1024 * 1024)
}
