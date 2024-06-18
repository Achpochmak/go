package models

type Packaging struct {
	Name      string
	Price     float64
	MaxWeight float64
}

var (
	PackagingBag  = Packaging{Name: "bag", Price: 5, MaxWeight: 10}
	PackagingBox  = Packaging{Name: "box", Price: 20, MaxWeight: 30}
	PackagingFilm = Packaging{Name: "film", Price: 1, MaxWeight: -1}
	PackagingNone = Packaging{Name: "None", Price: 0, MaxWeight: -1}
)
