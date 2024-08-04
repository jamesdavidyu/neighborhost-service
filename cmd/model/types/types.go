package types

import (
	"time"
)

type ZipcodeStore interface {
	GetZipcodeData(zipcode string) (*Zipcodes, error)
	ValidateZipcode(city string, state string, zipcode string) (*Zipcodes, error)
}

type NeighborStore interface {
	GetNeighborWithEmailOrUsername(emailOrUsername string) (*Neighbors, error)
	GetNeighborById(id int) (*Neighbors, error)
	CreateNeighbor(Neighbors) error
	GetNeighborWithEmail(email string) (*Neighbors, error)
	GetNeighborWithUsername(username string) (*Neighbors, error)
	UpdateZipcodeWithId(Neighbors) error
	UpdatePasswordWithId(Neighbors) error
}

type AddressStore interface {
	CreateAddress(Addresses) error
	GetAddressesByZipcode(zipcode string) (*Addresses, error)
	GetAddressIdByAddress(
		address string,
		city string,
		state string,
		zipcode string,
		Type string,
		neighborId int,
	) (*Addresses, error)
	GetAddressByNeighborId(id int) (*Addresses, error)
	GetAddressesByNeighborId(id int) (*Addresses, error)
}

type EventStore interface {
	GetPublicEvents() ([]Events, error)
	GetEventsByZipcode(zipcode string, dateTime time.Time) ([]EventAddresses, error)
	GetZipcodeEventsOnDate(zipcode string, dateTime time.Time) ([]EventAddresses, error)
	GetZipcodeEventsBeforeDate(zipcode string, dateTime time.Time) ([]EventAddresses, error)
	GetZipcodeEventsAfterDate(zipcode string, dateTime time.Time) ([]EventAddresses, error)
	GetEventsByNeighborhoodId(neighborhoodId int, dateTime time.Time) ([]EventAddresses, error)
	GetNeighborhoodEventsOnDate(neighborhoodId int, dateTime time.Time) ([]EventAddresses, error)
	GetNeighborhoodEventsBeforeDate(neighborhoodId int, dateTime time.Time) ([]EventAddresses, error)
	GetNeighborhoodEventsAfterDate(neighborhoodId int, dateTime time.Time) ([]EventAddresses, error)
	GetEventsByCity(city string, state string, zipcode string, dateTime time.Time) ([]EventAddresses, error)
	GetCityEventsOnDate(city string, state string, zipcode string, dateTime time.Time) ([]EventAddresses, error)
	GetCityEventsBeforeDate(city string, state string, zipcode string, dateTime time.Time) ([]EventAddresses, error)
	GetCityEventsAfterDate(city string, state string, zipcode string, dateTime time.Time) ([]EventAddresses, error)
	// GetAllEvents(dateTime time.Time) ([]EventAddresses, error)
	CreateEvent(Events) error
}

type FriendStore interface {
	GetFriendsByNeighborId(neighborId int) ([]Friends, error)
	GetFriendRequestsByNeighborId(neighborId int) ([]FriendRequests, error)
	CreateFriendRequest(FriendRequests) error
}

type ProfileStore interface {
	GetProfileByNeighborId(neighborId int) (*Profiles, error)
}

type NeighborhoodStore interface {
	GetNeighborhoods() ([]Neighborhoods, error)
	CreateNeighborhood(Neighborhoods) error
}

type Zipcodes struct {
	Zipcode  string `json:"zipcode"`
	City     string `json:"city"`
	State    string `json:"state"`
	Timezone string `json:"timezone"`
}

type Neighborhoods struct {
	Id           int       `json:"id"`
	Neighborhood string    `json:"neighborhood"`
	CreatedAt    time.Time `json:"createdAt"`
}

type Neighbors struct {
	Id             int       `json:"id"`
	Email          string    `json:"email"`
	Username       string    `json:"username"`
	Zipcode        string    `json:"zipcode"`
	Password       string    `json:"-"`
	Verified       bool      `json:"verified"`
	Ip             string    `json:"ip"`
	NeighborhoodId int       `json:"neighborhoodId"`
	CreatedAt      time.Time `json:"createdAt"`
	// role
}

type Register struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required"`
	Zipcode  string `json:"zipcode" validate:"required,min=5,max=5"`
	Password string `json:"password" validate:"required,min=8"`
}

