package azurestorageendpoints

type CreateContainerRequest struct {
	ContainerName        string `json:"containerName"`
	IsPubliclyAccessible bool   `json:"isPubliclyAccessible"`
}

type CreateContainerResponse struct {
	IsCreated bool  `json:"isCreated"`
	Err       error `json:"-"`
}

func (c CreateContainerResponse) Failed() error { return c.Err }

type DeleteContainerRequest struct {
	ContainerName string `json:"containerName"`
}

type DeleteContainerResponse struct {
	IsCreated bool  `json:"isCreated"`
	Err       error `json:"-"`
}

func (d DeleteContainerResponse) Failed() error { return d.Err }