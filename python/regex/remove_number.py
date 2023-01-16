import re

text = """
Relata a requerente que, após o julgamento dos Embargos de
Declaração interpostos ao acórdão por meio do qual foram julgados
improcedentes os pedidos deduzidos na Ação Rescisória de sua
autoria, “a empresa se vê hoje sem provimento judicial que a proteja
do iminente levantamento de aproximadamente R$ 40 MILHÕES,
por uma única pessoa física, antes do julgamento de seu recurso
extraordinário”.
Postula a concessão de efeito suspensivo ao Recurso
Código para aferir autenticidade deste caderno: 194487
3637/2023
Tribunal Superior do Trabalho
2
Data da Disponibilização: Segunda-feira, 09 de Janeiro de 2023
Extraordinário, visando obter a paralisação da “execução que se
processa nos autos da Reclamação Trabalhista 0025800-
58.2009.5.24.0022, em trâmite perante a 2ª Vara do Trabalho de
Dourados, até que o e. STF julgue o recurso”.
Argumenta que, se houver o levantamento do importe de cerca de
R$ 40.000.000,00 (quarenta milhões de reais), dificilmente poderá
"""

result = re.sub(r"^\d+$", "", text, flags=re.MULTILINE)
print(result)
