package main

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"sync"
	"testing"
)

func TestProbes(t *testing.T) {
	attempts := 50
	url_template := fmt.Sprintf("http://%s/probe", config.ApiAddr()) + "/%s"
	probes := []string{
		"psql", "kafka", "redis",
	}

	var wg sync.WaitGroup

	for _, probe := range probes {
		for i := range attempts {
			wg.Add(1)
			go func() {
				defer wg.Done()

				response, err := http.Get(fmt.Sprintf(url_template, probe))
				if err != nil {
					t.Errorf("probe \"%s\" - %d fail with err: %s", probe, i+1, err.Error())
				}
				defer response.Body.Close()

				responseContent, err := io.ReadAll(response.Body)
				if err != nil {
					t.Errorf("cannot read from response: %s", err.Error())
				}
				logContentSnapshot := string(responseContent[:int(math.Min(float64(len(responseContent)), 100))])
				if response.StatusCode != http.StatusOK {
					t.Errorf("probe \"%s\" - %d fail due non success status code: %d. Content: %s", probe, i+1, response.StatusCode, logContentSnapshot)
				}

				fmt.Printf("Successfully complete probe \"%s\" - %d with content: %s\n", probe, i+1, logContentSnapshot)
			}()
		}

		wg.Wait()
	}
}

func TestUpdateProfilingInfo(t *testing.T) {
	if !config.PROF_ENABLE {
		return
	}

	targetFolder := "./.prof" // using dot only because running from tmp forlder but real cwd from exec point
	profiles := []string{
		"goroutine",    // stack traces of all current goroutines
		"heap",         // a sampling of memory allocations of live objects
		"allocs",       // a sampling of all past memory allocations
		"threadcreate", // stack traces that led to the creation of new OS threads
		"block",        // stack traces that led to blocking on synchronization primitives
		"mutex",        // stack traces of holders of contended mutexes
	}

	if err := os.Mkdir(targetFolder, os.ModePerm); err != nil && !os.IsExist(err) {
		t.Error(err)
		return
	}

	for _, profile := range profiles {
		response, err := http.Get(fmt.Sprintf("http://%s/debug/pprof/%s", config.ProfAddr(), profile))
		if err != nil {
			t.Error(err)
			return
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			t.Errorf("invalid response status code for \"%s\" profile: %d", profile, response.StatusCode)
			return
		}
		responseContent, err := io.ReadAll(response.Body)
		if err != nil {
			t.Error(err)
			return
		}
		targetProfilePath := "./" + targetFolder + "/" + profile + ".prof"
		if err := os.WriteFile(targetProfilePath, responseContent, os.ModePerm); err != nil {
			t.Error(err)
			return
		}
	}
}
