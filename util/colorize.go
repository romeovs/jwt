package util

import (
	"regexp"

	"github.com/mgutz/ansi"
)

func bite(f func(string) string) func([]byte) []byte {
	return func(in []byte) []byte {
		return []byte(f(string(in)))
	}
}

var none = []byte(ansi.ColorCode("reset"))

var colors = regexp.MustCompile("\033\\[[0-9;]+m")

func clear(in []byte) []byte {
	return []byte(colors.ReplaceAllString(string(in), ""))
}

func colorizer(reg, color string) func([]byte) []byte {
	r := regexp.MustCompile(reg)
	return func(in []byte) []byte {
		return r.ReplaceAllFunc(in, func(m []byte) []byte {
			return []byte(ansi.Color(string(clear(m)), color))
		})
	}
}

var colorizers = []func([]byte) []byte{
	colorizer(`{|}|\[|\]`, "black+b"),
	colorizer(`[,:]`, "black"),
	colorizer(`[0-9]+\.?[0-9]*[eE][-+]?[0-9]*`, "yellow"),
	colorizer(`"[^"]*"`, "green"),
	colorizer(`false|true`, "red"),
}

func Colorize(json string) string {
	res := []byte(json)
	for _, colorize := range colorizers {
		res = colorize(res)
	}
	return string(res)
}
