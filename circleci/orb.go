package circleci

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"gopkg.in/yaml.v3"
)

// Orb is a parsed orb reference e.g. "circleci/node@5.1.0".
type Orb struct {
	Name    string
	Version string
}

func (o *Orb) String() string {
	return o.Name + "@" + o.Version
}

func (o *Orb) GoString() string {
	return "&Uses{" + o.String() + "}"
}

func (o *Orb) MarshalJSON() ([]byte, error) {
	return []byte("\"" + o.String() + "\""), nil
}

func Parse(orb string) (*Orb, error) {
	parts := strings.Split(orb, "@")
	if len(parts) != 2 {
		return nil, fmt.Errorf("parse orb: %s", orb)
	}

	return &Orb{
		Name:    parts[0],
		Version: parts[1],
	}, nil
}

func GetOrbSource(ctx context.Context, o *Orb) (*Source, error) {
	var (
		body    = new(bytes.Buffer)
		details = &struct {
			Data *struct {
				OrbVersion *struct {
					Source string `json:"source"`
				} `json:"orbVersion"`
			}
		}{
			Data: &struct {
				OrbVersion *struct {
					Source string `json:"source"`
				} `json:"orbVersion"`
			}{
				OrbVersion: &struct {
					Source string `json:"source"`
				}{},
			},
		}
		source = &Source{}
	)

	if err := json.NewEncoder(body).Encode(map[string]any{
		"operationName": "OrbDetailsQuery",
		"variables": map[string]string{
			"name":          o.Name,
			"orbVersionRef": o.String(),
		},
		"query": `query OrbDetailsQuery($name: String, $orbVersionRef: String) {
			orbVersion(orbVersionRef: $orbVersionRef) {
				source
			}
		}
		`,
	}); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://circleci.com/graphql-unstable", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(details); err != nil {
		return nil, err
	}

	return source, yaml.NewDecoder(bytes.NewReader([]byte(details.Data.OrbVersion.Source))).Decode(source)
}
