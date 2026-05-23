# Argos

## *All Eyez on You*

![argos](/Users/jonathansilva/playground/argos/assets/juno-and-argus.jpg)

### Maçaneta Eletrônica com Reconhecimento Facial e Matching de Celebridades


**Anteprojeto**
Autor: Jonathan Silva
Data: 2026-05-23

---

## 1. Introdução

Nos últimos anos, técnicas de visão computacional e reconhecimento facial deixaram de habitar exclusivamente laboratórios de pesquisa e passaram a integrar dispositivos do cotidiano — do desbloqueio do celular ao controle automatizado de fronteiras, passando por câmeras de vigilância residenciais. Em paralelo, fechaduras inteligentes residenciais (Yale, August, Intelbras) tornaram-se uma categoria de produto consolidada, com arquiteturas conhecidas (motor DC com gearbox acionando o thumbturn interno) e custo acessível ao consumidor final.

Esses dois movimentos, no entanto, raramente se cruzam em projetos pessoais e de aprendizado. As soluções comerciais que combinam reconhecimento facial e fechaduras tendem a ser caixas-pretas, e os projetos DIY documentados na literatura técnica se concentram em uma das duas pontas: ou no atuador mecânico (com servos hobby e fechaduras improvisadas) ou no modelo de visão (com câmeras conectadas a desktops via cabo USB).

Este anteprojeto descreve o **Argos**, uma maçaneta eletrônica residencial não-crítica que integra os dois universos. Ao detectar uma face na porta, o sistema destrava a fechadura e exibe, em um pequeno display embarcado, qual celebridade do banco VGGFace2 mais se assemelha à pessoa identificada. O propósito é deliberadamente lúdico: a maçaneta não é um sistema de segurança, mas sim um objeto de demonstração e diversão social.

A questão central que o projeto se propõe a investigar é:

> **É viável, com hardware de baixo custo e operação por bateria, construir um sistema integrado de visão computacional, atuação mecânica e fabricação por impressão 3D para uma aplicação residencial não-crítica, sem recorrer a serviços comerciais fechados de reconhecimento facial?**

A resposta a essa pergunta se desdobra em três sub-investigações práticas: (i) que arquitetura de hardware permite operação contínua por bateria com latência aceitável; (ii) que arquitetura de software (local vs. cloud, modelos pré-treinados vs. customizados) viabiliza o reconhecimento em tempo razoável; e (iii) que abordagem de fabricação caseira (CAD + impressão 3D) permite produzir um invólucro funcional sem comprometer a fechadura existente.

## 2. Objetivos

### 2.1 Objetivo geral

Projetar e construir um protótipo funcional de maçaneta eletrônica retrofit, operada por bateria, com reconhecimento facial integrado a um banco de celebridades, instalável em uma porta residencial interna sem substituir a fechadura original.

### 2.2 Objetivos específicos

1. **Modelagem mecânica em CAD**: produzir, em Fusion 360 ou FreeCAD, um invólucro 3D-imprimível que acopla ao thumbturn da fechadura existente e abriga motor DC, driver, ESP32-CAM, sensor PIR, display e baterias.
2. **Firmware embarcado**: implementar firmware em ESP32-CAM que opere em deep sleep, acorde por sensor PIR, capture imagem, envie via WiFi ao servidor de inferência, controle o motor com feedback de encoder e atualize o display com o resultado retornado.
3. **Serviço de inferência em cloud**: disponibilizar um endpoint HTTP em HuggingFace Spaces (ou equivalente gratuito) que recebe uma imagem, executa detecção facial, gera o embedding e retorna o nome e a foto da celebridade mais próxima do banco VGGFace2.
4. **Metas mensuráveis de aceite**:
  - Autonomia de bateria ≥ 30 dias em regime típico de 10 ativações por dia.
  - Latência ponta-a-ponta (acionamento do PIR até porta destravada + resultado no display) ≤ 5 segundos.
  - Taxa de detecção de face em condições residenciais normais (luz diurna, frente da porta) ≥ 90%.
  - Custo total de hardware ≤ R$ 400 (excluindo PC de desenvolvimento e impressora 3D).

## 3. Justificativa

### 3.1 Relevância para aprendizado e formação técnica

A construção do Argos demanda integração de três domínios técnicos que tendem a ser ensinados e praticados isoladamente:

