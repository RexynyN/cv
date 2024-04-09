# CV - Computer Vision
Computer Vision shenanigans and doohickeys

## ffmpeg - Some Commands 

Compression by codec:

```bash 
ffmpeg -i input.mp4 -vcodec h264 -acodec mp2 output.mp4
```

Compression by crf:

```bash 
ffmpeg -i input.mp4  -vcodec libx265 -crf 28 output.mp4
```

