package cause_test

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/alextanhongpin/errors/cause"
)

// Real-world example: Healthcare Patient Record System
type PatientRecord struct {
	PatientID        string                  `json:"patient_id"`
	MedicalNumber    string                  `json:"medical_number"`
	PersonalInfo     PersonalInfo            `json:"personal_info"`
	MedicalHistory   []MedicalEntry          `json:"medical_history"`
	Allergies        []Allergy               `json:"allergies"`
	Medications      []Medication            `json:"current_medications"`
	EmergencyContact MedicalEmergencyContact `json:"emergency_contact"`
	Insurance        InsuranceInfo           `json:"insurance"`
	Vitals           VitalSigns              `json:"vitals,omitempty"`
}

type PersonalInfo struct {
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	DateOfBirth time.Time `json:"date_of_birth"`
	Gender      string    `json:"gender"`
	BloodType   string    `json:"blood_type"`
	Height      float64   `json:"height_cm"`
	Weight      float64   `json:"weight_kg"`
	PhoneNumber string    `json:"phone_number"`
	Email       string    `json:"email,omitempty"`
	Address     Address   `json:"address"`
}

type MedicalEntry struct {
	Date         time.Time  `json:"date"`
	Condition    string     `json:"condition"`
	Diagnosis    string     `json:"diagnosis"`
	Treatment    string     `json:"treatment"`
	DoctorID     string     `json:"doctor_id"`
	Severity     string     `json:"severity"`
	Notes        string     `json:"notes,omitempty"`
	FollowUpDate *time.Time `json:"follow_up_date,omitempty"`
}

type Allergy struct {
	Allergen      string    `json:"allergen"`
	Reaction      string    `json:"reaction"`
	Severity      string    `json:"severity"`
	DiagnosedDate time.Time `json:"diagnosed_date"`
	Notes         string    `json:"notes,omitempty"`
}

type Medication struct {
	Name         string     `json:"name"`
	Dosage       string     `json:"dosage"`
	Frequency    string     `json:"frequency"`
	StartDate    time.Time  `json:"start_date"`
	EndDate      *time.Time `json:"end_date,omitempty"`
	PrescribedBy string     `json:"prescribed_by"`
	Purpose      string     `json:"purpose"`
	SideEffects  []string   `json:"side_effects,omitempty"`
}

type MedicalEmergencyContact struct {
	Name         string  `json:"name"`
	Relationship string  `json:"relationship"`
	PhoneNumber  string  `json:"phone_number"`
	Email        string  `json:"email,omitempty"`
	Address      Address `json:"address"`
}

type InsuranceInfo struct {
	Provider     string    `json:"provider"`
	PolicyNumber string    `json:"policy_number"`
	GroupNumber  string    `json:"group_number,omitempty"`
	ValidFrom    time.Time `json:"valid_from"`
	ValidTo      time.Time `json:"valid_to"`
	Copay        float64   `json:"copay"`
	Deductible   float64   `json:"deductible"`
}

type VitalSigns struct {
	BloodPressureSystolic  int       `json:"bp_systolic"`
	BloodPressureDiastolic int       `json:"bp_diastolic"`
	HeartRate              int       `json:"heart_rate"`
	Temperature            float64   `json:"temperature_celsius"`
	OxygenSaturation       int       `json:"oxygen_saturation"`
	RespiratoryRate        int       `json:"respiratory_rate"`
	RecordedAt             time.Time `json:"recorded_at"`
	RecordedBy             string    `json:"recorded_by"`
}

func (pr *PatientRecord) Validate() error {
	return cause.Map{
		"patient_id": cause.Required(pr.PatientID).
			Select(map[string]bool{
				"invalid patient ID format": !isValidPatientID(pr.PatientID),
				"patient ID too long":       len(pr.PatientID) > 20,
			}),

		"medical_number": cause.Required(pr.MedicalNumber).
			Select(map[string]bool{
				"invalid medical record number format": !isValidMedicalNumber(pr.MedicalNumber),
				"medical number too short":             len(pr.MedicalNumber) < 8,
				"medical number too long":              len(pr.MedicalNumber) > 15,
			}),

		"personal_info":       cause.Required(pr.PersonalInfo).Err(),
		"medical_history":     cause.Optional(pr.MedicalHistory).Err(),
		"allergies":           cause.Optional(pr.Allergies).Err(),
		"current_medications": cause.Optional(pr.Medications).Err(),
		"emergency_contact":   cause.Required(pr.EmergencyContact).Err(),
		"insurance":           cause.Required(pr.Insurance).Err(),
		"vitals":              cause.Optional(pr.Vitals).Err(),
	}.Err()
}

