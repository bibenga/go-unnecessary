package universe

import (
	"fmt"
	"testing"
)

type FakeObject struct {
	id        uint64
	universe  *Universe
	processed int
}

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
	universe := NewUniverse(NewRect(0, 0, 200, 200))
	if universe == nil {
		t.Errorf("universe is nil")
	}
}

func TestAdd(t *testing.T) {
	universe := NewUniverse(NewRect(0, 0, 200, 200))

	obj := FakeObject{id: 1}
	universe.Add(&obj)
	if obj.universe == nil {
		t.Errorf("obj.universe is nil")
	}
}

func TestProcessPhysics(t *testing.T) {
	universe := NewUniverse(NewRect(0, 0, 200, 200))
	universe.ProcessPhysics()
	if universe.tik != 1 {
		t.Errorf("universe.ProcessPhysics is not called")
	}

}

func TestObjectsProcessPhysics(t *testing.T) {
	universe := NewUniverse(NewRect(0, 0, 200, 200))

	obj := FakeObject{id: 1}
	universe.Add(&obj)

	universe.ProcessPhysics()
	if universe.tik != 1 {
		t.Errorf("universe.ProcessPhysics is not called")
	}

	if obj.processed != 1 {
		t.Errorf("obj.ProcessPhysics is not called")
	}
}
