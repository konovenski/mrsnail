# Turbo Snail
### The best speed test library and cli app* 


#### *according to Dumitru Conovenschi rating

![Mr. Snail](doc/snail.jpg)

# Using as library
## Installation:
```bash
$ go get github.com/konovenski/turbo-snail
```
## Usage:
```go
import "github.com/konovenski/turbo-snail"
// Pick provider (or write your own, using our interface)
// and initialize the test
ookla := speedtest.NewSpeedtest(&providers.OoklaProvider{})
// Perform test. 
dbits, err := ookla.DownloadTest()
ubits, err := ookla.UploadTest()
// We return results in bits per second.
// You can convert it to MBits
mbits := float64(bits) / 1024 / 1024
// You can also perform Complete test using oneliner
dbits, ubits, err := ookla.CompleteTest()
//That it!
```

# Using as standalone app
## Installation
```bash
$ git clone git@github.com:konovenski/turbo-snail.git
$ cd turbo-snail && make build
```

## Usage
```bash
$ ./turbo-snail
Welcome to the turbo-snail showcase

ookla provider initialized
Speedtest results for 'ookla' provider: 
Download speed: 36.76 MB/s
Upload speed: 57.92 MB/s

fast provider initialized
Speedtest results for 'fast' provider: 
Download speed: 47.91 MB/s
Upload speed: 70.64 MB/s
```

## LICENSE

[ISC License](https://github.com/konovenskki/turbo-snail/blob/main/LICENSE)