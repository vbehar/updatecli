package branch

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCondition(t *testing.T) {

	tests := []struct {
		name     string
		manifest struct {
			URL        string
			Token      string
			Owner      string
			Repository string
			Branch     string
		}
		wantResult bool
		wantErr    bool
	}{
		{
			name: "repository olblak/updatecli-donotexist should not exist",
			manifest: struct {
				URL        string
				Token      string
				Owner      string
				Repository string
				Branch     string
			}{
				URL:        "gitlab.com",
				Token:      "",
				Owner:      "olblak",
				Repository: "updatecli-donotexist",
				Branch:     "v2",
			},
			wantResult: false,
			wantErr:    true,
		},
		{
			name: "repository olblak/updatecli-mirror should exist with branches",
			manifest: struct {
				URL        string
				Token      string
				Owner      string
				Repository string
				Branch     string
			}{
				URL:        "gitlab.com",
				Token:      "",
				Owner:      "olblak",
				Repository: "updatecli",
				Branch:     "main",
			},
			wantResult: true,
			wantErr:    false,
		},
		{
			name: "repository olblak/updatecli should not have branch donotexist",
			manifest: struct {
				URL        string
				Token      string
				Owner      string
				Repository string
				Branch     string
			}{
				URL:        "gitlab.com",
				Token:      "",
				Owner:      "olblak",
				Repository: "updatecli",
				Branch:     "donotexist",
			},
			wantResult: false,
			wantErr:    false,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			g, gotErr := New(tt.manifest)
			require.NoError(t, gotErr)

			gotResult, gotErr := g.Condition("")

			if tt.wantErr {
				require.Error(t, gotErr)
			} else {
				require.NoError(t, gotErr)
			}

			assert.Equal(t, tt.wantResult, gotResult)

		})

	}
}
