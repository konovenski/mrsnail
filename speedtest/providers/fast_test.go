package providers

import (
	"testing"
	"time"
)

func TestFastProvider(t *testing.T) {
	f := FastProvider{}

	_, _, err := f.CompleteTest()
	if err == nil {
		t.Error("Should not be executed before initialization", err)
	}

	err = f.Init()
	if err != nil {
		t.Error("Got an error:", err)
	}

	if f.targetAmount == 0 || f.fileSizeInMBytes == 0 {
		t.Error("Target amount and file size should be > 0")
	}

	_, err = f.runSpeedTest("unknown", 1)
	if err == nil {
		t.Error("Speedtest should return error in case of unknown mode")
	}

	_, err = f.runSpeedTest("Download", -1)
	if err == nil {
		t.Error("Speedtest should return error in case of 0 payload")
	}

	mBits := float64(f.fileSizeInMBytes*f.targetAmount*8) / 1
	bits := uint64(mBits * 1024 * 1024)
	if f.calculateSpeed(time.Second) != bits {
		t.Error("Speed should be calculates in bits, using the corresponding parameters from provider config")
	}

	bits, err = f.DownloadTest()
	if err != nil {
		t.Error("An error occured: ", err)
	}

	if bits == 0 {
		t.Error("bits should not be 0 when no error triggered")
	}

	bits, err = f.UploadTest()
	if err != nil {
		t.Error("An error occured: ", err)
	}

	if bits == 0 {
		t.Error("bits should not be 0 when no error triggered")
	}

	ubits, dbits, err := f.CompleteTest()
	if err != nil {
		t.Error("An error occurred: ", err)
	}

	if dbits == 0 || ubits == 0 {
		t.Error("bits should not be 0 when no error triggered")
	}
}

func BenchmarkFastProvider(b *testing.B) {
	f := FastProvider{}
	err := f.Init()
	if err != nil {
		b.Error("Can't initialize provider:", err)
	}

	for i := 0; i < b.N; i++ {
		_, _, err = f.CompleteTest()
		if err != nil {
			b.Error("Error during benchmarking:", err)
		}
	}
}
