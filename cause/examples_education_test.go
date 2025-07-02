package cause_test

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/alextanhongpin/errors/cause"
)

// Real-world example: Educational Institution Management System
type StudentEnrollment struct {
	StudentID      string             `json:"student_id"`
	StudentNumber  string             `json:"student_number"`
	PersonalInfo   StudentProfile     `json:"personal_info"`
	AcademicInfo   AcademicProfile    `json:"academic_info"`
	Courses        []CourseEnrollment `json:"courses"`
	Guardian       GuardianInfo       `json:"guardian,omitempty"`
	FinancialAid   *FinancialAid      `json:"financial_aid,omitempty"`
	Transcripts    []Transcript       `json:"transcripts"`
	Status         string             `json:"status"`
	EnrollmentDate time.Time          `json:"enrollment_date"`
}

type StudentProfile struct {
	FirstName        string       `json:"first_name"`
	LastName         string       `json:"last_name"`
	MiddleName       string       `json:"middle_name,omitempty"`
	DateOfBirth      time.Time    `json:"date_of_birth"`
	Gender           string       `json:"gender"`
	Nationality      string       `json:"nationality"`
	Email            string       `json:"email"`
	PhoneNumber      string       `json:"phone_number"`
	Address          Address      `json:"address"`
	EmergencyContact GuardianInfo `json:"emergency_contact"`
}

type AcademicProfile struct {
	Program            string    `json:"program"`
	Major              string    `json:"major"`
	Minor              string    `json:"minor,omitempty"`
	Year               int       `json:"year"`
	Semester           string    `json:"semester"`
	Credits            int       `json:"credits"`
	GPA                float64   `json:"gpa"`
	ExpectedGraduation time.Time `json:"expected_graduation"`
	AdvisorID          string    `json:"advisor_id"`
	DegreeType         string    `json:"degree_type"`
}

type CourseEnrollment struct {
	CourseID      string    `json:"course_id"`
	CourseName    string    `json:"course_name"`
	Credits       int       `json:"credits"`
	InstructorID  string    `json:"instructor_id"`
	Schedule      Schedule  `json:"schedule"`
	EnrollDate    time.Time `json:"enroll_date"`
	Grade         string    `json:"grade,omitempty"`
	Status        string    `json:"status"`
	Prerequisites []string  `json:"prerequisites,omitempty"`
}

type Schedule struct {
	DaysOfWeek []string `json:"days_of_week"`
	StartTime  string   `json:"start_time"`
	EndTime    string   `json:"end_time"`
	Room       string   `json:"room"`
	Building   string   `json:"building"`
}

type GuardianInfo struct {
	Name         string  `json:"name"`
	Relationship string  `json:"relationship"`
	Phone        string  `json:"phone"`
	Email        string  `json:"email,omitempty"`
	Address      Address `json:"address"`
	IsEmergency  bool    `json:"is_emergency"`
}

type FinancialAid struct {
	Type            string    `json:"type"`
	Amount          float64   `json:"amount"`
	Semester        string    `json:"semester"`
	AcademicYear    string    `json:"academic_year"`
	Requirements    []string  `json:"requirements"`
	Status          string    `json:"status"`
	ApplicationDate time.Time `json:"application_date"`
	ExpiryDate      time.Time `json:"expiry_date"`
}

type Transcript struct {
	CourseID     string  `json:"course_id"`
	CourseName   string  `json:"course_name"`
	Credits      int     `json:"credits"`
	Grade        string  `json:"grade"`
	GradePoints  float64 `json:"grade_points"`
	Semester     string  `json:"semester"`
	AcademicYear string  `json:"academic_year"`
	InstructorID string  `json:"instructor_id"`
}

func (se *StudentEnrollment) Validate() error {
	return cause.Map{
		"student_id": cause.Required(se.StudentID).
			When(!isValidStudentID(se.StudentID), "invalid student ID format").
			When(len(se.StudentID) > 15, "student ID too long"),

		"student_number": cause.Required(se.StudentNumber).
			When(!isValidStudentNumber(se.StudentNumber), "invalid student number format").
			When(len(se.StudentNumber) < 6, "student number too short").
			When(len(se.StudentNumber) > 12, "student number too long"),

		"personal_info": cause.Required(se.PersonalInfo),
		"academic_info": cause.Required(se.AcademicInfo),
		"courses":       cause.Optional(se.Courses),
		"guardian":      cause.Optional(se.Guardian),
		"financial_aid": cause.Optional(se.FinancialAid),
		"transcripts":   cause.Optional(se.Transcripts),

		"status": cause.Required(se.Status).
			When(!isValidEnrollmentStatus(se.Status), "invalid enrollment status"),

		"enrollment_date": cause.Required(se.EnrollmentDate).
			When(se.EnrollmentDate.After(time.Now()), "enrollment date cannot be in the future").
			When(se.EnrollmentDate.Before(time.Date(1950, 1, 1, 0, 0, 0, 0, time.UTC)), "enrollment date too old"),
	}.Err()
}

