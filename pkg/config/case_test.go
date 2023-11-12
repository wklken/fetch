package config_test

import (
	"testing"

	"github.com/wklken/fetch/pkg/config"
)

func TestCaseID(t *testing.T) {
	c := config.Case{
		Title: "Test Case",
		Path:  "/path/to/test",
		Index: 1,
	}

	want := "/path/to/test | Test Case"
	if got := c.ID(); got != want {
		t.Errorf("Case.ID() = %v, want %v", got, want)
	}

	c.Index = 2
	want = "/path/to/test[2] | Test Case"
	if got := c.ID(); got != want {
		t.Errorf("Case.ID() = %v, want %v", got, want)
	}

	c.Title = ""
	want = "/path/to/test[2] | -"
	if got := c.ID(); got != want {
		t.Errorf("Case.ID() = %v, want %v", got, want)
	}
}

func TestCaseRender(t *testing.T) {
	c := config.Case{
		Title: "Test Case",
		Request: config.Request{
			Method: "GET",
			URL:    "http://example.com",
			Body:   `{"name": "{{.name}}"}`,
		},
		Assert: config.Assert{
			StatusCode: 200,
			Body:       `{"name": "{{.name}}"}`,
		},
	}

	ctx := map[string]interface{}{
		"name": "John",
	}

	c.Render(ctx)

	want := `{"name": "John"}`
	if got := c.Request.Body; got != want {
		t.Errorf("Case.Render() Request.Body = %v, want %v", got, want)
	}

	want = `{"name": "John"}`
	if got := c.Assert.Body; got != want {
		t.Errorf("Case.Render() Assert.Body = %v, want %v", got, want)
	}
}

func TestCaseGuessAssertLineNumber(t *testing.T) {
	c := config.Case{
		FileLines: map[int]map[int]string{
			1: {
				1: `{"name": "John"}`,
				2: `{"age": 30}`,
			},
			2: {
				1: `{"name": "Jane"}`,
				2: `{"age": 25}`,
			},
		},
	}

	want := 1
	if got := c.GuessAssertLineNumber(1, "name"); got != want {
		t.Errorf("Case.GuessAssertLineNumber() = %v, want %v", got, want)
	}

	want = 2
	if got := c.GuessAssertLineNumber(1, "age"); got != want {
		t.Errorf("Case.GuessAssertLineNumber() = %v, want %v", got, want)
	}

	want = 1
	if got := c.GuessAssertLineNumber(2, "name"); got != want {
		t.Errorf("Case.GuessAssertLineNumber() = %v, want %v", got, want)
	}

	want = 2
	if got := c.GuessAssertLineNumber(2, "age"); got != want {
		t.Errorf("Case.GuessAssertLineNumber() = %v, want %v", got, want)
	}

	// want = -1
	// if got := c.GuessAssertLineNumber(3, "name"); got != want {
	// 	t.Errorf("Case.GuessAssertLineNumber() = %v, want %v", got, want)
	// }
}
