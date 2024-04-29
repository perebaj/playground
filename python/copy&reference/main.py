class Copy:
    def __init__(self, value):
        self.value = value

    def changeReference(self, value):
        self.value = value

    def changeCopy(self, value):
        value = value


c = Copy(1)
print(c.value)  # 1
c.change(2)
print(c.value)  # 2
c.changeCopy(150)
print(c.value)  # 2
