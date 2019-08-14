package auth

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

const uri = "https://people.googleapis.com/v1/people/me?personFields=names,emailAddresses"

type GooglePeople struct {
	client *http.Client
}

func NewGooglePeople() *GooglePeople {
	return &GooglePeople{client: &http.Client{}}
}

type meData struct {
	ResourceName string `json:"resourceName"`
	Etag         string `json:"etag"`
	Names        []struct {
		Metadata struct {
			Primary bool `json:"primary"`
			Source  struct {
				Type string `json:"type"`
				Id   string `json:"id"`
			} `json:"source"`
		} `json:"metadata"`
		DisplayName          string `json:"displayName"`
		FamilyName           string `json:"familyName"`
		GivenName            string `json:"givenName"`
		DisplayNameLastFirst string `json:"displayNameLastFirst"`
	} `json:"names"`
	EmailAddresses []struct {
		Metadata struct {
			Primary  bool `json:"primary"`
			Verified bool `json:"verified"`
			Source   struct {
				Type string `json:"type"`
				Id   string `json:"id"`
			} `json:"source"`
		} `json:"metadata"`
		Value string `json:"value"`
	} `json:"emailAddresses"`
}

func (gp *GooglePeople) Me(accessCode string) (string, error) {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", "Bearer "+accessCode)

	resp, err := gp.client.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("invalid access code")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	data := meData{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}

	for _, emailAddress := range data.EmailAddresses {
		if emailAddress.Metadata.Primary {
			return emailAddress.Value, nil
		}
	}

	return "", errors.New("no email address")
}
