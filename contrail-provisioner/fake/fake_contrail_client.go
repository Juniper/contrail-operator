package fake

import (
	contrail "github.com/Juniper/contrail-go-api"
)

// FakeContrailClient enables unit tests for code that call contrail.Client
// methods. User can provide each method separately for each test.
type FakeContrailClient struct {
	CreateFake             func(contrail.IObject) error
	UpdateFake             func(contrail.IObject) error
	DeleteFake             func(contrail.IObject) error
	DeleteByUuidFake       func(string, string) error
	FindByUuidFake         func(string, string) (contrail.IObject, error)
	UuidByNameFake         func(string, string) (string, error)
	FqNameByUuidFake       func(string) ([]string, error)
	FindByNameFake         func(string, string) (contrail.IObject, error)
	ListFake               func(string) ([]contrail.ListResult, error)
	ListByParentFake       func(string, string) ([]contrail.ListResult, error)
	ListDetailFake         func(string, []string) ([]contrail.IObject, error)
	ListDetailByParentFake func(string, string, []string) ([]contrail.IObject, error)
	ReadListResultFake     func(string, *contrail.ListResult) (contrail.IObject, error)
}

// GetDefaultFakeContrailClient returns pointer to FakeContrailClient instance
// that has default implementations of each method, that return no errors, empty
// lists, nil pointers etc.
func GetDefaultFakeContrailClient() *FakeContrailClient {
	return &FakeContrailClient{
		CreateFake:             func(contrail.IObject) error { return nil },
		UpdateFake:             func(contrail.IObject) error { return nil },
		DeleteFake:             func(contrail.IObject) error { return nil },
		DeleteByUuidFake:       func(string, string) error { return nil },
		FindByUuidFake:         func(string, string) (contrail.IObject, error) { return nil, nil },
		UuidByNameFake:         func(string, string) (string, error) { return "", nil },
		FqNameByUuidFake:       func(string) ([]string, error) { return []string{}, nil },
		FindByNameFake:         func(string, string) (contrail.IObject, error) { return nil, nil },
		ListFake:               func(string) ([]contrail.ListResult, error) { return []contrail.ListResult{}, nil },
		ListByParentFake:       func(string, string) ([]contrail.ListResult, error) { return []contrail.ListResult{}, nil },
		ListDetailFake:         func(string, []string) ([]contrail.IObject, error) { return []contrail.IObject{}, nil },
		ListDetailByParentFake: func(string, string, []string) ([]contrail.IObject, error) { return []contrail.IObject{}, nil },
		ReadListResultFake:     func(string, *contrail.ListResult) (contrail.IObject, error) { return nil, nil },
	}
}

func (c *FakeContrailClient) Create(ptr contrail.IObject) error {
	return c.CreateFake(ptr)
}
func (c *FakeContrailClient) Update(ptr contrail.IObject) error {
	return c.UpdateFake(ptr)
}
func (c *FakeContrailClient) DeleteByUuid(typename, uuid string) error {
	return c.DeleteByUuidFake(typename, uuid)
}
func (c *FakeContrailClient) Delete(ptr contrail.IObject) error {
	return c.DeleteFake(ptr)
}
func (c *FakeContrailClient) FindByUuid(typename string, uuid string) (contrail.IObject, error) {
	return c.FindByUuidFake(typename, uuid)
}
func (c *FakeContrailClient) UuidByName(typename string, fqn string) (string, error) {
	return c.UuidByNameFake(typename, fqn)
}
func (c *FakeContrailClient) FQNameByUuid(uuid string) ([]string, error) {
	return c.FqNameByUuidFake(uuid)
}
func (c *FakeContrailClient) FindByName(typename string, fqn string) (contrail.IObject, error) {
	return c.FindByNameFake(typename, fqn)
}
func (c *FakeContrailClient) List(typename string) ([]contrail.ListResult, error) {
	return c.ListFake(typename)
}
func (c *FakeContrailClient) ListByParent(typename string, parentID string) ([]contrail.ListResult, error) {
	return c.ListByParentFake(typename, parentID)
}
func (c *FakeContrailClient) ListDetail(typename string, fields []string) ([]contrail.IObject, error) {
	return c.ListDetailFake(typename, fields)
}
func (c *FakeContrailClient) ListDetailByParent(typename string, parentID string, fields []string) ([]contrail.IObject, error) {
	return c.ListDetailByParentFake(typename, parentID, fields)
}
func (c *FakeContrailClient) ReadListResult(typename string, result *contrail.ListResult) (contrail.IObject, error) {
	return c.ReadListResultFake(typename, result)
}
