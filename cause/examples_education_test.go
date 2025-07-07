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
			When(len(se.StudentID) > 15, "student ID too long").Err(),

		"student_number": cause.Required(se.StudentNumber).
			When(!isValidStudentNumber(se.StudentNumber), "invalid student number format").
			When(len(se.StudentNumber) < 6, "student number too short").
			When(len(se.StudentNumber) > 12, "student number too long").Err(),

		"personal_info": cause.Required(se.PersonalInfo).Err(),
		"academic_info": cause.Required(se.AcademicInfo).Err(),
		"courses":       cause.Optional(se.Courses).Err(),
		"guardian":      cause.Optional(se.Guardian).Err(),
		"financial_aid": cause.Optional(se.FinancialAid).Err(),
		"transcripts":   cause.Optional(se.Transcripts).Err(),

		"status": cause.Required(se.Status).
			When(!isValidEnrollmentStatus(se.Status), "invalid enrollment status").Err(),

		"enrollment_date": cause.Required(se.EnrollmentDate).
			When(se.EnrollmentDate.After(time.Now()), "enrollment date cannot be in the future").
			When(se.EnrollmentDate.Before(time.Date(1950, 1, 1, 0, 0, 0, 0, time.UTC)), "enrollment date too old").Err(),
	}.Err()
}

func (sp *StudentProfile) Validate() error {
	age := getAge(sp.DateOfBirth)
	return cause.Map{
		"first_name": cause.Required(sp.FirstName).
			Select(map[string]bool{
				"first name is required":            len(sp.FirstName) < 1,
				"first name too long":               len(sp.FirstName) > 50,
				"first name cannot contain numbers": containsNumbers(sp.FirstName),
			}),

		"last_name": cause.Required(sp.LastName).
			Select(map[string]bool{
				"last name is required":            len(sp.LastName) < 1,
				"last name too long":               len(sp.LastName) > 50,
				"last name cannot contain numbers": containsNumbers(sp.LastName),
			}),

		"middle_name": cause.Optional(sp.MiddleName).
			Select(map[string]bool{
				"middle name too long":               len(sp.MiddleName) > 50,
				"middle name cannot contain numbers": containsNumbers(sp.MiddleName),
			}),

		"date_of_birth": cause.Required(sp.DateOfBirth).
			Select(map[string]bool{
				"date of birth cannot be in the future": sp.DateOfBirth.After(time.Now()),
				"student too young":                     age < 5,
				"invalid date of birth":                 age > 120,
			}),

		"gender": cause.Required(sp.Gender).
			Select(map[string]bool{
				"invalid gender value": !isValidGender(sp.Gender),
			}),

		"nationality": cause.Required(sp.Nationality).
			Select(map[string]bool{
				"nationality too short": len(sp.Nationality) < 2,
				"nationality too long":  len(sp.Nationality) > 50,
			}),

		"email": cause.Required(sp.Email).
			Select(map[string]bool{
				"invalid email format":         !isValidEmail(sp.Email),
				"non-educational email domain": !isEducationalEmail(sp.Email),
			}),

		"phone_number": cause.Required(sp.PhoneNumber).
			Select(map[string]bool{
				"invalid phone number format": !isValidPhoneNumber(sp.PhoneNumber),
			}),

		"address":           cause.Required(sp.Address).Err(),
		"emergency_contact": cause.Required(sp.EmergencyContact).Err(),
	}.Err()
}

