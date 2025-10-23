package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/B1NARY-GR0UP/violin"
	hasher "github.com/corona10/goimagehash"
	cv "gocv.io/x/gocv"
)

const PATH string = "./data/original"

func CalcHash(pathfile string) {
	path := filepath.Join(PATH, pathfile)
	read := cv.IMRead(path, cv.IMReadColor)
	img, _ := read.ToImage()

	hash, _ := hasher.PerceptionHash(img)
	fmt.Printf("Hash '%s' -> %d\n", path, hash.GetHash())
}

func main() {
	grayImage := cv.IMRead("image.jpg", cv.IMReadGrayScale)
	cv.Corn
	cv.ORBScoreTypeHarris

	cv.Dilate()


	window := cv.NewWindow("Trov")
	window.IMShow(image)
	window.WaitKey(-1)
}


# Applying the function 
dst = cv2.cornerHarris(gray_image, blockSize=2, ksize=3, k=0.04) 
  
# dilate to mark the corners 
dst = cv2.dilate(dst, None) 
image[dst > 0.01 * dst.max()] = [0, 255, 0] 
  
cv2.imshow('haris_corner', image) 
cv2.waitKey() 