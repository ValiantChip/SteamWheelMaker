package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
)

func main() {
	os.Exit(ReturnWithCode())
}

func ReturnWithCode() int {
	var filename string
	var maxHours int
	var outname string
	flag.StringVar(&filename, "f", "/home/valiantchip/Documents/steamgamelist", "filename (string - required)")
	flag.IntVar(&maxHours, "t", 0, "maximum amount of hours for a game to be selected (int >= 0 - default:0)")
	flag.StringVar(&outname, "o", "./f.out", "output filename (string - default:./f.out)")
	flag.Parse()
	if filename == "" {
		flag.PrintDefaults()
		return 1
	}

	fl, err := os.Open(filename)
	if err != nil {
		fmt.Printf("unable to open file: %s\n", err)
		return 1
	}
	out, err := os.Create(outname)
	if err != nil {
		fmt.Printf("unable to create file: %s\n", err)
		return 1
	}
	defer out.Close()

	defer fl.Close()

	writer := bufio.NewWriter(out)

	scanner := bufio.NewScanner(fl)
	var prev string
	var needachievements bool
	for scanner.Scan() {
		line := scanner.Text()
		if line == `TOTAL PLAYED` {
			scanner.Scan()
			var playtime float64
			var unit string
			txt := scanner.Text()
			fmt.Sscanf(txt, "%f %s", &playtime, &unit)
			if unit != `hours` {
				playtime = 1
			}
			playtime = math.Ceil(playtime)

			if playtime <= float64(maxHours) {
				writer.WriteString(prev + "\n")
			}
			needachievements = true
			continue

		} else if line == `ACHIEVEMENTS` {
			if needachievements {
				needachievements = false
				continue
			}
			writer.WriteString(prev + "\n")
		}
		prev = line
	}

	writer.Flush()

	return 0
}
