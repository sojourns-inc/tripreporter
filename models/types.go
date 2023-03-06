package models

import (
	"time"
)

type RouteOfAdministration int64

const ( // TODO: check if complete
	RoaOther RouteOfAdministration = iota
	RoaOral
	RoaBuccal
	RoaRectal
	RoaInhaled
	RoaSublabial
	RoaIntranasal
	RoaSublingual
	RoaOtherInjection
	RoaBuccalInjection
	RoaIntravenousInjection
	RoaSubcutanousInjection
	RoaIntramuscularInjection
)

type DosageUnit int64

const (
	DosageUnknown DosageUnit = iota
	DosageMicrograms
	DosageMilligrams
	DosageGrams
)

type DisplayUnit int64

const (
	UnitUnknown DisplayUnit = iota
	UnitMetric
	UnitImperial
)

type UserMedication struct { // TODO: How do we feel about this being unencrypted?
	Medications []Medication `json:"medications,omitempty"`
}

type UserPreferences struct { // TODO: How do we feel about this being unencrypted?
	Timezone     time.Location `json:"timezone,omitempty"`      // Default: Europe/London
	HeightFormat DisplayUnit   `json:"height_format,omitempty"` // Default: UnitMetric // Display height in centimeters or feet + inches
	WeightFormat DisplayUnit   `json:"weight_format,omitempty"` // Default: UnitMetric // Display weight in kilograms or pounds
}
