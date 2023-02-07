import re

text = """
Processo Nº AR-1000480-72.2019.5.00.0000
Relator
LELIO BENTES CORRÊA
AUTOR
TERRA SANTA AGRO S.A.
ADVOGADO
MARIA ADRIANNA LOBO LEAO DE
MATTOS(OAB: 47607/DF)
ADVOGADO
RAFAEL DE ALENCAR ARARIPE
CARNEIRO(OAB: 25120/DF)
ADVOGADO
ANTONIO CARLOS PAULA DE
OLIVEIRA(OAB: 12884/BA)
ADVOGADO
FERNANDA CURY MICHALANY(OAB:
314205/SP)
ADVOGADO
VALTON DORIA PESSOA(OAB:
317623/SP)
ADVOGADO
DANIELA YUASSA(OAB: 189774/SP)
RÉU
MARCOS CESAR DE MORAES
ADVOGADO
RENATA ARCOVERDE
HELCIAS(OAB: 38655/DF)
ADVOGADO
MAURICIO DE FIGUEIREDO
CORREA DA VEIGA(OAB: 21934/DF)
ADVOGADO
LUCAS BARBOSA DE ARAUJO(OAB:
60706/DF)
Intimado(s)/Citado(s):
  - TERRA SANTA AGRO S.A.
            PODER JUDICIÁRIO
            JUSTIÇA DO
"""

entities = {}
for line in text.split("\n"):
    match = re.match(r"(\w+):\s+(.+)", line)
    if match:
        key, value = match.groups()
        entities[key] = value

print(entities)