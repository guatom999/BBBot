package botRepositories

type (
	IBotRepository interface {
	}

	botRepository struct {
	}
)

func NewBotRepository() IBotRepository {
	return &botRepository{}
}
