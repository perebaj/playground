# argos-vision

Captura uma face pela webcam, gera embedding com FaceNet (VGGFace2) e retorna a celebridade do dataset escolhido mais parecida.

## Requisitos

- Python **3.10–3.12**
- [`uv`](https://docs.astral.sh/uv/) — instale antes:
  ```bash
  brew install uv
  ```

## Setup

```bash
make install               # default: instala dependências base (PINS)
uv sync --extra hf         # opcional: adiciona suporte ao dataset hf1000
```

## Rodar

O fluxo tem dois passos: **construir o cache** (uma vez por dataset) e **fazer matching** (toda vez que quiser identificar uma face).

### 1. Construir o cache de embeddings

```bash
make build                          # default: PINS
make build dataset=hf1000           # ou outro dataset
```

Baixa o dataset escolhido, roda o modelo em cada foto e salva embeddings em `data/<dataset>_embeddings.npz`. Demora 10–30min na primeira vez, mas só precisa rodar uma vez por dataset.

### 2. Rodar matching

```bash
make run                            # usa cache do PINS
make run dataset=hf1000             # usa cache do hf1000
```

`SPACE` captura, `ESC` sai. Se o cache não foi construído pra esse dataset, o programa avisa pra rodar `make build` primeiro.

### Ver status dos caches

```bash
make list
```

## Datasets

| Dataset | Pessoas | Fotos/pessoa | Tamanho | Como acessar |
|---|---|---|---|---|
| **pins** (default) | 105 | ~170 (cap 25) | ~1 GB | Via `kagglehub`, sem login |
| **hf1000** | 1.000 | ~18 | ~3 GB | Via `datasets`, precisa do extra `hf` |

Comandos equivalentes via CLI direta:

```bash
uv run argos-vision build --dataset pins
uv run argos-vision build --dataset hf1000

uv run argos-vision match --dataset pins
uv run argos-vision match --dataset hf1000 --image foto.jpg

uv run argos-vision list
```

### Recomendações

- **`pins`** — celebridades populares atuais (Robert Downey Jr., Margot Robbie), fotos em alta qualidade. **Default**.
- **`hf1000`** — variedade maior (1000 celebs) com qualidade boa. Use se quiser matches mais "surpresa".

## Flags

### `build`

| Flag | Default | Notas |
|---|---|---|
| `--dataset {pins,hf1000}` | `pins` | qual catalog construir |
| `--max-per-person N` | `25` | limita fotos por pessoa pra acelerar build |
| `--force` | — | reconstrói mesmo se o cache existir |

### `match`

| Flag | Default | Notas |
|---|---|---|
| `--dataset {pins,hf1000}` | `pins` | qual catalog usar (precisa ter o cache pronto) |
| `--image PATH` | — | usa imagem estática em vez da webcam |
| `--camera N` | `0` | índice da câmera |

## Testes

```bash
make test                          # roda tudo
make test testcase=test_search     # filtra por padrão de nome
```

## Lint

```bash
make lint                          # check (não modifica arquivos)
make format                        # aplica fixes automáticos
```
