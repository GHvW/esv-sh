package lib

type ESVClient struct {
	baseURL string
	apiKey  string
}

// func (client *ESVClient) GetPassage(query string) (string, error) {
// 	req, err := http.NewRequest("GET", client.baseURL, nil)
// 	if err != nil {
// 		return "", err
// 	}

// 	res, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return "", err
// 	}

// }
