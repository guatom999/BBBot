package botREST

type (
	ToSendMessage struct {
		ProjectName string `json:"subject"`
		Value       string `json:"value"`
	}

	InCommingMessage struct {
		ProjectName string `json:"project_name"`
		Color       int    `json:"color"`
		Status      string `json:"status"`
		Message     string `json:"message"`
	}
)
