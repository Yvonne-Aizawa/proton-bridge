// Copyright (c) 2023 Proton AG
//
// This file is part of Proton Mail Bridge.
//
// Proton Mail Bridge is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Proton Mail Bridge is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Proton Mail Bridge.  If not, see <https://www.gnu.org/licenses/>.

package user

import (
	"sync"

	"github.com/ProtonMail/go-proton-api"
)

type SyncDownloadCache struct {
	messageLock    sync.RWMutex
	messages       map[string]proton.Message
	attachmentLock sync.RWMutex
	attachments    map[string][]byte
}

func newSyncDownloadCache() *SyncDownloadCache {
	return &SyncDownloadCache{
		messages:    make(map[string]proton.Message, 64),
		attachments: make(map[string][]byte, 64),
	}
}

func (s *SyncDownloadCache) StoreMessage(message proton.Message) {
	s.messageLock.Lock()
	defer s.messageLock.Unlock()

	s.messages[message.ID] = message
}

func (s *SyncDownloadCache) StoreAttachment(id string, data []byte) {
	s.attachmentLock.Lock()
	defer s.attachmentLock.Unlock()

	s.attachments[id] = data
}

func (s *SyncDownloadCache) DeleteMessages(id ...string) {
	s.messageLock.Lock()
	defer s.messageLock.Unlock()

	for _, id := range id {
		delete(s.messages, id)
	}
}

func (s *SyncDownloadCache) DeleteAttachments(id ...string) {
	s.attachmentLock.Lock()
	defer s.attachmentLock.Unlock()

	for _, id := range id {
		delete(s.attachments, id)
	}
}

func (s *SyncDownloadCache) GetMessage(id string) (proton.Message, bool) {
	s.messageLock.RLock()
	defer s.messageLock.RUnlock()

	v, ok := s.messages[id]

	return v, ok
}

func (s *SyncDownloadCache) GetAttachment(id string) ([]byte, bool) {
	s.attachmentLock.RLock()
	defer s.attachmentLock.RUnlock()

	v, ok := s.attachments[id]

	return v, ok
}

func (s *SyncDownloadCache) Clear() {
	s.messageLock.Lock()
	s.messages = make(map[string]proton.Message, 64)
	s.messageLock.Unlock()

	s.attachmentLock.Lock()
	s.attachments = make(map[string][]byte, 64)
	s.attachmentLock.Unlock()
}
