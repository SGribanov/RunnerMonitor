from __future__ import annotations

from pathlib import Path
from math import cos, pi

from PIL import Image, ImageChops, ImageEnhance, ImageFilter, ImageOps


ROOT = Path(__file__).resolve().parents[1]
ASSETS = ROOT / "assets"
SOURCE = ASSETS / "hourglass-source-chromakey.png"
PNG_OUT = ASSETS / "runner-monitor-hourglass.png"
GIF_OUT = ASSETS / "runner-monitor-hourglass-spin.gif"
ICO_OUT = ASSETS / "runner-monitor-hourglass.ico"


def remove_chroma_key(source: Image.Image) -> Image.Image:
    image = source.convert("RGBA")
    pixels = image.load()
    width, height = image.size
    for y in range(height):
        for x in range(width):
            r, g, b, a = pixels[x, y]
            green_delta = g - max(r, b)
            if g > 140 and green_delta > 45:
                pixels[x, y] = (r, g, b, 0)
            elif g > 110 and green_delta > 25:
                alpha = max(0, min(a, int((45 - green_delta) / 20 * 255)))
                pixels[x, y] = (r, g, b, alpha)

    alpha = image.getchannel("A")
    contracted = alpha.filter(ImageFilter.MinFilter(3))
    soft = contracted.filter(ImageFilter.GaussianBlur(0.35))
    image.putalpha(soft)

    pixels = image.load()
    for y in range(height):
        for x in range(width):
            r, g, b, a = pixels[x, y]
            if a < 12:
                pixels[x, y] = (0, 0, 0, 0)
                continue
            if g > r and g > b:
                despill = int((g - max(r, b)) * 0.82)
                g = max(max(r, b), g - despill)
                pixels[x, y] = (r, g, b, a)
    return image


def crop_and_pad(image: Image.Image, size: int = 512, padding: int = 28) -> Image.Image:
    alpha = ImageChops.multiply(image.getchannel("A"), image.getchannel("A"))
    bbox = alpha.getbbox()
    if bbox is None:
        raise RuntimeError("No non-transparent subject found")

    subject = image.crop(bbox)
    subject.thumbnail((size - padding * 2, size - padding * 2), Image.Resampling.LANCZOS)
    canvas = Image.new("RGBA", (size, size), (0, 0, 0, 0))
    x = (size - subject.width) // 2
    y = (size - subject.height) // 2
    canvas.alpha_composite(subject, (x, y))
    return canvas


def spin_frames(icon: Image.Image, frame_count: int = 24) -> list[Image.Image]:
    frames: list[Image.Image] = []
    for index in range(frame_count):
        angle = 2 * pi * index / frame_count
        width_scale = max(0.12, abs(cos(angle)))
        frame = icon
        if cos(angle) < 0:
            frame = ImageOps.mirror(frame)

        target_width = max(32, int(icon.width * width_scale))
        squeezed = frame.resize((target_width, icon.height), Image.Resampling.BICUBIC)

        brightness = 0.78 + 0.22 * width_scale
        squeezed = ImageEnhance.Brightness(squeezed).enhance(brightness)
        contrast = 0.92 + 0.08 * width_scale
        squeezed = ImageEnhance.Contrast(squeezed).enhance(contrast)

        canvas = Image.new("RGBA", icon.size, (0, 0, 0, 0))
        canvas.alpha_composite(squeezed, ((icon.width - target_width) // 2, 0))
        frames.append(canvas)
    return frames


def save_gif(frames: list[Image.Image]) -> None:
    frames[0].save(
        GIF_OUT,
        save_all=True,
        append_images=frames[1:],
        duration=55,
        loop=0,
        disposal=2,
        transparency=0,
    )


def save_ico(icon: Image.Image) -> None:
    sizes = [(256, 256), (128, 128), (64, 64), (48, 48), (32, 32), (16, 16)]
    icon.save(ICO_OUT, sizes=sizes)


def main() -> None:
    ASSETS.mkdir(parents=True, exist_ok=True)
    source = Image.open(SOURCE)
    transparent = crop_and_pad(remove_chroma_key(source))
    transparent.save(PNG_OUT)
    save_gif(spin_frames(transparent))
    save_ico(transparent)
    print(PNG_OUT)
    print(GIF_OUT)
    print(ICO_OUT)


if __name__ == "__main__":
    main()
