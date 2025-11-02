package botREST

type (
	ToSendMessage struct {
		ProjectName string `json:"subject"`
		Value       string `json:"value"`
	}

	InCommingMessage struct {
		ProjectName string `json:"project_name"`
		Status      string `json:"status"`
	}
)
