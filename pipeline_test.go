package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

// func TestMain(m *testing.M) {
// 	log.SetOutput(ioutil.Discard)
// 	os.Exit(m.Run())
// }

func TestUploadPipelineCallsBuildkiteAgentCommand(t *testing.T) {
	plugin := Plugin{Diff: "echo ./foo-service"}
	p := PipelineUploader{generatePipeline: generatePipeline}
	p.uploadPipeline(plugin)

	// m.AssertCalled(t, mock.Anything)
}

func TestDiff(t *testing.T) {
	want := []string{
		"services/foo/serverless.yml",
		"services/bar/config.yml",
		"ops/bar/config.yml",
		"README.md",
	}

	got := diff("cat ./tests/mocks/diff1")

	assert.Equal(t, want, got)
}

func TestPipelinesToTriggerGetsListOfPipelines(t *testing.T) {
	want := []string{"service-1", "service-2", "service-4"}

	watch := []WatchConfig{
		{
			Paths: []string{"watch-path-1"},
			Step:  Step{Trigger: "service-1"},
		},
		{
			Paths: []string{"watch-path-2/", "watch-path-3/", "watch-path-4"},
			Step:  Step{Trigger: "service-2"},
		},
		{
			Paths: []string{"watch-path-5"},
			Step:  Step{Trigger: "service-3"},
		},
		{
			Paths: []string{"watch-path-2"},
			Step:  Step{Trigger: "service-4"},
		},
	}

	changedFiles := []string{
		"watch-path-1/text.txt",
		"watch-path-2/.gitignore",
		"watch-path-2/src/index.go",
		"watch-path-4/test/index_test.go",
	}

	pipelines := stepsToTrigger(changedFiles, watch)
	var got []string

	for _, v := range pipelines {
		got = append(got, v.Trigger)
	}

	assert.Equal(t, want, got)
}

func TestGeneratePipeline(t *testing.T) {
	steps := []Step{
		{
			Trigger: "foo-service-pipeline",
			Build:   Build{Message: "build message"},
		},
	}

	want := Pipeline{Steps: steps}
	got := Pipeline{}

	pipeline, err := generatePipeline(steps)
	defer os.Remove(pipeline.Name())

	if err != nil {
		// Failed to close the temp pipeline file
		assert.Equal(t, true, false)
	}

	file, _ := ioutil.ReadFile(pipeline.Name())

	if err = yaml.Unmarshal(file, &got); err != nil {
		// Failed to unmarshal temporary pipeline file
		assert.Equal(t, true, false)
	}

	assert.Equal(t, want, got)
}

// func TestUploadPipeline(t *testing.T) {
// 	want := "uploading pipelines"
// 	got := uploadPipeline()

// 	if want != got {
// 		t.Errorf(`uploadPipeline(), got %q, want "%v"`, got, want)
// 	}
// }
