package service

import "testing"

func TestService_RunningContainerName(t *testing.T) {
	s := Service{
		Host: "example.com",
	}

	if s.RunningContainerName() != "example_com" {
		t.Errorf("Expected running container name to be 'example_com', got '%s'", s.RunningContainerName())
	}
}

func TestService_NextContainerName(t *testing.T) {
	s := Service{
		Host: "example.co.uk",
	}

	if s.NextContainerName() != "next_example_co_uk" {
		t.Errorf("Expected next container name to be 'next_example_co_uk', got '%s'", s.NextContainerName())
	}
}
