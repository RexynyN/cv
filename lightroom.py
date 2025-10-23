from PIL import Image, ImageEnhance, ImageOps
import numpy as np


def adjust_brightness(image_path, output_path, factor):
    """
    Ajusta o brilho da imagem.
    """
    image = Image.open(image_path)
    enhancer = ImageEnhance.Brightness(image)
    brightened_image = enhancer.enhance(factor)
    brightened_image.save(output_path)
    print(f"Imagem com brilho ajustado salva em {output_path}")


def adjust_color(image_path, output_path, factor):
    """
    Ajusta as cores da imagem (saturação).
    """
    image = Image.open(image_path)
    enhancer = ImageEnhance.Color(image)
    color_adjusted_image = enhancer.enhance(factor)
    color_adjusted_image.save(output_path)
    print(f"Imagem com cores ajustadas salva em {output_path}")


def adjust_contrast(image_path, output_path, factor):
    """
    Ajusta o contraste da imagem.
    """
    image = Image.open(image_path)
    enhancer = ImageEnhance.Contrast(image)
    contrast_image = enhancer.enhance(factor)
    contrast_image.save(output_path)
    print(f"Imagem com contraste ajustado salva em {output_path}")


def adjust_sharpness(image_path, output_path, factor):
    """
    Ajusta a nitidez da imagem.
    """
    image = Image.open(image_path)
    enhancer = ImageEnhance.Sharpness(image)
    sharp_image = enhancer.enhance(factor)
    sharp_image.save(output_path)
    print(f"Imagem com nitidez ajustada salva em {output_path}")


def apply_sepia(image_path, output_path):
    """
    Aplica o filtro sépia à imagem.
    """
    image = Image.open(image_path).convert("RGB")
    np_image = np.array(image)
    tr = [112, 66, 20]  # Sepia tone factors (red, green, blue adjustment)
    sepia_image = np.dot(np_image[..., :3], [0.393, 0.769, 0.189])
    sepia_image = np.clip(sepia_image + tr, 0, 255).astype(np.uint8)
    Image.fromarray(sepia_image).save(output_path)
    print(f"Imagem com filtro sépia salva em {output_path}")


def apply_grayscale(image_path, output_path):
    """
    Converte a imagem para tons de cinza.
    """
    image = Image.open(image_path)
    grayscale_image = ImageOps.grayscale(image)
    grayscale_image.save(output_path)
    print(f"Imagem em tons de cinza salva em {output_path}")


def apply_negative(image_path, output_path):
    """
    Aplica o filtro negativo à imagem.
    """
    image = Image.open(image_path)
    inverted_image = ImageOps.invert(image.convert("RGB"))
    inverted_image.save(output_path)
    print(f"Imagem com filtro negativo salva em {output_path}")


def apply_custom_filter(image_path, output_path, filter_matrix):
    """
    Aplica um filtro personalizado utilizando um kernel de convolução.
    """
    image = Image.open(image_path).convert("RGB")
    np_image = np.array(image, dtype=np.float32)
    
    # Aplica um kernel simples (e.g., detecção de bordas)
    kernel = np.array(filter_matrix, dtype=np.float32)
    filtered_image = np.zeros_like(np_image)
    for i in range(3):  # Processar cada canal de cor separadamente
        filtered_image[..., i] = np.convolve(np_image[..., i].flatten(), kernel.flatten(), 'same')
    
    filtered_image = np.clip(filtered_image, 0, 255).astype(np.uint8)
    Image.fromarray(filtered_image).save(output_path)
    print(f"Imagem com filtro personalizado salva em {output_path}")


# Exemplos de uso:

# Ajustar brilho
adjust_brightness("input.jpg", "output_brightness.jpg", factor=1.2)

# Ajustar cores
adjust_color("input.jpg", "output_color.jpg", factor=1.5)

# Ajustar contraste
adjust_contrast("input.jpg", "output_contrast.jpg", factor=1.3)

# Ajustar nitidez
adjust_sharpness("input.jpg", "output_sharpness.jpg", factor=2.0)

# Aplicar filtro sépia
apply_sepia("input.jpg", "output_sepia.jpg")

# Converter para tons de cinza
apply_grayscale("input.jpg", "output_grayscale.jpg")

# Aplicar filtro negativo
apply_negative("input.jpg", "output_negative.jpg")

# Aplicar filtro personalizado (exemplo: detecção de bordas)
edge_detection_kernel = [[-1, -1, -1], [-1, 8, -1], [-1, -1, -1]]
apply_custom_filter("input.jpg", "output_custom.jpg", edge_detection_kernel)
