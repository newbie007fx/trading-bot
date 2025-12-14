package secret

import (
	"context"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

type Loader struct {
	projectID string
	client    *secretmanager.Client
}

func NewLoader(ctx context.Context, projectID string) (*Loader, error) {
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	return &Loader{
		projectID: projectID,
		client:    client,
	}, nil
}

func (l *Loader) Get(ctx context.Context, name string) (string, error) {
	secretName := fmt.Sprintf(
		"projects/%s/locations/asia-southeast2/secrets/%s/versions/latest",
		l.projectID,
		name,
	)

	req := &secretpb.AccessSecretVersionRequest{
		Name: secretName,
	}

	resp, err := l.client.AccessSecretVersion(ctx, req)
	if err != nil {
		return "", err
	}

	return string(resp.Payload.Data), nil
}

func (l *Loader) Close() error {
	return l.client.Close()
}
