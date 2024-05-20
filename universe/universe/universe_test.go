package universe

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type FakeObject struct {
	id        uint64
	universe  *Universe
	processed int
}

var _ fmt.Stringer = &FakeObject{}
var _ IObject = &FakeObject{}

func (fake FakeObject) String() string {
	return fmt.Sprintf("FakeObject-%d", fake.id)
}

func (fake *FakeObject) GetId() uint64 {
	return fake.id
}

func (fake *FakeObject) GetUniverse() *Universe {
	return fake.universe
}

func (fake *FakeObject) SetUniverse(universe *Universe) {
	fake.universe = universe
}

func (fake *FakeObject) ProcessPhysics() {
	fake.processed++
}

type MockObject struct {
	mock.Mock
}

var _ fmt.Stringer = &MockObject{}
var _ IObject = &MockObject{}

func (f *MockObject) GetId() uint64 {
	args := f.Called()
	return uint64(args.Int(0))
}

func (f *MockObject) GetUniverse() *Universe {
	args := f.Called()
	u := args.Get(0)
	if u == nil {
		return nil
	}
	return u.(*Universe)
}

func (f *MockObject) SetUniverse(universe *Universe) {
	f.Called(universe)
}

func (f *MockObject) ProcessPhysics() {
	f.Called()
}

func TestNew(t *testing.T) {
	universe := NewUniverse(NewRect(0, 0, 200, 200))
	assert.NotNil(t, universe)
}

func TestAdd(t *testing.T) {
	universe := NewUniverse(NewRect(0, 0, 200, 200))

	obj := FakeObject{id: 1}
	universe.Add(&obj)

	assert.Len(t, universe.objects, 1)
	assert.NotNil(t, obj.universe)
}

func TestDel(t *testing.T) {
	universe := NewUniverse(NewRect(0, 0, 200, 200))

	obj1 := FakeObject{id: 1}
	obj2 := FakeObject{id: 2}
	obj3 := FakeObject{id: 3}
	universe.Add(&obj1)
	universe.Add(&obj2)
	universe.Add(&obj3)
	assert.Len(t, universe.objects, 3)

	universe.Del(&obj2)
	assert.Len(t, universe.objects, 2)
	assert.Nil(t, obj2.universe)

	assert.Equal(t, universe.objects[0].GetId(), uint64(1))
	assert.Equal(t, universe.objects[1].GetId(), uint64(3))

	universe.Del(&obj1)
	assert.Len(t, universe.objects, 1)
	universe.Del(&obj3)
	assert.Len(t, universe.objects, 0)
}

func TestProcessPhysics(t *testing.T) {
	universe := NewUniverse(NewRect(0, 0, 200, 200))
	universe.ProcessPhysics()
	assert.Equal(t, universe.tik, uint64(1))
}

func TestObjectsProcessPhysics(t *testing.T) {
	universe := NewUniverse(NewRect(0, 0, 200, 200))

	obj := new(MockObject)
	obj.On("GetUniverse").Return(nil)
	obj.On("SetUniverse", universe).Return()
	obj.On("ProcessPhysics").Return()

	universe.Add(obj)

	universe.ProcessPhysics()
	assert.Equal(t, universe.tik, uint64(1))
	obj.AssertExpectations(t)
}
