package main

import (
	"context"
	"flag"
	"log"
	"terraform-provider-ome/ome"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --rendered-website-dir docs --provider-name terraform-provider-ome
// var (
// 	version string = "dev"
// )

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	err := providerserver.Serve(context.Background(), ome.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/dell/ome",
		Debug:   debug,
	})

	if err != nil {
		log.Fatal(err.Error())
	}
}
