package universe

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
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

func TestNew(t *testing.T) {
	assert := require.New(t)

	universe := NewUniverse(NewRect(0, 0, 200, 200))
	assert.NotNil(universe)
}

func TestAdd(t *testing.T) {
	assert := require.New(t)

	universe := NewUniverse(NewRect(0, 0, 200, 200))

	obj := FakeObject{id: 1}
	universe.Add(&obj)

	assert.Len(universe.objects, 1)
	assert.NotNil(obj.universe)
}

func TestDel(t *testing.T) {
	assert := require.New(t)

	universe := NewUniverse(NewRect(0, 0, 200, 200))

	obj1 := FakeObject{id: 1}
	obj2 := FakeObject{id: 2}
	obj3 := FakeObject{id: 3}
	universe.Add(&obj1)
	universe.Add(&obj2)
	universe.Add(&obj3)
	assert.Len(universe.objects, 3)

	universe.Del(&obj2)
	assert.Len(universe.objects, 2)
	assert.Nil(obj2.universe)

	assert.Equal(universe.objects[0].GetId(), uint64(1))
	assert.Equal(universe.objects[1].GetId(), uint64(3))

	universe.Del(&obj1)
	assert.Len(universe.objects, 1)
	universe.Del(&obj3)
	assert.Len(universe.objects, 0)
}

func TestProcessPhysics(t *testing.T) {
	assert := require.New(t)
	universe := NewUniverse(NewRect(0, 0, 200, 200))
	universe.ProcessPhysics()
	assert.Equal(universe.tik, uint64(1))
}

func TestObjectsProcessPhysics(t *testing.T) {
	assert := require.New(t)

	universe := NewUniverse(NewRect(0, 0, 200, 200))

	obj := FakeObject{id: 1}
	universe.Add(&obj)

	universe.ProcessPhysics()
	assert.Equal(universe.tik, uint64(1))
	assert.Equal(obj.processed, 1)
}