func (sp *StudentProfile) Validate() error {
	age := getAge(sp.DateOfBirth)
	return cause.Map{
		"first_name": cause.Required(sp.FirstName).
			When(len(sp.FirstName) < 1, "first name is required").
			When(len(sp.FirstName) > 50, "first name too long").
			When(containsNumbers(sp.FirstName), "first name cannot contain numbers"),

		"last_name": cause.Required(sp.LastName).
			When(len(sp.LastName) < 1, "last name is required").
			When(len(sp.LastName) > 50, "last name too long").
			When(containsNumbers(sp.LastName), "last name cannot contain numbers"),

		"middle_name": cause.Optional(sp.MiddleName).
			When(len(sp.MiddleName) > 50, "middle name too long").
			When(containsNumbers(sp.MiddleName), "middle name cannot contain numbers"),

		"date_of_birth": cause.Required(sp.DateOfBirth).
			When(sp.DateOfBirth.After(time.Now()), "date of birth cannot be in the future").
			When(age < 5, "student too young").
			When(age > 120, "invalid date of birth"),

		"gender": cause.Required(sp.Gender).
			When(!isValidGender(sp.Gender), "invalid gender value"),

		"nationality": cause.Required(sp.Nationality).
			When(len(sp.Nationality) < 2, "nationality too short").
			When(len(sp.Nationality) > 50, "nationality too long"),

		"email": cause.Required(sp.Email).
			When(!isValidEmail(sp.Email), "invalid email format").
			When(!isEducationalEmail(sp.Email), "non-educational email domain"),

		"phone_number": cause.Required(sp.PhoneNumber).
			When(!isValidPhoneNumber(sp.PhoneNumber), "invalid phone number format"),

		"address":           cause.Required(sp.Address),
		"emergency_contact": cause.Required(sp.EmergencyContact),
	}.Err()
}

func (ap *AcademicProfile) Validate() error {
	return cause.Map{
		"program": cause.Required(ap.Program).
			When(!isValidProgram(ap.Program), "invalid program").
			When(len(ap.Program) > 100, "program name too long"),

		"major": cause.Required(ap.Major).
			When(!isValidMajor(ap.Major), "invalid major").
			When(len(ap.Major) > 100, "major name too long"),

		"minor": cause.Optional(ap.Minor).
			When(!isValidMajor(ap.Minor), "invalid minor").
			When(len(ap.Minor) > 100, "minor name too long").
			When(ap.Minor == ap.Major, "minor cannot be the same as major"),

		"year": cause.Required(ap.Year).
			When(ap.Year < 1 || ap.Year > 8, "academic year must be between 1 and 8"),

		"semester": cause.Required(ap.Semester).
			When(!isValidSemester(ap.Semester), "invalid semester"),

		"credits": cause.Required(ap.Credits).
			When(ap.Credits < 0, "credits cannot be negative").
			When(ap.Credits > 200, "credits exceeds maximum allowed"),

		"gpa": cause.Required(ap.GPA).
			When(ap.GPA < 0.0 || ap.GPA > 4.0, "GPA must be between 0.0 and 4.0"),

		"expected_graduation": cause.Required(ap.ExpectedGraduation).
			When(ap.ExpectedGraduation.Before(time.Now()), "expected graduation cannot be in the past").
			When(ap.ExpectedGraduation.After(time.Now().AddDate(10, 0, 0)), "expected graduation too far in future"),

		"advisor_id": cause.Required(ap.AdvisorID).
			When(!isValidFacultyID(ap.AdvisorID), "invalid advisor ID"),

		"degree_type": cause.Required(ap.DegreeType).
			When(!isValidDegreeType(ap.DegreeType), "invalid degree type"),
	}.Err()
}