- **Mecânica e fabricação digital**: modelagem paramétrica em CAD, impressão FDM, ajustes dimensionais para encaixes mecânicos com peças comerciais (rolamentos, parafusos, thumbturn da fechadura).
- **Sistemas embarcados**: programação de microcontroladores com restrições de energia, controle de atuadores com feedback, deep sleep, comunicação via WiFi.
- **Deep learning aplicado**: pipeline de detecção, alinhamento e embedding de faces; uso de modelos pré-treinados; indexação vetorial; deploy em ambiente cloud.

Projetos que cruzam os três domínios são raros em currículos acadêmicos e exigem desenvolvimento autônomo, o que reforça o caráter formativo do trabalho.

### 3.2 Relevância técnica

A literatura DIY de fechaduras inteligentes (Instructables, Adafruit, Hackaday) concentra-se em mecanismos de atuação (servos, solenoides, motores) e em controles de acesso simples (senha, RFID, app via Bluetooth). Por outro lado, a literatura de reconhecimento facial (FaceNet, ArcFace, DeepFace) foca em arquiteturas de rede neural e métricas de desempenho em benchmarks acadêmicos. **Pouco se publica sobre o encaixe prático entre ambos sob a restrição combinada de baixo custo, operação por bateria e instalação residencial.** Este projeto explora justamente essa interface, documentando trade-offs reais entre latência, consumo, qualidade de imagem e acurácia de match.

### 3.3 Relevância lúdica e social

O Argos é deliberadamente uma aplicação **não-crítica**. A porta continua tendo sua fechadura original com chave externa funcional, e o sistema atua apenas pelo lado interno, como um atuador adicional. Isso elimina o risco real de "ficar trancado fora de casa por causa de um bug" e abre espaço para experimentação livre, com baixo custo de falha. A aplicação funciona, em última instância, como objeto de conversa e demonstração técnica — uma "maçaneta que diz com qual celebridade você se parece" é um artefato social.

## 4. Fundamentação Teórica

### 4.1 Reconhecimento facial: pipeline e estado da arte

O reconhecimento facial moderno segue um pipeline em quatro estágios:

1. **Detecção facial**: localizar regiões de face em uma imagem. Métodos clássicos como Viola-Jones (2001) baseados em cascatas de Haar features foram superados por métodos baseados em CNN, sendo MTCNN (Zhang et al. 2016) e RetinaFace (Deng et al. 2019) as referências atuais.
2. **Alinhamento**: normalizar a face detectada (rotação, escala, recorte) usando landmarks (olhos, nariz, boca), de forma que faces de pessoas diferentes fiquem em poses comparáveis.
3. **Extração de embedding**: transformar a face alinhada em um vetor de dimensão fixa (tipicamente 128 ou 512). A propriedade central buscada é que duas imagens da mesma pessoa produzam embeddings próximos no espaço vetorial, enquanto pessoas diferentes produzam embeddings distantes. Os modelos de referência são:
  - **FaceNet** (Schroff et al. 2015): treinado com triplet loss para empurrar embeddings de mesma identidade para perto e diferentes para longe.
  - **ArcFace** (Deng et al. 2019): usa additive angular margin loss, atualmente referência em benchmarks como LFW e MegaFace.
4. **Comparação**: medir distância (cosseno ou euclidiana) entre o embedding da consulta e os embeddings de um banco de referência, retornando o mais próximo.

### 4.2 Dataset VGGFace2

O banco VGGFace2 (Cao et al. 2018) contém aproximadamente 9.131 identidades e 3.31 milhões de imagens, com variação significativa de pose, idade e iluminação. Foi originalmente publicado pela Universidade de Oxford e, embora tenha sofrido restrições de distribuição por questões de privacidade em 2020, continua amplamente disponível em mirrors (Kaggle, repositórios acadêmicos). É uma das referências mais usadas para treinamento e avaliação de modelos de reconhecimento facial.

### 4.3 Busca de vizinho mais próximo em espaços vetoriais

Com ~~9.000 embeddings de referência, a busca linear é viável (~~10ms), mas a literatura padroniza o uso de índices vetoriais como **FAISS** (Johnson et al. 2017) que aceleram a busca via clustering hierárquico (IVF) ou quantização de produto (PQ). Para o escopo deste projeto, FAISS oferece API simples em Python e desempenho mais que adequado.

### 4.4 Sistemas embarcados de baixo consumo

O ESP32 (Espressif) é o microcontrolador padrão para projetos IoT residenciais por ter WiFi/Bluetooth integrados, ampla comunidade e modos de deep sleep com consumo na ordem de microamperes. A variante ESP32-CAM adiciona câmera OV2640 (2MP) ao mesmo SoC.

