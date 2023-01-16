import json
import re

_PROCESS_IDENTIFIER = "PROCESSO.*\d{7}-\d{2}.\d{4}.\d{1}.\d{2}.\d{4}"

with open("data/processed.json", "r", encoding="utf8") as f:
    data = json.load(f)


table_list = []
for process in data:
    dict_table = {}
    match = re.search(_PROCESS_IDENTIFIER, process, re.IGNORECASE)
    if match:
        dict_table["procces_id"] = match.group()
    else:
        dict_table["procces_id"] = None
    dict_table["process"] = process
    dict_table["extract_at"] = "2023-01-09"
    dict_table["title"] = "Edição 3637/2023 - Caderno do Tribunal Superior do Trabalho - Judiciário"
    if len(dict_table["process"]) > 1600:
        table_list.append(dict_table)

with open("data/processed_table.json", "w", encoding="utf8") as f:
    json.dump(table_list, f, ensure_ascii=False)
