package main

import (
	"context"
	"errors"
	"github.com/okcredit/go-common/log"
)

func main() {
	ctx := context.Background()
	ctx = log.WithLogger(
		ctx,
		log.NewZapLogger(
			log.WithNamespace("example"),
			log.WithLevel(log.INFO),
			log.WithFields("global_field", "qwerty"),
		),
	)

	log.Printf("hello, world")

	log.Info(ctx, "hello, world")
	doSomething(ctx)

	// add context
	ctx2 := log.Derive(ctx, log.WithNamespace("derive"), log.WithFields("name", "Aditya"))

	log.Info(ctx, "this is the original logger")
	log.Info(ctx2, "this is the logger with additional context")

	log.Info(ctx, "this is an info msg with extra fields", "mobile", "7760747507", "merchant_id", "aditya")
	log.Debug(ctx, "this is a debug message")
	log.Error(ctx, errors.New("chaos of the universe"), "this is an error message")
}

func doSomething(ctx context.Context) {
	ctx = log.Derive(ctx, log.WithNamespace("something"), log.WithLevel(log.DEBUG), log.WithFields("global_field", "qwerty"))
	log.Debug(ctx, "doing something")
}
