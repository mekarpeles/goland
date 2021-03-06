// Unit is a non-static game object
package game

import (
	"fmt"
)

const (
	DEFAULT_HP = 10
)

type Unit struct {
	*GameObject
	*Inventory

	Level     int
	Hp, HpMax int
}

func NewUnit(name string) *Unit {
	u := &Unit{Level: 1,
		Hp:    DEFAULT_HP,
		HpMax: DEFAULT_HP,
	}
	u.Inventory = NewInventory()
	u.GameObject = NewGameObject(name)

	return u
}

func (u Unit) String() string {
	return fmt.Sprintf("%s: Hp: %d(%d) %s", u.Name, u.Hp, u.HpMax, u.GameObject)
}
