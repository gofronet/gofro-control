package bootstrap

import (
	"encoding/json"
	"errors"
	"gofronet-foundation/gofro-control/internal/constants"
	"gofronet-foundation/gofro-control/internal/nodes/bootstrap/models"
	"os"
	"path/filepath"
	"sync"
)

const (
	invitesPath = constants.AppDataDir + "/invites.json"
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

func (s *InviteStore) DoneInvite(inviteID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	file, err := readInvitesFile()
	if err != nil {
		return err
	}

	var invites map[string]*models.InviteRecord
	if err := json.Unmarshal(file, &invites); err != nil {
		return err
	}

	_, ok := invites[inviteID]
	if !ok {
		return errors.New("invite not found")
	}

	delete(invites, inviteID)

	serializedInvites, err := json.Marshal(invites)
	if err != nil {
		return err
	}

	if err := writeInvitesFile(serializedInvites); err != nil {
		return err
	}

	return nil
}

func (s *InviteStore) GetInvite(inviteID string) (*models.InviteRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	file, err := readInvitesFile()
	if err != nil {
		return nil, err
	}

	var invites map[string]*models.InviteRecord
	if err := json.Unmarshal(file, &invites); err != nil {
		return nil, err
	}

	invite, ok := invites[inviteID]
	if !ok {
		return nil, errors.New("invite not found")
	}

	return invite, nil

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
	return os.WriteFile(invitesPath, content, 0o600)
}

func readInvitesFile() ([]byte, error) {

	invitesFile, err := os.ReadFile(invitesPath)
	if err != nil {
		return nil, err
	}
	return invitesFile, nil
}
