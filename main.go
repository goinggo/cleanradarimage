// Copyright 2013 Ardan Studios. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
	Ardan Studios
	12973 SW 112 ST, Suite 153
	Miami, FL 33186
	bill@ardanstudios.com

	// Install Binary Package C Library Procedure
	mkdir ~/temp
	cd ~/temp
	curl -O http://www.imagemagick.org/download/ImageMagick.tar.gz
	tar -xzf ImageMagick.tar.gz
	rm -f ImageMagick.tar.gz
	cd ImageMagick-6.8.6-9/
	./configure
	make
	sudo make install
	sudo ldconfig /usr/local/lib   ** LINUX ONLY
	remove ImageMagick-6.8.6-9
	pkg-config --cflags --libs MagickWand

	-- For Development Environment Only
	export GOPATH=$HOME/<My New Folder Location>
	go get github.com/gographics/imagick/imagick

	// Make sure these environment variables are set
	MAGICK_HOME=$HOME/Spaces/PublicPackages/ImageMagick-6.8.6
	DYLD_LIBRARY_PATH=$MAGICK_HOME/lib/
	PKG_CONFIG_PATH=$HOME/Spaces/PublicPackages/ImageMagick-6.8.6/lib/pkgconfig

	cleanradarimage source.gif out.gif
*/

// Package main will take a NOAA radar image and remove the noise colors.
package main

import (
	"fmt"
	"github.com/gographics/imagick/imagick"
	"io/ioutil"
	"log"
	"os"
)

// main removes noise from the image.
func main() {
	if len(os.Args) != 3 {
		fmt.Println("cleanradarimage source.gif out.gif")
		return
	}

	if fileExists(os.Args[2]) == true {
		os.Remove(os.Args[2])
	}

	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	imageBinary, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("Read File: %s", err)
	}

	//if err := mw.ReadImage(os.Args[1]); err != nil {
	//	log.Fatal(err)
	//}

	err = mw.ReadImageBlob(imageBinary)
	if err != nil {
		log.Fatalf("Read Image: %s", err)
	}

	fuzz := float64(10) // should be 10%

	colors := []string{
		"#3030CE",
		"#04e9e7",
		"#019ff4",
		"#0300f4",
		"#a9a879",
		"#777777",
		"#7a4679",
		"#aa7ca9",
		"#d7acd6",
		"#cccc99",
		"#999966",
		"#646464",
		"#663366",
	}

	pixelWand := imagick.NewPixelWand()
	defer func() {
		pixelWand.Destroy()
	}()

	// Remove unwanted colors
	for _, color := range colors {
		pixelWand.SetColor(color)
		mw.TransparentPaintImage(pixelWand, 0, fuzz, false)
	}

	// Wave the image
	if err := mw.WaveImage(5, 100); err != nil {
		return
	}

	// Crop the image
	if err := mw.CropImage(600, 540, 0, 10); err != nil {
		return
	}

	// Resize the image
	if err := mw.ResizeImage(600, 530, imagick.FILTER_BOX, 0); err != nil {
		return
	}

	// Equalize
	if err := mw.EqualizeImage(); err != nil {
		return
	}

	// Blue the image
	if err := mw.GaussianBlurImage(4, 2); err != nil {
		return
	}

	// Brightness
	if err := mw.BrightnessContrastImage(-20, 30); err != nil {
		return
	}

	// Reset the iterator for the write
	mw.ResetIterator()

	mw.WriteImage(os.Args[2])
}

// fileExists tests for the existance of the specified file on disk.
func fileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}

	return false
}
