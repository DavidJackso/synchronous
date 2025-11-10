package entity

type LeaderboardEntry struct {
	Rank           int     `json:"rank"`
	UserID         string  `json:"userId"`
	UserName       string  `json:"userName"`
	AvatarURL      *string `json:"avatarUrl"`
	TasksCompleted int     `json:"tasksCompleted"`
	FocusTime      int     `json:"focusTime"` // в минутах
	Score          int     `json:"score"`
}

type LeaderboardPeriod string

const (
	LeaderboardPeriodDay   LeaderboardPeriod = "day"
	LeaderboardPeriodWeek  LeaderboardPeriod = "week"
	LeaderboardPeriodMonth LeaderboardPeriod = "month"
	LeaderboardPeriodAll   LeaderboardPeriod = "all"
)
