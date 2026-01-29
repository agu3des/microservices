# 📦 Microservices Order System

Sistema distribuído de pedidos composto por 3 serviços: Order, Payment e Shipping.

## 📋 Pré-requisitos
- Go instalado (1.21+)
- Docker e Kubernetes (Kind ou Minikube) rodando
- Git configurado

## 🚀 Como Rodar o Projeto

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

## 📊 Regras de Negócio Implementadas
1. **Pagamento:** Pedidos acima de R$ 1.000,00 são recusados.
2. **Quantidade:** Pedidos acima de 50 itens são bloqueados.
3. **Entrega:**
   - Base: 1 dia.
   - Adicional: +1 dia a cada 5 itens.