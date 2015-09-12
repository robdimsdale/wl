package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/robdimsdale/wundergo"
)

// Memberships returns the memberships the client can access.
func (c oauthClient) Memberships() (*[]wundergo.Membership, error) {
	url := fmt.Sprintf(
		"%s/memberships",
		c.apiURL,
	)

	req, err := c.newGetRequest(url)
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	memberships := []wundergo.Membership{}
	err = json.NewDecoder(resp.Body).Decode(&memberships)
	if err != nil {
		c.logger.Println(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	return &memberships, nil
}

// Membership returns the Membership associated with the provided membershipID.
func (c oauthClient) Membership(membershipID uint) (*wundergo.Membership, error) {
	if membershipID == 0 {
		return nil, errors.New("membershipID must be > 0")
	}

	url := fmt.Sprintf(
		"%s/memberships/%d",
		c.apiURL,
		membershipID,
	)

	req, err := c.newGetRequest(url)
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	membership := wundergo.Membership{}
	err = json.NewDecoder(resp.Body).Decode(&membership)
	if err != nil {
		c.logger.Println(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	return &membership, nil
}

// MembershipsForListID returns the Memberships for the List associated with
// the provided listID.
func (c oauthClient) MembershipsForListID(listID uint) (*[]wundergo.Membership, error) {
	if listID == 0 {
		return nil, errors.New("listID must be > 0")
	}

	url := fmt.Sprintf(
		"%s/memberships?list_id=%d",
		c.apiURL,
		listID,
	)

	req, err := c.newGetRequest(url)
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	memberships := []wundergo.Membership{}
	err = json.NewDecoder(resp.Body).Decode(&memberships)
	if err != nil {
		c.logger.Println(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	return &memberships, nil
}

// AddMemberToListViaUserID creates a new Membership associating the User with
// the List.
func (c oauthClient) AddMemberToListViaUserID(userID uint, listID uint, muted bool) (*wundergo.Membership, error) {
	if userID == 0 {
		return nil, errors.New("userID must be > 0")
	}

	if listID == 0 {
		return nil, errors.New("listID must be > 0")
	}

	url := fmt.Sprintf("%s/memberships", c.apiURL)
	body := []byte(
		fmt.Sprintf(
			`{"user_id":%d,"list_id":%d,"muted":%t}`,
			userID,
			listID,
			muted,
		),
	)

	req, err := c.newPostRequest(url, body)
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if err != nil {
		c.logger.Println(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusCreated)
	}

	membership := wundergo.Membership{}
	err = json.NewDecoder(resp.Body).Decode(&membership)
	if err != nil {
		c.logger.Println(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	return &membership, nil
}

// AddMemberToListViaEmailAddress creates a new Membership joining the List
// with the user associated with the provided email address.
func (c oauthClient) AddMemberToListViaEmailAddress(emailAddress string, listID uint, muted bool) (*wundergo.Membership, error) {
	if emailAddress == "" {
		return nil, errors.New("emailAddress must not be empty")
	}

	if listID == 0 {
		return nil, errors.New("listID must be > 0")
	}

	url := fmt.Sprintf("%s/memberships", c.apiURL)
	body := []byte(
		fmt.Sprintf(
			`{"email":"%s","list_id":%d,"muted":%t}`,
			emailAddress,
			listID,
			muted,
		),
	)

	req, err := c.newPostRequest(url, body)
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if err != nil {
		c.logger.Println(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusCreated)
	}

	membership := wundergo.Membership{}
	err = json.NewDecoder(resp.Body).Decode(&membership)
	if err != nil {
		c.logger.Println(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	return &membership, nil
}

// AcceptMember updates the provided Membership to reflect the User has
// accepted the Membership request.
func (c oauthClient) AcceptMember(membership wundergo.Membership) (*wundergo.Membership, error) {
	membership.State = "accepted"
	body, err := json.Marshal(membership)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf(
		"%s/memberships/%d",
		c.apiURL,
		membership.ID,
	)

	req, err := c.newPatchRequest(url, body)
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	err = json.NewDecoder(resp.Body).Decode(&membership)
	if err != nil {
		c.logger.Println(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	return &membership, nil
}

// RejectInvite deletes the provided Membership.
func (c oauthClient) RejectInvite(membership wundergo.Membership) error {
	url := fmt.Sprintf(
		"%s/memberships/%d?revision=%d",
		c.apiURL,
		membership.ID,
		membership.Revision,
	)

	req, err := c.newDeleteRequest(url)
	if err != nil {
		return err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusNoContent)
	}

	return nil
}

// RemoveMemberFromList deletes the provided Membership.
func (c oauthClient) RemoveMemberFromList(membership wundergo.Membership) error {
	url := fmt.Sprintf(
		"%s/memberships/%d?revision=%d",
		c.apiURL,
		membership.ID,
		membership.Revision,
	)

	req, err := c.newDeleteRequest(url)
	if err != nil {
		return err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusNoContent)
	}

	return nil
}
