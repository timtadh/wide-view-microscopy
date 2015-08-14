package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

import (
	"github.com/timtadh/getopt"
)

import (
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
                                    '$(slide) $(subject) $(region) $(stain).tif'

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

func main() {
	args, optargs, err := getopt.GetOpt(
		os.Args[1:],
		"hl:d:o:f:",
		[]string{ "help", "directory=", "output=", "format="},
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error parsing command line flags", err)
		Usage(1)
	}
	if len(args) > 0 {
		fmt.Fprintf(os.Stderr, "unexpected trailing args `%v`\n", strings.Join(args, " "))
		Usage(1)
	}

	for _, oa := range optargs {
		switch oa.Opt() {
		case "-h", "--help":
			Usage(0)
			os.Exit(0)
		default:
			fmt.Fprintf(os.Stderr, "Unknown flag '%v'\n", oa.Opt())
			Usage(1)
		}
	}

	log.Println("done")
}

