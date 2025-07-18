package entity

// MemberLevel 會員等級
type MemberLevel int

const (
	LevelNormal   MemberLevel = iota // 0: 普通會員
	LevelSilver                      // 1: 白銀會員
	LevelGold                        // 2: 黃金會員
	LevelPlatinum                    // 3: 白金會員
)

// Member 會員實體
type Member struct {
	ID    int64
	Name  string
	Level MemberLevel
}
