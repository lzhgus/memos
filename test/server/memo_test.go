package testserver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/usememos/memos/api"
)

func TestMemoServer(t *testing.T) {
	ctx := context.Background()
	s, err := NewTestingServer(ctx, t)
	require.NoError(t, err)
	defer s.Shutdown(ctx)

	signup := &api.SignUp{
		Username: "testuser",
		Password: "testpassword",
	}
	user, err := s.postAuthSignup(signup)
	require.NoError(t, err)
	require.Equal(t, signup.Username, user.Username)
	memoList, err := s.getMemoList()
	require.NoError(t, err)
	require.Len(t, memoList, 0)
	memo, err := s.postMemoCreate(&api.MemoCreate{
		Content: "test memo",
	})
	require.NoError(t, err)
	require.Equal(t, "test memo", memo.Content)
	memoList, err = s.getMemoList()
	require.NoError(t, err)
	require.Len(t, memoList, 1)
	updatedContent := "updated memo"
	memo, err = s.patchMemo(&api.MemoPatch{
		ID:      memo.ID,
		Content: &updatedContent,
	})
	require.NoError(t, err)
	require.Equal(t, updatedContent, memo.Content)
	err = s.deleteMemo(&api.MemoDelete{
		ID: memo.ID,
	})
	require.NoError(t, err)
	memoList, err = s.getMemoList()
	require.NoError(t, err)
	require.Len(t, memoList, 0)
}

func (s *TestingServer) getMemo(memoID int) (*api.Memo, error) {
	body, err := s.get(fmt.Sprintf("/api/memo/%d", memoID), nil)
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	_, err = buf.ReadFrom(body)
	if err != nil {
		return nil, errors.Wrap(err, "fail to read response body")
	}

	type MemoCreateResponse struct {
		Data *api.Memo `json:"data"`
	}
	res := new(MemoCreateResponse)
	if err = json.Unmarshal(buf.Bytes(), res); err != nil {
		return nil, errors.Wrap(err, "fail to unmarshal get memo response")
	}
	return res.Data, nil
}

func (s *TestingServer) getMemoList() ([]*api.Memo, error) {
	body, err := s.get("/api/memo", nil)
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	_, err = buf.ReadFrom(body)
	if err != nil {
		return nil, errors.Wrap(err, "fail to read response body")
	}

	type MemoCreateResponse struct {
		Data []*api.Memo `json:"data"`
	}
	res := new(MemoCreateResponse)
	if err = json.Unmarshal(buf.Bytes(), res); err != nil {
		return nil, errors.Wrap(err, "fail to unmarshal get memo list response")
	}
	return res.Data, nil
}

func (s *TestingServer) postMemoCreate(memoCreate *api.MemoCreate) (*api.Memo, error) {
	rawData, err := json.Marshal(&memoCreate)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal memo create")
	}
	reader := bytes.NewReader(rawData)
	body, err := s.post("/api/memo", reader, nil)
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	_, err = buf.ReadFrom(body)
	if err != nil {
		return nil, errors.Wrap(err, "fail to read response body")
	}

	type MemoCreateResponse struct {
		Data *api.Memo `json:"data"`
	}
	res := new(MemoCreateResponse)
	if err = json.Unmarshal(buf.Bytes(), res); err != nil {
		return nil, errors.Wrap(err, "fail to unmarshal post memo create response")
	}
	return res.Data, nil
}

func (s *TestingServer) patchMemo(memoPatch *api.MemoPatch) (*api.Memo, error) {
	rawData, err := json.Marshal(&memoPatch)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal memo patch")
	}
	reader := bytes.NewReader(rawData)
	body, err := s.patch(fmt.Sprintf("/api/memo/%d", memoPatch.ID), reader, nil)
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	_, err = buf.ReadFrom(body)
	if err != nil {
		return nil, errors.Wrap(err, "fail to read response body")
	}

	type MemoPatchResponse struct {
		Data *api.Memo `json:"data"`
	}
	res := new(MemoPatchResponse)
	if err = json.Unmarshal(buf.Bytes(), res); err != nil {
		return nil, errors.Wrap(err, "fail to unmarshal patch memo response")
	}
	return res.Data, nil
}

func (s *TestingServer) deleteMemo(memoDelete *api.MemoDelete) error {
	_, err := s.delete(fmt.Sprintf("/api/memo/%d", memoDelete.ID), nil)
	return err
}
