import json

f = open("/root/playground/etl_dejt/clean/entire_table_clean3.json", "r")
data = json.load(f)
print(f"Len of data: {len(data)}")
result = list({dictionary["procces_id"]: dictionary for dictionary in data}.values())
print(f"Len of result: {len(result)}")

with open("dedup.json", "w") as f:
    json.dump(result, f, ensure_ascii=False)
