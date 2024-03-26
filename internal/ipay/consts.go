package ipay

const (
	ApiVersion = "1.28"

	baseUrl = "https://tokly.ipay.ua"

	ApiUrl    = baseUrl + "/api"
	ApiXMLUrl = baseUrl + "/api302"
)

type Lang string

const (
	LangUk Lang = "ua"
	LangEn Lang = "en"
)

type CurrencyKey string

const (
	CurrencyUAH CurrencyKey = "UAH"
	CurrencyUSD CurrencyKey = "USD"
	CurrencyEUR CurrencyKey = "EUR"
)
