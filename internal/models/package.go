package models

// интерфейс для упаковки заказа
type Packaging interface {
	GetName() string
	GetPrice() float64
	GetMaxWeight() float64
}

// пакет
type Bag struct{}

func (b Bag) GetName() string       { return "bag" }
func (b Bag) GetPrice() float64     { return 5 }
func (b Bag) GetMaxWeight() float64 { return 10 }

// коробка
type Box struct{}

func (bx Box) GetName() string       { return "box" }
func (bx Box) GetPrice() float64     { return 20 }
func (bx Box) GetMaxWeight() float64 { return 30 }

// пленка
type Film struct{}

func (f Film) GetName() string       { return "film" }
func (f Film) GetPrice() float64     { return 1 }
func (f Film) GetMaxWeight() float64 { return -1 }

// без упаковки
type NoPackaging struct{}

func (np NoPackaging) GetName() string       { return "none" }
func (np NoPackaging) GetPrice() float64     { return 0 }
func (np NoPackaging) GetMaxWeight() float64 { return -1 }
