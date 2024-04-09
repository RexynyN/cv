package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/corona10/goimagehash"
	hasher "github.com/corona10/goimagehash"
	"gocv.io/x/gocv"
)

func Benchmark(link1 string, link2 string) {
	fmt.Println(link1)
	fmt.Println(link2)

	file1, _ := os.Open(link1)
	file2, _ := os.Open(link2)
	defer file1.Close()
	defer file2.Close()

	img1, _ := jpeg.Decode(file1)
	img2, _ := jpeg.Decode(file2)
	hash1, _ := hasher.AverageHash(img1)
	hash2, _ := hasher.AverageHash(img2)
	distance, _ := hash1.Distance(hash2)
	fmt.Printf("Average: %v\n", distance)

	hash1, _ = hasher.DifferenceHash(img1)
	hash2, _ = hasher.DifferenceHash(img2)
	distance, _ = hash1.Distance(hash2)
	fmt.Printf("Difference: %v\n", distance)

	hash3, _ := hasher.PerceptionHash(img1)
	hash4, _ := hasher.PerceptionHash(img2)
	distance, _ = hash3.Distance(hash4)
	fmt.Printf("Perception: %v\n", distance)

	os.Exit(0)
}

func main() {
	type Foo struct {
		Number int    `json:"number"`
		Title  string `json:"title"`
	}

	jsonMarshalled, err := json.Marshal(Foo{Number: 1, Title: "test"})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(jsonMarshalled))
	err = os.WriteFile("test.json", jsonMarshalled, 0644)
	if err != err {
		log.Fatal(err)
	}

	Benchmark("frames/frame-0.jpg", "frames/frame-3.jpg")

	start := time.Now()
	PATH := path.Join("/mnt/c/Users/Admin/Downloads/Imersão vídeo04 - Desenvolvimento do Home Broker com Nextjs.mp4")

	// hashes := PartialFrameHashes(PATH)
	hashes := FrameHashes(PATH)
	elapsed := time.Since(start)

	fmt.Println(hashes)
	fmt.Println("Elapsed time: ", elapsed)
}

func FrameHashes(path string) []hasher.ImageHash {
	tokens := strings.Split(path, "/")
	file := tokens[len(tokens)-1]
	vidcap, err := gocv.VideoCaptureFile(path)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(file)
	frame, success := gocv.NewMat(), true
	count, hashes, timer := 0., make([]hasher.ImageHash, 0), 5000.
	for {
		fmt.Println(count)
		// Grab the next frame in the interval
		vidcap.Set(gocv.VideoCapturePosMsec, count*timer)
		success = vidcap.Read(&frame)
		if !success {
			break
		}
		gocv.IMWrite(fmt.Sprintf("frames/frame-%d.jpg", int(count)), frame)

		// Transform the map into an image
		imgMat, err := frame.ToImage()
		if err != nil {
			log.Fatal(err)
		}
		// Calculate the hash and append
		hash, err := hasher.PerceptionHash(imgMat)
		if err != nil {
			log.Fatal(err)
		}

		hashes = append(hashes, *hash)
		count++
	}

	return hashes
}

func PartialFrameHashes(path string, startPercent, endPercent float32) ([]goimagehash.ImageHash, error) {
	if startPercent > endPercent {
		return nil, errors.New("The starting percentage of the frame's height is bigger than the end percentage")
	}

	tokens := strings.Split(path, "/")
	file := tokens[len(tokens)-1]
	vidcap, err := gocv.VideoCaptureFile(path)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(file)
	frame, success := gocv.NewMat(), true
	count, hashes, timer := 0., make([]goimagehash.ImageHash, 0), 5000.
	for {
		fmt.Println(count)
		// Grab the next frame in the interval
		vidcap.Set(gocv.VideoCapturePosMsec, count*timer)
		success = vidcap.Read(&frame)
		if !success {
			break
		}

		// Calculate the two cropping points
		width, height := frame.Cols(), frame.Rows()
		x0, x1, fHeight := 0, width, float64(height)
		y0, y1 := int(fHeight*0.1), int(fHeight*0.8)

		// Save the image
		croppedMat := frame.Region(image.Rect(x0, y0, x1, y1))
		gocv.IMWrite(fmt.Sprintf("frames/frame-%d.jpg", int(count)), croppedMat)

		// Transform the map into an image
		imgMat, err := croppedMat.ToImage()
		if err != nil {
			log.Fatal(err)
		}
		// Calculate the hash and append
		hash, err := goimagehash.PerceptionHash(imgMat)
		if err != nil {
			log.Fatal(err)
		}

		hashes = append(hashes, *hash)
		count++
	}

	return hashes, nil
}

// func main() {
// 	file1, _ := os.Open("sample1.jpg")
// 	file2, _ := os.Open("sample2.jpg")
// 	defer file1.Close()
// 	defer file2.Close()

// 	img1, _ := jpeg.Decode(file1)
// 	img2, _ := jpeg.Decode(file2)
// 	hash1, _ := goimagehash.AverageHash(img1)
// 	hash2, _ := goimagehash.AverageHash(img2)
// 	distance, _ := hash1.Distance(hash2)
// 	fmt.Printf("Distance between images: %v\n", distance)

// 	hash1, _ = goimagehash.DifferenceHash(img1)
// 	hash2, _ = goimagehash.DifferenceHash(img2)
// 	distance, _ = hash1.Distance(hash2)
// 	fmt.Printf("Distance between images: %v\n", distance)
// 	width, height := 8, 8
// 	hash3, _ := goimagehash.ExtAverageHash(img1, width, height)
// 	hash4, _ := goimagehash.ExtAverageHash(img2, width, height)
// 	distance, _ = hash3.Distance(hash4)
// 	fmt.Printf("Distance between images: %v\n", distance)
// 	fmt.Printf("hash3 bit size: %v\n", hash3.Bits())
// 	fmt.Printf("hash4 bit size: %v\n", hash4.Bits())

// 	var b bytes.Buffer
// 	foo := bufio.NewWriter(&b)
// 	_ = hash4.Dump(foo)
// 	foo.Flush()
// 	bar := bufio.NewReader(&b)
// 	hash5, _ := goimagehash.LoadExtImageHash(bar)
// }
