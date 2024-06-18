package cli

import (
	"HOMEWORK-1/internal/models"
	"HOMEWORK-1/internal/models/customErrors"
)

func parsePackaging(packagingType string) (models.Packaging, error) {
	switch packagingType {
	case "bag":
		return models.PackagingBag, nil
	case "box":
		return models.PackagingBox, nil
	case "film":
		return models.PackagingFilm, nil
	case "":
		return models.PackagingNone, nil
	default:
		return models.PackagingNone, customErrors.ErrInvalidPackaging
	}
}
