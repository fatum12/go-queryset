package main

import (
	"context"
	"flag"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/jirfag/go-queryset/internal/parser"
	"github.com/jirfag/go-queryset/internal/queryset/generator"
)

func main() {
	inFile := flag.String("in", "models.go", "path to input file")
	outFile := flag.String("out", "", "path to output file")
	timeout := flag.Duration("timeout", time.Minute, "timeout for generation")
	flag.Parse()

	if *outFile == "" {
		inExt := filepath.Ext(*inFile)
		*outFile = strings.TrimSuffix(*inFile, inExt) + "_queryset" + inExt
	}

	g := generator.Generator{
		StructsParser: &parser.Structs{},
	}

	ctx, finish := context.WithTimeout(context.Background(), *timeout)
	defer finish()

	if err := g.Generate(ctx, *inFile, *outFile); err != nil {
		log.Fatalf("can't generate query sets: %s", err)
	}
}
