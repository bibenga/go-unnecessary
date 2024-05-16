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

	if len(universe.objects) != 1 {
		t.Errorf("object is not added")
	}
	if obj.universe == nil {
		t.Errorf("object.universe is nil")
	}
}

func TestDel(t *testing.T) {
	universe := NewUniverse(NewRect(0, 0, 200, 200))

	obj1 := FakeObject{id: 1}
	obj2 := FakeObject{id: 2}
	obj3 := FakeObject{id: 3}
	universe.Add(&obj1)
	universe.Add(&obj2)
	universe.Add(&obj3)
	if len(universe.objects) != 3 {
		t.Errorf("object is not added")
	}

	universe.Del(&obj2)
	if len(universe.objects) != 2 {
		t.Errorf("object is not removed")
	}
	if obj2.universe != nil {
		t.Errorf("object.universe is not nil")
	}
	if universe.objects[0].GetId() != 1 {
		t.Errorf("object1 is removed")
	}
	if universe.objects[1].GetId() != 3 {
		t.Errorf("object3 is removed")
	}

	universe.Del(&obj1)
	if len(universe.objects) != 1 {
		t.Errorf("object1 is not removed")
	}
	universe.Del(&obj3)
	if len(universe.objects) != 0 {
		t.Errorf("object3 is not removed")
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
