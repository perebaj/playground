import re

# define the regular expression pattern
pattern = r"\d+"

# define the string to search
string = "I have 3 cats and 42 dogs."

# search for a match using the re.search() function
match = re.search(pattern, string)

# check if a match was found
if match:
    print(match.group())
else:
    print("No match found.")