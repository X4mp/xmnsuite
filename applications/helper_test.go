package applications

import (
	"testing"
)

func TestFromURLPatternToRegex_withOneVariable_Success(t *testing.T) {
	//variables:
	urlPattern := "/videos/<id|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>"
	url := "/videos/70de0f1a-0623-4bf6-ac6c-384f56321ec0"

	//execute:
	pattern, variableNames, err := fromURLPatternToRegex(urlPattern)
	if err != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", err.Error())
		return
	}

	output := pattern.FindStringSubmatch(url)
	if len(output) != 2 {
		t.Errorf("there was expected to be 2 elements in the output, %d returned", len(output))
		return
	}

	if output[0] != url {
		t.Errorf("the first element in the output should be the url.  Expected: %s, Returned: %s", url, output[0])
		return
	}

	if output[1] != "70de0f1a-0623-4bf6-ac6c-384f56321ec0" {
		t.Errorf("the second element in the output should be the variable value.  Expected: %s, Returned: %s", "70de0f1a-0623-4bf6-ac6c-384f56321ec0", output[1])
		return
	}

	if len(variableNames) != 1 {
		t.Errorf("there should be 1 variable name.  Returned: %d", len(variableNames))
		return
	}

	if variableNames[0] != "id" {
		t.Errorf("the first element in the variableNames should be the variable name.  Expected: %s, Returned: %s", "id", variableNames[0])
		return
	}
}

func TestFromURLPatternToRegex_withTwoVariable_Success(t *testing.T) {
	//variables:
	urlPattern := "/this/is/a/<slug|[a-zA-Z0-9-]+>/oh-yes/<another|[a-z0-9-]+>"
	url := "/this/is/a/my-slug/oh-yes/just-a-name-09"

	//execute:
	pattern, variableNames, err := fromURLPatternToRegex(urlPattern)
	if err != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", err.Error())
		return
	}

	output := pattern.FindStringSubmatch(url)
	if len(output) != 3 {
		t.Errorf("there was expected to be 3 elements in the output, %d returned", len(output))
		return
	}

	if output[0] != url {
		t.Errorf("the first element in the output should be the url.  Expected: %s, Returned: %s", url, output[0])
		return
	}

	if output[1] != "my-slug" {
		t.Errorf("the second element in the output should be the variable value.  Expected: %s, Returned: %s", "my-slug", output[1])
		return
	}

	if output[2] != "just-a-name-09" {
		t.Errorf("the second element in the output should be the variable value.  Expected: %s, Returned: %s", "just-a-name-09", output[1])
		return
	}

	if len(variableNames) != 2 {
		t.Errorf("there should be 2 variable name.  Returned: %d", len(variableNames))
		return
	}

	if variableNames[0] != "slug" {
		t.Errorf("the first element in the variableNames should be the variable name.  Expected: %s, Returned: %s", "slug", variableNames[0])
		return
	}

	if variableNames[1] != "another" {
		t.Errorf("the second element in the variableNames should be the variable name.  Expected: %s, Returned: %s", "another", variableNames[1])
		return
	}
}
