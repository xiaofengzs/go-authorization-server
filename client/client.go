package client

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"github.com/ory/fosite"
	"log"
	"time"
)

var bCrypt = fosite.BCrypt{
	Config: &fosite.Config{HashCost: 6},
}

type Management struct {
	*Store
}

func NewClientManagement(store *Store) *Management {
	return &Management{
		store,
	}
}

// GetClient ...
func (clientManagement *Management) GetClient(ctx context.Context, id string) (fosite.Client, error) {
	client, err := clientManagement.Get(id)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (clientManagement *Management) ClientAssertionJWTValid(ctx context.Context, jti string) error {
	return nil
}

// SetClientAssertionJWT marks a JTI as known for the given
// expiry time. Before inserting the new JTI, it will clean
// up any existing JTIs that have expired as those tokens can
// not be replayed due to the expiry.
func (clientManagement *Management) SetClientAssertionJWT(ctx context.Context, jti string, exp time.Time) error {
	return nil
}

func (clientManagement *Management) SaveClient(ctx context.Context, id string) error {

	clientSecret, err := bCrypt.Hash(ctx, []byte(generateClientSecret()))
	if err != nil {
		log.Fatalf("hash client sceret failed, err: %s", err)
	}
	client := fosite.DefaultClient{
		ID:             id,
		Secret:         clientSecret,
		RotatedSecrets: [][]byte{clientSecret},
		RedirectURIs:   []string{"http://localhost:3846/callback"},
		ResponseTypes:  []string{"id_token", "code", "token", "id_token token", "code id_token", "code token", "code id_token token"},
		GrantTypes:     []string{"implicit", "refresh_token", "authorization_code", "password", "client_credentials"},
		Scopes:         []string{"fosite", "openid", "photos", "offline"},
	}

	clientManagement.Store.Save(client)

	return nil
}

func generateClientSecret() string {
	const length = 10

	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		panic(err) // Handle errors appropriately
	}

	return base64.URLEncoding.EncodeToString(b)
}
