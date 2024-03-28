# GO-iPay - iPay.ua Payment Gate Client

[![GoDoc](https://godoc.org/github.com/stremovskyy/go-ipay?status.svg)](https://godoc.org/github.com/stremovskyy/go-ipay)

**go-ipay** is GO client for the iPay.ua Payment Gate API.

## Jump to

* [Installation](#Installation)
* [Usage](#Usage)
* [API examples](#API-examples)
* [License](#License)
* [Contributing](#Contributing)
* [Authors](#Authors)
* [Acknowledgments](#Acknowledgments)
* [TODO](#TODO)

## Installation

```bash
go get github.com/stremovskyy/go-ipay
```

## Usage

### Create client

create a new iPay.ua client with default options

```go
client := go_ipay.NewDefaultClient()
```

### Using a merchant account

Merchant account is a structure that contains all the necessary information for the client to work with the iPay.ua API.

```go
	merchant := &go_ipay.Merchant{
		Name:        private.MerchantName,
		MerchantID:  private.MerchantID,
		MerchantKey: private.MerchantKey,
	}
```

## API examples

See [examples](examples) for more examples.

## License

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

## Acknowledgments

- [GoDoc](https://godoc.org/github.com/stremovskyy/gofondy)

## Author

* **Anton Stremovskyy** - *Initial work* - [stremovskyy](https://github.com/stremovskyy)

### Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.
Please make sure to update tests as appropriate.

## TODO

- [ ] Add tests
- [ ] Add more API methods
