package providers

import (
	"errors"
	ookla "github.com/showwin/speedtest-go/speedtest"
)

// OoklaProvider is being used for ookla based speedtest.
type OoklaProvider struct {
	initialized bool
	target      *ookla.Server
}

// Name returns human-readable name
func (o *OoklaProvider) Name() string {
	return "ookla"
}

// Init receives test URL for ookla provider.
func (o *OoklaProvider) Init() error {
	user, err := ookla.FetchUserInfo()
	if err != nil {
		return err
	}

	serverList, err := ookla.FetchServers(user)
	if err != nil {
		return err
	}

	targets, err := serverList.FindServer([]int{})
	if err != nil {
		return err
	}

	// Slice couldn't be empty corresponding to the implementation.
	//See FindServer implementation for more info
	o.target = targets[0]
	o.initialized = true
	return nil
}

// DownloadTest performs download speedtest.
func (o *OoklaProvider) DownloadTest() (bits uint64, err error) {
	if !o.initialized {
		return 0, errors.New("provider was not initialized")
	}

	err = o.target.DownloadTest(false)
	if err != nil {
		return 0, err
	}

	return convertOoklaResult(o.target.DLSpeed), nil
}

// UploadTest performs upload speedtest.
func (o *OoklaProvider) UploadTest() (bits uint64, err error) {
	if !o.initialized {
		return 0, errors.New("provider was not initialized")
	}

	err = o.target.UploadTest(false)
	if err != nil {
		return 0, err
	}

	return convertOoklaResult(o.target.ULSpeed), nil
}

// CompleteTest performs both download and upload speedtest.
func (o *OoklaProvider) CompleteTest() (dBits uint64, uBits uint64, err error) {
	dBits, err = o.DownloadTest()
	if err != nil {
		return 0, 0, err
	}
	uBits, err = o.UploadTest()
	if err != nil {
		return 0, 0, err
	}

	return dBits, uBits, nil
}

// convertOoklaResult converts ooklas format (MB) to our format (bits)
func convertOoklaResult(mBytes float64) (bits uint64) {
	return uint64(mBytes * 1024 * 1024)
}
