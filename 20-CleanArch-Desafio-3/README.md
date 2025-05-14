# Sistema de Pedidos - Clean Architecture

Este projeto implementa - a título de exercício - um sistema simples de pedidos utilizando os princípios da Clean Architecture, com múltiplas interfaces de comunicação: REST, GraphQL e gRPC.

## Estrutura do Projeto

O projeto segue a estrutura da Clean Architecture:

- **Entity**: Contém as regras de negócio e entidades principais
- **Use Cases**: Implementa os casos de uso da aplicação
- **Adapters/Interface**: Implementa as interfaces de comunicação (REST, GraphQL, gRPC)
- **Frameworks & Drivers**: Implementa a infraestrutura (banco de dados, mensageria)

## Funcionalidades

- Criação de pedidos (REST, GraphQL, gRPC)
- Listagem de pedidos (REST, GraphQL, gRPC)
- Eventos publicados no RabbitMQ quando um pedido é criado ou listado

## Pré-requisitos

- Docker e Docker Compose
- Go 1.19+ (para desenvolvimento local)

## Como Executar

### Usando Docker Compose

```bash
# Clone o repositório
git clone https://github.com/dudapinto/pos-golang.git
cd pos-golang/20-CleanArch-Desafio-3

# Inicie os containers
docker-compose up -d
```

A aplicação estará disponível em:
- API REST: http://localhost:8000
- API GraphQL: http://localhost:8080
- API gRPC: localhost:50051

### Desenvolvimento Local

```bash
# Configure as variáveis de ambiente
export DB_DRIVER=mysql
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=root
export DB_NAME=orders
export WEB_SERVER_PORT=:8000
export GRPC_SERVER_PORT=50051
export GRAPHQL_SERVER_PORT=8080

# Execute a aplicação
go run cmd/ordersystem/main.go
```

## Endpoints

### REST API

#### Criar Pedido
```http
POST http://localhost:8000/order
Content-Type: application/json

{
    "id": "123",
    "price": 100.5,
    "tax": 0.5
}
```

#### Listar Pedidos
```http
GET http://localhost:8000/list_orders
Content-Type: application/json
```

### GraphQL API

Acesse o playground em: http://localhost:8080

#### Criar Pedido
```graphql
mutation createOrder {
  createOrder(input: {id: "123", Price: 100.5, Tax: 0.5}) {
    id
    Price
    Tax
    FinalPrice
  }
}
```

#### Listar Pedidos
```graphql
query listOrders {
  orders {
    id
    Price
    Tax
    FinalPrice
  }
}
```

### gRPC API

Utilize o cliente gRPC [grpcurl](https://github.com/fullstorydev/grpcurl).

#### Criar Pedido
```
grpcurl -plaintext -d '{"id": "123", "price": 100.5, "tax": 0.5}' localhost:50051 pb.OrderService/CreateOrder
```

#### Listar Pedidos
```
grpcurl -plaintext localhost:50051 pb.ListOrdersService/ListOrders
```

## Banco de Dados

O projeto utiliza MySQL como banco de dados. A migração inicial é executada automaticamente quando o container é iniciado.

## Mensageria

O projeto utiliza RabbitMQ para publicar eventos quando um pedido é criado ou listado. 
Você pode acessar o painel de administração do RabbitMQ em http://localhost:15672 (usuário: guest, senha: guest).

## Testes

Para executar os testes:

```bash
go test ./...
```