type Login struct {
	EmailOrUsername string `json:"emailOrUsername" validate:"required"`
	Password        string `json:"password" validate:"required,min=8"`
}

type UpdatePassword struct {
	Password string `json:"password" validate:"required,min=8"`
}

type Addresses struct {
	Id             int       `json:"id"`
	FirstName      string    `json:"firstName"`
	LastName       string    `json:"lastName"`
	Address        string    `json:"address"`
	City           string    `json:"city"`
	State          string    `json:"state"`
	Zipcode        string    `json:"zipcode"`
	Type           string    `json:"type"`
	NeighborId     int       `json:"neighborId"`
	NeighborhoodId int       `json:"neighborhoodId"`
	RecordedAt     time.Time `json:"recordedAt"`
}

type AddressPayload struct {
	FirstName      string `json:"firstName"` // need to add validation
	LastName       string `json:"lastName"`
	Address        string `json:"address"`
	City           string `json:"city"`
	State          string `json:"state"`
	Zipcode        string `json:"zipcode"`
	Type           string `json:"type"`
	NeighborId     int    `json:"neighborId"`     // need?
	NeighborhoodId int    `json:"neighborhoodId"` // need?
}

type Profiles struct {
	NeighborId               int       `json:"id"`
	Bio                      string    `json:"bio"`
	DateOfBirth              time.Time `json:"dateOfBirth"`
	DateOfBirthPublic        bool      `json:"dateOfBirthPublic"`
	Gender                   string    `json:"gender"`
	GenderPublic             bool      `json:"genderPublic"`
	Race                     string    `json:"race"`
	RacePublic               bool      `json:"racePublic"`
	Ethnicity                string    `json:"ethnicity"`
	EthnicityPublic          bool      `json:"ethnicityPublic"`
	RelationshipStatus       string    `json:"relationshipStatus"`
	RelationshipStatusPublic bool      `json:"relationshipStatusPublic"`
	Religion                 string    `json:"religion"`
	ReligionPublic           bool      `json:"religionPublic"`
	Politics                 string    `json:"politics"`
	PoliticsPublic           bool      `json:"politicsPublic"`
}

type Bios struct {
	Id         int       `json:"id"`
	NeighborId int       `json:"neighborId"`
	Bio        string    `json:"bio"`
	CreatedAt  time.Time `json:"createdAt"`
}

type DatesOfBirth struct {
	Id          int       `json:"id"`
	NeighborId  int       `json:"neighborId"`
	DateOfBirth time.Time `json:"dateOfBirth"`
	Public      bool      `json:"public"`
	RecordedAt  time.Time `json:"recordedAt"`
}

type Genders struct {
	Id         int       `json:"id"`
	NeighborId int       `json:"neighborId"`
	Gender     string    `json:"gender"`
	Public     bool      `json:"public"`
	RecordedAt time.Time `json:"recordedAt"`
}

type Races struct {
	Id         int       `json:"id"`
	NeighborId int       `json:"neighborId"`
	Race       string    `json:"race"`
	Public     bool      `json:"public"`
	RecordedAt time.Time `json:"recordedAt"`
}

type Ethnicities struct {
	Id         int       `json:"id"`
	NeighborId int       `json:"neighborId"`
	Ethnicity  string    `json:"ethnicity"`
	Public     bool      `json:"public"`
	RecordedAt time.Time `json:"recordedAt"`
}

type RelationshipStatuses struct {
	Id                 int       `json:"id"`
	NeighborId         int       `json:"neighborId"`
	RelationshipStatus string    `json:"relationshipStatus"`
	Public             bool      `json:"public"`
	RecordedAt         time.Time `json:"recordedAt"`
}

type Religions struct {
	Id         int       `json:"id"`
	NeighborId int       `json:"neighborId"`
	Religion   string    `json:"religion"`
	Public     bool      `json:"public"`
	RecordedAt time.Time `json:"recordedAt"`
}

type Politics struct {
	Id         int       `json:"id"`
	NeighborId int       `json:"neighborId"`
	Politics   string    `json:"politics"`
	Public     bool      `json:"public"`
	RecordedAt time.Time `json:"recordedAt"`
}

