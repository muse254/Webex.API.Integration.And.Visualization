package types

import (
	"errors"
)

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
	MeetingID     string                `json:"meeting_id"`
	MediaSessions []MediaSessionQuality `json:"items"`
}

type MediaSessionQuality struct {
	MeetingID        string `json:"meetingId"`
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
	PacketLoss       []float32 `json:"packetLoss,omitempty"`
	Latency          []float32 `json:"latency,omitempty"`
	ResolutionHeight []float32 `json:"resolutionHeight,omitempty"`
	FrameRate        []float32 `json:"frameRate,omitempty"`
	MediaBitRate     []float32 `json:"mediaBitRate"`
	Codec            string    `json:"codec"`
	Jitter           []float32 `json:"jitter,omitempty"`
	TransportType    string    `json:"transportType"`
}

type Resources struct {
	ProcessAverageCPU []float32 `json:"processAverageCPU"`
	ProcessMaxCPU     []float32 `json:"processMaxCPU"`
	SystemAverageCPU  []float32 `json:"systemAverageCPU"`
	SystemMaxCPU      []float32 `json:"systemMaxCPU"`
}

type VisualData struct {
	MeetingID string          `json:"meeting_id"`
	DataPoint string          `json:"data_point"`
	Sessions  []VisualSession `json:"sessions"`
}

type VisualSession struct {
	StartTime  string    `json:"start_time"`
	EndTime    string    `json:"end_time"`
	PacketLoss []float32 `json:"packet_loss"`
	Latency    []float32 `json:"latency"`
	Jitter     []float32 `json:"jitter"`
}

func GetAllVisualData(qualities *MeetingQualities) ([]VisualData, error) {
	vData := func(dp string) VisualData {
		return VisualData{
			MeetingID: qualities.MeetingID,
			DataPoint: dp,
			Sessions: func() []VisualSession {
				var visualSessions []VisualSession
				for _, session := range qualities.MediaSessions {
					switch dp {
					case "audio_in":
						visualSessions = append(visualSessions, *populateVisualSession(session.AudioIn))
					case "audio_out":
						visualSessions = append(visualSessions, *populateVisualSession(session.AudioOut))
					case "video_in":
						visualSessions = append(visualSessions, *populateVisualSession(session.VideoIn))
					case "video_out":
						visualSessions = append(visualSessions, *populateVisualSession(session.VideoOut))
					case "share_in":
						visualSessions = append(visualSessions, *populateVisualSession(session.ShareIn))
					case "share_out":
						visualSessions = append(visualSessions, *populateVisualSession(session.ShareOut))
					}
				}
				return visualSessions
			}(),
		}
	}

	return []VisualData{vData("audio_in"), vData("audio_out"), vData("video_in"),
		vData("video_out"), vData("share_in"), vData("share_out"),
	}, nil
}

func GetVisualData(qualities *MeetingQualities, dp string) (*VisualData, error) {
	var visualSessions []VisualSession
	for _, session := range qualities.MediaSessions {
		switch dp {
		case "video_in":
			visualSessions = append(visualSessions, *populateVisualSession(session.VideoIn))
		case "video_out":
			visualSessions = append(visualSessions, *populateVisualSession(session.VideoOut))
		case "audio_in":
			visualSessions = append(visualSessions, *populateVisualSession(session.VideoOut))
		case "audio_out":
			visualSessions = append(visualSessions, *populateVisualSession(session.VideoOut))
		case "share_in":
			visualSessions = append(visualSessions, *populateVisualSession(session.VideoOut))
		case "share_out":
			visualSessions = append(visualSessions, *populateVisualSession(session.VideoOut))
		default:
			return nil, errors.New(`invalid request, "dp" parameter not recognized`)
		}
	}

	return &VisualData{
		MeetingID: qualities.MeetingID,
		DataPoint: dp,
		Sessions:  visualSessions,
	}, nil
}

func populateVisualSession(data []MediaQualityData) *VisualSession {
	var packetLoss, latency, jitter []float32
	for _, val := range data {
		packetLoss = append(packetLoss, val.PacketLoss...)
		latency = append(latency, val.Latency...)
		jitter = append(jitter, val.Jitter...)
	}

	set := len(data)
	return &VisualSession{
		StartTime:  data[0].StartTime,
		EndTime:    data[set-1].EndTime,
		PacketLoss: packetLoss,
		Latency:    latency,
		Jitter:     jitter,
	}
}
