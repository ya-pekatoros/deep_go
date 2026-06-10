package main

import (
	"bytes"
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type Option func(*GamePerson)

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		clear(person.name[:])
		copy(person.name[:], name)
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.x = int32(x)
		person.y = int32(y)
		person.z = int32(z)
	}
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.gold = uint32(gold)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.stats = (person.stats &^ (GamePersonStats(stats10BitMask) << statsManaShift)) |
			(GamePersonStats(mana&stats10BitMask) << statsManaShift)
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.stats = (person.stats &^ (GamePersonStats(stats10BitMask) << statsHealthShift)) |
			(GamePersonStats(health&stats10BitMask) << statsHealthShift)
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.socialStats = (person.socialStats &^ (GamePersonSocialStats(social4BitMask) << socialRespectShift)) |
			(GamePersonSocialStats(respect&social4BitMask) << socialRespectShift)
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.stats = (person.stats &^ (GamePersonStats(stats4BitMask) << statsStrengthShift)) |
			(GamePersonStats(strength&stats4BitMask) << statsStrengthShift)
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.stats = (person.stats &^ (GamePersonStats(stats4BitMask) << statsExperienceShift)) |
			(GamePersonStats(experience&stats4BitMask) << statsExperienceShift)
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.stats = (person.stats &^ (GamePersonStats(stats4BitMask) << statsLevelShift)) |
			(GamePersonStats(level&stats4BitMask) << statsLevelShift)
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.socialStats |= GamePersonSocialStats(social1BitMask) << socialHouseShift
	}
}

func WithWeapon() func(*GamePerson) {
	return func(person *GamePerson) {
		person.hasWeapon = true
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.socialStats |= GamePersonSocialStats(social1BitMask) << socialFamilyShift
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.socialStats = (person.socialStats &^ (GamePersonSocialStats(social2BitMask) << socialTypeShift)) |
			(GamePersonSocialStats(personType&social2BitMask) << socialTypeShift)
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

const (
	stats10BitMask = 0x3FF
	stats4BitMask  = 0xF

	statsHealthShift     = 0
	statsManaShift       = 10
	statsExperienceShift = 20
	statsLevelShift      = 24
	statsStrengthShift   = 28
)

const (
	social4BitMask = 0xF
	social2BitMask = 0x3
	social1BitMask = 0x1

	socialRespectShift = 0
	socialHouseShift   = 4
	socialFamilyShift  = 5
	socialTypeShift    = 6
)

// GamePersonStats битовая маска для хранения всех характеристик персонажа в одном uint32
// в обратном порядке
// [0-1000] - Health (10 бит)
// [0-1000] - Mana (10 бит)
// [0-10] - Experience (4 бита)
// [0-10] - Level (4 бита)
// [0-10] - Strength (4 бита)

type GamePersonStats uint32

// GamePersonSocialStats битовая маска для хранения всех социальных характеристик персонажа в одном uint8
// в обратном порядке
// [0-10] - Respect (4 бита)
// [0-1] - HasHouse (1 бит)
// [0-1] - HasFamily (1 бит)
// [0-3] - Type (2 бита)

type GamePersonSocialStats uint8

type GamePerson struct {
	x           int32                 // 4 байта
	y           int32                 // 4 байта
	z           int32                 // 4 байта
	gold        uint32                // 4 байта
	name        [42]byte              // 42 байта
	socialStats GamePersonSocialStats // 1 байт
	hasWeapon   bool                  // 1 байт
	stats       GamePersonStats       // 4 байта
	// 4 + 4 + 4 + 4 + 42 + 1 + 1 + 4 = 64 байта
	// выравнивание по 4 байтам
}

func NewGamePerson(options ...Option) GamePerson {
	person := GamePerson{}
	for _, option := range options {
		option(&person)
	}
	return person
}

func (p *GamePerson) Name() string {
	n := bytes.IndexByte(p.name[:], 0)
	if n == -1 {
		n = len(p.name)
	}
	return string(p.name[:n])
}

func (p *GamePerson) X() int {
	return int(p.x)
}

func (p *GamePerson) Y() int {
	return int(p.y)
}

func (p *GamePerson) Z() int {
	return int(p.z)
}

func (p *GamePerson) Gold() int {
	return int(p.gold)
}

func (p *GamePerson) Mana() int {
	return int((p.stats >> statsManaShift) & stats10BitMask)
}

func (p *GamePerson) Health() int {
	return int((p.stats >> statsHealthShift) & stats10BitMask)
}

func (p *GamePerson) Respect() int {
	return int((p.socialStats >> socialRespectShift) & social4BitMask)
}

func (p *GamePerson) Strength() int {
	return int((p.stats >> statsStrengthShift) & stats4BitMask)
}

func (p *GamePerson) Experience() int {
	return int((p.stats >> statsExperienceShift) & stats4BitMask)
}

func (p *GamePerson) Level() int {
	return int((p.stats >> statsLevelShift) & stats4BitMask)
}

func (p *GamePerson) HasHouse() bool {
	return ((p.socialStats >> socialHouseShift) & social1BitMask) != 0
}

func (p *GamePerson) HasWeapon() bool {
	return p.hasWeapon
}

func (p *GamePerson) HasFamily() bool {
	return ((p.socialStats >> socialFamilyShift) & social1BitMask) != 0
}

func (p *GamePerson) Type() int {
	return int((p.socialStats >> socialTypeShift) & social2BitMask)
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BuilderGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)
	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamily())
	assert.False(t, person.HasWeapon())
	assert.Equal(t, personType, person.Type())
}
