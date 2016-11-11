//
// sitemapper generates a sitemap.xml file by crawling the content generate with caitpage
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
//
// Copyright (c) 2016, Caltech
// All rights not granted herein are expressly reserved by Caltech.
//
// Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	// Caltech Library packages
	"github.com/caltechlibrary/cait"
)

type locInfo struct {
	Loc     string
	LastMod string
}

var (
	showHelp    bool
	showVersion bool
	changefreq  string
	locList     []*locInfo
)

func usage(appName, version string) {
	fmt.Println(`
 USAGE: %s [OPTIONS] HTDOCS_PATH MAP_FILENAME PUBLIC_BASE_URL

 OVERVIEW

 Generates a sitemap for the accession pages.

`, appName)
	flag.VisitAll(func(f *flag.Flag) {
		fmt.Printf("\t-%s\t%s\n", f.Name, f.Usage)
	})
	fmt.Printf(`
 EXAMPLE

    %s htdocs htdocs/sitemap.xml http://archives.example.edu

`, appName)
	fmt.Printf("\n%s %s\n", appName, version)
	os.Exit(0)
}

func init() {
	flag.BoolVar(&showHelp, "h", false, "display this help message")
	flag.BoolVar(&showHelp, "help", false, "display this help message")
	flag.BoolVar(&showVersion, "v", false, "display version info")
	flag.BoolVar(&showVersion, "version", false, "display version info")
	flag.StringVar(&changefreq, "c", "daily", "Set the change frequencely value, e.g. daily, weekly, monthly")
}

func main() {
	appName := path.Base(os.Args[0])
	flag.Parse()
	args := flag.Args()
	if showHelp == true {
		usage(appName, cait.Version)
	}
	if showVersion == true {
		fmt.Printf("%s %s\n", appName, cait.Version)
		os.Exit(0)
	}
	if len(args) != 3 {
		fmt.Printf("%s requires 3 parameters, see %s --help", appName, appName)
		os.Exit(1)
	}
	if changefreq == "" {
		changefreq = "daily"
	}

	log.Printf("Starting map of %s\n", args[0])
	filepath.Walk(args[0], func(p string, info os.FileInfo, err error) error {
		if strings.HasSuffix(p, ".html") {
			fname := path.Base(p)
			//NOTE: You can skip the eror pages in the sitemap
			if strings.HasPrefix(fname, "50") == false && strings.HasPrefix(p, "40") == false {
				finfo := new(locInfo)
				finfo.Loc = fmt.Sprintf("%s%s", args[2], strings.TrimPrefix(p, args[0]))
				yr, mn, dy := info.ModTime().Date()
				finfo.LastMod = fmt.Sprintf("%d-%0.2d-%0.2d", yr, mn, dy)
				log.Printf("Adding %s\n", finfo.Loc)
				locList = append(locList, finfo)
			}
		}
		return nil
	})
	log.Printf("Writing %s\n", args[1])
	fp, err := os.OpenFile(args[1], os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0664)
	if err != nil {
		log.Fatalf("Can't create %s, %s\n", args[1], err)
	}
	defer fp.Close()
	fp.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
`))
	for _, item := range locList {
		fp.WriteString(fmt.Sprintf(`
    <url>
            <loc>%s</loc>
            <lastmod>%s</lastmod>
            <changefreq>%s</changefreq>
    </url>
`, item.Loc, item.LastMod, changefreq))
	}
	fp.Write([]byte(`
</urlset>
`))
	log.Printf("Writing %s complete\n", args[1])
}
