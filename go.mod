module github.com/libsv/go-bn

go 1.17

require (
	github.com/libsv/go-bc v0.1.7
	github.com/libsv/go-bk v0.1.4
	github.com/libsv/go-bt/v2 v2.0.0-beta.9.0.20211021120434-0cd048b7ca09
	github.com/pkg/errors v0.9.1
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
)

require (
	github.com/theflyingcodr/govalidator v0.0.2 // indirect
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
)

replace github.com/libsv/go-bt/v2 => ../go-bt
