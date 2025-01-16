package asana

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reco-test-task/internal/common"
	"strconv"
)

type AsanaClient interface {
	fetchResources(url string, queryParams map[string]string) (ListResourceResponse, error)
	FetchUsers(limit int, offset string) (ListResourceResponse, error)
	FetchAllUsers() (ListResourceResponse, error)
	FetchProjects(limit int, offset string) (ListResourceResponse, error)
	FetchAllProjects() (ListResourceResponse, error)
}

type _asanaClient struct {
	config     *common.Config
	baseClient *http.Client
}

// Fetch resources from Asana API
func (ac *_asanaClient) fetchResources(
	url string,
	queryParams map[string]string,
) (ListResourceResponse, error) {
	req, err := http.NewRequest(
		"GET",
		url+"?"+formQueryParams(queryParams),
		nil,
	)
	if err != nil {
		return ListResourceResponse{}, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", ac.config.Asana.AccessToken))

	fmt.Printf("Sending request to url: %s\n", req.URL.String())
	resp, err := ac.baseClient.Do(req)
	if err != nil {
		return ListResourceResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		retryAfter, err := strconv.Atoi(resp.Header.Get("Retry-After"))
		if err != nil {
			return ListResourceResponse{}, err
		}

		return ListResourceResponse{}, TooManyRequestsError{
			RetryAfter: retryAfter,
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ListResourceResponse{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return ListResourceResponse{}, fmt.Errorf("users fetching failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result ListResourceResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return ListResourceResponse{}, err
	}

	return result, nil
}

// Fetch users with pagination
func (ac *_asanaClient) FetchUsers(limit int, offset string) (ListResourceResponse, error) {
	return ac.fetchResources(
		common.ASANA_URL+common.USERS_ENDPOINT,
		map[string]string{
			"limit":     strconv.Itoa(limit),
			"offset":    offset,
			"workspace": ac.config.Asana.WorkspaceId,
		},
	)
}

func (ac *_asanaClient) FetchAllUsers() (ListResourceResponse, error) {
	users, err := ac.FetchUsers(common.PAGE_LIMIT, "")
	if err != nil {
		return ListResourceResponse{}, err
	}

	for users.NextPage.Offset != "" {
		newUsers, err := ac.FetchUsers(common.PAGE_LIMIT, users.NextPage.Offset)
		if err != nil {
			return ListResourceResponse{}, err
		}
		users.Data = append(users.Data, newUsers.Data...)
		users.NextPage = newUsers.NextPage
	}

	return users, nil
}

func (ac *_asanaClient) FetchProjects(limit int, offset string) (ListResourceResponse, error) {
	return ac.fetchResources(
		common.ASANA_URL+common.PROJECTS_ENDPOINT,
		map[string]string{
			"limit":     strconv.Itoa(limit),
			"offset":    offset,
			"workspace": ac.config.Asana.WorkspaceId,
		},
	)
}

func (ac *_asanaClient) FetchAllProjects() (ListResourceResponse, error) {
	projects, err := ac.FetchProjects(common.PAGE_LIMIT, "")
	if err != nil {
		return ListResourceResponse{}, err
	}

	for projects.NextPage.Offset != "" {
		newProjects, err := ac.FetchProjects(common.PAGE_LIMIT, projects.NextPage.Offset)
		if err != nil {
			return ListResourceResponse{}, err
		}
		projects.Data = append(projects.Data, newProjects.Data...)
		projects.NextPage = newProjects.NextPage
	}

	return projects, nil
}

func formQueryParams(values map[string]string) string {
	queryParamsString := ""
	first := true

	for key, value := range values {
		if value == "" {
			continue
		}

		if first {
			queryParamsString += fmt.Sprintf("%s=%s", key, value)
			first = false
		} else {
			queryParamsString += fmt.Sprintf("&%s=%s", key, value)
		}
	}

	return queryParamsString
}

func NewClient(config *common.Config) AsanaClient {
	return &_asanaClient{
		config:     config,
		baseClient: &http.Client{},
	}
}
