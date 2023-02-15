#read json
import json
import re

f = open('/root/playground/etl_dejt/clean/entire_table_clean3.json', 'r')
data = json.load(f)

approved = []
for index, d in enumerate(data):
    if re.findall('dou provimento ao agravo de instrumento', d["process"], flags=re.MULTILINE) != []:
        d["sentiment"] = "Aprovado"
        approved.append(d)

with open("approved.json", "w") as f:
        json.dump(approved, f, ensure_ascii=False)