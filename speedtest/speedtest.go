package speedtest

type Provider interface {
	Name() string
	Init() error
	DownloadTest() (bits uint64, err error)
	UploadTest() (bits uint64, err error)
	CompleteTest() (dBits uint64, uBits uint64, err error)
}

func NewSpeedtest(provider Provider) (Provider, error) {
	err := provider.Init()
	return provider, err
}