O ciclo típico de baixo consumo combina:

- **Deep sleep**: ESP32 desligado quase por completo, mantendo apenas RTC e GPIOs configurados como wakeup source.
- **Wakeup por GPIO externa**: sensor PIR (HC-SR501) detecta movimento e levanta um pino, acordando o ESP32.
- **Driver de motor H-bridge** (DRV8833 ou equivalente): controla direção e velocidade do motor DC via PWM, com proteção contra correntes reversas.
- **Encoder rotativo no eixo do motor**: fornece feedback de posição (graus girados) para o microcontrolador, permitindo controle de malha fechada na rotação do thumbturn.

### 4.5 Fechaduras eletrônicas retrofit

A arquitetura de retrofit popularizada por August Home e Yale consiste em **acoplar um motor por dentro da porta, sem modificar a fechadura existente**. O motor gira o thumbturn (a borboleta interna que normalmente se opera com a mão), preservando o uso da chave pelo lado externo. Teardowns públicos da August Smart Lock revelam o uso de motor DC com gearbox controlado por microcontrolador STM32 via driver DRV8833 e feedback de posição por acelerômetro. Essa é a arquitetura adotada como referência neste projeto.

## 5. Metodologia

### 5.1 Abordagem geral: vertical-slice iterativo

O projeto adota uma metodologia de **vertical-slice iterativo**: cada iteração entrega um sistema funcional ponta-a-ponta, ainda que rudimentar, e iterações subsequentes substituem componentes específicos por versões mais sofisticadas. A justificativa é descobrir problemas de integração cedo, manter motivação alta e permitir parada precoce com um artefato demonstrável.

Quatro iterações são planejadas. Cada iteração termina com um conjunto de **critérios de aceite observáveis** que devem estar verdes antes de avançar.

### 5.2 Iteração 1 — MVP horroroso (~2 semanas)

**Objetivo**: validar o pipeline ponta-a-ponta com o mínimo de componentes reais.

**Componentes**:

- ESP32-CAM em protoboard, sem case.
- Botão físico como trigger (em vez de PIR).
- LED como "trava simulada" (em vez de motor).
- Display OLED 0.96" I2C.
- Servidor Python local (FastAPI) que retorna sempre o mesmo nome de celebridade.

**Aceite**: ao apertar o botão, o ESP32 captura uma foto, envia ao servidor, recebe resposta e mostra "Brad Pitt" no display em ≤ 5s; o LED acende por 3s simulando a destrava.

### 5.3 Iteração 2 — ML real (~3 semanas)

**Objetivo**: substituir o servidor dummy por um pipeline real de reconhecimento facial.

**Componentes adicionais**:

- Modelo: MTCNN para detecção + FaceNet (`facenet-pytorch`) para embedding.
- Banco de referência: 100 celebridades top do VGGFace2 (curadas manualmente por popularidade), com embeddings pré-computados.
- FAISS para busca do vizinho mais próximo.
- Servidor migrado para HuggingFace Spaces (FastAPI + Gradio).
- Servo motor SG90 acionando uma "porta de gaveta" (não a porta real), substituindo o LED.

**Aceite**: foto enviada de fato retorna uma celebridade plausível (validada manualmente em 20 testes com fotos próprias e de amigos); servo gira 90° quando a face é detectada com confiança suficiente.

### 5.4 Iteração 3 — Hardware definitivo (~3 semanas)

**Objetivo**: substituir componentes-rascunho pelos definitivos e adicionar gestão de energia.

**Componentes adicionais**:

- Motor DC com gearbox (6V, com torque suficiente para girar thumbturn de fechadura residencial padrão).
- Driver DRV8833.
- Encoder rotativo no eixo (ou acelerômetro MPU6050 como alternativa).
- Sensor PIR HC-SR501 substituindo botão.
- Deep sleep do ESP32 implementado, wakeup por GPIO do PIR.
- Banco VGGFace2 ampliado para o conjunto completo (~9.000 identidades) com índice FAISS.
- Primeira versão do case impresso em 3D (ainda funcional, sem refinamento estético).

**Aceite**: sistema fica em standby por 24h consumindo apenas o que PIR + ESP32 em deep sleep consomem; ao detectar movimento, executa o ciclo completo (foto → ML → motor) em ≤ 5s; motor gira o thumbturn de uma fechadura real montada em bancada.

