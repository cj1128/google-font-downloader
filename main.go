/*
* @Author: CJ Ting
* @Date: 2016-12-24 10:01:36
* @Email: fatelovely1128@gmail.com
 */

package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/briandowns/spinner"

	"gopkg.in/alecthomas/kingpin.v2"
)

const fontsURL = "https://fonts.googleapis.com/css"

var appVersion string // set by linker with -ldflags

var httpClient = http.Client{
	Timeout: 30 * time.Second,
}

var (
	input  string
	output string
	lang   string
	format string // "eot", "woff", "woff2", "svg", "ttf"
	fonts  []string
)

// we are gonna use these UAs to trick Google Fonts Server
// into serving us the correct css file
var formatUA = map[string]string{
	"eot":   "Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 6.1; Trident/4.0)",
	"woff":  "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:27.0) Gecko/20100101 Firefox/27.0",
	"woff2": "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:40.0) Gecko/20100101 Firefox/40.0",
	"svg":   "Mozilla/4.0 (iPad; CPU OS 4_0_1 like Mac OS X) AppleWebKit/534.46 (KHTML, like Gecko) Version/4.1 Mobile/9A405 Safari/7534.48.3",
	"ttf":   "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/534.54.16 (KHTML, like Gecko) Version/5.1.4 Safari/534.54.16",
}

func parseFlags() {
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Flag("input", "specify Google Fonts css file").
		Short('i').
		StringVar(&input)
	kingpin.Flag("output", "specify output css file").
		Short('o').
		Default("output.css").
		StringVar(&output)
	kingpin.Flag("lang", "specify language subsets, comma separated").
		Short('l').
		Default("latin").
		StringVar(&lang)
	kingpin.Flag("format", "specify webfont formats, can be eot,woff,woff2,svg,ttf").
		Short('f').
		Default("woff").
		EnumVar(&format, "eot", "woff", "woff2", "svg", "ttf")

	kingpin.Arg("fonts", "specify font specs").
		StringsVar(&fonts)

	kingpin.Version(appVersion)

	kingpin.Parse()
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
	}()
	parseFlags()
	if input != "" {
		buf, err := ioutil.ReadFile(input)
		if err != nil {
			panic(err)
		}
		transformStylesheet(buf)
	} else {
		if len(fonts) == 0 {
			kingpin.Usage()
			os.Exit(0)
		}
		css := getStylesheet()
		transformStylesheet(css)
	}
}

func transformStylesheet(buf []byte) {
	r := regexp.MustCompile(`@font-face \{[\s\S]+?font-family: '(.*?)'[\s\S]+?font-style: (.*?);[\s\S]+?font-weight: (\d+)[\s\S]+?url\((.*?)\)[\s\S]+?\}`)
	fontDownloaded := make(map[string]bool)
	for _, m := range r.FindAllSubmatch(buf, -1) {
		fontName := string(m[1])
		fontStyle := string(m[2])
		fontWeight := string(m[3])
		fontURL := string(m[4])
		fontFilename := fmt.Sprintf("%s_%s_%s.%s", slugify(fontName), fontWeight, fontStyle, format)

		// in browsers that support unicode-range, stylesheet may
		// contain multiple @font-face for same weight
		if !fontDownloaded[fontFilename] {
			fontDownloaded[fontFilename] = true
			downloadFont(fontURL, fontFilename)
		}
		buf = bytes.Replace(buf, []byte(fontURL), []byte("'"+fontFilename+"'"), -1)
	}

	f, err := os.Create(output)
	if err != nil {
		panic(err)
	}
	f.Write(buf)
}

// download individual font
func downloadFont(url, fontFilename string) {
	s := spinner.New(spinner.CharSets[26], 300*time.Millisecond)
	s.Prefix = fmt.Sprintf("download %s", fontFilename)
	s.FinalMSG = fmt.Sprintf("finish downloading %s", fontFilename)
	s.Start()
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := httpClient.Do(req)
	if err != nil {
		panic(fmt.Errorf("network error: %s", err))
	}
	defer resp.Body.Close()
	f, err := os.Create(fontFilename)
	if err != nil {
		panic(err)
	}
	io.Copy(f, resp.Body)
	s.Stop()
	fmt.Println()
}

// request Google to get corresponding
// stylesheet
func getStylesheet() []byte {
	query := url.Values{}
	query.Add("subset", lang)
	query.Add("family", strings.Join(fonts, "|"))

	finalURL := fontsURL + "?" + query.Encode()
	fmt.Println("request", finalURL)

	req, _ := http.NewRequest("GET", finalURL, nil)
	req.Header.Set("User-Agent", formatUA[format])
	resp, err := httpClient.Do(req)
	if err != nil {
		panic(fmt.Errorf("network error: %s", err))
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		panic(fmt.Errorf("can't get fonts, please check your font specs"))
	}
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return buf
}

func slugify(str string) string {
	spacereg := regexp.MustCompile(`\s+`)
	return spacereg.ReplaceAllString(str, "-")
}
