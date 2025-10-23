# CV - Computer Vision
Computer Vision shenanigans and doohickeys

## ffmpeg - Some Commands 

Compression by codec:

```bash 
ffmpeg -i input.mp4 -vcodec h264 -acodec mp2 output.mp4
```

Compression by crf:

```bash 
ffmpeg -i input.mp4 -vcodec libx265 -crf 28 output.mp4
```

Transform the aspect ratio to a vertical cellphone-like video

```bash
ffmpeg -i input.mp4 -vf "crop=ih*9/16:ih:(iw-ih*9/16)/2:0" -preset medium -crf 30 output.mp4
```


## Links for me

https://www.geeksforgeeks.org/change-image-resolution-using-pillow-in-python/
https://www.geeksforgeeks.org/feature-detection-and-matching-with-opencv-python/
https://www.geeksforgeeks.org/feature-extraction-and-image-classification-using-opencv/
https://www.geeksforgeeks.org/image-feature-extraction-using-python/

https://docs.opencv.org/4.x/db/d27/tutorial_py_table_of_contents_feature2d.html
https://docs.opencv.org/4.x/d8/d4b/tutorial_py_knn_opencv.html

https://www.youtube.com/watch?v=oEKg_jiV1Ng

https://www.quora.com/What-are-some-measures-of-image-quality