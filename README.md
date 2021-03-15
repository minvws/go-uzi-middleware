# UZI middleware for go

Proficient UZI pass reader for go projects.

The UZI card is part of an authentication mechanism for medical staff and doctors working in the Netherlands. The cards are distributed by the CIBG. More information and the relevant client software can be found at www.uziregister.nl (in Dutch).

go-uzi-middleware is a simple and functional module which allows you to use the UZI cards as authentication mechanism. 

go-uzi-middleware is available under the EU PL licence. It was created early 2021 during the COVID19 campaign as part of the vaccination registration project BRBA for the ‘Ministerie van Volksgezondheid, Welzijn & Sport, programma Realisatie Digitale Ondersteuning.’

Questions and contributions are welcome via [GitHub](https://github.com/minvws/pUZI-php/issues).

## Requirements

* go >= 1.13 or higher

## Installation

### Composer

```sh
go get github.com/minvws/go-uzi-middleware
```

## Usage

See the `examples` directory for usage. Basically you will need to add the middleware to `net/http` or `negroni`. If the middelware 
validates the UZI certificate correctly, it will return an UZI user with the card's information inside the request context:

```go
func handlerFunc(w http.ResponseWriter, r *http.Request) {
	uziUser := r.Context().Value(uzi.UziContext("uzi")).(*uzi.UziUser)

	fmt.Println("The current logged in UZI user is: ", uziUser.SurName)
}
```

Important:

 - The UZI middleware can ONLY be used on TLS connections.
 - In order to use the middelware, you will need to have a UZI CA certifcate.

## Options
You can add several options when instantiating the UZI middleware.

```go
type Options struct {
	// Strict check on the card
	StrictCACheck bool      
	
	// A list of allowed card types. Note that leaving this empty will NOT authenticate ANY card.
	AllowedTypes  []UziType
    
	// A list of allowed roles. Note that leaving this empty will NOT authenticate ANY card.
	AllowedRoles  []UziRole

	// When set, the authentication will call this function when authentication fails for any reason.
	ErrorHandler errorHandler 
	
	// When set to true, debug output will be send to the log package.
	Debug        bool         
}
```



## Contributing

1. Fork the Project

2. Create a Feature Branch

3. (Recommended) Run the Test Suite

    ```sh
    make test
    ```
4. Send us a Pull Request

