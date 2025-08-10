package models

import "errors"

var (
	// General-type error in case if item price doesn't comply with format: single number or <multiplier>*<price>.
	ErrItemPriceWrongFormat = errors.New("price should be a single number or '<multiplier>*<price>'")

	// For cases when price isn't a single int or float number
	ErrItemPriceNotSingleIntOrFloat = errors.New("item price must be a single number of int or float types")

	// For cases when price-multiplier isn't a single integer
	ErrItemPriceMultiplierNotSingleInt = errors.New("item price multiplier must be a single integer")
)
