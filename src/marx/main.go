package main

import (
	"encoding/csv"
	"flag"
	"log"
	"os"
	"path"
	"runtime"
	"strconv"
	"sync"

	"github.com/paidright/datalab/util"
)

var version = flag.Bool("version", false, "Just print the version and exit")
var input = flag.String("input", ".", "The directory to read source CSVs from")
var quiet = flag.Bool("quiet", false, "Tone down the output noise")

func main() {
	flag.Parse()
	if *version {
		log.Println(currentVersion)
		os.Exit(0)
	}

	files, err := util.ListFiles(*input, []string{}, []string{".csv"})
	if err != nil {
		log.Fatal(err)
	}

	for i, file := range files {
		files[i] = path.Join(*input, file)
	}

	if len(files) == 0 {
		log.Fatal("ERROR no files found in input directory")
	}

	log.Println("INFO unionising the following files:")
	for _, file := range files {
		log.Println(file)
	}

	headers, err := enumerateHeaders(files)
	if err != nil {
		log.Fatal(err)
	}

	headers = append(headers, "original_file_name", "original_row_number")

	output := csv.NewWriter(os.Stdout)

	if err := output.Write(headers); err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if err := processFile(file, headers, output); err != nil {
			log.Fatal(err)
		}
	}

	logDone()
}

func processFile(file string, headers []string, output *csv.Writer) error {
	log.Println("INFO working on file:", file)

	work, errors := util.ReadFileAsync(file)

	mutex := sync.Mutex{}

	workers := sync.WaitGroup{}
	for _, _ = range make([]bool, runtime.GOMAXPROCS(0)) {
		workers.Add(1)
		go (func() {
			for line := range work {
				if line.Number%100000 == 0 {
					log.Printf("INFO marx up to line number: %+v \n", line.Number)
				}
				record := []string{}
				for _, col := range headers {
					if col == "original_file_name" {
						line.Data[col] = file
					}
					if col == "original_row_number" {
						line.Data[col] = strconv.Itoa(line.Number)
					}
					record = append(record, line.Data[col])
				}

				mutex.Lock()
				if err := output.Write(record); err != nil {
					log.Fatal(err)
				}
				mutex.Unlock()
			}
			workers.Done()
		})()
	}

	workers.Wait()
	output.Flush()

	var cachedErr error

	for err := range errors {
		log.Println("ERROR", err)
		cachedErr = err
	}

	return cachedErr
}

func enumerateHeaders(names []string) ([]string, error) {
	cols := []string{}

	for _, name := range names {
		headers, err := util.ReadHeaders(name)
		if err != nil {
			return cols, err
		}

		cols = append(cols, headers...)
	}

	return util.Uniq(cols), nil
}

