package asana

type NextPageResponse struct {
	Offset string `json:"offset"`
	Path   string `json:"path"`
	URI    string `json:"uri"`
}

type ResourceResponse struct {
	GID          string `json:"gid"`
	ResourceType string `json:"resource_type"`
	Name         string `json:"name"`
}

type ListResourceResponse struct {
	Data     []ResourceResponse `json:"data"`
	NextPage NextPageResponse   `json:"next_page"`
}
