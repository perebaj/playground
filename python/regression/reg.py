import numpy as np

arr1 = [19, 12, 13, 0]
arr2 = [100, 28, 71, 6]

print(np.polyfit(arr1, arr2, deg=1))