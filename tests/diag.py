import base64
import numpy as np
import requests
import json
import io
import base64
from tqdm import tqdm
from PIL import Image

def get_image():
    get_image_url = "http://localhost:3001/json/image"
    image_response = requests.get(get_image_url)
    b64_image = json.loads(image_response.content)["data"]
    image = Image.open(io.BytesIO(base64.b64decode(b64_image)))

    return np.array(image)


if __name__ == "__main__":

    size = 100
    seed = 1000

    # for determinism
    np.random.seed(seed)

    servers = ["localhost:3001", "localhost:3002", "localhost:3003"]

    template = "http://{server}/update_pixel?X={x}&Y={y}&R={r}&G={g}&B={b}&A={a}"
    img_truth = get_image()

    for row in tqdm(range(size)):
        pixel = [100, 100, 100, 255]
        img_truth[row, row] = pixel
        url = template.format(
            server=np.random.choice(servers),
            x=row,
            y=row,
            r=pixel[0],
            g=pixel[1],
            b=pixel[2],
            a=pixel[3],
        )
        requests.get(url)

    assert np.allclose(img_truth, get_image())
