# read json file
import json

import matplotlib.pyplot as plt
import numpy as np

with open("/root/playground/etl_dejt/clean/entire_table_clean.json") as f:
    data = json.load(f)


# for key in data:
#     if len(key["process"]) > 2000 and len(key["process"]) < 3000:
#         print(f'LEN : {len(key["process"])}')

#         print(f'PROCESSOS: {key["process"]}')
#         break
process_text_len_list = []


for key in data:
    process_text_len_list.append(len(key["process"]))

#The "Freedman-Diaconis" rule is a method for determining the optimal number of bins for a histogram

#calculate the IQR
iqr = np.subtract(*np.percentile(process_text_len_list, [75, 25]))

# calculate the bin width
bin_width = 2 * iqr * len(process_text_len_list) ** (-1/3)

# calculate the number of bins
bins = int((np.max(process_text_len_list) - np.min(process_text_len_list)) / bin_width)
print(f'Number of bins: {bins}')
plt.title("Historgrama da distribuição do tamanho dos processos")
plt.xlabel("Tamanho do processo(em número de caracteres)")
plt.ylabel("Número de processos")
# plt.xticks(np.arange(min(process_text_len_list), max(process_text_len_list)+1, 5000))
counts, edges, bars = plt.hist(process_text_len_list, bins=bins)
# plt.xlim(300, 5000) #To increase the x axis range
plt.bar_label(bars)
plt.show()
