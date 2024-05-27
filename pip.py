import sqlite3

dbfile = 'VideoHashesV1.db'
con = sqlite3.connect(dbfile)

cur = con.cursor()
counter = cur.execute("SELECT COUNT(*) FROM video_hashes")
for count in counter:    
    print(count)

con.close()