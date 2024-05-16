di = {}

di["a"] = 0

# when we are doing this type of validation, python considers 0 as false, so, it will print "a" when di["a"] is 0
# otherwise, it will print "b"
if not di["a"]:
    print("a")
else:
    print("b")
