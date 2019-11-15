import base64
import numpy as np
import requests
import json
import io
import base64
from PIL import Image

if __name__ == "__main__":

    size = 1000
    seed = 1000

    # for determinism
    np.random.seed(seed)

    servers = ["localhost:3001", "localhost:3002", "localhost:3003"]

    template = "http://{server}/update_pixel?X={x}&Y={y}&R={r}&G={g}&B={b}"
    img_truth = np.zeros((size, size, 3))

    for row in range(size):
        pixel = [100, 100, 100]
        img_truth[row, row] = pixel
        url = template.format(
            server=np.random.choice(servers),
            x=row,
            y=row,
            r=pixel[0],
            g=pixel[1],
            b=pixel[2],
        )
        requests.get(url)

    get_image_url = "http://localhost:3001/json/image"
    image_response = requests.get(get_image_url)
    b64_image = json.loads(image_response.content)["data"]
    image = Image.open(io.BytesIO(base64.b64decode(b64_image)))
    img_recovered = np.array(image)[:, :, :3]

    print(img_truth.shape)
    print(img_recovered.shape)

    assert np.allclose(img_truth, img_recovered)
