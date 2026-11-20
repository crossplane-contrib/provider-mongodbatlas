package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/crossplane/upjet/v2/pkg/pipeline"
	"github.com/mongodb/terraform-provider-mongodbatlas/xpshim"

	"github.com/crossplane-contrib/provider-mongodbatlas/config"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] == "" {
		panic("root directory is required to be given as argument")
	}
	rootDir := os.Args[1]
	absRootDir, err := filepath.Abs(rootDir)
	if err != nil {
		panic(fmt.Sprintf("cannot calculate the absolute path with %s", rootDir))
	}
	sdk, fw := xpshim.GetSDKProvider(), xpshim.GetFrameworkProvider()
	pipeline.Run(config.GetProvider(sdk, fw), config.GetProviderNamespaced(sdk, fw), absRootDir)
}
