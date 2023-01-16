import json

with open("/root/playground/ocr/data/processed_table.json") as f:
    data = json.load(f)


for key in data:
    if len(key["process"]) > 1600 and len(key["process"]) < 2000:
        print(f'LEN : {len(key["process"])}')

        print(f'PROCESSOS: {key["process"]}')
        break

better_cases = sum([1 for key in data if len(key["process"]) > 1600 ])
print(f'Number of cases with more than 1600 characters: {better_cases}')