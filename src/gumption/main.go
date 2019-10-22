package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/paidright/datalab/util"
)

var version = flag.Bool("version", false, "Just print the version and exit")
var quiet = flag.Bool("quiet", false, "Tone down the output noise")
var colsString = flag.String("columns", "", "A comma separated list of columns to target. Leaving this blank will operate on all columns")

var stripLeadingZeroes = flag.Bool("strip-leading-zeroes", false, "Strip leading zeroes")
var unquote = flag.Bool("unquote", false, "Strip quotation marks from strings")
var commasToPoints = flag.Bool("commas-to-points", false, "Replace all commas with full stops")
var addMissing = flag.String("add-missing", "", "String with which to replace blank fields")
var replaceCell = flag.String("replace-cell", "", "Take any cells that match X and replace it with Y eg: X,Y. You may specify multiple tuples, ie: A,B,X,Y")
var rename = flag.String("rename", "", "New name to assign to the column(s)")
var splitOnDelim = flag.String("split", "", "Delimiter on which to split the column(s)")
var cp = flag.Bool("copy", false, "Whether to copy the column(s)")
var drop = flag.Bool("drop", false, "Whether to drop the column(s)")

var columns []string

type replacement struct {
	from string
	to   string
}

type flagval struct {
	active       bool
	value        string
	replacements []replacement
}

var logger = util.Logger{}

func main() {
	flag.Parse()
	columns = strings.Split(*colsString, ",")

	if *version {
		logger.Info(currentVersion)
		os.Exit(0)
	}

	flags := map[string]flagval{
		"stripLeadingZeroes": flagval{
			active: *stripLeadingZeroes,
		},
		"unquote": flagval{
			active: *unquote,
		},
		"commasToPoints": flagval{
			active: *commasToPoints,
		},
		"addMissing": flagval{
			active: *addMissing != "",
			value:  *addMissing,
		},
		"replaceCell": flagval{
			active:       *replaceCell != "",
			replacements: parseReplacements(*replaceCell),
		},
		"rename": flagval{
			active: *rename != "",
			value:  *rename,
		},
		"splitOnDelim": flagval{
			active: *splitOnDelim != "",
			value:  *splitOnDelim,
		},
		"cp": flagval{
			active: *cp,
		},
		"drop": flagval{
			active: *drop,
		},
	}

	for k, flag := range flags {
		if flag.active {
			logger.Info("flag", k, "is set")
		}
	}

	output := csv.NewWriter(os.Stdout)

	if err := gumption(os.Stdin, *output, columns, flags); err != nil {
		logger.Fatal(err)
	}

	output.Flush()

	logDone()
}

func gumption(input io.Reader, output csv.Writer, columns []string, flags map[string]flagval) error {
	cachedHeaders := []string{}

	handleHeaders := func(headers []string) ([]string, error) {
		if len(cachedHeaders) > 0 {
			return cachedHeaders, nil
		}

		cachedHeaders = append([]string{}, headers...)

		if len(columns) == 0 {
			columns = cachedHeaders
		}

		if flags["rename"].active {
			if len(columns) > 1 {
				return []string{}, fmt.Errorf("Can only rename one column at at a time")
			}
			if len(columns) == 0 {
				return []string{}, fmt.Errorf("Cannot rename without setting a single target column")
			}
			for i, header := range cachedHeaders {
				if header == columns[0] {
					cachedHeaders[i] = flags["rename"].value
				}
			}
		}

		if flags["splitOnDelim"].active {
			for _, col := range columns {
				cachedHeaders = append(cachedHeaders, suffixed(col, columns, 1))
			}
		}

		if flags["cp"].active {
			for _, col := range columns {
				cachedHeaders = append(cachedHeaders, suffixed(col, columns, 1))
			}
		}

		if flags["drop"].active {
			newHeaders := []string{}
			for _, header := range cachedHeaders {
				shouldDrop := false
				for _, col := range columns {
					if col == header {
						shouldDrop = true
					}
				}
				if !shouldDrop {
					newHeaders = append(newHeaders, header)
				}
			}
			cachedHeaders = newHeaders
		}

		if err := output.Write(cachedHeaders); err != nil {
			return []string{}, err
		}
		output.Flush()

		return cachedHeaders, nil
	}

	work, errors := util.ReadSourceAsync(input)

	for line := range work {
		headers, err := handleHeaders(line.Headers)
		if err != nil {
			return fmt.Errorf("Error handling headers %w", err)
		}

		for _, col := range columns {
			cell := line.Data[col]
			if flags["stripLeadingZeroes"].active {
				cell = strings.TrimLeft(cell, "0")
			}

			if flags["unquote"].active {
				cell = strings.Trim(cell, `"`)
				cell = strings.Trim(cell, `'`)
			}

			if flags["commasToPoints"].active {
				cell = strings.ReplaceAll(cell, ",", ".")
			}

			if flags["addMissing"].active {
				if cell == "" {
					cell = flags["addMissing"].value
				}
			}

			if flags["replaceCell"].active {
				for _, rep := range flags["replaceCell"].replacements {
					if cell == rep.from {
						cell = rep.to
					}
				}
			}

			line.Data[col] = cell

			if flags["rename"].active {
				line.Data[flags["rename"].value] = cell
			}

			if flags["splitOnDelim"].active {
				parts := strings.SplitN(cell, flags["splitOnDelim"].value, 2)
				if len(parts) > 1 {
					line.Data[col] = parts[0]
					line.Data[suffixed(col, columns, 1)] = parts[1]
				}
			}

			if flags["cp"].active {
				line.Data[suffixed(col, columns, 1)] = cell
			}
		}

		newLine := []string{}
		for _, header := range headers {
			newLine = append(newLine, line.Data[header])
		}
		if err = output.Write(newLine); err != nil {
			return err
		}
		output.Flush()
	}

	var cachedErr error

	for err := range errors {
		log.Println("ERROR", err)
		cachedErr = err
	}

	return cachedErr
}

