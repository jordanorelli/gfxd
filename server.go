package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"net/http"
	"strconv"
	"strings"
)

var (
	white = color.RGBA{0xff, 0xff, 0xff, 0xff}
	red   = color.RGBA{0xff, 0, 0, 0xff}
	green = color.RGBA{0, 0xff, 0, 0xff}
	blue  = color.RGBA{0, 0, 0xff, 0xff}
)

type server struct {
	out    io.Writer
	errors io.Writer
}

func (s *server) logReceived(r *http.Request) {
	fmt.Fprintf(s.out, "> %s\n", r.URL.String())
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.logReceived(r)

	encoder, err := getEncoding(r)
	if err != nil {
		s.writeError(w, err)
		return
	}

	width := parseInt(r.URL.Query().Get("w"), 400)
	height := parseInt(r.URL.Query().Get("h"), 400)
	bg := parseColor(r.URL.Query().Get("bg"), white)
	data := parseSeries(r.URL.Query().Get("s"), series{})

	// create a new canvas
	m := image.NewRGBA(image.Rect(0, 0, width, height))

	// paint the background
	draw.Draw(m, m.Bounds(), &image.Uniform{bg}, m.Bounds().Min, draw.Src)

	// compute column dimensions
	num_cols := data.len()
	col_width := width
	if num_cols > 0 {
		col_width = width / num_cols // float here?
	}

	fg := &image.Uniform{blue}
	for idx, val := range data {
		col_height_n := norm(val, data.min(), data.max())
		col_height := int(col_height_n * float64(height))
		x := idx * col_width
		rect := image.Rect(x, height, x+col_width, height-col_height)
		draw.Draw(m, rect, fg, image.ZP, draw.Src)
	}
	encoder.WriteImage(w, m)
}

// parses an integer from a user string. If the user string is invalid, return
// the suuplied default int.
func parseInt(i_s string, i int) int {
	n, err := strconv.Atoi(i_s)
	if err != nil {
		return i
	}
	return n
}

// parses a color from a user string. If the user string is invalid, return the
// supplied default color.
func parseColor(c_s string, d color.RGBA) color.RGBA {
	var err error
	f := func(s string) (n uint8) {
		if len(s) != 2 {
			err = coalesce(err, fmt.Errorf("input too long"))
			return 0
		}
		i, e := strconv.ParseUint(s, 16, 0)
		if e != nil {
			err = coalesce(err, e)
			return 0
		}
		return uint8(i)
	}

	var c color.RGBA
	c_s = strings.ToLower(c_s)
	switch len(c_s) {
	case 2:
		n := f(c_s[0:2])
		c = color.RGBA{n, n, n, 0xff}
	case 4:
		n, a := f(c_s[0:2]), f(c_s[2:4])
		c = color.RGBA{n, n, n, a}
	case 6:
		c = color.RGBA{f(c_s[0:2]), f(c_s[2:4]), f(c_s[4:6]), 0xff}
	case 8:
		c = color.RGBA{f(c_s[0:2]), f(c_s[2:4]), f(c_s[4:6]), f(c_s[6:8])}
	default:
		return d
	}
	if err != nil {
		return d
	}
	return c
}

type series []int

// minimum value in the series
func (s series) min() int {
	switch len(s) {
	case 0:
		return 0
	case 1:
		return s[0]
	default:
	}
	m := s[0]
	for _, i := range s[1:] {
		if i < m {
			m = i
		}
	}
	return m
}

// maximum value in the series
func (s series) max() int {
	switch len(s) {
	case 0:
		return 0
	case 1:
		return s[0]
	default:
	}
	m := s[0]
	for _, i := range s[1:] {
		if i > m {
			m = i
		}
	}
	return m
}

func (s series) len() int {
	return len(s)
}

// parses a user-supplied series
func parseSeries(s string, missing series) series {
	parts := strings.Split(s, ",")
	out := make(series, 0, len(parts))
	for _, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil {
			return missing
		}
		out = append(out, n)
	}
	return out
}

type imageWriter interface {
	WriteImage(w http.ResponseWriter, m image.Image)
}

type pngWriter struct {
	CompressionLevel png.CompressionLevel
}

func (w pngWriter) WriteImage(rw http.ResponseWriter, m image.Image) {
	enc := png.Encoder{CompressionLevel: w.CompressionLevel}
	rw.Header().Add("Content-Type", "image/png")
	enc.Encode(rw, m)
}

func getEncoding(r *http.Request) (imageWriter, error) {
	parts := strings.Split(r.URL.Path, ".")
	if len(parts) == 0 {
		return nil, fmt.Errorf("no encoding specified")
	}

	switch e := parts[len(parts)-1]; e {
	case "png":
		return pngWriter{CompressionLevel: png.DefaultCompression}, nil
	default:
		return nil, fmt.Errorf("invalid encoding: %s", e)
	}
}

func (s *server) writeError(w http.ResponseWriter, err error) {
	fmt.Fprintf(s.errors, "E %s\n", err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}
