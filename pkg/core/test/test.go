package test

import (
	"testing"

	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/cidverse/cid-sdk-go/mocks"
)

func Setup(t *testing.T) *mocks.SDKClient {
	cidsdk.JoinSeparator = "/"
	sdk := mocks.NewSDKClient(t)
	return sdk
}