func (pi *PersonalInfo) Validate() error {
	age := getAge(pi.DateOfBirth)
	return cause.Map{
		"first_name": cause.Required(pi.FirstName).
			Select(map[string]bool{
				"first name is required":            len(pi.FirstName) < 1,
				"first name too long":               len(pi.FirstName) > 50,
				"first name cannot contain numbers": containsNumbers(pi.FirstName),
			}),

		"last_name": cause.Required(pi.LastName).
			Select(map[string]bool{
				"last name is required":            len(pi.LastName) < 1,
				"last name too long":               len(pi.LastName) > 50,
				"last name cannot contain numbers": containsNumbers(pi.LastName),
			}),

		"date_of_birth": cause.Required(pi.DateOfBirth).
			Select(map[string]bool{
				"date of birth cannot be in the future": pi.DateOfBirth.After(time.Now()),
				"invalid date of birth - age too high":  age > 150,
				"invalid date of birth":                 age < 0,
			}),

		"gender": cause.Required(pi.Gender).
			When(!isValidGender(pi.Gender), "invalid gender value").Err(),

		"blood_type": cause.Required(pi.BloodType).
			When(!isValidBloodType(pi.BloodType), "invalid blood type").Err(),

		"height_cm": cause.Required(pi.Height).
			Select(map[string]bool{
				"height must be positive":        pi.Height <= 0,
				"height outside normal range":    pi.Height < 30 || pi.Height > 300,
				"height unusually low for adult": age >= 18 && pi.Height < 120,
			}),

		"weight_kg": cause.Required(pi.Weight).
			Select(map[string]bool{
				"weight must be positive":        pi.Weight <= 0,
				"weight outside normal range":    pi.Weight < 1 || pi.Weight > 500,
				"weight unusually low for adult": age >= 18 && pi.Weight < 30,
			}),

		"phone_number": cause.Required(pi.PhoneNumber).
			When(!isValidPhoneNumber(pi.PhoneNumber), "invalid phone number format").Err(),

		"email": cause.Optional(pi.Email).
			When(!isValidEmail(pi.Email), "invalid email format").Err(),

		"address": cause.Required(pi.Address).Err(),
	}.Err()
}

func (me *MedicalEntry) Validate() error {
	return cause.Map{
		"date": cause.Required(me.Date).
			Select(map[string]bool{
				"medical entry date cannot be in the future": me.Date.After(time.Now()),
				"medical entry date too old":                 me.Date.Before(time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)),
			}),

		"condition": cause.Required(me.Condition).
			Select(map[string]bool{
				"condition description too short": len(me.Condition) < 3,
				"condition description too long":  len(me.Condition) > 100,
			}),

		"diagnosis": cause.Required(me.Diagnosis).
			Select(map[string]bool{
				"diagnosis too short": len(me.Diagnosis) < 5,
				"diagnosis too long":  len(me.Diagnosis) > 500,
			}),

		"treatment": cause.Required(me.Treatment).
			Select(map[string]bool{
				"treatment description too short": len(me.Treatment) < 3,
				"treatment description too long":  len(me.Treatment) > 1000,
			}),

		"doctor_id": cause.Required(me.DoctorID).
			When(!isValidDoctorID(me.DoctorID), "invalid doctor ID format").Err(),

		"severity": cause.Required(me.Severity).
			When(!isValidSeverity(me.Severity), "invalid severity level").Err(),

		"notes": cause.Optional(me.Notes).
			When(len(me.Notes) > 2000, "notes too long").Err(),

		"follow_up_date": cause.Optional(me.FollowUpDate).
			Select(map[string]bool{
				"follow-up date cannot be before entry date": me.FollowUpDate != nil && me.FollowUpDate.Before(me.Date),
				"follow-up date too far in future":           me.FollowUpDate != nil && me.FollowUpDate.After(time.Now().AddDate(2, 0, 0)),
			}),
	}.Err()
}

