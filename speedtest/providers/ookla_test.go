package providers

import "testing"

func TestOoklaProvider(t *testing.T) {
	o := OoklaProvider{}

	_, _, err := o.CompleteTest()
	if err == nil {
		t.Error("Should not be executed before initialization", err)
	}

	err = o.Init()
	if err != nil {
		t.Error("Got an error:", err)
	}

	if o.target == nil {
		t.Error("Target should be provided")
	}

	bits, err := o.DownloadTest()
	if err != nil {
		t.Error("An error occurred: ", err)
	}

	if bits == 0 {
		t.Error("bits should not be 0 when no error triggered")
	}

	bits, err = o.UploadTest()
	if err != nil {
		t.Error("An error occurred: ", err)
	}

	if bits == 0 {
		t.Error("bits should not be 0 when no error triggered")
	}

	ubits, dbits, err := o.CompleteTest()
	if err != nil {
		t.Error("An error occurred: ", err)
	}

	if dbits == 0 || ubits == 0 {
		t.Error("bits should not be 0 when no error triggered")
	}
}

func BenchmarkOoklaProvider(b *testing.B) {
	o := OoklaProvider{}
	err := o.Init()
	if err != nil {
		b.Error("Can't initialize provider:", err)
	}

	for i := 0; i < b.N; i++ {
		_, _, err = o.CompleteTest()
		if err != nil {
			b.Error("Error during benchmarking:", err)
		}
	}
}
