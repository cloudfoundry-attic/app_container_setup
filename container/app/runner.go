package main

import (
	"io/ioutil"
	"os"

	"github.com/cloudfoundry/app_container_setup/container"
)

func main() {
	inputJson, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic("failed to read")
	}

	err = container.Main(inputJson)
	if err != nil {
		panic("Main failed")
	}
}