func (a *Allergy) Validate() error {
	return cause.Map{
		"allergen": cause.Required(a.Allergen).
			Select(map[string]bool{
				"allergen name too short": len(a.Allergen) < 2,
				"allergen name too long":  len(a.Allergen) > 100,
			}),

		"reaction": cause.Required(a.Reaction).
			Select(map[string]bool{
				"reaction description too short": len(a.Reaction) < 3,
				"reaction description too long":  len(a.Reaction) > 200,
			}),

		"severity": cause.Required(a.Severity).
			When(!isValidAllergySeverity(a.Severity), "invalid allergy severity level").Err(),

		"diagnosed_date": cause.Required(a.DiagnosedDate).
			Select(map[string]bool{
				"diagnosed date cannot be in the future": a.DiagnosedDate.After(time.Now()),
				"diagnosed date too old":                 a.DiagnosedDate.Before(time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)),
			}),

		"notes": cause.Optional(a.Notes).
			When(len(a.Notes) > 500, "notes too long").Err(),
	}.Err()
}

func (m *Medication) Validate() error {
	return cause.Map{
		"name": cause.Required(m.Name).
			Select(map[string]bool{
				"medication name too short": len(m.Name) < 2,
				"medication name too long":  len(m.Name) > 100,
			}),

		"dosage": cause.Required(m.Dosage).
			Select(map[string]bool{
				"invalid dosage format":       !isValidDosage(m.Dosage),
				"dosage description too long": len(m.Dosage) > 50,
			}),

		"frequency": cause.Required(m.Frequency).
			Select(map[string]bool{
				"invalid medication frequency": !isValidMedicationFrequency(m.Frequency),
			}),

		"start_date": cause.Required(m.StartDate).
			Select(map[string]bool{
				"start date cannot be more than 1 day in future": m.StartDate.After(time.Now().AddDate(0, 0, 1)),
				"start date too old":                             m.StartDate.Before(time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)),
			}),

		"end_date": cause.Optional(m.EndDate).
			Select(map[string]bool{
				"end date cannot be before start date": m.EndDate != nil && m.EndDate.Before(m.StartDate),
				"end date too far in future":           m.EndDate != nil && m.EndDate.After(time.Now().AddDate(10, 0, 0)),
			}),

		"prescribed_by": cause.Required(m.PrescribedBy).
			Select(map[string]bool{
				"invalid prescribing doctor ID": !isValidDoctorID(m.PrescribedBy),
			}),

		"purpose": cause.Required(m.Purpose).
			Select(map[string]bool{
				"medication purpose too short": len(m.Purpose) < 5,
				"medication purpose too long":  len(m.Purpose) > 200,
			}),

		"side_effects": cause.Optional(m.SideEffects).
			Select(map[string]bool{
				"too many side effects listed": len(m.SideEffects) > 20,
			}),
	}.Err()
}

func (ec *MedicalEmergencyContact) Validate() error {
	return cause.Map{
		"name": cause.Required(ec.Name).
			Select(map[string]bool{
				"emergency contact name too short": len(ec.Name) < 2,
				"emergency contact name too long":  len(ec.Name) > 100,
			}),

		"relationship": cause.Required(ec.Relationship).
			Select(map[string]bool{
				"invalid relationship type": !isValidMedicalRelationship(ec.Relationship),
			}),

		"phone_number": cause.Required(ec.PhoneNumber).
			Select(map[string]bool{
				"invalid emergency contact phone number": !isValidPhoneNumber(ec.PhoneNumber),
			}),

		"email": cause.Optional(ec.Email).
			Select(map[string]bool{
				"invalid emergency contact email": !isValidEmail(ec.Email),
			}),

		"address": cause.Required(ec.Address).Err(),
	}.Err()
}

func (ii *InsuranceInfo) Validate() error {
	return cause.Map{
		"provider": cause.Required(ii.Provider).
			Select(map[string]bool{
				"insurance provider name too short": len(ii.Provider) < 2,
				"insurance provider name too long":  len(ii.Provider) > 100,
			}),

		"policy_number": cause.Required(ii.PolicyNumber).
			Select(map[string]bool{
				"invalid policy number format": !isValidPolicyNumber(ii.PolicyNumber),
				"policy number too short":      len(ii.PolicyNumber) < 5,
				"policy number too long":       len(ii.PolicyNumber) > 30,
			}),

		"group_number": cause.Optional(ii.GroupNumber).
			Select(map[string]bool{
				"group number too long": len(ii.GroupNumber) > 20,
			}),

		"valid_from": cause.Required(ii.ValidFrom).
			Select(map[string]bool{
				"insurance valid from date too far in future": ii.ValidFrom.After(time.Now().AddDate(0, 1, 0)),
			}),

		"valid_to": cause.Required(ii.ValidTo).
			Select(map[string]bool{
				"insurance valid to date cannot be before valid from date": ii.ValidTo.Before(ii.ValidFrom),
				"insurance has been expired for too long":                  ii.ValidTo.Before(time.Now().AddDate(0, -1, 0)),
			}),

		"copay": cause.Required(ii.Copay).
			Select(map[string]bool{
				"copay cannot be negative":    ii.Copay < 0,
				"copay amount unusually high": ii.Copay > 10000,
			}),

		"deductible": cause.Required(ii.Deductible).
			Select(map[string]bool{
				"deductible cannot be negative":    ii.Deductible < 0,
				"deductible amount unusually high": ii.Deductible > 100000,
			}),
	}.Err()
}