func suffixed(target string, cols []string, i int) string {
	candidate := fmt.Sprintf("%s_%d", target, i)
	for _, col := range cols {
		if col == candidate {
			return suffixed(target, cols, i+1)
		}
	}
	return candidate
}

func parseReplacements(input string) []replacement {
	parts := strings.Split(input, ",")
	replacements := []replacement{}
	for i, part := range parts {
		if (i+1)%2 == 0 {
			rep := replacement{
				from: parts[i-1],
				to:   part,
			}
			replacements = append(replacements, rep)
		}
	}
	return replacements
}

func logDone() {
	if *quiet {
		return
	}
	logger.Info(`
                                     WWWWWWWWNNNNNNNNNNNNNNNNNNNNNNNNNNNNNWWWWWWWWWWWWW
                   WWWWNNNNXXXXKKKK0000OOOOOOkkkkkkkkkkkkkkkkkkkkkkkkOOOOOOOOOOOOO0000000KKKXXXXNNNNNWWWWW
         WWWNNXXK000OOOkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkOOOOO000KKKXNNWWWW
    WNNXK0OOkkxxxxxxkkkkkkkkOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOkkkkkkkkkxxxxkkkOO00KXXNWW
  WX0OkkxxxxkkkkkOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO0OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOkkkkkxxxxxxxkkO0KXNW
 NK0OxxxxkOOOOOOOOOO00OO00000000000000000000000000000000000000000000000000000000000000OOOOOOO00000000000000OOOOOOOOkkkxxxxxkkOKXW 
WKO00OOOO0OOOOO000000OOOOO000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000OOOOkkxddkOOX 
XOkkkOOO00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000OOOOOkkkkk0W
KkkkkxxxkkkOOOO000000000000000000000000000000000000000000K000000000000000000000000000000000000000000000000OO000000000OOOOkkkxxxxkX
0kkkkkkkkkxxxxkkkkkkOOOOO000000000000000000000000000000KKK00000000000000000000000000KKKK00000000000000000OOOOOOOOOOkkkxxxxxxxxxxkX
KkxxxxxkkkkkkkkkkkkkkkkkkkkkkkkkkkOOOOOOOOOO00O00000000000000000000000000000000000000000000OOOOOOOkkkkkxxxxxxxxxxxxxxxxxxxxxxxxxkX
WKOkxxxxxxxxxxkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkxxkkkkkkkkkkkkkxxxxxxxxxxxxxxxxxxxxxkkkkkkkkkxxxxddddxxk0W
  N0OOkkxxxxxxxxxxxxxxxkkkkkkkkkkkkkkkkkOOOOOOOOOOkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkxxxxxxxxxdddddxxxkOXW
  WKKK00OOOkkkkkxxxxxxxxxxxxxxxxxxxxxxxxxxkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkxxxxxxxxxxxxxxxxxxxddddddxxxxxxxxkkkO0N
  WXKKKKKKK0000OOOOOkkkkkkkkkxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxddddddddxxxdddddddddddddddddxxxxxxxxxxxxxxxxxkkkkkkOOO000KW
   XKKKKKKKKKKK00000000000OOOOOOOkOkkkkkkkkkkkkkkkxxxxxxxxxkxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxkkkkkkkkkkkOOOOOOO000000000XW
   NKKXXXKKKKKKK00KKKKKKKKKKK000000000000000OOOOOOOOOOOOOOOOOOkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkOOOOOOOOOOOOOO0O0000000000000000KN
   WKKKKXKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKK00000000000000000000000000000O00OOOOOOOOOOOOOOOOOOOO0OOO00000000000000000000000KKKXNW
    WXKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKK0000000000000000000000OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO0000000000000000KKKKKKKN
     WXKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKK0000000000000000000OOOOOOOOOOOOOOOOOOOOOOOOOOOO00000000000000000000KKKK0000KKKN
     WXKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKK00000000000000000000000000000000000000000000000KKK00000000000000000K0KN
     WNXKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKK000000000000000000000000000000000000000000000XW
      NXX0dxKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKK0000KKKKKKKKKKKKKKKKKKKKKKKKKK00000000000000000000000000000000000000000000KKX
      NXKd;:kKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKK0000000KKKKKKKK00000000000000000000000000000000000000000000000OO000000000000KKN
      NX0o;,:kKKKKKKKKKKKKKKKKKK000000000000000000000000000000000000000000000000000000OOOOOOOOOOOO00000000000OkOxx00000KKKN
      WX0dolclOK0Okxdooolllcccccc::::ccccccllcccllllooooodddxxxkkkkkOOOOOOOOOOOOOOOOOOOOOOOOOO000000000000OOx:;okk00000KKKN
      WX0lcolcdkc,''.......................;,................'''',,,;;:::cclloodxkkOO000000000000000000Oo:;cc',d000000KK0KW
      WN0ococ';d:.'....................''.'lc..................................'',;:cloddxkOO000OxolldOx,..,,.;k000000KKOKW
      WXOdlc:';dc......................':::xd;:c,.........''''''''',,;;::clloooddddddddxxkkxk0kl;'...'oo'..''.cO000000K0kK
       Kkdlcc::oc                 ',;;::cxO00Okl,,,,,;;;::cllooodddxxxkkkkkOOOOOOOOOOkkOxc;;co;  ;:  'oc     'o000000K0kkX
       Kkxlcc:;ol'                ';:lodxk0KKK0OkxdddddxxxkkkkkkkkkkkkkkOOOO0OOkxolc;,;dc   ::  'oc  ,o:     ,d00000K0kdkN
       Xkxl:c:;ll'          '',;:c:::cllok0KK00OOOOOOkkkkkkkkOOOO0000OOOOO00kc,'    ';cd:  ':;  ,o:  ;l, ,,  ;x0O0000kooOW
       NOkOkdl;cl,''   ',,;;:clokOkxxkOO0OO00OOOOOkkOOOO000OO00Odllc:;;,,;cxx:,,   'lk0k;  ':,  ;o;  :l' :c;:oOO0000kocoK
       NOOKOk0kxo,',,;::clodxxkOO00OOOOOOOO00OOOOOkO000xlc:;:oOd'    ,,'   ;xOOl   ,d00x,  ,;'  :l' ,ol,;dkkkOOO000kl:cdX
       W0kKOxk0Kkc:ccodxkO000000OOOOkkOOkdddO0xc:;,;oOk;    'cko    'ld;   :x0Oc   ;x00d'  ,;   '' ,oOxdkOOOOOOO0Oxc;:cxN
        0kOOkOKX0doxk0000OdlccoxO0x:;;lkl' 'lkl     ;do'    'ckl     ;;  ':d00O:   :k00l   :o:'',;coddccxOOOOOO0ko:;;:lOW
        KkkxxkOO00KKK000Ol' '' 'ckx,  :xl  ':xc     ':;     'ckl     '';cdk000k;   ck0Oc'';oOOdolcc::c;:xOOOO0Odc;,,;:l0W
        XkkkxdooOKKKKK0Kk: ':l, 'od,  :xl   :dc  ''  '  ''  'lkc    ,okO00000Kx,  'lO00xodxlol,;:;:;cooxkOO0Oxl;'',;;:dK
        NOkkxxxdkKKKKK0KO: 'cxoclxx,  ;xl   ;d:  ',    ','  'lk:    ;x0000000KOolllddxkl:ol,;;';clodxkOOO0Oxl,''',,;:ckN
        W0OOkxxdkKKKKKKK0c ':;,',od;  ;xl   ;o:  ';,   ,:,  'lx:   ':xK000OOxdk0kccl;lx::lc:cldxkkOOOOO0Odc,   '',,;:oKW
        WKOkkkkxkKXKKKKX0c 'cl, 'cd;  ;dl  ':d:  'cc' ':o;',:dOxllloddoooccl;;d0o;ll:colodxkkkkkkkOO0Oxl;'    ''',,;lON
         KOOkkkdxKXKKKKX0l'':l, 'cxc  ';, ';dkl,,:okdllxkxooooodl;::::;;::::,;d0klldxxkkkkkkkkOOOOkxl;'        '',;lOKN
         XOOxxkxx0XKKKKKKk:''',,,cOOoc;;:clxO0xoddkK0dcc:cc;,,,:; ;ol;;cllldddkOOOOkkkkkkkOOOOOkdc;,'         ''';d0KXW
         X00kxxdd0XKKKKKXK0xddxOOOK0xxolool:okl:oclkx:;l:;;;::clolxKOxxO00OOOOkkkkkOOOOOOOkxoc:,..,ll,........',lkKKKXW
         NK0kdddx0XKKKKXXXXXXXXKKKXx;,,;:::,ckolxlokxodKOddkOO0000000OOOkkkkkkkOOOO000Oxoc,'......:oc'.......;lk0KKKKXW
         WXK0xllxKXKKKKKXXXXXXXXXKXkcclooloddxkkOO0KK0000000OOOkkkkkkkkkkOOkkkkxddxkkl;'..................,cdOKKKKKKKN
          NKKOdd0KKKKKKKKKKXXXXXXKXK00KKKKKKKK0OOkkkkOOkkkkkkkxxxxxxdddoolc::;,''',c:.................';cdk0KKKKKKKKKN
          WXKKO0KKKKKKKKKKKXXXXXXXXXXKKKKKKKKKK0Okkxdolc:;;;;;,,,,,'''.............;,............',:ldkO00OOO0KKKKKKXW
          WNKKKKKKKKKKKKKKKKXXXXXXKKXXKKKXKKKKKKKKKKK0Okxdolc::;,,,''...................'',;;:cldkOOkkxxkkxkO0K0KKKKN
           WKKKKKKKKKKKKKKKKKXXXXXXKXKkdxxxxxkxxkkkOkOOOO0K000Okxxdddooolccccllcllcc:clc:cclllodOK0Oxxkkk00000KKKKKXW
            WXKKKKKKKKKKKKKKKXXXKXXKKKOxxddxddlloolllllclxOdxkllolllllllc:clcxkdxkoolloooooddxkO0K0000KKKKKKKKKKKKKN
             WNKKKKKKKKKKXXXXKXKKKKKKKKXXXKKKKK00000OOOOO00000OkkOOOOkkOkkkOOO0000OO00000KKK0000KKKKKKKKKKKKKKKKKXN
               WNXXKKKKXXXXXXXXXXXXXXXKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKK0000000KKKKKKKKKKKKKKKKKKKKKKKKKKKKKKXNW
                   WNNXXKKKXXXXXXXXXXXXXXXXXXXKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKXKKKKKXXNW
                       WWNNNXXKKKXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXKKKKKKKKKKKKKKKKXKXXKXXKKXXXKKKKKXXNNWW
                             WWNXXXXXKKKKKKXXXXXXXXXXXXXXXXXXXXXXXXXXKKKKKKKKKXXXXXXXXXKKKKKKXXXXNNNWW
                                    WWWNNNNNNNXXXXXXXXXXXXXXKKKKKKKKKKKKKKKXXXXXXXXXXNNNNNNWWWW
                                                   WWWWWWWWWNNNNNNNNNNNNWWWWWWWWW`)
}
