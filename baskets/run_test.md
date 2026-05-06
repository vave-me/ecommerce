To run all the application tests for the Shopping Baskets module, you would use the following command:

`go test ./baskets/internal/application`

To run only one test eg. RemoveItem you need to add `-run "RemoveItem$"` to the command

`go test ./baskets/internal/application -run "RemoveItem$"`


To run only the `NoProduct` subtest for the `RemoveItem` test we need to add `RemoveItem/NoProduct$`

`go test ./baskets/internal/application -run "RemoveItem/NoProduct$"`