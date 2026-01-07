package shared

// Gender represents gender for measurement purposes.
type Gender string

const (
	GenderMen   Gender = "men"
	GenderWomen Gender = "women"
)

// BodyMeasurement represents body measurements for tailoring.
// All measurements are in centimeters unless otherwise specified.
type BodyMeasurement struct {
	gender Gender
	name   string // optional name like "My Baju Kurung Size"

	// Upper body (cm)
	bust          *float64
	chest         *float64
	waist         *float64
	hip           *float64
	shoulderWidth *float64
	armLength     *float64

	// Lower body (cm)
	inseam  *float64
	outseam *float64
	thigh   *float64

	// Additional (cm/kg)
	neck   *float64
	wrist  *float64
	height *float64
	weight *float64

	notes string
}

// BodyMeasurementParams contains parameters for creating a BodyMeasurement.
type BodyMeasurementParams struct {
	Gender        Gender
	Name          string
	Bust          *float64
	Chest         *float64
	Waist         *float64
	Hip           *float64
	ShoulderWidth *float64
	ArmLength     *float64
	Inseam        *float64
	Outseam       *float64
	Thigh         *float64
	Neck          *float64
	Wrist         *float64
	Height        *float64
	Weight        *float64
	Notes         string
}

// NewBodyMeasurement creates a new BodyMeasurement.
func NewBodyMeasurement(params BodyMeasurementParams) BodyMeasurement {
	return BodyMeasurement{
		gender:        params.Gender,
		name:          params.Name,
		bust:          params.Bust,
		chest:         params.Chest,
		waist:         params.Waist,
		hip:           params.Hip,
		shoulderWidth: params.ShoulderWidth,
		armLength:     params.ArmLength,
		inseam:        params.Inseam,
		outseam:       params.Outseam,
		thigh:         params.Thigh,
		neck:          params.Neck,
		wrist:         params.Wrist,
		height:        params.Height,
		weight:        params.Weight,
		notes:         params.Notes,
	}
}

// Getters
func (m BodyMeasurement) Gender() Gender          { return m.gender }
func (m BodyMeasurement) Name() string            { return m.name }
func (m BodyMeasurement) Bust() *float64          { return m.bust }
func (m BodyMeasurement) Chest() *float64         { return m.chest }
func (m BodyMeasurement) Waist() *float64         { return m.waist }
func (m BodyMeasurement) Hip() *float64           { return m.hip }
func (m BodyMeasurement) ShoulderWidth() *float64 { return m.shoulderWidth }
func (m BodyMeasurement) ArmLength() *float64     { return m.armLength }
func (m BodyMeasurement) Inseam() *float64        { return m.inseam }
func (m BodyMeasurement) Outseam() *float64       { return m.outseam }
func (m BodyMeasurement) Thigh() *float64         { return m.thigh }
func (m BodyMeasurement) Neck() *float64          { return m.neck }
func (m BodyMeasurement) Wrist() *float64         { return m.wrist }
func (m BodyMeasurement) Height() *float64        { return m.height }
func (m BodyMeasurement) Weight() *float64        { return m.weight }
func (m BodyMeasurement) Notes() string           { return m.notes }

// HasUpperBodyMeasurements returns true if upper body measurements are set.
func (m BodyMeasurement) HasUpperBodyMeasurements() bool {
	return m.bust != nil || m.chest != nil || m.waist != nil || m.hip != nil ||
		m.shoulderWidth != nil || m.armLength != nil
}

// HasLowerBodyMeasurements returns true if lower body measurements are set.
func (m BodyMeasurement) HasLowerBodyMeasurements() bool {
	return m.inseam != nil || m.outseam != nil || m.thigh != nil
}

// IsComplete returns true if all standard measurements are set.
func (m BodyMeasurement) IsComplete() bool {
	return m.bust != nil && m.waist != nil && m.hip != nil && m.height != nil
}

// BMI calculates BMI if height and weight are available.
// Returns nil if either is not set.
func (m BodyMeasurement) BMI() *float64 {
	if m.height == nil || m.weight == nil {
		return nil
	}
	heightM := *m.height / 100 // convert cm to m
	bmi := *m.weight / (heightM * heightM)
	return &bmi
}

// UpperBodyDifference returns bust-waist difference (useful for garment sizing).
func (m BodyMeasurement) UpperBodyDifference() *float64 {
	if m.bust == nil || m.waist == nil {
		return nil
	}
	diff := *m.bust - *m.waist
	return &diff
}

// HipToWaistRatio calculates hip to waist ratio.
func (m BodyMeasurement) HipToWaistRatio() *float64 {
	if m.hip == nil || m.waist == nil {
		return nil
	}
	ratio := *m.waist / *m.hip
	return &ratio
}
