import sqlite3
import imagehash as ih

def hash_diff(this: str, other: str):
    if not this or not other:
        raise ValueError("One of the operand are null or empty")
    
    if "," in this or "," in other:
        return ih.ImageMultiHash.__sub__(ih.hex_to_multihash(this), ih.hex_to_multihash(other))

    return ih.ImageHash.__sub__(ih.hex_to_hash(this), ih.hex_to_hash(other))

dbfile = 'VideoHashes.db'

conn = sqlite3.connect(dbfile)
cursor = conn.cursor()

# Creating the hamming distance for sqlite (the imagehash implementation, anyway)
conn.create_function("hash_diff", 2, hash_diff)

# Executar um SELECT
query = """
SELECT vh.id, vh.path, fh.frame_order, fh.perc_hash, fh.crop_hash
FROM video_hashes vh
JOIN frame_hashes fh 
ON vh.id = fh.video_id
WHERE fh.frame_order < 10
ORDER BY vh.id, fh.frame_order;
"""


query = """
SELECT DISTINCT vh.id, vh.path
FROM video_hashes AS vh
JOIN frame_hashes AS fh 
    ON vh.id = fh.video_id
WHERE HASH_DIFF(perc_hash, 'f2c213d966b90dc6') < 100
"""

cursor.execute(query)

# Iterar pelos resultados
results = cursor.fetchall()
for row in results:
    print(f"Video ID: {row[0]}, Path: {row[1]}")

