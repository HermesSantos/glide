package main

import (
	"fmt"
	"os"

	"glide"
)

func main() {
	err := glide.New().
		AddElement(glide.Element{
			Shape: glide.ShapeBox,
			X:     5,
			Y:     3,
			Label: "alpha",
		}).
		AddElement(glide.Element{
			Shape: glide.ShapeBox,
			X:     25,
			Y:     8,
			Label: "beta",
		}).
		AddElement(glide.Element{
			Shape:  glide.ShapeBox,
			Center: true,
			Label:  "center",
		}).
		Run()

	if err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}
}
