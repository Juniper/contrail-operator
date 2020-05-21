package mock

import (
	contrail "github.com/Juniper/contrail-go-api"
)

type MockContrailClient struct {
	createMock             func(contrail.IObject) error
	updateMock             func(contrail.IObject) error
	deleteMock             func(contrail.IObject) error
	deleteByUuidMock       func(string, string) error
	findByUuidMock         func(string, string) (contrail.IObject, error)
	uuidByNameMock         func(string, string) (string, error)
	fqNameByUuidMock       func(string) ([]string, error)
	findByNameMock         func(string, string) (contrail.IObject, error)
	listMock               func(string) ([]contrail.ListResult, error)
	listByParentMock       func(string, string) ([]contrail.ListResult, error)
	listDetailMock         func(string, []string) ([]contrail.IObject, error)
	listDetailByParentMock func(string, string, []string) ([]contrail.IObject, error)
	readListResultMock     func(string, *contrail.ListResult) (contrail.IObject, error)
}

func GetDefaultMockContrailClient() MockContrailClient {
	return MockContrailClient{
		createMock:             func(contrail.IObject) error { return nil },
		updateMock:             func(contrail.IObject) error { return nil },
		deleteMock:             func(contrail.IObject) error { return nil },
		deleteByUuidMock:       func(string, string) error { return nil },
		findByUuidMock:         func(string, string) (contrail.IObject, error) { return nil, nil },
		uuidByNameMock:         func(string, string) (string, error) { return "", nil },
		fqNameByUuidMock:       func(string) ([]string, error) { return []string{}, nil },
		findByNameMock:         func(string, string) (contrail.IObject, error) { return nil, nil },
		listMock:               func(string) ([]contrail.ListResult, error) { return []contrail.ListResult{}, nil },
		listByParentMock:       func(string, string) ([]contrail.ListResult, error) { return []contrail.ListResult{}, nil },
		listDetailMock:         func(string, []string) ([]contrail.IObject, error) { return []contrail.IObject{}, nil },
		listDetailByParentMock: func(string, string, []string) ([]contrail.IObject, error) { return []contrail.IObject{}, nil },
		readListResultMock:     func(string, *contrail.ListResult) (contrail.IObject, error) { return nil, nil },
	}
}

func (c *MockContrailClient) Create(ptr contrail.IObject) error {
	return c.createMock(ptr)
}
func (c *MockContrailClient) Update(ptr contrail.IObject) error {
	return c.updateMock(ptr)
}
func (c *MockContrailClient) DeleteByUuid(typename, uuid string) error {
	return c.deleteByUuidMock(typename, uuid)
}
func (c *MockContrailClient) Delete(ptr contrail.IObject) error {
	return c.deleteMock(ptr)
}
func (c *MockContrailClient) FindByUuid(typename string, uuid string) (contrail.IObject, error) {
	return c.findByUuidMock(typename, uuid)
}
func (c *MockContrailClient) UuidByName(typename string, fqn string) (string, error) {
	return c.uuidByNameMock(typename, fqn)
}
func (c *MockContrailClient) FQNameByUuid(uuid string) ([]string, error) {
	return c.fqNameByUuidMock(uuid)
}
func (c *MockContrailClient) FindByName(typename string, fqn string) (contrail.IObject, error) {
	return c.findByNameMock(typename, fqn)
}
func (c *MockContrailClient) List(typename string) ([]contrail.ListResult, error) {
	return c.listMock(typename)
}
func (c *MockContrailClient) ListByParent(typename string, parentID string) ([]contrail.ListResult, error) {
	return c.listByParentMock(typename, parentID)
}
func (c *MockContrailClient) ListDetail(typename string, fields []string) ([]contrail.IObject, error) {
	return c.listDetailMock(typename, fields)
}
func (c *MockContrailClient) ListDetailByParent(typename string, parentID string, fields []string) ([]contrail.IObject, error) {
	return c.listDetailByParentMock(typename, parentID, fields)
}
func (c *MockContrailClient) ReadListResult(typename string, result *contrail.ListResult) (contrail.IObject, error) {
	return c.readListResultMock(typename, result)
}