func (ap *AcademicProfile) Validate() error {
	return cause.Map{
		"program": cause.Required(ap.Program).
			Select(map[string]bool{
				"invalid program":       !isValidProgram(ap.Program),
				"program name too long": len(ap.Program) > 100,
			}),

		"major": cause.Required(ap.Major).
			Select(map[string]bool{
				"invalid major":       !isValidMajor(ap.Major),
				"major name too long": len(ap.Major) > 100,
			}),

		"minor": cause.Optional(ap.Minor).
			Select(map[string]bool{
				"invalid minor":                     !isValidMajor(ap.Minor),
				"minor name too long":               len(ap.Minor) > 100,
				"minor cannot be the same as major": ap.Minor == ap.Major,
			}),

		"year": cause.Required(ap.Year).
			Select(map[string]bool{
				"academic year must be between 1 and 8": ap.Year < 1 || ap.Year > 8,
			}),

		"semester": cause.Required(ap.Semester).
			Select(map[string]bool{
				"invalid semester": !isValidSemester(ap.Semester),
			}),

		"credits": cause.Required(ap.Credits).
			Select(map[string]bool{
				"credits cannot be negative":      ap.Credits < 0,
				"credits exceeds maximum allowed": ap.Credits > 200,
			}),

		"gpa": cause.Required(ap.GPA).
			Select(map[string]bool{
				"GPA must be between 0.0 and 4.0": ap.GPA < 0.0 || ap.GPA > 4.0,
			}),

		"expected_graduation": cause.Required(ap.ExpectedGraduation).
			Select(map[string]bool{
				"expected graduation cannot be in the past": ap.ExpectedGraduation.Before(time.Now()),
				"expected graduation too far in future":     ap.ExpectedGraduation.After(time.Now().AddDate(10, 0, 0)),
			}),

		"advisor_id": cause.Required(ap.AdvisorID).
			Select(map[string]bool{
				"invalid advisor ID": !isValidFacultyID(ap.AdvisorID),
			}),

		"degree_type": cause.Required(ap.DegreeType).
			Select(map[string]bool{
				"invalid degree type": !isValidDegreeType(ap.DegreeType),
			}),
	}.Err()
}

