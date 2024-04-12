package main

import "github.com/mikezzb/steam-trading-crawler/types"

var ListingsHandler = &types.Handler{
	OnResult: func(result interface{}) {
		// save to db
		// save preview urls
	},
	OnError: func(err error) {
	},
	OnComplete: func() {
	}}

var TransactionsHandler = &types.Handler{
	OnResult: func(result interface{}) {
		// save to db
	},
	OnError: func(err error) {
	},
	OnComplete: func() {
	}}
