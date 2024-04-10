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

	hasher "github.com/corona10/goimagehash"
	"gocv.io/x/gocv"
)

type VideoHash struct {
	Path   string      `json:"path"`
	Frames []FrameHash `json:"frames"`
}

type FrameHash struct {
	PerceptionHash uint64 `json:"percep_hash"`
	AverageHash    uint64 `json:"avg_hash"`
	DifferenceHash uint64 `json:"diff_hash"`
}

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
	Benchmark("frames/frame-0.jpg", "frames/frame-3.jpg")

	start := time.Now()
	PATH := path.Join("/mnt/c/Users/Admin/Downloads/Imersão vídeo04 - Desenvolvimento do Home Broker com Nextjs.mp4")

	// hashes := PartialFrameHashes(PATH)
	hashes := FrameHashes(PATH)
	elapsed := time.Since(start)

	fmt.Println(hashes)
	fmt.Println("Elapsed time: ", elapsed)
}

func LoadJsonHashes() {
	return
}

func SaveJsonHashes(hashes VideoHash) {
	jsonMarshalled, err := json.Marshal(hashes)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(jsonMarshalled))
	err = os.WriteFile("test.json", jsonMarshalled, 0644)
	if err != err {
		log.Fatal(err)
	}
}

func FrameHashes(path string) VideoHash {
	tokens := strings.Split(path, "/")
	file := tokens[len(tokens)-1]
	vidcap, err := gocv.VideoCaptureFile(path)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(file)
	frame, success := gocv.NewMat(), true
	count, hashes, timer := 0., make([]FrameHash, 0), 5000.
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

		frameHashes := ComputeHashes(imgMat)

		hashes = append(hashes, frameHashes)
		count++
	}
	videoHash := VideoHash{
		Path:   path,
		Frames: hashes,
	}
	return videoHash
}

// Compute all three hashes for the given frame
func ComputeHashes(frame image.Image) FrameHash {
	avg, err := hasher.AverageHash(frame)
	if err != nil {
		log.Fatal(err)
	}

	diff, err := hasher.DifferenceHash(frame)
	if err != nil {
		log.Fatal(err)
	}

	perc, err := hasher.PerceptionHash(frame)
	if err != nil {
		log.Fatal(err)
	}

	frameHashes := FrameHash{
		PerceptionHash: perc.GetHash(),
		DifferenceHash: diff.GetHash(),
		AverageHash:    avg.GetHash(),
	}
	return frameHashes
}

func PartialFrameHashes(path string, startPercent, endPercent float32) ([]hasher.ImageHash, error) {
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
	count, hashes, timer := 0., make([]hasher.ImageHash, 0), 5000.
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
		hash, err := hasher.PerceptionHash(imgMat)
		if err != nil {
			log.Fatal(err)
		}

		hashes = append(hashes, *hash)
		count++
	}

	return hashes, nil
}
