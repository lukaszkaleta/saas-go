module github.com/lukaszkaleta/saas-go/payment/integration/stripe

go 1.25.0

require (
	github.com/lukaszkaleta/saas-go/payment v0.2.284
	github.com/stripe/stripe-go/v85 v85.0.0
)

require github.com/lukaszkaleta/saas-go/universal v0.2.284 // indirect

replace github.com/lukaszkaleta/saas-go/payment => ../../../payment

replace github.com/lukaszkaleta/saas-go/universal => ../../../universal
