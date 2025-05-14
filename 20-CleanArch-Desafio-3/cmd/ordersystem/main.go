package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/dudapinto/pos-golang/20-CleanArch-Desafio-3/configs"
	"github.com/dudapinto/pos-golang/20-CleanArch-Desafio-3/internal/event"
	"github.com/dudapinto/pos-golang/20-CleanArch-Desafio-3/internal/event/handler"
	"github.com/dudapinto/pos-golang/20-CleanArch-Desafio-3/internal/infra/graph"
	"github.com/dudapinto/pos-golang/20-CleanArch-Desafio-3/internal/infra/grpc/pb"
	"github.com/dudapinto/pos-golang/20-CleanArch-Desafio-3/internal/infra/grpc/service"
	"github.com/dudapinto/pos-golang/20-CleanArch-Desafio-3/internal/infra/web/webserver"
	"github.com/dudapinto/pos-golang/20-CleanArch-Desafio-3/pkg/events"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rabbitMQChannel := getRabbitMQChannel()

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})
	eventDispatcher.Register("OrdersListed", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	listOrdersUseCase := NewListOrdersUseCase(db, eventDispatcher)

	webserver := webserver.NewWebServer(configs.WebServerPort)
	orderCreatedEvent := &event.OrderCreated{}
	ordersListedEvent := &event.OrdersListed{}
	webOrderHandler := NewWebOrderHandler(eventDispatcher, db, orderCreatedEvent, ordersListedEvent)
	webserver.AddHandler("/order", webOrderHandler.Create)
	webserver.AddHandler("/list_orders", webOrderHandler.List)

	fmt.Println("Starting web server on port", configs.WebServerPort)
	go webserver.Start()

	grpcServer := grpc.NewServer()
	createOrderService := service.NewOrderService(*createOrderUseCase)
	listOrdersService := service.NewListOrdersService(*listOrdersUseCase)
	pb.RegisterOrderServiceServer(grpcServer, createOrderService)
	pb.RegisterListOrdersServiceServer(grpcServer, listOrdersService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", configs.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
		ListOrdersUseCase:  *listOrdersUseCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", configs.GraphQLServerPort)
	http.ListenAndServe(":"+configs.GraphQLServerPort, nil)
}

func getRabbitMQChannel() *amqp.Channel {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}
