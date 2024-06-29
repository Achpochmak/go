package models

// Упаковка
type Packaging struct {
	Name      string
	Price     float64
	MaxWeight float64
}

// Пакет
func NewBag() Packaging {
	return Packaging{
		Name:      "bag",
		Price:     5,
		MaxWeight: 10,
	}
}

// Коробка
func NewBox() Packaging {
	return Packaging{
		Name:      "box",
		Price:     20,
		MaxWeight: 30,
	}
}

// Пленка
func NewFilm() Packaging {
	return Packaging{
		Name:      "film",
		Price:     1,
		MaxWeight: -1,
	}
}

// нет упаковки
func NewNoPackaging() Packaging {
	return Packaging{
		Name:      "none",
		Price:     0,
		MaxWeight: -1,
	}
}

func (p Packaging) GetName() string {
	return p.Name
}

func (p Packaging) GetPrice() float64 {
	return p.Price
}

func (p Packaging) GetMaxWeight() float64 {
	return p.MaxWeight
}
