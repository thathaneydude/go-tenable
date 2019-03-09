# Go-Tenable

This is a simple library for interacting with Tenable's suite of products.  
* Tenable.sc
* Tenable.io
* Nessus

## Installation

    go get github.com/thathaneydude/go-tenable

## Examples

### Tenable.sc

```go
package main

import (
	"crypto/tls"
	"fmt"
	"github.com/thathaneydude/go-tenable"
	"net/http"
)

func main() {
	transport :=  &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	scClient := go_tenable.NewTenableSCClient("sc-address", transport)
	scClient.Login("sc-user", "sc-password")
	for _, asset := range scClient.ListAssets().Response.Usable {
		fmt.Printf("Asset %v (%v): %v\n", asset.Name, asset.ID, asset.Type)
	}

	scClient.Logout()
}
```

### Tenable.io
```go
package main

import (
	"crypto/tls"
	"fmt"
	"github.com/thathaneydude/go-tenable"
	"net/http"
	"time"
)

func main(){
	transport :=  &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	tio := go_tenable.NewTenableIOClient("access-key", "secret-key", transport)

	var Payload = go_tenable.AssetRequestBody{ChunkSize: 10000}
	export := tio.NewAssetExport()
	export.RequestAssetExport(Payload)
	for {
		status := export.RequestStatus()
		if status != "FINISHED" {
			time.Sleep(5)
		} else {
			break
		}
	}

	for _, asset := range export.DownloadAssetChunk(1) {
		fmt.Printf("Asset %v IPs: %v\n", asset.ID, asset.Ipv4S)
	}
}
```

### Nessus
```go
package main

import (
	"crypto/tls"
	"fmt"
	"github.com/thathaneydude/go-tenable"
	"net/http"
)

func main(){
	transport :=  &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	var nessusScanner = "nessus-address"
	nessusClient := go_tenable.NewNessusClient("access-key", "secret-key", nessusScanner, 8834, transport)
	status := nessusClient.GetStatus()
	properties := nessusClient.GetProperties()
	fmt.Printf("Scanner %v\n", nessusScanner)
	fmt.Printf("-Status: %v\n", status.Status)
	fmt.Printf("-Version: %v\n", properties.ServerVersion)
	fmt.Printf("-Licensed IP Count: %v\n", properties.License.Ips)
	fmt.Printf("-Licensed Expiration: %v\n", properties.License.ExpirationDate)
}
```
## Authors
Ryan Haney [@thathaneydude](https://twitter.com/thathaneydude)  
John Lampe [@f00dikator](https://twitter.com/f00dikator)   
Julian Davies 