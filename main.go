package main

import (
	"context"
	"flag"
	"log"
	"terraform-provider-ome/ome"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var (
	version string = "dev"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	err := providerserver.Serve(context.Background(), ome.New(version), providerserver.ServeOpts{
		Address: "registry.terraform.io/hashicorp/ome",
		Debug:   debug,
	})

	if err != nil {
		log.Fatal(err.Error())
	}
}
