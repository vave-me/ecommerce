package shipping

//go:generate buf generate

//go:generate mockery --quiet --dir ./shippingpb -r --all --inpackage --case underscore
//go:generate mockery --quiet --dir ./internal -r --all --inpackage --case underscore

//go:generate swagger generate client -q -f ./internal/rest/api.swagger.json -c shippingclient -m shippingclient/models --with-flatten=remove-unused
