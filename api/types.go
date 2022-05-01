package api

// AuthResponse is returned on successful authorization. The access token is to be used in susbsequent requests.
type AuthResponse struct {
	AccessToken           string `json:"access_token"`
	ExpiresIn             int    `json:"expires_in"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresIn int    `json:"refresh_token_expires_in"`
}

// HTTP4XXError is returned when the access token is expired or invalid.
// To recover, a new access token must be generated using the refresh token.
type HTTP4XXError struct {
	Message string `json:"message"`
	Errors  []struct {
		Description string `json:"description"`
	} `json:"errors"`
	TrackingID string `json:"trackingId"`
}

type OAuthRequest struct {
	ClientID     string
	ClientSecret string
	Scope        string
}

type RefreshTokenRequest struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RefreshToken string `json:"refresh_token"`
}

type GenericPage struct {
	Heading          string
	Message          string
	ShowHomeRedirect bool
	ShowAPIRedirect  bool
}

type MeetingsList struct {
	Items []MeetingSeries `json:"items"`
}

type MeetingSeries struct {
	ID                                  string                 `json:"id"`
	MeetingNumber                       string                 `json:"meetingNumber"`
	Title                               string                 `json:"title"`
	Agenda                              string                 `json:"agenda"`
	Password                            string                 `json:"password"`
	PhoneAndVideoSystemPassword         string                 `json:"phoneAndVideoSystemPassword"`
	MeetingType                         string                 `json:"meetingType"`
	State                               string                 `json:"state"`
	Timezone                            string                 `json:"timezone"`
	Start                               string                 `json:"start"`
	End                                 string                 `json:"end"`
	Recurrence                          string                 `json:"recurrence"`
	HostUserID                          string                 `json:"hostUserId"`
	HostDisplayName                     string                 `json:"hostDisplayName"`
	HostEmail                           string                 `json:"hostEmail"`
	HostKey                             string                 `json:"hostKey"`
	SiteURL                             string                 `json:"siteUrl"`
	WebLink                             string                 `json:"webLink"`
	SipAddress                          string                 `json:"sipAddress"`
	DialInIPAddress                     string                 `json:"dialInIpAddress"`
	RoomID                              string                 `json:"roomId"`
	EnableAutoRecordMeeting             bool                   `json:"enableAutoRecordMeeting"`
	AllowUserToBeCoHost                 bool                   `json:"allowUserToBeCoHost"`
	EnabledJoinBeforeHost               bool                   `json:"enabledJoinBeforeHost"`
	EnableConnectAudioBeforeHost        bool                   `json:"enableConnectAudioBeforeHost"`
	JoinBeforeHostMinutes               int                    `json:"joinBeforeHostMinutes"`
	ExcludePassword                     bool                   `json:"excludePassword"`
	PublicMeeting                       bool                   `json:"publicMeeting"`
	ReminderTime                        int                    `json:"reminderTime"`
	UnlockedMeetingJoinSecurity         string                 `json:"unlockedMeetingJoinSecurity"`
	SessionTypeID                       int                    `json:"sessionTypeId"`
	ScheduledType                       string                 `json:"scheduledType"`
	EnabledWebcastView                  bool                   `json:"enabledWebcastView"`
	PanelistPassword                    string                 `json:"panelistPassword"`
	PhoneAndVideoSystemPanelistPassword string                 `json:"phoneAndVideoSystemPanelistPassword"`
	EnableAutomaticLock                 bool                   `json:"enableAutomaticLock"`
	AutomaticLockMinutes                int                    `json:"automaticLockMinutes"`
	AllowFirstUserToBeCoHost            bool                   `json:"allowFirstUserToBeCoHost"`
	AllowAuthenticatedDevices           bool                   `json:"allowAuthenticatedDevices"`
	Telephony                           map[string]interface{} `json:"telephony"`
	Registration                        map[string]interface{} `json:"registration"`
	IntegrationTags                     []string               `json:"integrationTags"`
}
