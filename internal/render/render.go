// Package render provides render  ->  Renders the image as ANSI
package render

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"golang.org/x/term"
)

func cmdExec(width, height int, path string) (*exec.Cmd, error) {
	cmdScale := fmt.Sprintf("scale=%d:%d", width, height)
	cmd := exec.Command("ffmpeg", "-re", //#nosec
		"-i", path, "-vf", cmdScale,
		"-pix_fmt", "rgb24", "-f",
		"rawvideo", "-",
	)

	return cmd, nil
}

func Process(path string) error {
	fmt.Print("\033[?25l")
	defer fmt.Print("\033[?25h")

	termWidth, termHeight, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return err
	}

	cmd, err := cmdExec(termWidth, termHeight, path)
	if err != nil {
		return err
	}

	pipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	buf := make([]byte, termWidth*termHeight*3)

	if err := cmd.Start(); err != nil {
		return err
	}

	for {
		_, err := io.ReadFull(pipe, buf)
		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
				break
			}

			return err
		}

		fmt.Print("\033[H")
		if _, err = os.Stdout.Write(ansi(termWidth, termHeight, buf)); err != nil {
			return err
		}
	}

	return nil
}

func ansi(width, height int, rgb []byte) []byte {
	output := make([]byte, width*height*3)

	for y := 0; y < height; y += 2 {
		if y+1 >= height {
			break
		}
		for x := 0; x < width; x++ {
			top := (y*width + x) * 3
			bottom := ((y+1)*width + x) * 3

			tR, tG, tB := rgb[top], rgb[top+1], rgb[top+2]
			bR, bG, bB := rgb[bottom], rgb[bottom+1], rgb[bottom+2]

			ansi := fmt.Sprintf("\033[38;2;%d;%d;%dm\033[48;2;%d;%d;%dm▀", tR, tG, tB, bR, bG, bB)
			output = append(output, []byte(ansi)...)
		}

		output = append(output, '\n')
	}

	return output
}
