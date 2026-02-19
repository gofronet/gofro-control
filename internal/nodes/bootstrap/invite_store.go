package bootstrap

import (
	"encoding/json"
	"gofronet-foundation/gofro-control/internal"
	"gofronet-foundation/gofro-control/internal/nodes/bootstrap/models"
	"os"
	"path/filepath"
	"sync"
)

const (
	invitesPath = internal.AppDataDir + "/invites.json"
)

type InviteStore struct {
	mu sync.Mutex
}

func NewInviteStore() *InviteStore {

	if err := os.MkdirAll(filepath.Dir(invitesPath), 0o700); err != nil {
		panic(err)
	}

	if _, err := os.Stat(invitesPath); err != nil {
		if os.IsNotExist(err) {
			if err := os.WriteFile(invitesPath, []byte("{}"), 0o600); err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	return &InviteStore{}
}

func (s *InviteStore) AddInvite(invite *models.InviteRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	invitesFile, err := readInvitesFile()
	if err != nil {
		return err
	}

	var invites map[string]*models.InviteRecord
	if err := json.Unmarshal(invitesFile, &invites); err != nil {
		return err
	}

	invites[invite.InviteID] = invite

	serializedInvites, err := json.Marshal(invites)
	if err != nil {
		return err
	}

	if err := writeInvitesFile(serializedInvites); err != nil {
		return err
	}

	return nil
}

func writeInvitesFile(content []byte) error {
	return os.WriteFile(invitesPath, content, 0o644)
}

func readInvitesFile() ([]byte, error) {
	invitesFile, err := os.ReadFile(invitesPath)
	if err != nil {
		return nil, err
	}
	return invitesFile, nil
}
