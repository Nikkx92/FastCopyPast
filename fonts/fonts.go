package fonts

import (
	"embed"
	"path/filepath"
	"strings"

	"gioui.org/font"
	"gioui.org/font/opentype"
)

//go:embed *
var SteticaFonts embed.FS

const (
	SteticaBlack        font.Typeface = "SteticaBlack"
	SteticaBold         font.Typeface = "SteticaBold"
	SteticaBoldItalic   font.Typeface = "SteticaBoldItalic"
	SteticaItalic       font.Typeface = "SteticaItalic"
	SteticaLight        font.Typeface = "SteticaLight"
	SteticaLightItalic  font.Typeface = "SteticaLightItalic"
	SteticaMedium       font.Typeface = "SteticaMedium"
	SteticaMediumItalic font.Typeface = "SteticaMediumItalic"
	SteticaRegular      font.Typeface = "SteticaRegular"
)

func ParseFont(fsys embed.FS) []font.FontFace {
	entries, err := fsys.ReadDir(".")
	if err != nil {
		return nil
	}

	var coll []font.FontFace
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if ext != ".ttf" && ext != ".otf" {
			continue
		}

		data, err := fsys.ReadFile(entry.Name())
		if err != nil {
			return nil
		}

		face, err := opentype.Parse(data)
		if err != nil {
			return nil
		}

		typefaceName := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))

		coll = append(coll, font.FontFace{
			Font: font.Font{
				Typeface: font.Typeface(typefaceName),
			},
			Face: face,
		})
	}
	return coll
}