### 5.5 Iteração 4 — Polimento (~2 semanas)

**Objetivo**: instalação em porta residencial real, refinamento estético e medição rigorosa.

**Componentes adicionais**:

- Case CAD refinado (paredes finas, encaixes precisos, cabeamento interno).
- Instalação física em uma porta interna residencial.
- Pack de baterias 18650 (2 ou 3 células em série) com proteção e carregador.
- Threshold de similaridade ajustado a partir de dados coletados nas iterações anteriores.
- Logging de eventos em arquivo local no servidor para análise posterior.

**Aceite**: autonomia medida ≥ 30 dias em uso típico (10 ativações/dia, simuladas via gatilhos manuais distribuídos no tempo); latência ponta-a-ponta ≤ 5s em pelo menos 90% das ativações; sistema instalado e usado durante o período de medição sem necessidade de reset manual.

### 5.6 Stack tecnológica


| Camada                | Tecnologia                                                                                |
| --------------------- | ----------------------------------------------------------------------------------------- |
| CAD                   | Fusion 360 (licença pessoal gratuita) ou FreeCAD                                          |
| Impressão 3D          | FDM com PLA ou PETG                                                                       |
| Microcontrolador      | ESP32-CAM (AI-Thinker)                                                                    |
| Framework embarcado   | Arduino (via PlatformIO) ou ESP-IDF                                                       |
| Linguagem do servidor | Python 3.11+                                                                              |
| Framework web         | FastAPI                                                                                   |
| Deep learning         | PyTorch + `facenet-pytorch`                                                               |
| Indexação vetorial    | FAISS                                                                                     |
| Hospedagem cloud      | HuggingFace Spaces (gratuito)                                                             |
| Versionamento         | Git, com repositório único (`argos`) contendo subpastas para firmware, server, cad e docs |


### 5.7 Critérios de avaliação transversais

Em todas as iterações, três métricas são acompanhadas:

- **Latência ponta-a-ponta**: medida em segundos do trigger ao resultado no display.
- **Consumo médio**: medido em mA com multímetro ou USB power meter, separadamente em standby e em ciclo ativo.
- **Acurácia perceptual do match**: avaliada qualitativamente — o sósia retornado é "razoável" segundo julgamento humano? Não há ground truth, então a métrica é necessariamente subjetiva.

## 6. Resultados Esperados

Por se tratar de um anteprojeto, esta seção antecipa os artefatos e medições esperadas ao fim das quatro iterações, com discussão preliminar de limitações previstas.

### 6.1 Artefatos esperados

- **Protótipo físico** instalado e funcional em uma porta interna residencial.
- **Repositório Git público** (`argos`) contendo:
  - Código de firmware (ESP32, C++/Arduino).
  - Código do servidor de inferência (Python).
  - Modelos CAD (arquivos `.f3d` ou `.FCStd` originais + STL exportado para impressão).
  - Scripts de processamento do dataset VGGFace2 (download, filtragem, geração de embeddings).
  - Este anteprojeto e um relatório final.
- **Vídeo de demonstração** mostrando o ciclo completo: pessoa se aproxima, PIR aciona, foto capturada, servidor responde, porta destrava, display mostra "você é igual ao [celebridade]".
- **Tabela de medições** com autonomia, latência e consumo por iteração.

### 6.2 Discussão antecipada de limitações

**Viés do dataset VGGFace2**: a base é fortemente enviesada para atores americanos brancos, com sub-representação de outras etnias, gêneros e regiões. O sósia retornado para pessoas fora desse perfil tende a ser sistematicamente "menos parecido" em termos perceptuais. Esse viés é uma propriedade conhecida da literatura e não será corrigido neste projeto.

**Sensibilidade a condições adversas**: iluminação noturna ou contraluz, óculos escuros, máscaras faciais e ângulos extremos degradam a qualidade do embedding. Em ambiente residencial interno, espera-se que essas condições sejam minoritárias, mas o sistema falhará em casos como "pessoa de boné e luz fraca atrás".

**Dependência de WiFi externa**: o servidor de inferência hospedado em cloud exige conexão à internet. Quedas de WiFi ou indisponibilidade do HuggingFace Spaces impedem o sistema de funcionar. Como a fechadura original está preservada, isso não impede o acesso à casa, mas elimina a função "smart".

**Privacidade**: fotos da pessoa que se aproxima da porta são enviadas a um servidor cloud. Embora o servidor seja controlado pelo autor do projeto, não há criptografia ponta-a-ponta nem mecanismo de exclusão automática, e o tráfego depende da confiança em HuggingFace Spaces. Em uma versão futura, processamento local seria mais privado.

