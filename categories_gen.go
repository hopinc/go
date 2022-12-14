package hop

// Code generated by generate_categories.go; DO NOT EDIT.

//go:generate go run generate_categories.go

// ClientCategoryChannelsTokens is an auto-generated struct which is used to allow for simple categorisation of the APIs.
// It is public since it may be desired to store a reference to this somewhere, however, do NOT create a instance of this
// directly. Instead, call NewClient and then go to the field Channels.Tokens.
type ClientCategoryChannelsTokens struct {
	c clientDoer
}

// ClientCategoryChannels is an auto-generated struct which is used to allow for simple categorisation of the APIs.
// It is public since it may be desired to store a reference to this somewhere, however, do NOT create a instance of this
// directly. Instead, call NewClient and then go to the field Channels.
type ClientCategoryChannels struct {
	c clientDoer

	Tokens *ClientCategoryChannelsTokens
}

func newChannels(c clientDoer) *ClientCategoryChannels {
	return &ClientCategoryChannels{
		c:      c,
		Tokens: &ClientCategoryChannelsTokens{c},
	}
}

// ClientCategoryIgniteGateways is an auto-generated struct which is used to allow for simple categorisation of the APIs.
// It is public since it may be desired to store a reference to this somewhere, however, do NOT create a instance of this
// directly. Instead, call NewClient and then go to the field Ignite.Gateways.
type ClientCategoryIgniteGateways struct {
	c clientDoer
}

// ClientCategoryIgniteDeployments is an auto-generated struct which is used to allow for simple categorisation of the APIs.
// It is public since it may be desired to store a reference to this somewhere, however, do NOT create a instance of this
// directly. Instead, call NewClient and then go to the field Ignite.Deployments.
type ClientCategoryIgniteDeployments struct {
	c clientDoer
}

// ClientCategoryIgniteContainers is an auto-generated struct which is used to allow for simple categorisation of the APIs.
// It is public since it may be desired to store a reference to this somewhere, however, do NOT create a instance of this
// directly. Instead, call NewClient and then go to the field Ignite.Containers.
type ClientCategoryIgniteContainers struct {
	c clientDoer
}

// ClientCategoryIgnite is an auto-generated struct which is used to allow for simple categorisation of the APIs.
// It is public since it may be desired to store a reference to this somewhere, however, do NOT create a instance of this
// directly. Instead, call NewClient and then go to the field Ignite.
type ClientCategoryIgnite struct {
	c clientDoer

	Gateways    *ClientCategoryIgniteGateways
	Deployments *ClientCategoryIgniteDeployments
	Containers  *ClientCategoryIgniteContainers
}

func newIgnite(c clientDoer) *ClientCategoryIgnite {
	return &ClientCategoryIgnite{
		c:           c,
		Gateways:    &ClientCategoryIgniteGateways{c},
		Deployments: &ClientCategoryIgniteDeployments{c},
		Containers:  &ClientCategoryIgniteContainers{c},
	}
}

// ClientCategoryPipeRooms is an auto-generated struct which is used to allow for simple categorisation of the APIs.
// It is public since it may be desired to store a reference to this somewhere, however, do NOT create a instance of this
// directly. Instead, call NewClient and then go to the field Pipe.Rooms.
type ClientCategoryPipeRooms struct {
	c clientDoer
}

// ClientCategoryPipe is an auto-generated struct which is used to allow for simple categorisation of the APIs.
// It is public since it may be desired to store a reference to this somewhere, however, do NOT create a instance of this
// directly. Instead, call NewClient and then go to the field Pipe.
type ClientCategoryPipe struct {
	c clientDoer

	Rooms *ClientCategoryPipeRooms
}

func newPipe(c clientDoer) *ClientCategoryPipe {
	return &ClientCategoryPipe{
		c:     c,
		Rooms: &ClientCategoryPipeRooms{c},
	}
}

// ClientCategoryProjectsTokens is an auto-generated struct which is used to allow for simple categorisation of the APIs.
// It is public since it may be desired to store a reference to this somewhere, however, do NOT create a instance of this
// directly. Instead, call NewClient and then go to the field Projects.Tokens.
type ClientCategoryProjectsTokens struct {
	c clientDoer
}

// ClientCategoryProjectsSecrets is an auto-generated struct which is used to allow for simple categorisation of the APIs.
// It is public since it may be desired to store a reference to this somewhere, however, do NOT create a instance of this
// directly. Instead, call NewClient and then go to the field Projects.Secrets.
type ClientCategoryProjectsSecrets struct {
	c clientDoer
}

// ClientCategoryProjects is an auto-generated struct which is used to allow for simple categorisation of the APIs.
// It is public since it may be desired to store a reference to this somewhere, however, do NOT create a instance of this
// directly. Instead, call NewClient and then go to the field Projects.
type ClientCategoryProjects struct {
	c clientDoer

	Tokens  *ClientCategoryProjectsTokens
	Secrets *ClientCategoryProjectsSecrets
}

func newProjects(c clientDoer) *ClientCategoryProjects {
	return &ClientCategoryProjects{
		c:       c,
		Tokens:  &ClientCategoryProjectsTokens{c},
		Secrets: &ClientCategoryProjectsSecrets{c},
	}
}

// ClientCategoryRegistryImages is an auto-generated struct which is used to allow for simple categorisation of the APIs.
// It is public since it may be desired to store a reference to this somewhere, however, do NOT create a instance of this
// directly. Instead, call NewClient and then go to the field Registry.Images.
type ClientCategoryRegistryImages struct {
	c clientDoer
}

// ClientCategoryRegistry is an auto-generated struct which is used to allow for simple categorisation of the APIs.
// It is public since it may be desired to store a reference to this somewhere, however, do NOT create a instance of this
// directly. Instead, call NewClient and then go to the field Registry.
type ClientCategoryRegistry struct {
	c clientDoer

	Images *ClientCategoryRegistryImages
}

func newRegistry(c clientDoer) *ClientCategoryRegistry {
	return &ClientCategoryRegistry{
		c:      c,
		Images: &ClientCategoryRegistryImages{c},
	}
}

// ClientCategoryUsersMe is an auto-generated struct which is used to allow for simple categorisation of the APIs.
// It is public since it may be desired to store a reference to this somewhere, however, do NOT create a instance of this
// directly. Instead, call NewClient and then go to the field Users.Me.
type ClientCategoryUsersMe struct {
	c clientDoer
}

// ClientCategoryUsers is an auto-generated struct which is used to allow for simple categorisation of the APIs.
// It is public since it may be desired to store a reference to this somewhere, however, do NOT create a instance of this
// directly. Instead, call NewClient and then go to the field Users.
type ClientCategoryUsers struct {
	c clientDoer

	Me *ClientCategoryUsersMe
}

func newUsers(c clientDoer) *ClientCategoryUsers {
	return &ClientCategoryUsers{
		c:  c,
		Me: &ClientCategoryUsersMe{c},
	}
}
