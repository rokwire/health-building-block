package rokmetro

//Adapter implements the Rokmetro interface
type Adapter struct {
	groupsBBHost string
	apiKey       string
}

//GetExtJoinExternalApproval loads the user data by uuid
func (a *Adapter) GetExtJoinExternalApproval(externalApproverID string) (interface{}, error) {
	//TODO
	return nil, nil
}

//NewRokmetroAdapter creates a new rokmetroadapter instance
func NewRokmetroAdapter(groupsBBHost string, apiKey string) *Adapter {
	return &Adapter{groupsBBHost: groupsBBHost, apiKey: apiKey}
}