**Não-segurança**: o sistema **não é** uma fechadura de segurança. Qualquer pessoa cuja face seja detectada destravará a porta. Essa propriedade é deliberada (escolha de escopo no início do projeto), mas limita o uso a portas internas não-críticas.

### 6.3 Métricas previstas

- Autonomia: 30–60 dias com 10 ativações/dia, usando 2x 18650 (~6.000mAh combinados a 3.7V).
- Latência: 3–5s ponta-a-ponta, dominada por (i) wakeup do ESP32 e estabilização da câmera (~~1s), (ii) captura e envio da foto (~~1s), (iii) inferência no servidor (~~1s), (iv) atualização de display e acionamento do motor (~~1s).
- Custo de hardware: R$ 250–400 dependendo de fornecedores e câmbio.

## 7. Considerações Finais

O Argos demonstra, como hipótese a ser verificada nas iterações, que a integração entre visão computacional, sistemas embarcados de baixo consumo e fabricação caseira por impressão 3D é hoje acessível a um indivíduo, dado o estado de comoditização dos componentes (~R$ 300 em hardware) e a disponibilidade pública de modelos pré-treinados e ferramentas de deploy.

A abordagem metodológica de **vertical-slice iterativo** é defendida como contrapeso ao risco mais comum em projetos multi-domínio: investir profundamente em uma das camadas antes de validar que ela se integra às demais. Cada iteração entrega um sistema demonstrável, e o conjunto de critérios de aceite por iteração funciona como um circuit breaker contra otimização prematura.

As limitações reconhecidas — viés do dataset, dependência de cloud, ausência de segurança real — são consequências diretas das escolhas de escopo e custo, não acidentes. Reconhecê-las explicitamente é parte do que diferencia um anteprojeto honesto de uma apresentação publicitária.

### 7.1 Trabalhos futuros

Após a conclusão das quatro iterações, vislumbram-se as seguintes extensões:

1. **Inferência on-device**: portar o modelo (quantizado, possivelmente MobileFaceNet) para o próprio ESP32-S3 ou para uma placa com NPU (Sipeed MaixCam, OpenMV), eliminando a dependência de WiFi e melhorando a privacidade.
2. **Modo de controle de acesso real**: adicionar enrollment de faces autorizadas e fazer a maçaneta destravar **apenas** para essas faces — passando do modo "diversão" ao modo "segurança", com todas as preocupações adicionais que isso implica.
3. **Feedback multimodal**: integrar um speaker pequeno com TTS para que o resultado seja anunciado em voz ("Você é igual ao Pedro Pascal"), adicionando dimensão sensorial ao artefato.
4. **Reconhecimento de múltiplas pessoas**: estender o pipeline para identificar todas as faces presentes no quadro simultaneamente, lidando com o caso de grupos chegando juntos à porta.
5. **Alimentação por painel solar**: substituir as baterias 18650 por um painel solar pequeno com supercapacitor, alcançando autonomia indefinida em residências com janela próxima.
6. **Dataset customizado**: complementar o VGGFace2 com celebridades brasileiras (atores, músicos, esportistas), reduzindo o viés cultural do banco original.

---

## Referências

- Cao, Q., Shen, L., Xie, W., Parkhi, O. M., & Zisserman, A. (2018). *VGGFace2: A dataset for recognising faces across pose and age*. IEEE International Conference on Automatic Face & Gesture Recognition.
- Deng, J., Guo, J., Niannan, X., & Zafeiriou, S. (2019). *ArcFace: Additive Angular Margin Loss for Deep Face Recognition*. CVPR.
- Johnson, J., Douze, M., & Jégou, H. (2017). *Billion-scale similarity search with GPUs*. arXiv:1702.08734 (FAISS).
- Schroff, F., Kalenichenko, D., & Philbin, J. (2015). *FaceNet: A Unified Embedding for Face Recognition and Clustering*. CVPR.
- Viola, P., & Jones, M. (2001). *Rapid Object Detection using a Boosted Cascade of Simple Features*. CVPR.
- Zhang, K., Zhang, Z., Li, Z., & Qiao, Y. (2016). *Joint Face Detection and Alignment using Multi-task Cascaded Convolutional Networks*. IEEE Signal Processing Letters.
- РУКИ. (2018). *August Smart Lock Teardown*. RUKI Journal, Medium.