import os
import subprocess
from uuid import uuid4

# YAFL - Yet Another Ffmpeg Library


FFMPEG = "ffmpeg"
DEFAULT_ARGS = ["-v", "error"]

DELETE_TEMP_FILES = True

def run_command(command: str|list[str], cwd: str=".") -> None:
    if isinstance(command, str):
        command = command.split(" ")

    pipe = subprocess.Popen(command, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
    # Returns tuple of (output, errors)
    return pipe.communicate()


def concat_inputs(inputs: list[str]):
    if isinstance(inputs, str):
        return ["-i", inputs]
    
    ins = []
    for i in inputs: ins.extend(['-i', i])
    return ins
    

def concat(videos: list, output: str, reencode: bool=False): 
    if reencode: 
        """
        Base Case: 

        ffmpeg -i opening.mkv -i episode.mkv -i ending.mkv 
        -filter_complex "[0:v] [0:a] [1:v] [1:a] [2:v] [2:a] concat=n=3:v=1:a=1 [v] [a]" 
        -map "[v]" -map "[a]" output.mkv 
        """
        n = len(videos)
        fc_streams = "".join([f"[{i}:v][{i}:a]" for i in range(n)])
        cmd = [
            FFMPEG, *DEFAULT_ARGS,
            *concat_inputs(videos),
            "-filter_complex",
            f'{fc_streams}concat=n={n}:v=1:a=1[v][a]',
            '-map', '[v]', '-map','[a]', 
            output
        ]
        print(" ".join(cmd))
        return run_command(cmd)

    # This doesn't fucking work, goddamn stack overflow fucker 
    # https://trac.ffmpeg.org/wiki/Concatenate
    "BASE CASE: ffmpeg -f concat -i opening.mkv -i episode.mkv -c copy output.mkv"
    cmd = [
        FFMPEG, *DEFAULT_ARGS,
        *concat_inputs(videos),
        "-f", "concat", 
        "-c", "copy", output
    ]
    print(" ".join(cmd))

    return run_command(cmd)


def clip(video: str, start: str, end: str,  output: str=None, reencode: bool=False):
    """ Create a clip of a video, given a start and end timestamp """
    if not output:
        output = video + "-trimmed"
    
    cmd = [
        FFMPEG, *DEFAULT_ARGS,
        *concat_inputs(video),
        "-ss", start, # "-accurate_seek",
        "-to", end,
        "-c:v", "libx264", "-c:a", "aac",
        output
    ]

    print(" ".join(cmd))

    "ffmpeg -i input.mp4 -ss 00:05:20 -accurate_seek -to 00:10:00 -c:v libx264 -c:a aac output7.mp4"
    return run_command(cmd) 

def clip_join(video: str, clips: list[tuple], output: str=None):
    """ Divide an video in clips, and join them consecutively """
    paths, ext = [], video.split(".")[-1]
    for clop in clips:
        start, end = clop 
        print(clop)
        video_hex = str(uuid4()) + f".{ext}"
        print(clip(video, start, end, output=video_hex))
        paths.append(video_hex)

    output, errors = concat(paths, output, reencode=True)
    print(errors)
    
    if DELETE_TEMP_FILES: 
        [os.remove(path) for path in paths]
    return output, errors

def reencode(video: str, output: str, vcodec: str="copy", acodec: str="copy", args: list[str]=None):
    cmd = [
        FFMPEG, *DEFAULT_ARGS,
        concat_inputs(video), 
        *args,
        f"-c:v {vcodec}", f"-c:a {acodec}",
        output
    ]

    "ffmpeg -i input.mp4 <<ARGS>> -c:v libx264 -c:a aac output.mp4"
    return run_command(cmd)
    


def compress(video: str, output: str, crf: int=None, preset: str="medium", tune: str=None, acodec:str="copy"): 
    """ Shorthand for reencode(video, 'libx264', 'copy', crf, args). 
    It reencodes the video with a different crf value, to lower the bitrate, and consequently lower the file size
    
    - CRF (0 - 51): 0 -> Lossless; 51 -> Maximum Compression 
    - Preset: More time, less file size
    https://trac.ffmpeg.org/wiki/Encode/H.264
        - ()"""
    
    preset, tune = preset.strip().lower(), tune.strip().lower()
    if preset not in ["ultrafast", "superfast", "veryfast", "faster", "fast", "medium", "slow", "slower", "veryslow"]:
        return 
    
    args = [
        "-preset", preset
    ]
    
    if tune and tune in ["film", "animation", "grain", "stillimage", "fastdecode", "zerolatency"]:
        args.extend(["-tune",  tune])

    if crf and isinstance(crf, (int, str)):
        args.extend(["-crf", crf])

    return reencode(video, output, 'libx264', acodec=acodec, preset=preset, args=args)




if __name__ == "__main__":
    PATH = "/mnt/c/users/breno/downloads/G2arEtna_720p.mp4"

    # print(concat([
    #     "78c8966b-5cc9-4a92-9740-67c2fda1493e.mp4",
    #     "2090f644-3155-4622-b1e9-b5ba3245fe67.mp4",
    #     "c811ccca-d4c3-4444-80cc-7fc5984d68ca.mp4"
    # ],
    
    # 'flopper.mp4', reencode=False)
    # )
    # exit()
    clips = [
        ("00:00:00", "00:01:00"),
        ("00:02:00", "00:03:00"),
        ("00:04:00", "00:05:00")
    ]
    clip_join(
        PATH, clips, "floog.mp4"
    )