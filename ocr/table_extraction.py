import json
import re

import fitz

_DEJT_PATH = "/root/playground/ocr/data/diario-09-01-2023.pdf"
_PROCESS_IDENTIFIER = "^PROCESSO.*\d{7}-\d{2}.\d{4}.\d{1}.\d{2}.\d{4}$"

file_name = _DEJT_PATH.split("/")[-1]
_OUTPUT_PATH = f"data/{file_name}_output.txt"


_REMOVE_REGEX_1 = "^Código para aferir autenticidade deste caderno: \d{6}$"
_REMOVE_REGEX_2 = "^\d{4}\/\d{4}$"
_REMOVE_REGEX_3 = "^\d+$"
_REMOVE_REGEX_4 = "^Tribunal Superior do Trabalho$"


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
procces_chunk_list = process_chunk(diario_list, _PROCESS_IDENTIFIER)

procces_chunk_list_before_treatment = []
for process in procces_chunk_list:
    proccess_after_treatment = process
    n_process = re.sub(_REMOVE_REGEX_1, "", proccess_after_treatment, flags=re.MULTILINE)
    n_process = re.sub(_REMOVE_REGEX_2, "", n_process, flags=re.MULTILINE)
    n_process = re.sub(_REMOVE_REGEX_3, "", n_process, flags=re.MULTILINE)
    n_process = re.sub(_REMOVE_REGEX_4, "", n_process, flags=re.MULTILINE)
    procces_chunk_list_before_treatment.append(n_process)

with open("data/processed.json", "w", encoding="utf8") as f:
    json.dump(procces_chunk_list_before_treatment, f, ensure_ascii=False)
