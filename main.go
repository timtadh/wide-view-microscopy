package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

import (
	"github.com/timtadh/getopt"
)

import (
	"github.com/timtadh/wide-view-microscopy/ingest"
	"github.com/timtadh/wide-view-microscopy/charts"
)


var UsageMessage string = "wide-view-microscopy --help"
var ExtendedMessage string = `
wide-view-microscopy -d <path> -o <out.html> \
                     -f '$(slide) $(subject) $(region) $(stain).tif'

+---------+
| Options |
+---------+

-h, --help                          view this message
-d, directory=<path>                the directory where the imanges are stored
-o, output=<path>                   output for the html
                                    (optional will go to stdout)
-f, format=<format-string>          a format for the names of the images
                                    default: '$(slide) $(subject) $(region) $(stain).tif'
-r, row-group=<vars>                variables to group row on
                                    default: 'region'
-c, chart-group=<vars>              variables to group charts on
                                    default: 'subject,slide'
-s, column-sort=<vars>              variables to sort columns on
                                    default: 'stain'

+-------+
| Specs |
+-------+

<format-string>     See below
<path>              A file system path
<vars>              A comma separated list of variables.

+---------------+
| Format Fields |
+---------------+

$(slide)    int     the slide identifier. should be a number (required)
$(subject)  string  the identifier of the subject the sample was taken from
                    (required)
$(region)   string  the identifier for the region of the slide the image is
                    from (required)
$(stain)    string  the stain type which was used for this image (required)


+----------------+
| Format Strings |
+----------------+

The format strings consist of utf-8 characters and variables enclosed in $(-). The

tiff: '$(slide) $(subject) $(region) $(stain).tif'
png:  '$(slide) $(subject) $(region) $(stain).png'
jpeg: '$(slide) $(subject) $(region) $(stain).jpg'
`

func Usage(code int) {
	fmt.Fprintln(os.Stderr, UsageMessage)
	if code == 0 {
		fmt.Fprintln(os.Stdout, ExtendedMessage)
	}
	os.Exit(code)
}

func Vars(str string) []string {
	split := strings.Split(str, ",")
	vars := make([]string, 0, len(split))
	for _, s := range split {
		s = strings.TrimSpace(s)
		if s != "" {
			vars = append(vars, s)
		}
	}
	return vars
}

func main() {
	args, optargs, err := getopt.GetOpt(
		os.Args[1:],
		"hl:d:o:f:s:r:c:",
		[]string{ "help", "directory=", "output=", "format=",
		          "column-sort=", "row-group=", "chart-group=" },
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error parsing command line flags", err)
		Usage(1)
	}
	if len(args) > 0 {
		fmt.Fprintf(os.Stderr, "unexpected trailing args `%v`\n", strings.Join(args, " "))
		Usage(1)
	}
	
	format, err := ingest.ParseFormatString("$(slide) $(subject) $(region) $(stain).tif")
	if err != nil {
		log.Fatal(err)
	}
	directory := ""
	rowGroup := Vars("region")
	chartGroup := Vars("subject,slide")
	columnSort := Vars("stain")
	for _, oa := range optargs {
		switch oa.Opt() {
		case "-h", "--help":
			Usage(0)
			os.Exit(0)
		case "-f", "--format":
			format, err = ingest.ParseFormatString(oa.Arg())
			if err != nil {
				fmt.Fprintf(os.Stderr, "Invalid format string (%v) '%v'\n", oa.Opt(), oa.Arg())
				fmt.Fprintln(os.Stderr, err)
				Usage(1)
			}
		case "-d", "--directory":
			directory, err = filepath.Abs(oa.Arg())
			if err != nil {
				fmt.Fprintf(os.Stderr, "Bad directory (%v) '%v' supplied\n", oa.Opt(), oa.Arg())
				Usage(1)
			}
			if _, err := os.Stat(directory); err != nil && os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "Bad directory (%v) '%v' supplied\n", oa.Opt(), oa.Arg())
				Usage(1)
			} else if err != nil {
				fmt.Fprintln(os.Stderr, err)
				fmt.Fprintf(os.Stderr, "Bad directory (%v) '%v' supplied\n", oa.Opt(), oa.Arg())
				Usage(1)
			}
		case "-s", "--column-sort":
			columnSort = Vars(oa.Arg())
		case "-r", "--row-group":
			rowGroup = Vars(oa.Arg())
		case "-c", "--chart-group":
			chartGroup = Vars(oa.Arg())
		default:
			fmt.Fprintf(os.Stderr, "Unknown flag '%v'\n", oa.Opt())
			Usage(1)
		}
	}

	log.Println(directory)

	files, err := ingest.Ingest(directory, format)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(files)
	for _, img := range files {
		log.Println(img)
	}

	C := charts.MakeCharts(files, chartGroup, rowGroup, columnSort) 
	for _, chart := range C {
		log.Println("chart", chart.Meta())
		for _, row := range chart.Rows() {
			log.Println("row", row.Meta())
			for _, img := range row.Images() {
				log.Println(img)
			}
			log.Println()
		}
		log.Println()
	}

	log.Println("done")
	fmt.Println(charts.ChartsHTML(C))
}

