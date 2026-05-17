# Welcome to CS224N!

We'll be using Python throughout the course. If you've got a good Python setup already, great! But make sure that it is at least Python version 3.8.

---

## Setup with pyenv + virtualenv (recommended)

> **Note:** gensim does not support Python 3.14+. Use Python 3.12.

### 1. Install Python 3.12 with pyenv:
    pyenv install 3.12.9

### 2. Create a virtual environment:
    pyenv virtualenv 3.12.9 cs224n

### 3. Activate the environment:
    pyenv activate cs224n

### 4. Install dependencies:
    pip install jupyter matplotlib numpy ipykernel scikit-learn gensim "datasets==2.18.0"

### 5. Register the kernel for Jupyter:
    python -m ipykernel install --user --name=cs224n --display-name "Python (cs224n)"

### 6. Open the notebook:
    jupyter notebook exploring_word_vectors.ipynb

### 7. Select the kernel:
Go to the toolbar of `exploring_word_vectors.ipynb`, click on **Kernel -> Change kernel**, and select **cs224n**.

### To deactivate the environment:
    pyenv deactivate

---

## Virtual environment path

```
/Users/jonathansilva/.pyenv/versions/cs224n
```

Python binary:

```
/Users/jonathansilva/.pyenv/versions/cs224n/bin/python
```

---

## Setup with Conda (alternative)

### 1. Create an environment with dependencies specified in env.yml:
    conda env create -f env.yml

### 2. Activate the new environment:
    conda activate cs224n

### 3. Install IPython kernel:
    python -m ipykernel install --user --name cs224n

### To deactivate:
    conda deactivate
