// Package gwc is simple HTTP client based on premise to use client side
// middlewares.
// It is based on cliware library that defines basic type for client side
// middlewares and series of middlewares implemented in cliware-middlewares
// package.
//
// Basic idea behind this client is to create simple, yet easily composable
// HTTP client. Client does provide simple Request and Response types, but
// they are only thin wrappers around same type in http package used for
// middleware composition.
package gwc // import "go.delic.rs/gwc"
