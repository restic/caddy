package restic

import (
	"testing"

	"github.com/caddyserver/caddy"
	"github.com/caddyserver/caddy/caddyhttp/httpserver"
)

func TestSetup(t *testing.T) {
	c := caddy.NewTestController("http", `restic /basepath`)
	err := setup(c)
	if err != nil {
		t.Fatalf("Expected no errors, got: %v", err)
	}

	mids := httpserver.GetConfig(c).Middleware()
	if len(mids) == 0 {
		t.Fatal("Expected middleware, had 0 instead")
	}

	handler := mids[0](httpserver.EmptyNext)
	myHandler, ok := handler.(ResticHandler)
	if !ok {
		t.Fatalf("Expected handler to be type ResticHandler, got: %#v", handler)
	}

	if myHandler.BasePath != "/basepath" {
		t.Error("Expected base path to be /basepath")
	}
	if !httpserver.SameNext(myHandler.Next, httpserver.EmptyNext) {
		t.Error("'Next' field of handler was not set properly")
	}

}

func TestExtParse(t *testing.T) {
	tests := []struct {
		inputStr         string
		shouldErr        bool
		expectedBasePath string
	}{
		{`restic`, false, "/"},
		{`restic /basepath`, false, "/basepath"},
		{`restic /basepath /backups`, false, "/basepath"},
		{`restic /basepath /backups extra`, true, "/basepath"},
	}
	for i, test := range tests {
		c := caddy.NewTestController("http", test.inputStr)
		err := setup(c)
		if err == nil && test.shouldErr {
			t.Errorf("Test %d didn't error, but it should have", i)
		} else if err != nil && !test.shouldErr {
			t.Errorf("Test %d errored, but it shouldn't have; got '%v'", i, err)
		}
		if test.shouldErr {
			continue
		}

		mids := httpserver.GetConfig(c).Middleware()
		if len(mids) == 0 {
			t.Fatalf("Test %d: Expected middleware, had 0 instead", i)
		}

		handler := mids[0](httpserver.EmptyNext)
		myHandler, ok := handler.(ResticHandler)
		if !ok {
			t.Fatalf("Expected handler to be type ResticHandler, got: %#v", handler)
		}

		if test.expectedBasePath != myHandler.BasePath {
			t.Errorf("Test %d: Expected base path to be %s but was %s", i, test.expectedBasePath, myHandler.BasePath)
		}
	}

}
