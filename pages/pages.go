package pages

import (
	"errors"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/biehlerj/tldr/config"
)

const (
	sep = string(os.PathSeparator)
	ext = ".md"
)

// Read finds and creates the Page, if it does not find, simply returns abstract
// contribution guide
func Read(seq []string) (p *Page, err error) {
	page := ""
	for i, l := range seq {
		if len(seq)-1 == i {
			page = page + l
			break
		} else {
			page = page + l + "-"
		}
	}
	// Common pages are more, so we have better luck there
	p, err = queryCommon(page)
	if err != nil {
		p, err = queryOS(page)
		if err != nil {
			return p, errors.New("This page (" + page + ") doesn't exist yet!\n" +
				"Submit new pages here: https://github.com/tldr-pages/tldr")
		}
	}
	return p, nil
}

// Queries from common folder
func queryCommon(page string) (p *Page, err error) {
	d := config.SourceDir + sep + "pages" + sep + "common" + sep
	b, err := ioutil.ReadFile(d + page + ".md")
	if err != nil {
		return p, err
	}
	p = ParsePage(string(b))
	return p, nil
}

// Queries from os specific folder
func queryOS(page string) (p *Page, err error) {
	d := config.SourceDir + sep + "pages" + sep + config.OSName() + sep
	b, err := ioutil.ReadFile(d + page + ".md")
	if err != nil {
		return p, err
	}
	p = ParsePage(string(b))
	return p, nil
}

// QueryRandom brings a random page from the source
func QueryRandom() (p *Page, err error) {
	dCmn := config.SourceDir + sep + "pages" + sep + "common" + sep
	dOs := config.SourceDir + sep + "pages" + sep + config.OSName() + sep
	paths := []string{dCmn, dOs}
	sources := make([]string, 0)
	for _, p := range paths {
		files, err := ioutil.ReadDir(p)
		if err != nil {
			break
		}
		for _, f := range files {
			if strings.HasSuffix(f.Name(), ".md") {
				sources = append(sources, f.Name()[:len(f.Name())-3])
			}
		}
	}
	rand.Seed(time.Now().UTC().UnixNano())
	page := sources[rand.Intn(len(sources))]
	return Read([]string{page})
}

// ReadAll reads every single page from the source inta single page
func ReadAll() (p *Page, err error) {
	dCmn := config.SourceDir + sep + "pages" + sep + "common" + sep
	dOs := config.SourceDir + sep + "pages" + sep + config.OSName() + sep
	paths := []string{dCmn, dOs}
	p = &Page{Name: "Search All"}
	p.Tips = make([]*Tip, 0)
	for _, pt := range paths {
		files, err := ioutil.ReadDir(pt)
		if err != nil {
			break
		}
		for _, f := range files {
			if strings.HasSuffix(f.Name(), ".md") {
				page, err := Read([]string{f.Name()[:len(f.Name())-3]})
				if err != nil {
					continue
				}
				p.Tips = append(p.Tips, page.Tips...)
			}
		}
	}
	return p, nil
}
