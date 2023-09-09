# Trafilea Challenge

[![Trafilea](https://mma.prnewswire.com/media/1959233/Trafilea_Logo.jpg?w=200)](https://www.trafilea.com/)

![Build Status](https://travis-ci.org/joemccann/dillinger.svg?branch=master)

A company is building an e-commerce platform for their coffee shop to sell products online and
needs to build a RESTful API to handle the user cart products and calculate the total and
discount of the order.

## Features

- Creating a cart
- Adding products to a cart
- Updating products quantities
- Create order applying discounts

## Installation

Trafilea Challenge requires [Go](https://go.dev/).

Install the dependencies:

```sh
make install
```

## Run app

In order to run the application, you have to run:

```sh
make run
```

## Test

In order to run tests:

```sh
make test
```

## Docker

In order to run the Dockerfile:

```sh
docker run -p 8080:8080 trafilea-challenge
```

## Considerations
- It is being used an inmemory storage represented by a map where the key is the UserID and value is the Cart. We assume that every user will have only ONE cart.
- More unit tests should be added to have a 100% coverage
