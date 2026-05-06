package notifications

//go:generate buf generate

//go:generate mockery --quiet --dir ./notificationspb -r --all --inpackage --case underscore
//go:generate mockery --quiet --dir ./internal -r --all --inpackage --case underscore

//go:generate swagger generate client -q -f ./internal/rest/api.swagger.json -c notificationsclient -m notificationsclient/models --with-flatten=remove-unused