type Education struct {
	Id         int       `json:"id"`
	NeighborId int       `json:"neighborId"`
	School     string    `json:"school"`
	Degree     string    `json:"degree"`
	Start      time.Time `json:"start"`
	End        time.Time `json:"end"`
	Public     bool      `json:"public"`
	RecordedAt time.Time `json:"recordedAt"`
}

type Occupations struct {
	Id         int       `json:"id"`
	NeighborId int       `json:"neighborId"`
	Role       string    `json:"role"`
	Employer   string    `json:"employer"`
	Start      time.Time `json:"start"`
	End        time.Time `json:"end"`
	Public     bool      `json:"public"`
	RecordedAt time.Time `json:"recordedAt"`
}

type Interests struct {
	Id         int       `json:"id"`
	NeighborId int       `json:"neighborId"`
	Interest   string    `json:"interest"`
	RecordedAt time.Time `json:"recordedAt"`
}

type Events struct {
	Id             int       `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Start          time.Time `json:"start"`
	End            time.Time `json:"end"`
	Reoccurrence   string    `json:"reoccurrence"`
	ForUnloggedins bool      `json:"forUnloggedins"`
	ForUnverifieds bool      `json:"forUnverifieds"`
	InviteOnly     bool      `json:"inviteOnly"`
	HostId         int       `json:"hostId"`
	AddressId      int       `json:"addressId"`
	CreatedAt      time.Time `json:"createdAt"`
	// add category
}

type CreateEventPayload struct {
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Start          time.Time `json:"start"`
	End            time.Time `json:"end"`
	Reoccurrence   string    `json:"reoccurrence"`
	ForUnloggedins bool      `json:"forUnloggedins"`
	ForUnverifieds bool      `json:"forUnverifieds"`
	InviteOnly     bool      `json:"inviteOnly"`
	Address        string    `json:"address"`
	City           string    `json:"city"`
	State          string    `json:"state"`
	Zipcode        string    `json:"zipcode"`
	Type           string    `json:"type"`
	HostId         int       `json:"hostId"`
	AddressId      int       `json:"addressId"`
}

type LocationFilterPayload struct {
	LocationFilter string `json:"locationFilter"`
}

type DateFilterPayload struct {
	DateFilter string    `json:"dateFilter"`
	DateTime   time.Time `json:"dateTime"`
}

type LocationDateFilterPayload struct {
	LocationFilter string    `json:"locationFilter"`
	DateFilter     string    `json:"dateFilter"`
	DateTime       time.Time `json:"dateTime"`
}

type EventAddresses struct {
	Id               int       `json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Start            time.Time `json:"start"`
	End              time.Time `json:"end"`
	Reoccurrence     string    `json:"reoccurrence"`
	ForUnloggedins   bool      `json:"forUnloggedins"`
	ForUnverifieds   bool      `json:"forUnverifieds"`
	InviteOnly       bool      `json:"inviteOnly"`
	HostId           int       `json:"hostId"`
	AddressId        int       `json:"addressId"`
	CreatedAt        time.Time `json:"createdAt"`
	AddressAddressId int       `json:"addressAddressId"`
	FirstName        string    `json:"firstName"`
	LastName         string    `json:"lastName"`
	Address          string    `json:"address"`
	City             string    `json:"city"`
	State            string    `json:"state"`
	Zipcode          string    `json:"zipcode"`
	Type             string    `json:"type"`
	NeighborId       int       `json:"neighborId"`
	NeighborhoodId   int       `json:"neighborhoodId"`
	RecordedAt       time.Time `json:"recordedAt"`
}

type EventInvites struct {
	Id                int       `json:"id"`
	EventId           int       `json:"eventId"`
	InvitedNeighborId int       `json:"invitedNeighborId"`
	InvitedAt         time.Time `json:"invitedAt"`
}

type Friends struct {
	Id                int       `json:"id"`
	NeighborId        int       `json:"neighborId"`
	NeighborsFriendId int       `json:"neighborsFriendId"`
	FriendedAt        time.Time `json:"friendedAt"`
}

type FriendRequests struct {
	Id                 int       `json:"id"`
	NeighborId         int       `json:"neighborId"`
	RequestingFriendId int       `json:"requestingFriendId"`
	Status             string    `json:"status"`
	FriendRequestedAt  time.Time `json:"friendRequestedAt"`
}
