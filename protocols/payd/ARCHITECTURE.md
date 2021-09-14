# Architecture

payd uses a standard 3 layer architecture with clear separation between these layers of the application:

![3 Tier Layout](https://aspblogs.blob.core.windows.net/media/fredriknormen/WindowsLiveWriter/UsingWebServicesina3tierarchitecture_134F6/3tier_thumb.jpg "3 tier layout")

All domain objects are stored in the top level of the application enforcing Hexagonal structuring of the dependencies.

## Transports

This can contain one or more transport and is responsible for parsing requests and delivering responses via those specific transports.

At the moment there is one supported transport, HTTP, which implements the BIP-270 protocol.

The API then passes data down to the service layer.

## Service

This is where business logic lives and is agnostic entirely to the Transport layer above as well as the data layer below it.

It is responsible for validating data and enforcing the protocol rules.

## Data

Data stores each have their own top level package, named to match the store.

At the moment there is one supported data store, sqlite, but this will be extended to support more.

The data layer knows only about how to interact with the data store to store or retrieve data. This will be called by the service layer, but the service layer doesn't know or care about what data store it is interacting with.

## Dependency Injection

We use dependency injection throughout the application to:

1) Improve testability
2) Improve flexibility - as mentioned we will be implementing more data store and transports as required.

