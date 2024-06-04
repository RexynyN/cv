import cv2
import imagehash as ih
import sqlalchemy as sqla
from sqlalchemy import String, Integer, Column
from sqlalchemy.orm import relationship, sessionmaker
from PIL import Image
from multiprocessing import Pool

# import sqlite3    
# con = sqlite3.connect(dbfile)

NUM_WORKERS = 20 
dbfile = 'VideoHashesV1.db'

engine = sqla.create_engine(f"sqlite://{dbfile}", echo=True)
Base = sqla.orm.declarative_base()

class FrameHash(Base):
    __tablename__ ="frame_hashes"
    id = Column(Integer, primary_key=True)
    perception_hash = Column(Integer)
    average_hash = Column(Integer)
    difference_hash = Column(Integer)
    color_hash = Column(Integer)
    wavelet_hash = Column(Integer)

    video = relationship("VideoHash", back_populates="video_hashes")


class VideoHash(Base):
    __tablename__ = "video_hashes"
    
    id = Column(Integer, primary_key=True)
    path = Column(String) 

    frames = relationship("FrameHash", order_by=FrameHash.id, back_populates="frame_hashes") 


Base.create_all(engine)
session = sessionmaker(bind=engine)

# # Create all records 
# session.add_all([
#     frame,
#     frame,
# ])

# Needs to commit an insertion
session.commit()

# Raw SQL
sql = sqla.text("SELECT product_id, BIT_COUNT(phash1 ^ phash2) as hd from A ORDER BY hd ASC;") 
rs = engine.execute(sql)

for row in rs:
    print(row)



def compute_hashes(frame: Image) -> FrameHash:
    return FrameHash(
        color_hash=ih.colorhash(frame),
        wavelet_hash=ih.whash(frame),
        perception_hash=ih.phash(frame),
        average_hash=ih.average_hash(frame),
        difference_hash=ih.dhash(frame)
    )

def add_video_hash(path: str) -> str:
    count = 0
    vidcap = cv2.VideoCapture(path)
    hashes = [] 
    while success:
        vidcap.set(cv2.CAP_PROP_POS_MSEC, (count * 5000))
        success, image = vidcap.read()
        frame = Image.fromarray(image)
        hashes.append(compute_hashes(frame))
        count += 1

    if len(hashes) == 0:
        return None

    video = VideoHash(
        path=path,
        frames=hashes
    )

    try:
        session.add(video)
        session.commit()
        return ""
    except Exception as e:
        return path


def run_workers(paths: list[str]):
    with Pool(processes=NUM_WORKERS) as pool:
        leftovers = pool.map(add_video_hash, paths)
        leftovers = [overs for overs in leftovers if overs]