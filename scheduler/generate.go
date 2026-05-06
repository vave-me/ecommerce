package scheduler

//go:generate buf generate

//go:generate mockery --quiet --dir ./schedulerspb -r --all --inpackage --case underscore
//go:generate mockery --quiet --dir ./internal -r --all --inpackage --case underscore

//go:generate swagger generate client -q -f ./internal/rest/api.swagger.json -c schedulerclient -m schedulerclient/models --with-flatten=remove-unused
