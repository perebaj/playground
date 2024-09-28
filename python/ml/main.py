import pickle


f = open("models/misc.pickle", "rb")

model = pickle.load(f)


print(model)
