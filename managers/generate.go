package managers

//go:generate buf generate

//go:generate mockery --quiet --dir ./managerspb -r --all --inpackage --case underscore
//go:generate mockery --quiet --dir ./internal -r --all --inpackage --case underscore

//go:generate swagger generate client -q -f ./internal/rest/api.swagger.json -c managersclient -m managersclient/models --with-flatten=remove-unused
