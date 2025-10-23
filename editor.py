import re
import subprocess
import os
from uuid import uuid4
from typing import Iterator, Literal

# 99:59:59.999
# I wrote this fucking regex, I do have some pride in this achievement
RE_TIMESTAMP = re.compile(r"([0-9]{2}:)+([0-5][0-9]:)+([0-5][0-9])+(.([0-9]{3}))?")



class OutputError(Exception):
    def __init__(self, message):            
        super().__init__(message)

class Timestamp:
    def __init__(self, duration: int|str) -> None:
        self.__set_duration(duration)
        
    def __set_duration(self, duration: int|str) -> None:
        if isinstance(duration, str):
            if self.__is_timestamp(duration):
                self.stamp = duration
                self.duration_milliseconds = self.stamp_to_milliseconds(duration)
                self.valid = True
            elif self.__is_number(duration):
                self.duration_milliseconds = int(duration)
                self.stamp = self.milliseconds_to_stamp(duration)
                self.valid = True
            else:
                self.duration_seconds = 0
                self.stamp = "00:00:00.000"
                self.valid = False
        elif isinstance(duration, (int, float)):
            self.duration_milliseconds = int(duration)
            self.stamp = self.milliseconds_to_stamp(duration)
            self.valid = True
        else: 
            self.duration_milliseconds = 0
            self.stamp = "00:00:00.000"
            self.valid = False

    @staticmethod
    def milliseconds_to_stamp(self, duration: int):
        # Really Really Lazy Solution™
        stamp = self.seconds_to_stamp(duration // 1000)
        mili = duration % 1000
        return stamp.replace(".000", f"{mili}")

    @staticmethod
    def seconds_to_stamp(self, duration: int|str) -> str:
        duration = int(duration)
        hours, minutes, seconds = (duration // 3600), (duration // 60) % 60, (duration % 60)

        return f"{hours}:{minutes}:{seconds}.000"
    
    @staticmethod
    def stamp_to_milliseconds(self, stamp: str) -> int:
        tokens = stamp.split(".")
        mili = tokens[-1] if len(tokens) > 1 else 0 

        tokens = tokens[0].split(":")
        tokens.reverse() # To stay in the order of SECS, MINS, HOURS
        for time, conv in zip(tokens, [1000, 60000, 3600000]):
            mili += time * conv

        return mili
    
    @staticmethod
    def stamp_to_seconds(self, stamp: str) -> int:
        return self.stamp_to_miliseconds(stamp) // 1000

    @staticmethod
    def is_timestamp(self, string: str) -> str | None:
        result = RE_TIMESTAMP.match(string)
        return result.group() if result else None
    
    @staticmethod
    def is_number(self, string: str):
        return string if string.isnumeric() else None

    def to_string(self) -> str:
        return self.__str__()
    
    def to_seconds(self) -> int:
        return self.duration_milliseconds // 1000
    
    def to_milliseconds(self) -> int:
        return self.duration_milliseconds
    
    def __add__(self, other: Timestamp) -> Timestamp: # type: ignore
        if not isinstance(other, Timestamp):
            raise ValueError("The given object to add is not a Timestamp or Number (as milliseconds)")
        
        return Timestamp(self.duration_milliseconds + other.duration_milliseconds)

    def __sub__(self, other: Timestamp) -> Timestamp:  # type: ignore
        if not isinstance(other, Timestamp):
            raise ValueError("The given object to subtract is not a Timestamp or Number (as milliseconds)")
        
        return Timestamp(self.duration_milliseconds - other.duration_milliseconds)

    def __mul__(self, other: float) -> Timestamp: # type: ignore
        if not isinstance(other, (int, float)):
            raise ValueError("Multiplication is only allowed with a Timestamp and a Number")
        
        return Timestamp(self.duration_milliseconds * other)

    def __truediv__(self, other: float) -> Timestamp: # type: ignore
        if not isinstance(other, (int, float)):
            raise ValueError("Dividing is only allowed with a Timestamp and a Number")
        
        return Timestamp(self.duration_milliseconds / other)

    def __str__(self) -> str:
        return self.stamp
    
    def __repr__(self) -> str:
        return self.stamp

class Editor:
    def __init__(self, path: str) -> None:
        self.path = path
        self.clip_num = 0 # If we have multiple clips, create numbered videos
        self.operations = []

    def trim(self, start: str, end: str, name: str=None, reencode: bool=True, vcodec: str="l"):
        comm = ["ffmpeg"]
        if not self.__is_number(start) and not self.__is_timestamp(start):
            raise ValueError(f"The start of the trimmed video is not a valid timestamp or seconds integer: {start}")
        
        comm.extend(["-ss", start, "-accurate_seek", "-i", self.path])                     
        if not self.__is_number(start) and not self.__is_timestamp(start):
            raise ValueError(f"The end of the trimmed video is not a valid timestamp or seconds integer: {end}")
        
        comm.extend(["-t", end])
        if reencode:
            comm.extend(["-c:v", "libx264", "-c:a", "aac", name])
            

    def trim_start(self, end: str):
        pass
        
    def trim_end(self, start: str):
        pass
    
    def cut(self):
        pass

    def reverse(self):
        pass

    def re_encode_video(self, vcodec: str, acodec: str, input: str=None):
        input, output = self.__treat_iofiles(input, output)

        
        comm = [
            "ffmpeg",
            "-i", input,
            "-vcodec", vcodec,
            "-acodec", acodec,
            output
        ]

        _, err = self.__run_command(comm)
        if err:
            raise OutputError("Re-encode operation failed:\n{err}")


        pass 

    def extract_frame(self):
        pass

    def extract_frames(self):
        pass

    def downscale_video(self):
        pass
    
    def clips(self, clips: list[tuple|str|int]):
        typel = self.__homogeneous_type(clips)
        if not typel:
            raise ValueError("The clips list has a mixture of types. It must be homogenous as tuples ou strings.")

        if typel is tuple:
            pass
        elif typel is str and len():
            pass
        elif typel is int:
            pass
        else:
            raise ValueError("The values inside the timestamps are wrong or the .")

    
    
    def __homogeneous_type(self, seq: Iterator) -> (type | Literal[False]):
        if len(seq) < 1:
            return False
            
        iseq = iter(seq)
        first_type = type(next(iseq))
        # Return the type if homogeneous, and false if not
        return first_type if all((type(x) is first_type) for x in iseq) else False
    

    def __treat_iofiles(self, input: str, output: str):
        input = self.path if not input else input
        ext = input.split(".")[-1] 
        output = f"{uuid4()}.{ext}" if not output else output

        return input, output
        
        

    def __run_command(self, command: str|list[str]) -> tuple[str, str]:
        pass
        if isinstance(command, str):
            command = command.split(" ")

        process = subprocess.Popen(command, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
        output, errors = process.communicate()

        return output, errors

    def show_result(self, path: str) -> None:
        cap = cv2.VideoCapture(path)

        if cap.isOpened() == False:
            print("Error opening video file")

        while cap.isOpened():
            ret, frame = cap.read()
            if ret == True:
                cv2.imshow('Frame', frame)
                
                # Press Q on keyboard to exit
                if cv2.waitKey(25) & 0xFF == ord('q'):
                    break
            else:
                break

        cap.release()
        cv2.destroyAllWindows()


new = Editor("Breno")