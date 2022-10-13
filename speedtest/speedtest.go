package speedtest

type Provider interface {
	//Name should return the human-readable name of the provider
	Name() string
	//Init should initialize all the provider resources. This should be done before test
	Init() error
	//DownloadTest performs download speedtest. Speed should be implemented as bits per second
	DownloadTest() (bits uint64, err error)
	//UploadTest performs upload speedtest. Speed should be implemented as bits per second
	UploadTest() (bits uint64, err error)
	//CompleteTest performs both download and upload speedtest. Just for the comfort
	CompleteTest() (dBits uint64, uBits uint64, err error)
}

//NewSpeedtest receives an instance of Provider and initiates it
func NewSpeedtest(provider Provider) (Provider, error) {
	err := provider.Init()
	return provider, err
}