func (vs *VitalSigns) Validate() error {
	return cause.Map{
		"bp_systolic": cause.Required(vs.BloodPressureSystolic).
			Select(map[string]bool{
				"systolic blood pressure outside normal range":    vs.BloodPressureSystolic < 50 || vs.BloodPressureSystolic > 300,
				"systolic pressure must be higher than diastolic": vs.BloodPressureSystolic <= vs.BloodPressureDiastolic,
			}),

		"bp_diastolic": cause.Required(vs.BloodPressureDiastolic).
			Select(map[string]bool{
				"diastolic blood pressure outside normal range": vs.BloodPressureDiastolic < 30 || vs.BloodPressureDiastolic > 200,
			}),

		"heart_rate": cause.Required(vs.HeartRate).
			Select(map[string]bool{
				"heart rate outside normal range": vs.HeartRate < 30 || vs.HeartRate > 250,
			}),

		"temperature_celsius": cause.Required(vs.Temperature).
			Select(map[string]bool{
				"body temperature outside viable range":        vs.Temperature < 30.0 || vs.Temperature > 45.0,
				"hypothermia detected - temperature too low":   vs.Temperature < 35.0,
				"hyperthermia detected - temperature too high": vs.Temperature > 42.0,
			}),

		"oxygen_saturation": cause.Required(vs.OxygenSaturation).
			Select(map[string]bool{
				"oxygen saturation outside normal range": vs.OxygenSaturation < 70 || vs.OxygenSaturation > 100,
				"critically low oxygen saturation":       vs.OxygenSaturation < 90,
			}),

		"respiratory_rate": cause.Required(vs.RespiratoryRate).
			Select(map[string]bool{
				"respiratory rate outside normal range": vs.RespiratoryRate < 5 || vs.RespiratoryRate > 60,
			}),

		"recorded_at": cause.Required(vs.RecordedAt).
			Select(map[string]bool{
				"vital signs cannot be recorded in the future": vs.RecordedAt.After(time.Now()),
				"vital signs too old":                          vs.RecordedAt.Before(time.Now().AddDate(0, 0, -30)),
			}),

		"recorded_by": cause.Required(vs.RecordedBy).
			Select(map[string]bool{
				"invalid staff ID for recorder": !isValidStaffID(vs.RecordedBy),
			}),
	}.Err()
}

// Healthcare-specific validation helper functions
func isValidPatientID(id string) bool {
	matched, _ := regexp.MatchString(`^P\d{6,10}$`, id)
	return matched
}

func isValidMedicalNumber(number string) bool {
	matched, _ := regexp.MatchString(`^MR\d{8,12}$`, number)
	return matched
}

func isValidGender(gender string) bool {
	validGenders := []string{"M", "F", "O", "U"} // Male, Female, Other, Unknown
	for _, valid := range validGenders {
		if gender == valid {
			return true
		}
	}
	return false
}

func isValidBloodType(bloodType string) bool {
	validTypes := []string{"A+", "A-", "B+", "B-", "AB+", "AB-", "O+", "O-"}
	for _, valid := range validTypes {
		if bloodType == valid {
			return true
		}
	}
	return false
}

func isValidDoctorID(id string) bool {
	matched, _ := regexp.MatchString(`^DR\d{4,8}$`, id)
	return matched
}

func isValidSeverity(severity string) bool {
	validSeverities := []string{"Low", "Moderate", "High", "Critical"}
	for _, valid := range validSeverities {
		if severity == valid {
			return true
		}
	}
	return false
}

