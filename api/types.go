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

type MeetingQualities struct {
	Items []MediaSessionQuality `json:"items"`
}

type MediaSessionQuality struct {
	MeetingId        string `json:"meetingId"`
	DisplayName      string `json:"displayName"`
	Email            string `json:"email"`
	Joined           string `json:"joined"`
	Client           string `json:"client"`
	ClientVersion    string `json:"clientVersion"`
	OsType           string `json:"osType"`
	OsVersion        string `json:"osVersion"`
	HardwareType     string `json:"hardwareType"`
	SpeakerName      string `json:"speakerName"`
	NetworkType      string `json:"networkType"`
	LocalIP          string `json:"localIP"`
	PublicIP         string `json:"publicIP"`
	MaskedLocalIP    string `json:"maskedLocalIP"`
	MaskedPublicIP   string `json:"maskedPublicIP"`
	Camera           string `json:"camera"`
	Microphone       string `json:"microphone"`
	ServerRegion     string `json:"serverRegion"`
	VideoMeshCluster string `json:"videoMeshCluster"`
	ParticipantID    string `json:"participantId"`
	// VideoIn is the collection of downstream (sent to the client) video quality data.
	VideoIn []MediaQualityData `json:"videoIn"`
	// VideoOut is the collection of upstream (sent from the client) video quality data.
	VideoOut []MediaQualityData `json:"videoOut"`
	// AudioIn is the collection of downstream (sent to the client) audio quality data.
	AudioIn []MediaQualityData `json:"audioIn"`
	// AudioOut is the collection of upstream (sent from the client) audio quality data.
	AudioOut []MediaQualityData `json:"audioOut"`
	// ShareIn is the collection of downstream (sent to the client) share quality data.
	ShareIn []MediaQualityData `json:"shareIn"`
	// ShareOut is the collection of upstream (sent from the client) share quality data.
	ShareOut []MediaQualityData `json:"shareOut"`
	// Resources sre devide resources such as CPU and memory.
	Resources []Resources `json:"resources"`
}

type MediaQualityData struct {
	SamplingInterval int       `json:"samplingInterval"`
	StartTime        string    `json:"startTime"`
	EndTime          string    `json:"endTime"`
	PacketLoss       []float32 `json:"packetLoss"`
	Latency          []int     `json:"latency"`
	ResolutionHeight []int     `json:"resolutionHeight,omitempty"`
	FrameRate        []int     `json:"frameRate,omitempty"`
	MediaBitRate     []int     `json:"mediaBitRate"`
	Codec            string    `json:"codec"`
	Jitter           []int     `json:"jitter"`
	TransportType    string    `json:"transportType"`
}

type Resources struct {
	ProcessAverageCPU []string `json:"processAverageCPU"`
	ProcessMaxCPU     []string `json:"processMaxCPU"`
	SystemAverageCPU  []string `json:"systemAverageCPU"`
	SystemMaxCPU      []string `json:"systemMaxCPU"`
}
