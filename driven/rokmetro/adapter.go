package rokmetro

import (
	"encoding/json"
	"errors"
	"fmt"
	"health/core"
	"io/ioutil"
	"log"
	"net/http"
)

//Adapter implements the Rokmetro interface
type Adapter struct {
	groupsBBHost string
	apiKey       string
}

//GetExtJoinExternalApproval loads the join groups external approvements
func (a *Adapter) GetExtJoinExternalApproval(externalApproverID string) ([]core.RokmetroJoinGroupExtApprovement, error) {
	url := fmt.Sprintf("%s/ext/join-external-approvements?external-approver-id=%s", a.groupsBBHost, externalApproverID)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("error getting join groups external approvements request - %s", err)
		return nil, err
	}
	req.Header.Set("ROKMETRO-EXTERNAL-API-KEY", a.apiKey)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error join groups external approvements data - %s", err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("getting join groups external approvements - error with response code - %d", resp.StatusCode)
		return nil, errors.New("getting join groups external approvements - error with response code != 200")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading the body data for getting join groups external approvements data request - %s", err)
		return nil, err
	}

	var result []core.RokmetroJoinGroupExtApprovement
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Printf("error converting data for getting join groups external approvements data request - %s", err)
		return nil, err
	}

	return result, nil
}

//UpdateExtJoinExternalApprovement approve/reject jea
func (a *Adapter) UpdateExtJoinExternalApprovement(jeaID string, status string) error {
	return nil
}

//NewRokmetroAdapter creates a new rokmetro adapter instance
func NewRokmetroAdapter(groupsBBHost string, apiKey string) *Adapter {
	return &Adapter{groupsBBHost: groupsBBHost, apiKey: apiKey}
}