func (ce *CourseEnrollment) Validate() error {
	return cause.Map{
		"course_id": cause.Required(ce.CourseID).
			When(!isValidCourseID(ce.CourseID), "invalid course ID format").
			When(len(ce.CourseID) > 15, "course ID too long"),

		"course_name": cause.Required(ce.CourseName).
			When(len(ce.CourseName) < 3, "course name too short").
			When(len(ce.CourseName) > 100, "course name too long"),

		"credits": cause.Required(ce.Credits).
			When(ce.Credits < 0, "credits cannot be negative").
			When(ce.Credits > 10, "credits exceeds maximum for a single course"),

		"instructor_id": cause.Required(ce.InstructorID).
			When(!isValidFacultyID(ce.InstructorID), "invalid instructor ID"),

		"schedule": cause.Required(ce.Schedule),

		"enroll_date": cause.Required(ce.EnrollDate).
			When(ce.EnrollDate.After(time.Now()), "enrollment date cannot be in the future").
			When(ce.EnrollDate.Before(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)), "enrollment date too old"),

		"grade": cause.Optional(ce.Grade).
			When(!isValidGrade(ce.Grade), "invalid grade"),

		"status": cause.Required(ce.Status).
			When(!isValidCourseStatus(ce.Status), "invalid course status"),

		"prerequisites": cause.Optional(ce.Prerequisites).
			When(len(ce.Prerequisites) > 10, "too many prerequisites listed"),
	}.Err()
}

func (s *Schedule) Validate() error {
	return cause.Map{
		"days_of_week": cause.Required(s.DaysOfWeek).
			When(len(s.DaysOfWeek) == 0, "at least one day of week required").
			When(len(s.DaysOfWeek) > 7, "cannot have more than 7 days").
			When(!areValidDays(s.DaysOfWeek), "invalid days of week"),

		"start_time": cause.Required(s.StartTime).
			When(!isValidTime(s.StartTime), "invalid start time format").
			When(!isValidClassTime(s.StartTime), "start time outside class hours"),

		"end_time": cause.Required(s.EndTime).
			When(!isValidTime(s.EndTime), "invalid end time format").
			When(!isValidClassTime(s.EndTime), "end time outside class hours").
			When(!isEndTimeAfterStart(s.StartTime, s.EndTime), "end time must be after start time"),

		"room": cause.Required(s.Room).
			When(!isValidRoomNumber(s.Room), "invalid room number format").
			When(len(s.Room) > 20, "room number too long"),

		"building": cause.Required(s.Building).
			When(len(s.Building) < 1, "building name required").
			When(len(s.Building) > 50, "building name too long"),
	}.Err()
}

func (gi *GuardianInfo) Validate() error {
	return cause.Map{
		"name": cause.Required(gi.Name).
			When(len(gi.Name) < 2, "guardian name too short").
			When(len(gi.Name) > 100, "guardian name too long"),

		"relationship": cause.Required(gi.Relationship).
			When(!isValidGuardianRelationship(gi.Relationship), "invalid guardian relationship"),

		"phone": cause.Required(gi.Phone).
			When(!isValidPhoneNumber(gi.Phone), "invalid guardian phone number"),

		"email": cause.Optional(gi.Email).
			When(!isValidEmail(gi.Email), "invalid guardian email"),

		"address": cause.Required(gi.Address),
	}.Err()
}

func (fa *FinancialAid) Validate() error {
	return cause.Map{
		"type": cause.Required(fa.Type).
			When(!isValidFinancialAidType(fa.Type), "invalid financial aid type"),

		"amount": cause.Required(fa.Amount).
			When(fa.Amount <= 0, "financial aid amount must be positive").
			When(fa.Amount > 100000, "financial aid amount exceeds maximum"),

		"semester": cause.Required(fa.Semester).
			When(!isValidSemester(fa.Semester), "invalid semester"),

		"academic_year": cause.Required(fa.AcademicYear).
			When(!isValidAcademicYear(fa.AcademicYear), "invalid academic year format"),

		"requirements": cause.Optional(fa.Requirements).
			When(len(fa.Requirements) > 20, "too many requirements listed"),

		"status": cause.Required(fa.Status).
			When(!isValidFinancialAidStatus(fa.Status), "invalid financial aid status"),

		"application_date": cause.Required(fa.ApplicationDate).
			When(fa.ApplicationDate.After(time.Now()), "application date cannot be in the future").
			When(fa.ApplicationDate.Before(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)), "application date too old"),

		"expiry_date": cause.Required(fa.ExpiryDate).
			When(fa.ExpiryDate.Before(fa.ApplicationDate), "expiry date cannot be before application date").
			When(fa.ExpiryDate.After(time.Now().AddDate(5, 0, 0)), "expiry date too far in future"),
	}.Err()
}

