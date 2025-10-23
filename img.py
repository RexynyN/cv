import os
import cv2
import numpy as np
from PIL import Image

class ImageQuality:
    def __init__(self, path: str):
        self.path = path 
        self.image = cv2.imread(path)
        self.extension = path.split(".")[-1]

    def get_size(self) -> float: 
        # Retuns size in KBs
        return os.path.getsize(self.path) / 1024

    def get_resolution(self) -> tuple[int, int]: 
        # Obter largura e altura
        return Image.open(self.path).size

    def laplacian_sharpness(self) -> float:
        # Calcular a variância do Laplaciano
        image = cv2.imread(self.path, cv2.IMREAD_GRAYSCALE)
        return cv2.Laplacian(image, cv2.CV_64F).var()

    def jpg_blockiness(self):
        if self.extension not in ("jpg", "jpeg", "jfif"):
            raise ValueError(f"{self.path} is not a jpeg file")
        
        # Convertendo para escala de cinza
        gray_image = cv2.cvtColor(self.image, cv2.COLOR_BGR2GRAY)
        h, w = gray_image.shape

        # Quantidade de blocos horizontais e verticais de 8x8
        h_blocks = h // 8
        w_blocks = w // 8

        blockiness = 0

        # Iterar sobre blocos e calcular variação entre eles
        for i in range(0, h_blocks * 8, 8):
            for j in range(0, w_blocks * 8, 8):
                block = gray_image[i:i+8, j:j+8]
                blockiness += np.std(block)

        return blockiness / (h_blocks * w_blocks)



quality = ImageQuality("image.jpeg")
blockiness_score = quality.jpg_blockiness()

print(f"Pontuação de blocos (blockiness): {blockiness_score}")

