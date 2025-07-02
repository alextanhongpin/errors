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
			When(!isValidPatientID(pr.PatientID), "invalid patient ID format").
			When(len(pr.PatientID) > 20, "patient ID too long"),

		"medical_number": cause.Required(pr.MedicalNumber).
			When(!isValidMedicalNumber(pr.MedicalNumber), "invalid medical record number format").
			When(len(pr.MedicalNumber) < 8, "medical number too short").
			When(len(pr.MedicalNumber) > 15, "medical number too long"),

		"personal_info":       cause.Required(pr.PersonalInfo),
		"medical_history":     cause.Optional(pr.MedicalHistory),
		"allergies":           cause.Optional(pr.Allergies),
		"current_medications": cause.Optional(pr.Medications),
		"emergency_contact":   cause.Required(pr.EmergencyContact),
		"insurance":           cause.Required(pr.Insurance),
		"vitals":              cause.Optional(pr.Vitals),
	}.Err()
}

func (pi *PersonalInfo) Validate() error {
	age := getAge(pi.DateOfBirth)
	return cause.Map{
		"first_name": cause.Required(pi.FirstName).
			When(len(pi.FirstName) < 1, "first name is required").
			When(len(pi.FirstName) > 50, "first name too long").
			When(containsNumbers(pi.FirstName), "first name cannot contain numbers"),

		"last_name": cause.Required(pi.LastName).
			When(len(pi.LastName) < 1, "last name is required").
			When(len(pi.LastName) > 50, "last name too long").
			When(containsNumbers(pi.LastName), "last name cannot contain numbers"),

		"date_of_birth": cause.Required(pi.DateOfBirth).
			When(pi.DateOfBirth.After(time.Now()), "date of birth cannot be in the future").
			When(age > 150, "invalid date of birth - age too high").
			When(age < 0, "invalid date of birth"),

		"gender": cause.Required(pi.Gender).
			When(!isValidGender(pi.Gender), "invalid gender value"),

		"blood_type": cause.Required(pi.BloodType).
			When(!isValidBloodType(pi.BloodType), "invalid blood type"),

		"height_cm": cause.Required(pi.Height).
			When(pi.Height <= 0, "height must be positive").
			When(pi.Height < 30 || pi.Height > 300, "height outside normal range").
			When(age >= 18 && pi.Height < 120, "height unusually low for adult"),

		"weight_kg": cause.Required(pi.Weight).
			When(pi.Weight <= 0, "weight must be positive").
			When(pi.Weight < 1 || pi.Weight > 500, "weight outside normal range").
			When(age >= 18 && pi.Weight < 30, "weight unusually low for adult"),

		"phone_number": cause.Required(pi.PhoneNumber).
			When(!isValidPhoneNumber(pi.PhoneNumber), "invalid phone number format"),

		"email": cause.Optional(pi.Email).
			When(!isValidEmail(pi.Email), "invalid email format"),

		"address": cause.Required(pi.Address),
	}.Err()
}

func (me *MedicalEntry) Validate() error {
	return cause.Map{
		"date": cause.Required(me.Date).
			When(me.Date.After(time.Now()), "medical entry date cannot be in the future").
			When(me.Date.Before(time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)), "medical entry date too old"),

		"condition": cause.Required(me.Condition).
			When(len(me.Condition) < 3, "condition description too short").
			When(len(me.Condition) > 100, "condition description too long"),

		"diagnosis": cause.Required(me.Diagnosis).
			When(len(me.Diagnosis) < 5, "diagnosis too short").
			When(len(me.Diagnosis) > 500, "diagnosis too long"),

		"treatment": cause.Required(me.Treatment).
			When(len(me.Treatment) < 3, "treatment description too short").
			When(len(me.Treatment) > 1000, "treatment description too long"),

		"doctor_id": cause.Required(me.DoctorID).
			When(!isValidDoctorID(me.DoctorID), "invalid doctor ID format"),

		"severity": cause.Required(me.Severity).
			When(!isValidSeverity(me.Severity), "invalid severity level"),

		"notes": cause.Optional(me.Notes).
			When(len(me.Notes) > 2000, "notes too long"),

		"follow_up_date": cause.Optional(me.FollowUpDate).
			When(me.FollowUpDate != nil && me.FollowUpDate.Before(me.Date), "follow-up date cannot be before entry date").
			When(me.FollowUpDate != nil && me.FollowUpDate.After(time.Now().AddDate(2, 0, 0)), "follow-up date too far in future"),
	}.Err()
}

func (a *Allergy) Validate() error {
	return cause.Map{
		"allergen": cause.Required(a.Allergen).
			When(len(a.Allergen) < 2, "allergen name too short").
			When(len(a.Allergen) > 100, "allergen name too long"),

		"reaction": cause.Required(a.Reaction).
			When(len(a.Reaction) < 3, "reaction description too short").
			When(len(a.Reaction) > 200, "reaction description too long"),

		"severity": cause.Required(a.Severity).
			When(!isValidAllergySeverity(a.Severity), "invalid allergy severity level"),

		"diagnosed_date": cause.Required(a.DiagnosedDate).
			When(a.DiagnosedDate.After(time.Now()), "diagnosed date cannot be in the future").
			When(a.DiagnosedDate.Before(time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)), "diagnosed date too old"),

		"notes": cause.Optional(a.Notes).
			When(len(a.Notes) > 500, "notes too long"),
	}.Err()
}

