import json

with open("flowtrack.json", "r") as file:
    data = json.load(file)

count_empty = 0
count_content = 0
for d in data["data"]:
    r = json.loads(d["output_data"])
    if r["module_output_data"]["bigquery_result"] == {}:
        count_empty += 1
    else:
        count_content += 1
print(count_empty, count_content)
