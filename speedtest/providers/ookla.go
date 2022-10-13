package providers

import (
	ookla "github.com/showwin/speedtest-go/speedtest"
)

type OoklaProvider struct {
	target *ookla.Server
}

func (o *OoklaProvider) Name() string {
	return "ookla"
}

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

	return nil
}

func (o *OoklaProvider) DownloadTest() (bits uint64, err error) {
	//err := t.PingTest()
	err = o.target.DownloadTest(false)
	if err != nil {
		return 0, err
	}

	return convertOoklaResult(o.target.DLSpeed), nil
}

func (o *OoklaProvider) UploadTest() (bits uint64, err error) {
	err = o.target.UploadTest(false)
	if err != nil {
		return 0, err
	}

	return convertOoklaResult(o.target.ULSpeed), nil
}

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

func convertOoklaResult(mBytes float64) (bits uint64) {
	return uint64(mBytes * 1024 * 1024 * 8)
}
