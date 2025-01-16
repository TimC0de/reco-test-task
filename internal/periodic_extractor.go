package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"reco-test-task/internal/asana"
	"reco-test-task/internal/common"
	"time"
)

type PeriodicExtractor interface {
	Start(period time.Duration)
}

type _periodicExtractor struct {
	config   *common.Config
	client   asana.AsanaClient
	uploader Uploader
}

func (pe *_periodicExtractor) Start(period time.Duration) {
	for {
		// Fetch users and upload them to folder
		users, err := pe.client.FetchAllUsers()

		// Check for TooManyRequests error
		var tooManyRequestsError asana.TooManyRequestsError
		if errors.As(err, &tooManyRequestsError) {
			time.Sleep(time.Second * time.Duration(tooManyRequestsError.RetryAfter))
			continue
		}

		// Process other errors
		if err != nil {
			fmt.Printf("Error fetching users: %s\n", err.Error())
		}
		fmt.Printf("New users data: %v\n", users.Data)

		err = pe.uploadResources(users.Data)
		if err != nil {
			fmt.Println(err.Error())
		}

		// Fetch projects and upload them to folder
		projects, err := pe.client.FetchAllProjects()

		// Check for TooManyRequests error
		if errors.As(err, &tooManyRequestsError) {
			time.Sleep(time.Second * time.Duration(tooManyRequestsError.RetryAfter))
			continue
		}

		// Process other errors
		if err != nil {
			fmt.Printf("Error fetching projects: %s\n", err.Error())
		}
		fmt.Printf("New projects data: %v\n", projects.Data)

		err = pe.uploadResources(projects.Data)
		if err != nil {
			fmt.Println(err.Error())
		}

		// Sleep for specified period
		time.Sleep(period)
	}
}

func (pe *_periodicExtractor) uploadResources(resources []asana.ResourceResponse) error {
	for _, resource := range resources {
		userSerialized, err := json.Marshal(resource)
		if err != nil {
			return fmt.Errorf("Error serializing resource: %s\n", err.Error())
		}

		fileName := fmt.Sprintf("%s_%s", resource.ResourceType, resource.GID)
		err = pe.uploader.UploadNewFile(
			fileName,
			userSerialized,
		)
		if err != nil {
			return fmt.Errorf("Error uploading resource: %s\n", err.Error())
		}
	}

	return nil
}

func NewPeriodicExtractor(config *common.Config) PeriodicExtractor {
	return &_periodicExtractor{
		config:   config,
		client:   asana.NewClient(config),
		uploader: NewUploader(),
	}
}
