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
	"sync"
	"time"

	hasher "github.com/corona10/goimagehash"
	"github.com/google/uuid"
	"gocv.io/x/gocv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type VideoHash struct {
	gorm.Model
	Path   string      `json:"path"`
	Frames []FrameHash `json:"frames" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type FrameHash struct {
	gorm.Model
	PerceptionHash int64 `json:"percep_hash"`
	AverageHash    int64 `json:"avg_hash"`
	DifferenceHash int64 `json:"diff_hash"`
	VideoHashID    uint
}

const (
	NUM_WORKERS = 20
	PATH        = "<path>"
	DB_NAME     = "VideoHashesV1.db"
)

var db *gorm.DB

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
	var err error
	db, err = gorm.Open(sqlite.Open(DB_NAME), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Error in creating the database: ", err))
	}

	db.AutoMigrate(&VideoHash{}, &FrameHash{})
	// Benchmark("frames/frame-0.jpg", "frames/frame-3.jpg")

	files, err := os.ReadDir(PATH)
	if err != nil {
		panic(err)
	}

	dirs := []string{PATH}
	for _, entry := range files {
		if entry.IsDir() {
			dirs = append(dirs, path.Join(PATH, entry.Name()))
		}
	}

	start := time.Now()

	wg := sync.WaitGroup{}
	for _, dir := range dirs {
		wg.Add(1)
		go GoVideoDirHashes(dir, &wg)
	}

	wg.Wait()
	fmt.Println("Passed the main waitgroup!")

	elapsed := time.Since(start)
	fmt.Println("Elapsed time: ", elapsed)
}

func SaveJsonHashes(hashes VideoHash) {
	jsonMarshalled, err := json.Marshal(hashes)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(path.Join("json", uuid.New().String()+".json"), jsonMarshalled, 0644)
	if err != err {
		log.Fatal(err)
	}
}

func clearFileExtensions(dirPath string) []string {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		panic(err)
	}

	// entries = slices.DeleteFunc(entries, func(entry fs.DirEntry) bool {
	// 	return !strings.HasSuffix(entry.Name(), ".mp4")
	// })

	files := []string{}
	for _, entry := range entries {
		fileName := path.Join(dirPath, entry.Name())
		files = append(files, fileName)
	}
	return files
}

func GoVideoDirHashes(dirPath string, wg *sync.WaitGroup) {
	files := clearFileExtensions(dirPath)
	innerWg := sync.WaitGroup{}

	stride := int(len(files) / NUM_WORKERS) // How many files in each worker
	ceil := stride * NUM_WORKERS            // Get the maximum if the len(files) is divisible by workers
	lowBound, upBound := 0, 0
	for count := range NUM_WORKERS {
		upBound = ((count + 1) * stride) + 1
		if upBound == ceil {
			// Include the rest of the files left if the # of files is not
			// Divisible by the # of workers
			upBound += len(files) - ceil
		}

		innerWg.Add(1)
		go GoSaveVideoHashes(files[lowBound:upBound], &innerWg)
		lowBound = upBound + 1
	}

	fmt.Println("All workers done!")
	innerWg.Wait()
	wg.Done() // We're done :D
	fmt.Println("Passed the WaitGroup!")
}

func GoSaveVideoHashes(paths []string, wg *sync.WaitGroup) {
	videoHashes := []VideoHash{}
	total := len(paths)
	for idx, file := range paths {
		// All of this just to log a damn file ffs
		tokens := strings.Split(file, "/")
		logF := tokens[len(tokens)-1]
		fmt.Println(fmt.Sprintf("%d/%d -> %s", (idx + 1), total, logF))

		// Do all the reading and hash calculation shenanigans
		hash, err := ReadFrameHashes(file)
		if err != nil {
			log.Println(err)
			continue
		}
		videoHashes = append(videoHashes, hash)
	}

	for _, hash := range videoHashes {
		// SaveJsonHashes(hash)
		db.Create(&hash)
	}
	fmt.Println("Worker Done!")
	wg.Done()
	fmt.Println("Passed the WaitGroup!")
}

func ReadFrameHashes(path string) (VideoHash, error) {
	// Open the video
	vidcap, err := gocv.VideoCaptureFile(path)
	if err != nil {
		log.Println("Error reading the video: ", err)
		return VideoHash{}, err
	}

	frame, success := gocv.NewMat(), true
	count, hashes, timer := 0., make([]FrameHash, 0), 5000.
	for {
		// Grab the next frame in the interval
		vidcap.Set(gocv.VideoCapturePosMsec, count*timer)
		success = vidcap.Read(&frame)
		// Usually it means the video is over, but if it was an error, I really don't care lmao
		if !success {
			break
		}

		// Save the Frame on disk
		// gocv.IMWrite(fmt.Sprintf("frames/frame-%d.jpg", int(count)), frame)

		// Transform the map into an image
		imgMat, err := frame.ToImage()
		if err != nil {
			log.Println("Error transforming in an image: ", err)
			continue
		}

		// Compute the Hashes and Append
		frameHashes := ComputeHashes(imgMat)
		hashes = append(hashes, frameHashes)
		count++
	}

	// At least one Frame is a must
	if len(hashes) == 0 {
		return VideoHash{}, errors.New("No frame was read, error.")
	}

	// Let it live its best life
	videoHash := VideoHash{
		Path:   path,
		Frames: hashes,
	}
	return videoHash, nil
}

// Compute all three hashes for the given frame
func ComputeHashes(frame image.Image) FrameHash {
	// If this is true, some shit went down
	if frame == nil {
		return FrameHash{}
	}

	// Calculate all of the hashes and return it
	avg, _ := hasher.AverageHash(frame)
	diff, _ := hasher.DifferenceHash(frame)
	perc, _ := hasher.PerceptionHash(frame)
	return FrameHash{
		PerceptionHash: int64(perc.GetHash()),
		DifferenceHash: int64(diff.GetHash()),
		AverageHash:    int64(avg.GetHash()),
	}
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
