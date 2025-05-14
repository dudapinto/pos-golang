//go:build wireinject
// +build wireinject

package main

import (
	"database/sql"

	"github.com/dudapinto/pos-golang/20-CleanArch-Desafio-3/internal/entity"
	"github.com/dudapinto/pos-golang/20-CleanArch-Desafio-3/internal/event"
	"github.com/dudapinto/pos-golang/20-CleanArch-Desafio-3/internal/infra/database"
	"github.com/dudapinto/pos-golang/20-CleanArch-Desafio-3/internal/infra/web"
	"github.com/dudapinto/pos-golang/20-CleanArch-Desafio-3/internal/usecase"
	"github.com/dudapinto/pos-golang/20-CleanArch-Desafio-3/pkg/events"
	"github.com/google/wire"
)

var setOrderRepositoryDependency = wire.NewSet(
	database.NewOrderRepository,
	wire.Bind(new(entity.OrderRepositoryInterface), new(*database.OrderRepository)),
)

var setEventDispatcherDependency = wire.NewSet(
	events.NewEventDispatcher,
	event.NewOrderCreated,
	event.NewOrdersListed,
	wire.Bind(new(events.EventInterface), new(*event.OrderCreated)),
	wire.Bind(new(events.EventDispatcherInterface), new(*events.EventDispatcher)),
)

var setOrderCreatedEvent = wire.NewSet(
	event.NewOrderCreated,
	wire.Bind(new(events.EventInterface), new(*event.OrderCreated)),
)

var setOrdersListedEvent = wire.NewSet(
	event.NewOrdersListed,
	wire.Bind(new(events.EventInterface), new(*event.OrdersListed)),
)

func NewCreateOrderUseCase(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *usecase.CreateOrderUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		setOrderCreatedEvent,
		usecase.NewCreateOrderUseCase,
	)
	return &usecase.CreateOrderUseCase{}
}

func NewListOrdersUseCase(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *usecase.ListOrdersUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		setOrdersListedEvent,
		usecase.NewListOrdersUseCase,
	)
	return &usecase.ListOrdersUseCase{}
}

func NewWebOrderHandler(eventDispatcher events.EventDispatcherInterface, db *sql.DB, orderCreatedEvent events.EventInterface, ordersListedEvent events.EventInterface) *web.WebOrderHandler {
	wire.Build(
		setOrderRepositoryDependency,
		web.NewWebOrderHandler,
	)
	return &web.WebOrderHandler{}
}
