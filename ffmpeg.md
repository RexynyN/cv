# ffmpeg - Course 

Fast Forward Moving Picture Expert Group 

- Sources
https://img.ly/blog/ultimate-guide-to-ffmpeg/

## ffprobe

Command: 

`ffprobe -v error video.mp4 -show_format -show_streams -print_format json`

v: verbose 
show_format: Show informations about the format of the video
show_streams: Show informations about the streams inside the video 
print_format: How this information is going to be shown (Can be piped)


`ffprobe -v error video.mp4 -select_streams v -show_entries stream=codec_name`

select_streams: Select a specific stream
show_entries: Filter single informations about the stream
stream=codec_name: Returns only the codec name

`ffprobe -v error video.mp4 v -show_entries format=format_long_name -print_format json`


NOTE: The input can be a HTTP url!


## ffplay 


- Simple video playing 

`ffplay video.mp4`

- Set the window size and border 

`ffplay video.mp4 -x 600 -y 600 -no_border`

- Set the window starting location 

`ffplay video.mp4 -x 600 -y 600 -no_border -top 0 -left 0`

- Playing Full Screen

`ffplay video.mp4 -no_border -fs`

- Playing With no Sound

`ffplay video.mp4 -no_border -an`

- Playing With no Video (just sound)

`ffplay video.mp4 -no_border -vn -showmode waves`


- Shortcuts 


Space: Play/Pause
M: Mute
F: Fullscreen
0 - 9: Audio levels
/ - *: Up/Down Audio Levels
W: Cycle Showmodes
S: Frame by Frame
<- - -> Keys: Back/Forward 10s
Up - Down Keys: Back/Forward 1min
Esc: Quit


# Media Concepts 

Video - Audio - Image


## Codecs - Containers

Coder-Decoder = CoDec

- Examples
    - H.264 
    - H.265
    - VP9
    - Prores -> + Quality - Compression
    - DNxHD

Container 

Wrapper for the media essence, the file format itself
How the media data is organized

- Examples
    - MP5
    - MXF 
    - QT/MOV
    - MKT

    - WAV
    - M4A

## Transcoding 

One codec to another 

Prores -> H.264

- Transmuxing 

From one container to another 

MXF -> MP4

- Thumbnail Generation 

Preview image, Search hit, Hover-scrubbing, Poster frame


- Frame rate conversion 

25fps -> 30fps

- Bitrate conversion 

50Mbps -> 2Mbps

+ Quality -> + Space Occupied


# Core ffmpeg concepts 

## Architecture 

Input -> Unpack Audio/Video -> Uncompress -> (Transformations) -> Compress -> Pack -> Output

More Sophisticated: 

Input -> Demuxing Audio/Video -> Decode -> (Transformations) -> Encode -> Muxing -> Output


## Streams 

- Generally 1 video streams and 1 or more audio streams (voices, background, music, etc)
- Can be gathered with ffprobe (see First section)


- Select streams 

- Required for filters and outputs (-map)
- Sintax 

[input-index]:[stream-type]:[stream-index]

`ffmpeg -i video.mp4 -i overlay.png -i sounds.wav`

- Each (-i) is a 0-index input-index


- 0:v:0 -> First video stream of the first input 
- 1:a:2 -> Third audio stream of the second input

`ffmpeg -v error -y -i multitrack.mp4 -to 1 -map 0 multitrack-1s.mp4`

Creates a 1 second snippet of the video. (-map 0) gathers all the audio streams, if ommited, only the first stream if gathered. 

`ffmpeg -v error -y -i multitrack.mp4 -to 1 -map 0:v multitrack-1s.mp4`

Same thing, but it only gathers the video streams. 

`ffmpeg -v error -y -i multitrack.mp4 -to 1 -map 0:a multitrack-1s.mp4`

Same thing, but now audio-only. 

`ffmpeg -v error -y -i multitrack.mp4 -to 1 -map 0:a:1 multitrack-1s.mp4`

Gathers only the first audio stream.