func isValidAllergySeverity(severity string) bool {
	validSeverities := []string{"Mild", "Moderate", "Severe", "Life-threatening"}
	for _, valid := range validSeverities {
		if severity == valid {
			return true
		}
	}
	return false
}

func isValidDosage(dosage string) bool {
	// Simple validation for dosage format (e.g., "500mg", "2.5ml", "1 tablet")
	matched, _ := regexp.MatchString(`^\d+(\.\d+)?\s*(mg|ml|tablet|capsule|drop|unit)s?$`, strings.ToLower(dosage))
	return matched
}

func isValidMedicationFrequency(frequency string) bool {
	validFrequencies := []string{
		"Once daily", "Twice daily", "Three times daily", "Four times daily",
		"Every 4 hours", "Every 6 hours", "Every 8 hours", "Every 12 hours",
		"As needed", "Before meals", "After meals", "At bedtime",
	}
	for _, valid := range validFrequencies {
		if frequency == valid {
			return true
		}
	}
	return false
}

func isValidMedicalRelationship(relationship string) bool {
	validRelationships := []string{
		"Spouse", "Parent", "Child", "Sibling", "Grandparent", "Grandchild",
		"Friend", "Partner", "Guardian", "Other",
	}
	for _, valid := range validRelationships {
		if relationship == valid {
			return true
		}
	}
	return false
}

func isValidPolicyNumber(policyNumber string) bool {
	// Policy numbers are typically alphanumeric
	matched, _ := regexp.MatchString(`^[A-Z0-9]{5,30}$`, policyNumber)
	return matched
}

func isValidStaffID(id string) bool {
	matched, _ := regexp.MatchString(`^(DR|NP|RN|LPN|ST)\d{4,8}$`, id)
	return matched
}

// Example test function
func ExamplePatientRecord_Validate() {
	// Valid patient record
	validPatient := &PatientRecord{
		PatientID:     "P12345678",
		MedicalNumber: "MR123456789",
		PersonalInfo: PersonalInfo{
			FirstName:   "John",
			LastName:    "Doe",
			DateOfBirth: time.Date(1980, 5, 15, 0, 0, 0, 0, time.UTC),
			Gender:      "M",
			BloodType:   "O+",
			Height:      175.5,
			Weight:      70.2,
			PhoneNumber: "+1-555-123-4567",
			Email:       "john.doe@email.com",
			Address: Address{
				Street:  "123 Main St",
				City:    "Springfield",
				State:   "IL",
				Country: "US",
				ZipCode: "62701",
			},
		},
		EmergencyContact: MedicalEmergencyContact{
			Name:         "Jane Doe",
			Relationship: "Spouse",
			PhoneNumber:  "+1-555-987-6543",
			Email:        "jane.doe@email.com",
			Address: Address{
				Street:  "123 Main St",
				City:    "Springfield",
				State:   "IL",
				Country: "US",
				ZipCode: "62701",
			},
		},
		Insurance: InsuranceInfo{
			Provider:     "HealthCare Plus",
			PolicyNumber: "HP123456789",
			ValidFrom:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			ValidTo:      time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
			Copay:        25.00,
			Deductible:   1000.00,
		},
	}

	err := validPatient.Validate()
	fmt.Printf("Valid patient error: %v\n", err)

	// Invalid patient record
	invalidPatient := &PatientRecord{
		PatientID:     "INVALID",
		MedicalNumber: "123",
		PersonalInfo: PersonalInfo{
			FirstName:   "",
			LastName:    "X",
			DateOfBirth: time.Now().AddDate(0, 0, 1), // Future date
			Gender:      "Invalid",
			BloodType:   "Z+",
			Height:      -10,
			Weight:      0,
			PhoneNumber: "invalid-phone",
		},
		EmergencyContact: MedicalEmergencyContact{
			Name:         "A",
			Relationship: "Invalid",
			PhoneNumber:  "bad-phone",
		},
		Insurance: InsuranceInfo{
			Provider:     "A",
			PolicyNumber: "123",
			ValidFrom:    time.Now().AddDate(0, 2, 0),  // Too far in future
			ValidTo:      time.Now().AddDate(0, -2, 0), // In the past
			Copay:        -100,
			Deductible:   -500,
		},
	}

	err = invalidPatient.Validate()
	if err != nil {
		fmt.Printf("Invalid patient has errors: %v\n", err != nil)
	}

	// Output:
	// Valid patient error: <nil>
	// Invalid patient has errors: true
}
