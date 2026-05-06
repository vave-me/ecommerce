package e2e

import (
	"github.com/cucumber/godog"
	"github.com/go-openapi/strfmt"
	"github.com/stackus/errors"
	"middleman/users/usersclient"
	"middleman/users/usersclient/models"
	"middleman/users/usersclient/user"
)

type usersContext struct {
	*suiteContext
	client       *usersclient.Users
	fetchedUsers bool
}

var _ featureContext = (*usersContext)(nil)

func newUsersContext(sc *suiteContext) featureContext {
	return &usersContext{
		suiteContext: sc,
		client:       usersclient.New(sc.transport, strfmt.Default),
	}
}

func (c *usersContext) register(ctx *godog.ScenarioContext) {
	ctx.Step(`^a valid user$`, c.aValidUser)
	ctx.Step(`^I create the user called "([^"]*)"$`, c.iCreateTheUserCalled)
	ctx.Step(`^(?:ensure |expect )?a user called "([^"]*)" (?:to )?exists?$`, c.expectAUserCalledToExist)
	ctx.Step(`^(?:ensure |expect )?no user called "([^"]*)" (?:to )?exists?$`, c.expectNoUserCalledToExist)
}

func (c *usersContext) reset() {
	c.users = make(map[string]string)
	c.fetchedUsers = false
	c.truncate("users.users")
	c.truncate("users.products")
	c.truncate("users.events")
	c.truncate("users.snapshots")
	c.truncate("users.inbox")
	c.truncate("users.outbox")

}

func (c *usersContext) aValidUser() {
	// noop
}

func (c *usersContext) expectAUserCalledToExist(name string) error {
	if !c.fetchedUsers {
		err := c.fetchUsers()
		if err != nil {
			return err
		}
	}

	if _, exists := c.users[name]; !exists {
		return errors.ErrNotFound.Msgf("the user `%s` does not exist", name)
	}
	return nil
}

func (c *usersContext) expectNoUserCalledToExist(name string) error {
	if !c.fetchedUsers {
		err := c.fetchUsers()
		if err != nil {
			return err
		}
	}

	if _, exists := c.users[name]; exists {
		return errors.ErrNotFound.Msgf("the user `%s` does exist", name)
	}
	return nil
}

func (c *usersContext) iCreateTheUserCalled(name string) {
	resp, err := c.client.User.CreateUser(user.NewCreateUserParams().WithBody(&models.UserspbCreateUserRequest{
		Email:    "anywhere",
		Password: name,
	}))
	if err != nil {
		c.lastErr = err
		return
	}

	c.users[name] = resp.Payload.ID
}

func (c *usersContext) fetchUsers() error {
	resp, err := c.client.User.GetUsers(user.NewGetUsersParams())
	if err != nil {
		return err
	}

	for _, s := range resp.Payload.Users {
		c.users[s.FirstName] = s.ID
	}

	c.fetchedUsers = true

	return nil
}