func (ce *CourseEnrollment) Validate() error {
	return cause.Map{
		"course_id": cause.Required(ce.CourseID).
			Select(map[string]bool{
				"invalid course ID format": !isValidCourseID(ce.CourseID),
				"course ID too long":       len(ce.CourseID) > 15,
			}),

		"course_name": cause.Required(ce.CourseName).
			Select(map[string]bool{
				"course name too short": len(ce.CourseName) < 3,
				"course name too long":  len(ce.CourseName) > 100,
			}),

		"credits": cause.Required(ce.Credits).
			Select(map[string]bool{
				"credits cannot be negative":                  ce.Credits < 0,
				"credits exceeds maximum for a single course": ce.Credits > 10,
			}),

		"instructor_id": cause.Required(ce.InstructorID).
			Select(map[string]bool{
				"invalid instructor ID": !isValidFacultyID(ce.InstructorID),
			}),

		"schedule": cause.Required(ce.Schedule).Err(),

		"enroll_date": cause.Required(ce.EnrollDate).
			Select(map[string]bool{
				"enrollment date cannot be in the future": ce.EnrollDate.After(time.Now()),
				"enrollment date too old":                 ce.EnrollDate.Before(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
			}),

		"grade": cause.Optional(ce.Grade).
			Select(map[string]bool{
				"invalid grade": !isValidGrade(ce.Grade),
			}),

		"status": cause.Required(ce.Status).
			Select(map[string]bool{
				"invalid course status": !isValidCourseStatus(ce.Status),
			}),

		"prerequisites": cause.Optional(ce.Prerequisites).
			Select(map[string]bool{
				"too many prerequisites listed": len(ce.Prerequisites) > 10,
			}),
	}.Err()
}

func (s *Schedule) Validate() error {
	return cause.Map{
		"days_of_week": cause.Required(s.DaysOfWeek).
			Select(map[string]bool{
				"at least one day of week required": len(s.DaysOfWeek) == 0,
				"cannot have more than 7 days":      len(s.DaysOfWeek) > 7,
				"invalid days of week":              !areValidDays(s.DaysOfWeek),
			}),

		"start_time": cause.Required(s.StartTime).
			Select(map[string]bool{
				"invalid start time format":      !isValidTime(s.StartTime),
				"start time outside class hours": !isValidClassTime(s.StartTime),
			}),

		"end_time": cause.Required(s.EndTime).
			Select(map[string]bool{
				"invalid end time format":           !isValidTime(s.EndTime),
				"end time outside class hours":      !isValidClassTime(s.EndTime),
				"end time must be after start time": !isEndTimeAfterStart(s.StartTime, s.EndTime),
			}),

		"room": cause.Required(s.Room).
			Select(map[string]bool{
				"invalid room number format": !isValidRoomNumber(s.Room),
				"room number too long":       len(s.Room) > 20,
			}),

		"building": cause.Required(s.Building).
			Select(map[string]bool{
				"building name required": len(s.Building) < 1,
				"building name too long": len(s.Building) > 50,
			}),
	}.Err()
}

func (gi *GuardianInfo) Validate() error {
	return cause.Map{
		"name": cause.Required(gi.Name).
			Select(map[string]bool{
				"guardian name too short": len(gi.Name) < 2,
				"guardian name too long":  len(gi.Name) > 100,
			}),

		"relationship": cause.Required(gi.Relationship).
			Select(map[string]bool{
				"invalid guardian relationship": !isValidGuardianRelationship(gi.Relationship),
			}),

		"phone": cause.Required(gi.Phone).
			Select(map[string]bool{
				"invalid guardian phone number": !isValidPhoneNumber(gi.Phone),
			}),

		"email": cause.Optional(gi.Email).
			Select(map[string]bool{
				"invalid guardian email": !isValidEmail(gi.Email),
			}),

		"address": cause.Required(gi.Address).Err(),
	}.Err()
}

func (fa *FinancialAid) Validate() error {
	return cause.Map{
		"type": cause.Required(fa.Type).
			Select(map[string]bool{
				"invalid financial aid type": !isValidFinancialAidType(fa.Type),
			}),

		"amount": cause.Required(fa.Amount).
			Select(map[string]bool{
				"financial aid amount must be positive": fa.Amount <= 0,
				"financial aid amount exceeds maximum":  fa.Amount > 100000,
			}),

		"semester": cause.Required(fa.Semester).
			Select(map[string]bool{
				"invalid semester": !isValidSemester(fa.Semester),
			}),

		"academic_year": cause.Required(fa.AcademicYear).
			Select(map[string]bool{
				"invalid academic year format": !isValidAcademicYear(fa.AcademicYear),
			}),

		"requirements": cause.Optional(fa.Requirements).
			Select(map[string]bool{
				"too many requirements listed": len(fa.Requirements) > 20,
			}),

		"status": cause.Required(fa.Status).
			Select(map[string]bool{
				"invalid financial aid status": !isValidFinancialAidStatus(fa.Status),
			}),

		"application_date": cause.Required(fa.ApplicationDate).
			Select(map[string]bool{
				"application date cannot be in the future": fa.ApplicationDate.After(time.Now()),
				"application date too old":                 fa.ApplicationDate.Before(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
			}),

		"expiry_date": cause.Required(fa.ExpiryDate).
			Select(map[string]bool{
				"expiry date cannot be before application date": fa.ExpiryDate.Before(fa.ApplicationDate),
				"expiry date too far in future":                 fa.ExpiryDate.After(time.Now().AddDate(5, 0, 0)),
			}),
	}.Err()
}

func (t *Transcript) Validate() error {
	return cause.Map{
		"course_id": cause.Required(t.CourseID).
			Select(map[string]bool{
				"invalid course ID": !isValidCourseID(t.CourseID),
			}),

		"course_name": cause.Required(t.CourseName).
			Select(map[string]bool{
				"course name too short": len(t.CourseName) < 3,
				"course name too long":  len(t.CourseName) > 100,
			}),

		"credits": cause.Required(t.Credits).
			Select(map[string]bool{
				"credits cannot be negative": t.Credits < 0,
				"credits exceeds maximum":    t.Credits > 10,
			}),

		"grade": cause.Required(t.Grade).
			Select(map[string]bool{
				"invalid grade": !isValidGrade(t.Grade),
			}),

		"grade_points": cause.Required(t.GradePoints).
			Select(map[string]bool{
				"grade points must be between 0.0 and 4.0": t.GradePoints < 0.0 || t.GradePoints > 4.0,
			}),

		"semester": cause.Required(t.Semester).
			Select(map[string]bool{
				"invalid semester": !isValidSemester(t.Semester),
			}),

		"academic_year": cause.Required(t.AcademicYear).
			Select(map[string]bool{
				"invalid academic year": !isValidAcademicYear(t.AcademicYear),
			}),

		"instructor_id": cause.Required(t.InstructorID).
			Select(map[string]bool{
				"invalid instructor ID": !isValidFacultyID(t.InstructorID),
			}),
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
