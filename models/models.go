package models

type Term struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

type School struct {
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

type Subject struct {
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

type SubjectAvailability struct {
	Term    int    `json:"term"`
	School  string `json:"school"`
	Subject string `json:"subject"`
}

type Instructor struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

type InstructorSubject struct {
	Instructor int    `json:"instructor"`
	Subject    string `json:"subject"`
}

type Building struct {
	Id   int     `json:"id"`
	Name string  `json:"name"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
}

type Room struct {
	Id         int    `json:"id"`
	BuildingId int    `json:"buildingId"`
	Name       string `json:"name"`
}

type Course struct {
	Id           int    `json:"id"`
	Title        string `json:"title"`
	Term         int    `json:"term"`
	School       string `json:"school"`
	Instructor   int    `json:"instructor"`
	Subject      string `json:"subject"`
	CatalogNum   string `json:"catalogNum"`
	Section      string `json:"section"`
	Room         int    `json:"room"`
	MeetingDays  string `json:"meetingDays"`
	StartTime    string `json:"startTime"`
	EndTime      string `json:"endTime"`
	StartDate    string `json:"startDate"`
	EndDate      string `json:"endDate"`
	Seats        int    `json:"seats"`
	Overview     string `json:"overview"`
	Topic        string `json:"topic"`
	Attributes   string `json:"attributes"`
	Requirements string `json:"requirements"`
	Component    string `json:"component"`
	ClassNum     int    `json:"classNum"`
	CourseId     int    `json:"courseId"`
}

type CourseDescription struct {
	Course int    `json:"course"`
	Name   string `json:"name"`
	Desc   string `json:"desc"`
}

type CourseComponent struct {
	Course      int    `json:"course"`
	Component   string `json:"component"`
	MeetingDays string `json:"meetingDays"`
	StartTime   string `json:"startTime"`
	EndTime     string `json:"endTime"`
	Section     string `json:"section"`
	Room        string `json:"room"`
}
