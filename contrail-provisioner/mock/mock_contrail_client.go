package mock

import (
	contrail "github.com/Juniper/contrail-go-api"
)

// MockContrailClient enables unit tests for code that call contrail.Client
// methods. User can provide each method separately for each test.
type MockContrailClient struct {
	CreateMock             func(contrail.IObject) error
	UpdateMock             func(contrail.IObject) error
	DeleteMock             func(contrail.IObject) error
	DeleteByUuidMock       func(string, string) error
	FindByUuidMock         func(string, string) (contrail.IObject, error)
	UuidByNameMock         func(string, string) (string, error)
	FqNameByUuidMock       func(string) ([]string, error)
	FindByNameMock         func(string, string) (contrail.IObject, error)
	ListMock               func(string) ([]contrail.ListResult, error)
	ListByParentMock       func(string, string) ([]contrail.ListResult, error)
	ListDetailMock         func(string, []string) ([]contrail.IObject, error)
	ListDetailByParentMock func(string, string, []string) ([]contrail.IObject, error)
	ReadListResultMock     func(string, *contrail.ListResult) (contrail.IObject, error)
}

// GetDefaultMockContrailClient returns pointer to MockContrailClient instance
// that has default implementations of each method, that return no errors, empty
// lists, nil pointers etc.
func GetDefaultMockContrailClient() *MockContrailClient {
	return &MockContrailClient{
		CreateMock:             func(contrail.IObject) error { return nil },
		UpdateMock:             func(contrail.IObject) error { return nil },
		DeleteMock:             func(contrail.IObject) error { return nil },
		DeleteByUuidMock:       func(string, string) error { return nil },
		FindByUuidMock:         func(string, string) (contrail.IObject, error) { return nil, nil },
		UuidByNameMock:         func(string, string) (string, error) { return "", nil },
		FqNameByUuidMock:       func(string) ([]string, error) { return []string{}, nil },
		FindByNameMock:         func(string, string) (contrail.IObject, error) { return nil, nil },
		ListMock:               func(string) ([]contrail.ListResult, error) { return []contrail.ListResult{}, nil },
		ListByParentMock:       func(string, string) ([]contrail.ListResult, error) { return []contrail.ListResult{}, nil },
		ListDetailMock:         func(string, []string) ([]contrail.IObject, error) { return []contrail.IObject{}, nil },
		ListDetailByParentMock: func(string, string, []string) ([]contrail.IObject, error) { return []contrail.IObject{}, nil },
		ReadListResultMock:     func(string, *contrail.ListResult) (contrail.IObject, error) { return nil, nil },
	}
}

func (c *MockContrailClient) Create(ptr contrail.IObject) error {
	return c.CreateMock(ptr)
}
func (c *MockContrailClient) Update(ptr contrail.IObject) error {
	return c.UpdateMock(ptr)
}
func (c *MockContrailClient) DeleteByUuid(typename, uuid string) error {
	return c.DeleteByUuidMock(typename, uuid)
}
func (c *MockContrailClient) Delete(ptr contrail.IObject) error {
	return c.DeleteMock(ptr)
}
func (c *MockContrailClient) FindByUuid(typename string, uuid string) (contrail.IObject, error) {
	return c.FindByUuidMock(typename, uuid)
}
func (c *MockContrailClient) UuidByName(typename string, fqn string) (string, error) {
	return c.UuidByNameMock(typename, fqn)
}
func (c *MockContrailClient) FQNameByUuid(uuid string) ([]string, error) {
	return c.FqNameByUuidMock(uuid)
}
func (c *MockContrailClient) FindByName(typename string, fqn string) (contrail.IObject, error) {
	return c.FindByNameMock(typename, fqn)
}
func (c *MockContrailClient) List(typename string) ([]contrail.ListResult, error) {
	return c.ListMock(typename)
}
func (c *MockContrailClient) ListByParent(typename string, parentID string) ([]contrail.ListResult, error) {
	return c.ListByParentMock(typename, parentID)
}
func (c *MockContrailClient) ListDetail(typename string, fields []string) ([]contrail.IObject, error) {
	return c.ListDetailMock(typename, fields)
}
func (c *MockContrailClient) ListDetailByParent(typename string, parentID string, fields []string) ([]contrail.IObject, error) {
	return c.ListDetailByParentMock(typename, parentID, fields)
}
func (c *MockContrailClient) ReadListResult(typename string, result *contrail.ListResult) (contrail.IObject, error) {
	return c.ReadListResultMock(typename, result)
}
