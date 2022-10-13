package main

import (
	"fmt"
	"github.com/konovenski/turbo-snail/speedtest"
	"github.com/konovenski/turbo-snail/speedtest/providers"
	"log"
)

func main() {
	fmt.Println("Welcome to the turbo-snail showcase")

	for _, provider := range []speedtest.Provider{
		&providers.OoklaProvider{},
		&providers.FastProvider{},
	} {
		fmt.Println()
		runTest(provider)
	}
}

const TestFailedTemplate = "An error occurred during speedtest for provider '%s': %s\n"
const TestFinishedTemplate = "Speedtest results for '%s' provider: \n" +
	"Download speed: %.2f MB/s\n" +
	"Upload speed: %.2f MB/s\n"

func runTest(provider speedtest.Provider) {
	sp, err := speedtest.NewSpeedtest(provider)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("%s provider initialized\n", sp.Name())

	downloadBits, uploadBits, err := sp.CompleteTest()
	if err != nil {
		fmt.Printf(TestFailedTemplate, sp.Name(), err)
		return
	}

	fmt.Printf(TestFinishedTemplate, sp.Name(), toMBits(downloadBits), toMBits(uploadBits))
}

func toMBits(bits uint64) float64 {
	return float64(bits) / 1024 / 1024
}
