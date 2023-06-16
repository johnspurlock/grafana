package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/grafana/thema"
	"github.com/grafana/thema/vmux"

	"github.com/grafana/grafana/pkg/cuectx"
	"github.com/grafana/grafana/pkg/kinds/dashboard"
)

func main() {
	var inputFile, outputFile string
	flag.StringVar(&inputFile, "in", "dashboard.json", "input json file to schematize")
	flag.StringVar(&outputFile, "out", "", "output file (default: sdtout)")
	flag.Parse()

	src, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("opening input file failed: %s", err.Error())
		os.Exit(1)
	}

	dashdata, err := vmux.NewJSONCodec(inputFile).Decode(cuectx.GrafanaCUEContext(), src)
	if err != nil {
		fmt.Printf("decoding input file failed: %s", err.Error())
		os.Exit(1)
	}

	dk, err := dashboard.NewKind(cuectx.GrafanaThemaRuntime())
	if err != nil {
		fmt.Printf("creating dashboard kind failed: %s", err.Error())
		os.Exit(1)
	}

	var i *thema.Instance
	// try to validate against the latest lineage. for this example we're going
	// to continue through and run ValidateAny regardless of the outcome.
	_, err = dk.Lineage().Latest().Validate(dashdata)
	if err != nil {
		fmt.Printf("validation against latest schema failed: %s", err.Error())
	}

	i = dk.Lineage().ValidateAny(dashdata)
	if i == nil {
		fmt.Println("ValidateAny returned a nil instance")
		os.Exit(1)
	}

	if i.Schema().Version() != i.Schema().LatestInMajor().Version() {
		i, _ = i.Translate(i.Schema().LatestInMajor().Version())
	}

	j, err := i.Underlying().MarshalJSON()
	if err != nil {
		fmt.Printf("marshaling output failed: %s", err.Error())
		os.Exit(1)
	}
	if outputFile == "" {
		fmt.Printf("writing to stdout\n")
		os.Stdout.Write(j)
	} else {
		fmt.Printf("writing to %s\n", outputFile)
		os.WriteFile(outputFile, j, 0644)
	}
}
