# Module Interfaces

## Prerequisites
* [Application CLI](.//interfaces/cli.md)
* [Application REST Interface](.//interfaces/cli.md)

## Synopsis

This document details how to build CLI and REST interfaces for a module.
- [CLI](#cli)
  + [Transaction Commands](#tx-commands)
  + [Query Commands](#query-commands)
- [REST](#rest)
  + [Request Types](#request-types)
  + [Request Handlers](#request-handlers)
  + [Register Routes](#register-routes)

## CLI

One of the main interfaces for an application is the [command-line interface](../interfaces/cli.md). This entrypoint created by the application developer will add commands from the application's modules to make [**transactions**](../core/transactions.md) and [**queries**]().  

### Transaction Commands 

### Query Commands

## REST

Application users may find the most intuitive methods of interfacing with the application are web services that use HTTP requests (e.g. a web wallet like [Lunie.io](lunie.io)). Thus, application developers will also use REST Routes to route HTTP requests to the application's modules.

### Request Types


### Request Handlers

### Register Routes