`ffmpeg -v error -y -i multitrack.mp4 -i input2.mp4 -map 1:v:0 -map 0:a:1 mishmash.mp4`

Amalgamation of the second input video stream and the first input audio stream. 


## Filters
Changes the input in some way 


- Sintax
filter=key1=value1:key2=value2

- Examples
    - scale=width=1920:height=1080 
    - scale=w=1920:h=1080
    - scale=1920:1080


- Labelling 
    - Labelling 2 inputs to the split filter
        - split=2[sd_in][hd_in] 
    - Labelling two inputs (video/image) to overlay filter
        - [background][overlay]overlay

- Filter Chain

Separated by comma (,) and they're sequential

- scale=width=1920:height=1080 , split=2[sd_in][hd_in] , [background][overlay]overlay

## Filter Graph 

Multiple filter chains that can be non-linear and have multiple input-outputs

Separated with a semicolon (;) or specified with -vf, -af, -filter_complex

- -vf

Simple Video Filter Graphs (1 input, 1 output)

- Example: 

`ffmpeg -v error -y -i video.mp4 -vf "split[bg][ol]; [bg]scale=width=1920:height=1080, format=gray[bg_out]; [ol]scale=-1:480, hflip[ol_out]; [bg_out][ol_out]overlay=x=W-w:y=(H-h)/2" overlayed_inversed_grayed.mp4`

Creates a video, with a smaller, flipped version overlayed, and the original video is grayed out. 

```bash
ffmpeg -v error -y -i video.mp4 -vf 
"split[bg][ol]; 
[bg]scale=width=1920:height=1080, format=gray[bg_out]; 
[ol]scale=-1:480, hflip[ol_out]; 
[bg_out][ol_out]overlay=x=W-w:y=(H-h)/2" 
overlayed_inversed_grayed.mp4
```
Same command, but with the filters separated;

- -af

Simple Audio Filter Graphs (1 input, 1 output)

- Example: 

`ffmpeg -v error -y -i audio.wav -af "asplit=2[voice][bg]; [voice]volume=volume=2, pan=mono|c0=c0+1[voice_out]; [bg]volume=volume=0.5, pan=mono|c0=c2+c3[bg_out];[voice_out][bg_out]amerge=input=2" processed.wav`

Creates a audio, with a quiter background audio and higher voice audio;

```bash
ffmpeg -v error -y -i audio.wav -af 
"asplit=2[voice][bg]; 
[voice]volume=volume=2, pan=mono|c0=c0+1[voice_out]; 
[bg]volume=volume=0.5, pan=mono|c0=c2+c3[bg_out];
[voice_out][bg_out]amerge=input=2" 
processed.wav
```
Same command, but with the filters separated;

- -filter_complex

Complex Filters (N inputs -> N Outputs)

```bash
ffmpeg -v error -y -i video.mp4 -i logo.png -filter_complex "[1:v]scale=-1:200[small_logo]; [0:v][small_logo]overlay=x=W-w-50:y=H-h-50, split=2[sd_in][hd_in]; [sd_in]scale-2:480[sd]; [hd_in]scale=-2:1080[hd]; [0:a]pan=stereo|FL=c0+c2|FR=c1+c3[stereo_mix]" -map [sd] sd.mp4 -map [hd] hd.mp4 -map [stereo_mix] stereo.mp3
```

Not gonna explain this one, just see the picture, homie. 


```bash
ffmpeg -v error -y -i video.mp4 -i logo.png -filter_complex 
"[1:v]scale=-1:200[small_logo];
[0:v][small_logo]overlay=x=W-w-50:y=H-h-50, split=2[sd_in][hd_in];
[sd_in]scale-2:480[sd]; 
[hd_in]scale=-2:1080[hd];
[0:a]pan=stereo|FL=c0+c2|FR=c1+c3[stereo_mix]" 
-map [sd] sd.mp4 -map [hd] hd.mp4 -map [stereo_mix] stereo.mp3
```

Same command, but with the filters separated;