func (m *Medication) Validate() error {
	return cause.Map{
		"name": cause.Required(m.Name).
			When(len(m.Name) < 2, "medication name too short").
			When(len(m.Name) > 100, "medication name too long"),

		"dosage": cause.Required(m.Dosage).
			When(!isValidDosage(m.Dosage), "invalid dosage format").
			When(len(m.Dosage) > 50, "dosage description too long"),

		"frequency": cause.Required(m.Frequency).
			When(!isValidMedicationFrequency(m.Frequency), "invalid medication frequency"),

		"start_date": cause.Required(m.StartDate).
			When(m.StartDate.After(time.Now().AddDate(0, 0, 1)), "start date cannot be more than 1 day in future").
			When(m.StartDate.Before(time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)), "start date too old"),

		"end_date": cause.Optional(m.EndDate).
			When(m.EndDate != nil && m.EndDate.Before(m.StartDate), "end date cannot be before start date").
			When(m.EndDate != nil && m.EndDate.After(time.Now().AddDate(10, 0, 0)), "end date too far in future"),

		"prescribed_by": cause.Required(m.PrescribedBy).
			When(!isValidDoctorID(m.PrescribedBy), "invalid prescribing doctor ID"),

		"purpose": cause.Required(m.Purpose).
			When(len(m.Purpose) < 5, "medication purpose too short").
			When(len(m.Purpose) > 200, "medication purpose too long"),

		"side_effects": cause.Optional(m.SideEffects).
			When(len(m.SideEffects) > 20, "too many side effects listed"),
	}.Err()
}

func (ec *MedicalEmergencyContact) Validate() error {
	return cause.Map{
		"name": cause.Required(ec.Name).
			When(len(ec.Name) < 2, "emergency contact name too short").
			When(len(ec.Name) > 100, "emergency contact name too long"),

		"relationship": cause.Required(ec.Relationship).
			When(!isValidMedicalRelationship(ec.Relationship), "invalid relationship type"),

		"phone_number": cause.Required(ec.PhoneNumber).
			When(!isValidPhoneNumber(ec.PhoneNumber), "invalid emergency contact phone number"),

		"email": cause.Optional(ec.Email).
			When(!isValidEmail(ec.Email), "invalid emergency contact email"),

		"address": cause.Required(ec.Address),
	}.Err()
}

func (ii *InsuranceInfo) Validate() error {
	return cause.Map{
		"provider": cause.Required(ii.Provider).
			When(len(ii.Provider) < 2, "insurance provider name too short").
			When(len(ii.Provider) > 100, "insurance provider name too long"),

		"policy_number": cause.Required(ii.PolicyNumber).
			When(!isValidPolicyNumber(ii.PolicyNumber), "invalid policy number format").
			When(len(ii.PolicyNumber) < 5, "policy number too short").
			When(len(ii.PolicyNumber) > 30, "policy number too long"),

		"group_number": cause.Optional(ii.GroupNumber).
			When(len(ii.GroupNumber) > 20, "group number too long"),

		"valid_from": cause.Required(ii.ValidFrom).
			When(ii.ValidFrom.After(time.Now().AddDate(0, 1, 0)), "insurance valid from date too far in future"),

		"valid_to": cause.Required(ii.ValidTo).
			When(ii.ValidTo.Before(ii.ValidFrom), "insurance valid to date cannot be before valid from date").
			When(ii.ValidTo.Before(time.Now().AddDate(0, -1, 0)), "insurance has been expired for too long"),

		"copay": cause.Required(ii.Copay).
			When(ii.Copay < 0, "copay cannot be negative").
			When(ii.Copay > 10000, "copay amount unusually high"),

		"deductible": cause.Required(ii.Deductible).
			When(ii.Deductible < 0, "deductible cannot be negative").
			When(ii.Deductible > 100000, "deductible amount unusually high"),
	}.Err()
}

func (vs *VitalSigns) Validate() error {
	return cause.Map{
		"bp_systolic": cause.Required(vs.BloodPressureSystolic).
			When(vs.BloodPressureSystolic < 50 || vs.BloodPressureSystolic > 300, "systolic blood pressure outside normal range").
			When(vs.BloodPressureSystolic <= vs.BloodPressureDiastolic, "systolic pressure must be higher than diastolic"),

		"bp_diastolic": cause.Required(vs.BloodPressureDiastolic).
			When(vs.BloodPressureDiastolic < 30 || vs.BloodPressureDiastolic > 200, "diastolic blood pressure outside normal range"),

		"heart_rate": cause.Required(vs.HeartRate).
			When(vs.HeartRate < 30 || vs.HeartRate > 250, "heart rate outside normal range"),

		"temperature_celsius": cause.Required(vs.Temperature).
			When(vs.Temperature < 30.0 || vs.Temperature > 45.0, "body temperature outside viable range").
			When(vs.Temperature < 35.0, "hypothermia detected - temperature too low").
			When(vs.Temperature > 42.0, "hyperthermia detected - temperature too high"),

		"oxygen_saturation": cause.Required(vs.OxygenSaturation).
			When(vs.OxygenSaturation < 70 || vs.OxygenSaturation > 100, "oxygen saturation outside normal range").
			When(vs.OxygenSaturation < 90, "critically low oxygen saturation"),

		"respiratory_rate": cause.Required(vs.RespiratoryRate).
			When(vs.RespiratoryRate < 5 || vs.RespiratoryRate > 60, "respiratory rate outside normal range"),

		"recorded_at": cause.Required(vs.RecordedAt).
			When(vs.RecordedAt.After(time.Now()), "vital signs cannot be recorded in the future").
			When(vs.RecordedAt.Before(time.Now().AddDate(0, 0, -30)), "vital signs too old"),

		"recorded_by": cause.Required(vs.RecordedBy).
			When(!isValidStaffID(vs.RecordedBy), "invalid staff ID for recorder"),
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