func (t *Transcript) Validate() error {
	return cause.Map{
		"course_id": cause.Required(t.CourseID).
			When(!isValidCourseID(t.CourseID), "invalid course ID"),

		"course_name": cause.Required(t.CourseName).
			When(len(t.CourseName) < 3, "course name too short").
			When(len(t.CourseName) > 100, "course name too long"),

		"credits": cause.Required(t.Credits).
			When(t.Credits < 0, "credits cannot be negative").
			When(t.Credits > 10, "credits exceeds maximum"),

		"grade": cause.Required(t.Grade).
			When(!isValidGrade(t.Grade), "invalid grade"),

		"grade_points": cause.Required(t.GradePoints).
			When(t.GradePoints < 0.0 || t.GradePoints > 4.0, "grade points must be between 0.0 and 4.0"),

		"semester": cause.Required(t.Semester).
			When(!isValidSemester(t.Semester), "invalid semester"),

		"academic_year": cause.Required(t.AcademicYear).
			When(!isValidAcademicYear(t.AcademicYear), "invalid academic year"),

		"instructor_id": cause.Required(t.InstructorID).
			When(!isValidFacultyID(t.InstructorID), "invalid instructor ID"),
	}.Err()
}

// Educational system validation helper functions
func isValidStudentID(id string) bool {
	matched, _ := regexp.MatchString(`^ST\d{6,10}$`, id)
	return matched
}

func isValidStudentNumber(number string) bool {
	matched, _ := regexp.MatchString(`^\d{6,12}$`, number)
	return matched
}

func isValidEnrollmentStatus(status string) bool {
	validStatuses := []string{"Active", "Inactive", "Graduated", "Transferred", "Suspended", "Withdrawn"}
	for _, valid := range validStatuses {
		if status == valid {
			return true
		}
	}
	return false
}

func isEducationalEmail(email string) bool {
	// Check for .edu domains or known educational domains
	eduDomains := []string{".edu", ".ac.", ".university"}
	emailLower := strings.ToLower(email)
	for _, domain := range eduDomains {
		if strings.Contains(emailLower, domain) {
			return true
		}
	}
	return false
}

func isValidProgram(program string) bool {
	validPrograms := []string{
		"Computer Science", "Engineering", "Business Administration", "Medicine",
		"Law", "Arts", "Education", "Psychology", "Biology", "Chemistry",
		"Physics", "Mathematics", "History", "Literature", "Economics",
	}
	for _, valid := range validPrograms {
		if program == valid {
			return true
		}
	}
	return false
}

func isValidMajor(major string) bool {
	// Simplified check - in real system, this would check against a database
	return len(major) >= 3 && len(major) <= 100
}

func isValidSemester(semester string) bool {
	validSemesters := []string{"Fall", "Spring", "Summer", "Winter"}
	for _, valid := range validSemesters {
		if semester == valid {
			return true
		}
	}
	return false
}

func isValidDegreeType(degreeType string) bool {
	validTypes := []string{"Bachelor", "Master", "PhD", "Associate", "Certificate", "Diploma"}
	for _, valid := range validTypes {
		if degreeType == valid {
			return true
		}
	}
	return false
}

func isValidFacultyID(id string) bool {
	matched, _ := regexp.MatchString(`^(FAC|PROF|INST)\d{4,8}$`, id)
	return matched
}

func isValidCourseID(id string) bool {
	matched, _ := regexp.MatchString(`^[A-Z]{2,4}\d{3,4}$`, id)
	return matched
}

func isValidGrade(grade string) bool {
	validGrades := []string{"A+", "A", "A-", "B+", "B", "B-", "C+", "C", "C-", "D+", "D", "F", "I", "W", "P", "NP"}
	for _, valid := range validGrades {
		if grade == valid {
			return true
		}
	}
	return false
}

func isValidCourseStatus(status string) bool {
	validStatuses := []string{"Enrolled", "Completed", "Dropped", "Withdrawn", "In Progress"}
	for _, valid := range validStatuses {
		if status == valid {
			return true
		}
	}
	return false
}

