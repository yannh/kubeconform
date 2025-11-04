package sarif

// NewLocationWithPhysicalLocation creates a Location with a PhysicalLocation
func NewLocationWithPhysicalLocation(physicalLocation *PhysicalLocation) *Location {
	return NewLocation().WithPhysicalLocation(physicalLocation)
}
