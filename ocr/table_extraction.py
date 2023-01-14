import json
import re

import fitz

_DEJT_PATH = "/root/playground/ocr/data/diario-09-01-2023.pdf"
_PROCESS_IDENTIFIER = "^PROCESSO.*\d{7}-\d{2}.\d{4}.\d{1}.\d{2}.\d{4}$"

file_name = _DEJT_PATH.split("/")[-1]
_OUTPUT_PATH = f"data/{file_name}_output.txt"


_REMOVE_REGEX_1 = "^Código para aferir autenticidade deste caderno: \d{6}$"
_REMOVE_REGEX_2 = "^(\d{4}\/\d{4})$"
_REMOVE_REGEX_3 = "^(\d)$"


def pdf_to_text(pdf_path):
    doc = fitz.open(pdf_path)

    f = open(_OUTPUT_PATH, "wb")
    text = ""
    for page in doc:
        text = page.get_text().encode("utf-8")
        f.write(text)
    print(f"Saved output:{_OUTPUT_PATH} ")


def get_output_txt() -> list[str]:
    diario_text = open(_OUTPUT_PATH, "r")
    diario_list = diario_text.readlines()
    return diario_list


def process_chunk(diario_list, regex):
    process_chunk_list = []
    process_chunk: str = ""
    for line in diario_list:
        if re.match(regex, line, re.IGNORECASE):
            # process_chunk = process_chunk + line
            process_chunk_list.append(process_chunk)
            process_chunk = "" + line
        else:
            process_chunk = process_chunk + line
    return process_chunk_list


# pdf_to_text(_DEJT_PATH)
diario_list = get_output_txt()
procces_chunk_list = process_chunk(diario_list[:5000], _PROCESS_IDENTIFIER)
print(f"Found {len(procces_chunk_list)} processos")
print(f"Type: {type(procces_chunk_list[2])}")
for item in procces_chunk_list:
    item_list = item.split("\n")
    for index, value in enumerate(item_list):
        if re.match(_REMOVE_REGEX_1, value):
            item_list.pop(index)
        if re.match(_REMOVE_REGEX_2, value):
            item_list.pop(index)
        if re.match(_REMOVE_REGEX_3, value):
            item_list.pop(index)
    print(item)
# print(procces_chunk_list[2].split("\n"))
# print(procces_chunk_list[0])
# print(re.match("DIÁRIO", "UM TESTE DIÀRIO DIÁRIO"))
# clean = re.findall(_REMOVE_REGEX_1, procces_chunk_list[2])
# print(clean)

# text_file = open("data/sample3.txt", "w")
# n = text_file.write(procces_chunk_list[2])
# text_file.close()
# print(get_output_txt(_OUTPUT_PATH))
# for page in doc.pages():
#     text = page.get_text().encode("utf-8")
#     f.write(text)
#     f.write(bytes((12,)))