func areValidDays(days []string) bool {
	validDays := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	for _, day := range days {
		found := false
		for _, valid := range validDays {
			if day == valid {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func isValidTime(timeStr string) bool {
	matched, _ := regexp.MatchString(`^([01]?[0-9]|2[0-3]):[0-5][0-9]$`, timeStr)
	return matched
}

func isValidClassTime(timeStr string) bool {
	// Class hours typically 7:00 AM to 10:00 PM
	matched, _ := regexp.MatchString(`^(0[7-9]|1[0-9]|2[0-2]):[0-5][0-9]$`, timeStr)
	return matched
}

func isEndTimeAfterStart(startTime, endTime string) bool {
	// Simple comparison - in real system would parse times properly
	return endTime > startTime
}

func isValidRoomNumber(room string) bool {
	matched, _ := regexp.MatchString(`^[A-Z]?\d{2,4}[A-Z]?$`, room)
	return matched
}

func isValidGuardianRelationship(relationship string) bool {
	validRelationships := []string{
		"Parent", "Father", "Mother", "Guardian", "Grandparent",
		"Uncle", "Aunt", "Sibling", "Spouse", "Other",
	}
	for _, valid := range validRelationships {
		if relationship == valid {
			return true
		}
	}
	return false
}

func isValidFinancialAidType(aidType string) bool {
	validTypes := []string{
		"Scholarship", "Grant", "Loan", "Work-Study", "Fellowship",
		"Assistantship", "Need-Based Aid", "Merit-Based Aid",
	}
	for _, valid := range validTypes {
		if aidType == valid {
			return true
		}
	}
	return false
}

func isValidFinancialAidStatus(status string) bool {
	validStatuses := []string{"Applied", "Approved", "Denied", "Pending", "Disbursed", "Cancelled"}
	for _, valid := range validStatuses {
		if status == valid {
			return true
		}
	}
	return false
}

func isValidAcademicYear(year string) bool {
	matched, _ := regexp.MatchString(`^\d{4}-\d{4}$`, year)
	return matched
}

// Example test function
func ExampleStudentEnrollment_Validate() {
	// Valid student enrollment
	validStudent := &StudentEnrollment{
		StudentID:     "ST12345678",
		StudentNumber: "202312345",
		PersonalInfo: StudentProfile{
			FirstName:   "Alice",
			LastName:    "Johnson",
			DateOfBirth: time.Date(2000, 3, 15, 0, 0, 0, 0, time.UTC),
			Gender:      "F",
			Nationality: "American",
			Email:       "alice.johnson@university.edu",
			PhoneNumber: "+1-555-123-4567",
			Address: Address{
				Street:  "123 College Ave",
				City:    "University City",
				State:   "CA",
				Country: "US",
				ZipCode: "90210",
			},
			EmergencyContact: GuardianInfo{
				Name:         "John Johnson",
				Relationship: "Father",
				Phone:        "+1-555-987-6543",
				Email:        "john.johnson@email.com",
				Address: Address{
					Street:  "456 Parent St",
					City:    "Hometown",
					State:   "CA",
					Country: "US",
					ZipCode: "90211",
				},
			},
		},
		AcademicInfo: AcademicProfile{
			Program:            "Computer Science",
			Major:              "Computer Science",
			Year:               2,
			Semester:           "Fall",
			Credits:            45,
			GPA:                3.7,
			ExpectedGraduation: time.Date(2026, 5, 15, 0, 0, 0, 0, time.UTC),
			AdvisorID:          "FAC12345",
			DegreeType:         "Bachelor",
		},
		Status:         "Active",
		EnrollmentDate: time.Date(2023, 8, 25, 0, 0, 0, 0, time.UTC),
	}

	err := validStudent.Validate()
	fmt.Printf("Valid student error: %v\n", err)

	// Invalid student enrollment
	invalidStudent := &StudentEnrollment{
		StudentID:     "INVALID",
		StudentNumber: "123",
		PersonalInfo: StudentProfile{
			FirstName:   "",
			LastName:    "X",
			DateOfBirth: time.Now().AddDate(0, 0, 1), // Future date
			Gender:      "Invalid",
			Email:       "invalid-email",
			PhoneNumber: "bad-phone",
		},
		AcademicInfo: AcademicProfile{
			Program:    "Invalid Program",
			Major:      "",
			Year:       15, // Invalid year
			Semester:   "Invalid",
			Credits:    -10,
			GPA:        5.0, // Above 4.0
			AdvisorID:  "INVALID",
			DegreeType: "Invalid",
		},
		Status:         "Invalid Status",
		EnrollmentDate: time.Now().AddDate(0, 0, 1), // Future date
	}

	err = invalidStudent.Validate()
	if err != nil {
		fmt.Printf("Invalid student has errors: %v\n", err != nil)
	}

	// Output:
	// Valid student error: <nil>
	// Invalid student has errors: true
}
