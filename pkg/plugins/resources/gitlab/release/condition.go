package release

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/updatecli/updatecli/pkg/core/pipeline/scm"
	"github.com/updatecli/updatecli/pkg/core/result"
)

func (g *Gitlab) Condition(source string) (bool, error) {
	releases, err := g.SearchReleases()

	if len(g.spec.Tag) == 0 {
		g.spec.Tag = source
	}

	if err != nil {
		logrus.Error(err)
		return false, err
	}

	if len(releases) == 0 {
		logrus.Infof("%s No Gitlab release found. As a fallback you may be looking for git tags", result.ATTENTION)
		return false, nil
	}

	for _, release := range releases {
		if release == g.spec.Tag {
			logrus.Infof("%s Gitlab Release tag %q found", result.SUCCESS, release)
			return true, nil
		}
	}

	logrus.Infof("%s No Gitlab Release tag found matching pattern %q", result.FAILURE, g.versionFilter.Pattern)
	return false, nil

}

func (g *Gitlab) ConditionFromSCM(source string, scm scm.ScmHandler) (bool, error) {
	return false, fmt.Errorf("Condition not supported for the plugin Gitlab Release")
}