func logDone() {
	if *quiet {
		return
	}
	log.Println(`
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%&@@@@@@@@@@@@@@@&%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%&@@@@@@@@@@@@&&&@%@@@@@@@@@%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%@@@@@@@&*&,           ,(@@@@@@@@%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%&@@@@@@%%..      ./*#%,,/* *(#%@@@@@@%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%&@@@@@@@@/#         *#(/#,,(,&@&/ ,#@@@@@@@@@@&%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%&@@@@@@@%/ .      .              .@@@@&#%@@@@@@@@@@@@&%%%%%%%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%@@@@@@*(,    #&&&&@&%(,             .     (%&&%/* .(@@@@@(@@%%%%%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%&@@@@/       ,(&@&                              ,/#%@@@/ ,% #@@%%%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%&@@@**        /@@/                                   ,*&@@@@/ /&/*@@&%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%%&@@*         *@%.                                       .%@@@@&. #,*&@@%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%@@&.   /    *@@@.                                         .@@@(@@   %.@@@%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%@@%  .*     &@@%*                                           .%@@*@/  ,**#@@%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%@@#   .      @@@%                                             ,@@%*%#  /#((@@&%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%@@@  (.      *@@@%                                              @@@(##. .#//%@@%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%@@#         (@%(#@,                               .*/./,*,#(/%. &@@&,(#  //(%@@&%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%@@@,        %@&%,*%%, .(               ...*//(,  **         *#@@@%@@@.*&  (#/*&@@%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%&@@@.       (&(@&,.,@#.(*                    ,,/%,           //&@@@@@@(,@@@@@&&%@@%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%@@@&       #@,@@@(  ,&%@,                      .              ,##&@@@@&*@@@@@@@@@@&%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%&@@@#      ,@,@@@@@,* ,*@%                                      ,, ((@@@,  .##@@@@@&%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%&@@@@@# ,/   *@,&@*#@@&,   ** /%&@&&%%(/                        .%@@#(%@@@@*%    (*%@@@&%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%@@@@%,  ,    @#%%.#@@@@,,.   *&@@@@@@@@@@@&(*,              .#@@@@@@@@@@@@@/%.     ,#&@@@%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%@@@#.#*      *&.(/,,@#&@% * %@@@%* /%@@@@@@@@@@@@@@@*  .%@@@@@@@@%*.(@@@@@@@,#.      (.%@@@&%%%%%%%%%%%
%%%%%%%%%%%%%%%@@&,  *,  .  ,@**,  (@#*@&%%@@@%.   ..  **%@@@@@@@@*    .@@@@@@@@,,*/@@@@@. ,(       .(*@@@@%%%%%%%%%%
%%%%%%%%%%%%%%@@&    ,.     #@.*.  &% *@@,#%%###@@&&&@@@@&@@@@@%      &@@@@(.,@@@&@@@@@@@.           */(@@@&%%%%%%%%%
%%%%%%%%%%%%%@@@.     .,    #@  .  @( %@,       (#.   *@&%( ,#%.        (@@&* .(##*  (@@@@@#,#&@@@@@&&(/%#%@@@%%%%%%%%%
%%%%%%%%%%%%@@&.       /    &% ., .@/,@#               ,%@@@@#   ..     ,@@@@,    *@%#,.*@@@.#&*%(@@@@@@@@@@@@%%%%%%%%%
%%%%%%%%%%%@@@.        (,   @#.,   @( #&*              .*(&/      .     .@@%*&           &@@,(@@@@%%%@@@@@@@@@%%%%%%%%%
%%%%%%%%%%@@@/          /   /#      .@#(                       /       @@#           #%@@@,/@@@@@@%*@@@@@@@&%%%%%%%%%
%%%%%%%%%&@@*   (#.    ,           *@. %@@#/.     ,.            ./&.     .&%.       .%@@@@@,  @@#%@@@@@@@@@@@%%%%%%%%%%
%%%%%%%%%@@@*#((#*(#  /&.           .@(#@@@@%#, @&*.          /&,         *@#&        #@@@@  *#.&@@@(/@@@@@@@%%%%%%%%%%
%%%%%%%%@@@@/@*//(/*(  /@.*,.       ..##%@@@@@@@,%         .@@.           /@@@(       ,/@@%,#@@@%*@@@@@@@@@@%%%%%%%%%%%
%%%%%%%%@@@@,@**(*&@@%.   /#      . *@##%#@@@@&/*    ..*/&@,.&   ,%&&/,./@@@@@@@&*   (/&@@%/(@/*@@@@@@@@@@@%%%%%%%%%%%%
%%%%%%%%@@@@*/@*&@@@@@@&*  .#&&, .*..@@(/&@@@@/,   (%@@@@#. .%@@@@@@@@@@@@@@@@@@@@%. //@@@@,&&@@@@@@@@@@@@%%%%%%%%%%%%%
%%%%%%%%&@@@&,@@*@@@@@@@@@&*     .   ./(&&@@@@/,#&@@&, .*&%,(%@#(./((@@@@@@@@@@@@@@@&%@@@@@@@@@@@@@@@@@@&%%%%%%%%%%%%%%
%%%%%%%%%%%@@@@@@@@@@@@@@@@@@%, ./*&&@#.,/@@@@@@@&*#.*&,.%. /(/,///*&@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@&%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%&@@@@@@@@@@@@*.**..&@@@&@@@@%/,@@@@&@@#   %@@%(%#@/##&@@@,#@%%@@@@@@@@@@@@@@@@@@@@@@@&%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%&@@@@@@@&//(&@@@@@@%,   /@@@@@@@&%,.&@@@%@@@@@@@@@@@@@%,(@(@@@@&&@@@@@@@@@#(%@@@@@@@%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%@@@@##@@&      /@@@@@@@@@@%@@@@@@@@@@@@@@@@@@@@@@@@@@@@@%@@@@@@@@,,*./%(@@@@@%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%@@@(.*@@*      .#@@@@@@@@@@@@@@@@@@@@@@@&(/,,,*&@@@@@@@@@@@@@@@@@@%#&%(/#/%@@@@%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%@@@,.#@        /#@@@@@@@(@@@@@@@@@@@@@/     ,(@@@@@@@@@@@@@@@@@@@@@(@##(&%#&@@@&%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%@@@  *.        #@@@@%@@@(&@@@@@@@@%*.        ,./##&@@@@@@@@@@@@@@@@@@@&@(##(@@@&%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%&@@% .,.%       *@@@@@@@@@,@@@@(,                  #,@@//*##@@@@@@@@@@@@@&&((@@@&%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%@@@&.  /(       .@&&/#.* .                          .*(@%,#@&@@@@@@@@@@@@@%##@@@%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%@@@@,              .                                   /((@&(.@@@@@@@@@@@@@@%@@@%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%@@@@@#.                                                    *#**@#,&@@@@@@@@@@@@@@@%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%@@@@@&%                                                      /&%.**.@@&@@@@@@@@@@@%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%@@@@&..                                                         (* .@#@@@@@@@@@@%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%&@@@@@&.(                                                       . ,/.*@/&@@@@@@&%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%@@@@@@#*/                                        , .          .%/.##@*(%@@@@@%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%@@@@@@&,,                                    .*,.*             .#@#/#,(@@@@&%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%@@@@@@#,/,@(                             ,#*                *#..*(((@@@@&%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%&@@@@@@#./&@@%/,.       .,(#       ./**&*                   ,%#/#%@@@@@%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%&@@@@@@@*.(@@@@@@@@@@@@@%##%&%*/#&%,,#%%(.%&             .(#./@@@@@%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%@@@@@@@@#*.*##&@@@@@@@@@@@#**/&@%%%%%&@@@@@(****/((*##&@@@@@@@&%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%@@@@@@@@@@@@@@@@@@@@@@@@@@&%%%%%%%%@@@@@@@@@@@@@@@@@@@@@&%%%%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%&@@@@@@@@@@@@@@@@&%%%%%%%%%%%%%&@@@@@@@@@@&&&%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%`)
}
