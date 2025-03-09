package types

type SseStatus string // status

const (
	Disconnected     SseStatus = "disconnected"
	Connected        SseStatus = "connected"
	UserJoin         SseStatus = "users"
	UserLeave        SseStatus = "user_leave"
	LiveContest      SseStatus = "live"
	StartContest     SseStatus = "start_contest"
	EndContest       SseStatus = "end_contest"
	CloseContest     SseStatus = "closed_contest"
	ContestInfo      SseStatus = "contest_info"
	ContestResults   SseStatus = "contest_results"
	ContestClosed    SseStatus = "contest_closed"
	ContestCountDown SseStatus = "contest_countdown"
	ErrContest                 = "err_contest"
)
