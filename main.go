package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

//go:embed latte.json
var latte []byte
//go:embed frappe.json
var frappe []byte
//go:embed macchiato.json
var macchiato []byte
//go:embed mocha.json
var mocha []byte

type Catpuccin []map[string][]string
type Palette []map[string]string

func hexToRgb(hex string) (int, int, int) {
	hex = strings.TrimPrefix(hex, "#")
	var r, g, b int
	fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	return r, g, b
}

func hexToHsl(hex string) (int, int, int) {
	r, g, b := hexToRgb(hex)
	red, green, blue := float64(r)/255, float64(g)/255, float64(b)/255
	max, min := red, red
	if green > max {
		max = green
	}
	if blue > max {
		max = blue
	}
	if green < min {
		min = green
	}
	if blue < min {
		min = blue
	}

	l := (max + min) / 2
	var h, s float64
	if max != min {
		delta := max - min
		if l > 0.5 {
			s = delta / (2 - max - min)
		} else {
			s = delta / (max + min)
		}
		switch max {
		case red:
			h = (green - blue) / delta
			if green < blue { h += 6 }
		case green:
			h = (blue - red) / delta + 2
		case blue:
			h = (red - green) / delta + 4
		}
		h /= 6
	}
	return int(h * 360), int(s * 100), int(l * 100)
}

func process(path string, catpuccin map[string][]string, target map[string]string) error {
	input, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	output := string(input)
	original := output

	for color, values := range catpuccin {
		newHex, exists := target[color]
		if !exists {
			continue
		}

		r, g, b := hexToRgb(newHex)
		h, s, l := hexToHsl(newHex)

		newRgb := fmt.Sprintf("rgb(%d, %d, %d)", r, g, b)
		newRgbArray := fmt.Sprintf("[%d, %d, %d]", r, g, b)
		newHsl := fmt.Sprintf("hsl(%ddeg, %d%%, %d%%)", h, s, l)

		for _, value := range values {
			if strings.HasPrefix(value, "#") {
				output = strings.ReplaceAll(output, value, newHex)
			} else if strings.HasPrefix(value, "rgb") {
				output = strings.ReplaceAll(output, value, newRgb)
				var rv, gv, bv int
				_, err := fmt.Sscanf(value, "rgb(%d, %d, %d)", &rv, &gv, &bv)
				if err == nil {
					pattern := fmt.Sprintf(`\[\s*%d\s*,\s*%d\s*,\s*%d\s*\]`, rv, gv, bv)
					regex := regexp.MustCompile(pattern)
					output = regex.ReplaceAllString(output, newRgbArray)
				}
			} else if strings.HasPrefix(value, "hsl") {
				output = strings.ReplaceAll(output, value, newHsl)
			}
		}
	}

	if output != original {
		err = os.WriteFile(path, []byte(output), 0644)
		if err == nil {
			fmt.Printf("Modified %s\n", filepath.Base(path))
		}
		return err
	}
	return nil
}

func main() {
	path := flag.String("path", "", "File/Dir to modify")
	newTheme := flag.String("theme", "", "The theme file")
	catpuccinType := flag.String("type", "macchiato", "The catppucin theme you're replacing")
	flag.Parse()

	if *path == "" || *newTheme == "" {
		log.Fatal("Usage: -path <file or dir> -theme <theme>.json")
	}

	newpath := *path
	if strings.HasPrefix(newpath, "~/") {
		home, _ := os.UserHomeDir()
		newpath = filepath.Join(home, newpath[2:])
	}

	catpuccinThemes := map[string][]byte{
		"latte":     latte,
		"frappe":    frappe,
		"macchiato": macchiato,
		"mocha":     mocha,
	}

	selected, exists := catpuccinThemes[*catpuccinType]
	if !exists {
		log.Fatalf("Invalid type: %s. Options are latte, frappe, macchiato, and mocha", *catpuccinType)
	}

	themeData, err := os.ReadFile(*newTheme)
	if err != nil {
		log.Fatalf("Error reading theme file: %v", err)
	}

	var catppuccinTheme Catpuccin
	var theme Palette

	json.Unmarshal(selected, &catppuccinTheme)
	json.Unmarshal(themeData, &theme)

	catpuccin, target := catppuccinTheme[0], theme[0]

	err = filepath.Walk(newpath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		return process(path, catpuccin, target)
	})

	if err != nil {
		log.Fatal(err)
	}
}