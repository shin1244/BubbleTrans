import cv2 as cv
from glob import glob
from PIL import Image
import deepl
from manga_ocr import MangaOcr
from ultralytics import YOLO
from PIL import Image
import os

if __name__ == '__main__':
  mocr = MangaOcr()
  model = YOLO("./best.pt")
  auth_key = os.getenv("DEEPL_AUTH_KEY")
  translator = deepl.Translator(auth_key)
  img_list = glob('./image/*')
  page_trans = []
  for image_path in img_list:
    if "txt" in image_path:
      continue
    img = cv.imread(image_path)
    result = model(img)
    with open(image_path[:-4]+".txt", 'w', encoding='utf-8') as file:
      for box in result[0].boxes:
        xA, yA, xB, yB = box.xyxy[0].tolist()
        xA, yA, xB, yB = map(int, [xA, yA, xB, yB])
        cropped_img = img[yA:yB, xA:xB]
        pil_img = Image.fromarray(cv.cvtColor(cropped_img, cv.COLOR_BGR2RGB))

        text = mocr(pil_img)
        tran = translator.translate_text(text, source_lang="JA", target_lang="KO")
        file.write(f"{yA} {yB} {xA} {xB} {tran.text}\n")