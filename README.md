### 1\. Resumo do Projeto: "E-commerce Microservices"

**Objetivo e Utilidade**
Este projeto simula o backend de um e-commerce utilizando **Arquitetura de Microsserviços**. Sua utilidade é demonstrar como sistemas distribuídos se comunicam de forma robusta e desacoplada.
O fluxo principal implementado é a **Criação de Pedidos**, onde o sistema:

1.  Recebe um pedido.
2.  Valida o pagamento (Microsserviço **Payment**).
3.  Calcula o prazo de entrega baseado na quantidade de itens (Microsserviço **Shipping**).
4.  Retorna o resultado consolidado ao cliente.

**Tecnologias Utilizadas**

  * **Linguagem:** Go (Golang) 1.22+.
  * **Comunicação:** gRPC (Protocol Buffers) para alta performance entre serviços.
  * **Arquitetura:** Hexagonal (Ports & Adapters) para isolar o domínio da infraestrutura.
  * **Infraestrutura:** Docker (Containerização) e Kubernetes (Orquestração).
  * **Banco de Dados:** MySQL (Persistência dos pedidos).

**Histórico de Alterações**
Desde o início da nossa sessão, realizamos as seguintes implementações:

1.  **Atualização de Contrato (Proto):** Adicionamos o campo `delivery_days` no arquivo `order.proto` e regeneramos os códigos Go.
2.  **Lógica de Shipping (Domínio):** Criamos a regra de negócio onde `Prazo = 1 dia + (TotalItens / 5)`.
3.  **Integração Order -\> Shipping:** Implementamos o Adapter gRPC no serviço de Order para consultar o Shipping e salvar o prazo.
4.  **Correção de Build:** Resolvemos problemas de cache do Go Proxy usando `direct` e configuramos o `main.go` para tratar erros de conexão.
5.  **Cliente de Teste Rico:** Criamos um CLI (`client.go`) colorido e inteligente que valida regras de negócio e exibe erros amigáveis.

-----

### 2\. Arquivo de Documentação: `README.md`

Salve este conteúdo como **README.md** na raiz do projeto. Ele contém tudo o que é necessário para rodar.

```markdown
# Microservices Order System

Sistema distribuído de pedidos composto por 3 serviços: Order, Payment e Shipping.

## Pré-requisitos
- Go instalado (1.21+)
- Docker e Kubernetes (Kind ou Minikube) rodando
- Git configurado

## Como Rodar o Projeto

### Passo 1: Atualizar Contratos (Proto)
Se houver mudanças no `.proto`, gere o código e suba para o Git:
1. Vá para a pasta `microservices-proto`.
2. Rode o comando `protoc` (ver lista de comandos).
3. Faça o `git push` para a branch main.

### Passo 2: Build e Deploy dos Serviços

**Serviço de Shipping (Calculadora de Frete):**
1. Navegue até `microservices/shipping`.
2. Atualize dependências: `go mod tidy`.
3. Build da imagem: `docker build -t shipping:latest .`
4. Deploy: `kubectl rollout restart deployment shipping-deployment`

**Serviço de Order (Orquestrador):**
*Nota: Este serviço depende do Payment e Shipping.*
1. Navegue até `microservices/order`.
2. Force a atualização do proto: `$env:GOPROXY="direct"; go get -u github.com/agu3des/microservices-proto@latest`.
3. Build da imagem: `docker build --no-cache -t order:latest .`
4. Deploy: `kubectl rollout restart deployment order-deployment`

### Passo 3: Executar Testes
Aguarde os pods reiniciarem (`kubectl get pods`).

Execute o cliente de teste automatizado:
`go run cmd/test_client/client.go`

## Regras de Negócio Implementadas
1. **Pagamento:** Pedidos acima de R$ 1.000,00 são recusados.
2. **Quantidade:** Pedidos acima de 50 itens são bloqueados.
3. **Entrega:**
   - Base: 1 dia.
   - Adicional: +1 dia a cada 5 itens.
```

-----

### 3\. Arquivo de Comandos: `comandos_uteis.sh`

Salve este arquivo como **comandos.txt** ou **cheatsheet.sh**. Ele contém os comandos exatos que usamos no terminal.

```bash
# ========================================================
# GERAÇÃO DE CÓDIGO GPRC (PROTOBUF)
# ========================================================
# Deve ser rodado na raiz da pasta 'microservices-proto'

# Gerar código para o serviço Order
protoc -I=order --go_out=golang/order --go_opt=paths=source_relative --go-grpc_out=golang/order --go-grpc_opt=paths=source_relative order.proto

# Subir alterações para o GitHub (Essencial para o Docker ver as mudanças)
git add .
git commit -m "Atualizando contratos proto"
git push origin main

# ========================================================
# GERENCIAMENTO DE DEPENDÊNCIAS (GO MODULES)
# ========================================================
# Deve ser rodado dentro da pasta de cada serviço (ex: microservices/order)

# Limpar cache local do Go (se der erro estranho)
go clean -modcache

# Forçar atualização ignorando o cache do Google (A "Opção Nuclear")
# No Windows Powershell:
$env:GOPROXY="direct"
go get -u github.com/agu3des/microservices-proto@latest
go mod tidy

# No Linux/Mac:
GOPROXY=direct go get -u github.com/agu3des/microservices-proto@latest
go mod tidy

# ========================================================
# DOCKER E KUBERNETES
# ========================================================

# Construir imagem ignorando cache (Para garantir que baixa o código novo)
docker build --no-cache -t order:latest .
docker build --no-cache -t shipping:latest .

# Reiniciar os Pods para pegar a nova imagem
kubectl rollout restart deployment order-deployment
kubectl rollout restart deployment shipping-deployment

# Verificar status dos pods
kubectl get pods

# Ver logs de um serviço específico (para debugar erros internos)
kubectl logs -f deployment/order-deployment
kubectl logs -f deployment/shipping-deployment

# ========================================================
# TESTES
# ========================================================
# Rodar na raiz de microservices/order

# Executar o cliente de teste rico
go run cmd/test_client/client.go
```