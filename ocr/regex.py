import re

pattern = "^(?!3637\/2023).*$"
string = "This string does not contain the pattern 3637/2023"

print(
    re.sub(
        pattern,
        "",
        string,
    )
)
