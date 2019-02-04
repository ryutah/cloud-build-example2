package function

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"golang.org/x/oauth2"

	"github.com/golang/glog"
	"github.com/google/go-github/v22/github"
)

var (
	authToken = os.Getenv("AUTH_TOKEN")
)

type buildStatus string

const (
	buildStatusSuccess       buildStatus = "SUCCESS"
	buildStatusFailure       buildStatus = "FAILURE"
	buildStatusInternalError buildStatus = "INTERNAL_ERROR"
	buildStatusTimeout       buildStatus = "TIMEOUT"
)

type githubMessages struct {
	status      string
	description string
}

var buildStatusGithubMsgMap = map[buildStatus]githubMessages{
	buildStatusSuccess: {
		status:      "success",
		description: "Your tests passed on Cloud Build",
	},
	buildStatusFailure: {
		status:      "failure",
		description: "Your tests failed on Cloud Build",
	},
	buildStatusInternalError: {
		status:      "success",
		description: "Your test failed caused by Cloud Build internal error",
	},
	buildStatusTimeout: {
		status:      "success",
		description: "Your test timeouted on Cloud Build",
	},
}

func (b buildStatus) githubStatus() string {
	s, ok := buildStatusGithubMsgMap[b]
	if !ok {
		return "pending"
	}
	return s.status
}

func (b buildStatus) githubDescription() string {
	s, ok := buildStatusGithubMsgMap[b]
	if !ok {
		return "Your test proceeded on Cloud Build"
	}
	return s.description
}

type PubSubMessage struct {
	Data []byte `json:"data"`
}

type (
	cloudBuildMessage struct {
		ID               string                     `json:"id"`
		ProjectID        string                     `json:"projectId"`
		Status           buildStatus                `json:"status"`
		SourceProvenance cloudBuildSourceProvenance `json:"sourceProvenance"`
	}

	cloudBuildSourceProvenance struct {
		ResolvedRepoSource cloudBuildSourceProvenanceResolvedRepoSource `json:"resolvedRepoSource"`
	}

	cloudBuildSourceProvenanceResolvedRepoSource struct {
		ProjectID string `json:"projectId"`
		RepoName  string `json:"repoName"`
		CommitSHA string `json:"commitSha"`
	}
)

func HelloPubSub(ctx context.Context, m PubSubMessage) error {
	if !flag.Parsed() {
		flag.Set("stderrthreshold", "INFO")
		flag.Parse()
	}

	if len(m.Data) == 0 {
		glog.Info("finish function")
		return nil
	}
	gcbData := new(cloudBuildMessage)
	if err := json.Unmarshal(m.Data, gcbData); err != nil {
		glog.Warningf("could not parse data: %v", err)
		return nil
	}
	glog.Infof("%#v", *gcbData)

	ts := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: authToken,
	})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	parts := strings.SplitN(gcbData.SourceProvenance.ResolvedRepoSource.RepoName, "_", 3)
	if len(parts) != 3 || parts[0] != "github" {
		glog.Infof("%s is not github resource", gcbData.SourceProvenance.ResolvedRepoSource.RepoName)
		return nil
	}
	if _, _, err := client.Repositories.CreateStatus(ctx, parts[1], parts[2], gcbData.SourceProvenance.ResolvedRepoSource.CommitSHA, &github.RepoStatus{
		State: strPtr(gcbData.Status.githubStatus()),
		TargetURL: strPtr(fmt.Sprintf(
			"https://console.cloud.google.com/gcr/builds/%s?project=%s",
			gcbData.ID, gcbData.ProjectID,
		)),
		Description: strPtr(gcbData.Status.githubDescription()),
		Context:     strPtr("ci/cloud-build"),
	}); err != nil {
		glog.Errorf("failed to create stats: %v", err)
		return err
	}
	return nil
}

func strPtr(s string) *string {
	return &s
}
