import os
import sqlite3
import cv2
import imagehash as ih
from PIL import Image
from multiprocessing import Pool

DIR_PATH = "/mnt/c/Users/Breno/Downloads/ghqe/nok/"

FRAME_INTERVAL = 1000 # 1 second
EXT_WHITELIST = ("mp4", "m4v", "mkv", "mov")

NUM_WORKERS = 20

dbfile = 'VideoHashes.db'
conn, cursor = None, None


def t(hour, min, sec):
    return hour * 3600 + min * 60 + sec

def create_tables():
    global conn, cursor 

    cursor.execute("""
    CREATE TABLE IF NOT EXISTS video_hashes (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        path TEXT NOT NULL
    )
    """)

    cursor.execute("""
    CREATE TABLE IF NOT EXISTS frame_hashes (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        video_id INTEGER NOT NULL,
        frame_order INTEGER NOT NULL,
        perc_hash TEXT,
        avg_hash TEXT,
        diff_hash TEXT,
        color_hash TEXT,
        wave_hash TEXT,
        crop_hash TEXT,
        FOREIGN KEY (video_id) REFERENCES video_hashes (id)
    )
    """)

    conn.commit()

def hash_diff(this: str, other: str):
    if not this or not other:
        raise ValueError("One of the operand are null or empty")
    
    if "," in this or "," in other:
        return ih.ImageMultiHash.__sub__(ih.hex_to_multihash(this), ih.hex_to_multihash(other))

    return ih.ImageHash.__sub__(ih.hex_to_hash(this), ih.hex_to_hash(other))


def compute_hashes(frame: Image):
    return {
        "avg": ih.average_hash(frame),
        "crop": ih.crop_resistant_hash(frame),
        "whash": ih.whash(frame),
        "phash": ih.phash(frame),
        "dhash": ih.dhash(frame),
        "color": ih.colorhash(frame)
    }

def video_hashes(path: str) -> str:
    vidcap = cv2.VideoCapture(path)

    frame_count, fps = vidcap.get(cv2.CAP_PROP_FRAME_COUNT), vidcap.get(cv2.CAP_PROP_FPS)
    duration = (frame_count / fps * 1000) if fps > 0 else 0 # In milliseconds

    count, hashes = 0, []
    while True:
        if count >= duration: 
            break
        vidcap.set(cv2.CAP_PROP_POS_MSEC, count)
        # print("Frame", count)
        success, image = vidcap.read()
        if not success:
            break

        hash = compute_hashes(Image.fromarray(image))
        hash["frame_order"] = count // FRAME_INTERVAL
        hashes.append(hash)
        
        count += FRAME_INTERVAL

    return hashes 

def add_video_hashes(path: str):
    global conn, cursor
    print(path)

    hashes = video_hashes(path)
    if len(hashes) == 0:
        return None

    try:
        cursor.execute("INSERT INTO video_hashes (path) VALUES (?)", (path,))
        video_id = cursor.lastrowid

        for h in hashes:
            cursor.execute("""
            INSERT INTO frame_hashes (
                video_id, frame_order, avg_hash, crop_hash, wave_hash, perc_hash, diff_hash, color_hash
            ) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
            """, (
                video_id, str(h["frame_order"]), str(h["avg"]), str(h["crop"]),
                str(h["whash"]),str(h["phash"]),str(h["dhash"]),str(h["color"])
            ))
        conn.commit()
        return ""
    except Exception as e:
        print(e)
        return path

def run_workers(paths: list[str]):
    with Pool(processes=NUM_WORKERS) as pool:
        leftovers = pool.map(add_video_hashes, paths)
        leftovers = [overs for overs in leftovers if overs]

def get_all_files(dir_path: str) -> list[str]:
    all_files = []
    for root, _, files in os.walk(dir_path):
        for file in files:
            all_files.append(os.path.join(root, file))

    return [p for p in all_files if p.split(".")[-1] in EXT_WHITELIST]


def main():
    global conn, cursor 
    conn = sqlite3.connect(dbfile)
    cursor = conn.cursor()

    # Creating the hamming distance for sqlite (the imagehash implementation, anyway)
    conn.create_function("hash_diff", 2, hash_diff)
    create_tables()

    video_paths = get_all_files(DIR_PATH)
    print(f"Scanned {len(video_paths)} videos to process")

    run_workers(video_paths)

if __name__ == "__main__":
    main()
