package deploy_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/a-novel/golib/deploy"
	"github.com/a-novel/golib/testutils"
)

func TestUnsetEnv(t *testing.T) {
	testutils.RunCMD(t, &testutils.CMDConfig{
		CmdFn: func(t *testing.T) {
			require.Equal(t, deploy.ENV, deploy.DevENV)
			require.False(t, deploy.IsReleaseEnv())
		},
		MainFn: func(t *testing.T, res *testutils.CMDResult) {
			require.True(t, res.Success)
		},
	})
}

func TestDevEnv(t *testing.T) {
	testutils.RunCMD(t, &testutils.CMDConfig{
		CmdFn: func(t *testing.T) {
			require.Equal(t, deploy.ENV, deploy.DevENV)
			require.False(t, deploy.IsReleaseEnv())
		},
		MainFn: func(t *testing.T, res *testutils.CMDResult) {
			require.True(t, res.Success)
		},
		Env: []string{"ENV=dev"},
	})
}

func TestStagingEnv(t *testing.T) {
	testutils.RunCMD(t, &testutils.CMDConfig{
		CmdFn: func(t *testing.T) {
			require.Equal(t, deploy.ENV, deploy.StagingEnv)
			require.True(t, deploy.IsReleaseEnv())
		},
		MainFn: func(t *testing.T, res *testutils.CMDResult) {
			require.True(t, res.Success)
		},
		Env: []string{"ENV=staging"},
	})
}

func TestProdEnv(t *testing.T) {
	testutils.RunCMD(t, &testutils.CMDConfig{
		CmdFn: func(t *testing.T) {
			require.Equal(t, deploy.ENV, deploy.ProdENV)
			require.True(t, deploy.IsReleaseEnv())
		},
		MainFn: func(t *testing.T, res *testutils.CMDResult) {
			require.True(t, res.Success)
		},
		Env: []string{"ENV=prod"},
	})
}

func TestInvalidEnv(t *testing.T) {
	testutils.RunCMD(t, &testutils.CMDConfig{
		CmdFn: func(t *testing.T) {
			fmt.Println(deploy.ENV)
		},
		MainFn: func(t *testing.T, res *testutils.CMDResult) {
			require.False(t, res.Success)
		},
		Env: []string{"ENV=foo"},
	})
}
