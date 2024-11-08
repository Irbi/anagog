package schema

type VisitReport struct {
	DataVer       int     `json:"DataVer"`
	UserId        int     `json:"UserId"`
	EnterTime     string  `json:"EnterTime"`
	ExitTime      string  `json:"ExitTime"`
	AlgorithmType int     `json:"AlgorithmType"`
	PoiId         int64   `json:"PoiId"`
	Latitude      float64 `json:"Latitude"`
	Longitude     float64 `json:"Longitude"`
}

type ActivityReport struct {
	DataVer        int     `json:"DataVer"`
	UserId         int     `json:"UserId"`
	StartTime      string  `json:"StartTime"`
	EndTime        string  `json:"EndTime"`
	ActivityType   int     `json:"ActivityType"`
	StartLatitude  float64 `json:"StartLatitude"`
	StartLongitude float64 `json:"StartLongitude"`
	EndLatitude    float64 `json:"EndLatitude"`
	EndLongitude   float64 `json:"EndLongitude"`
}