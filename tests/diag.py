import base64
import numpy as np
import requests
import json

if __name__ == "__main__":

    size = 1000
    template = "http://localhost:3002/update_pixel?x={x}&y={y}&r={r}&g={g}&b={b}"
    img_truth = np.zeros((size, size, 3))
    for row in range(size):
        pixel = [100, 100, 100]
        img_truth[row, row] = pixel
        url = template.format(
            x=row,
            y=row,
            r=pixel[0],
            g=pixel[1],
            b=pixel[2],
        )
        requests.get(url)

    get_image_url = "http://localhost:3001/get/image"
    image_response = requests.get(get_image_url)
    image_data = json.loads(image_response.content)

    # TODO(gabe): figure out how to go from base64 string to an array of ints
