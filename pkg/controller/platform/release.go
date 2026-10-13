package platform

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

const (
	defaultCoreVersion = "v3.4.0"
	autoVersion        = "auto"
	releaseCheckEvery  = 24 * time.Hour
	releaseRetry       = time.Hour
	coreReleasesURL    = "https://api.github.com/repos/infinimesh/core/releases?per_page=100"
)

var (
	stickyTag = regexp.MustCompile(`^v[0-9]+\.[0-9]+\.[0-9]+(-r[0-9]+|-rc\.[0-9]+)?$`)
	official  = regexp.MustCompile(`^v([0-9]+)\.([0-9]+)\.([0-9]+)(?:-r([0-9]+))?$`)
)

type ghRelease struct {
	TagName    string `json:"tag_name"`
	Draft      bool   `json:"draft"`
	Prerelease bool   `json:"prerelease"`
}

// resolveCoreVersion returns the image tag, whether status should be written, and when to look again.
func resolveCoreVersion(spec, status, checkedAt string, now time.Time, lookup func() (string, error)) (string, bool, time.Duration, error) {
	if spec == "" {
		spec = defaultCoreVersion
	}
	if spec != autoVersion {
		if !stickyTag.MatchString(spec) {
			return "", false, 0, fmt.Errorf("version %q is not vX.Y.Z, vX.Y.Z-rN, vX.Y.Z-rc.N, or auto", spec)
		}
		return spec, status != spec, 0, nil
	}
	checked, _ := time.Parse(time.RFC3339, checkedAt)
	if status != "" && !checked.IsZero() && now.Sub(checked) < releaseCheckEvery {
		return status, false, releaseCheckEvery - now.Sub(checked), nil
	}
	next, err := lookup()
	if err != nil {
		if status != "" {
			return status, false, releaseRetry, nil
		}
		return "", false, 0, err
	}
	if !stickyTag.MatchString(next) || !official.MatchString(next) {
		if status != "" {
			return status, false, releaseRetry, nil
		}
		return "", false, 0, fmt.Errorf("newest core release %q is not an official tag", next)
	}
	return next, true, releaseCheckEvery, nil
}

func newestOfficial(releases []ghRelease) (string, bool) {
	var best string
	var bestN [4]int
	found := false
	for _, rel := range releases {
		if rel.Draft || rel.Prerelease {
			continue
		}
		n, ok := officialNumbers(rel.TagName)
		if !ok {
			continue
		}
		if !found || lessOfficial(bestN, n) {
			best, bestN, found = rel.TagName, n, true
		}
	}
	return best, found
}

func officialNumbers(tag string) ([4]int, bool) {
	m := official.FindStringSubmatch(tag)
	if m == nil {
		return [4]int{}, false
	}
	var n [4]int
	for i := 0; i < 3; i++ {
		v, err := strconv.Atoi(m[i+1])
		if err != nil {
			return [4]int{}, false
		}
		n[i] = v
	}
	if m[4] != "" {
		v, err := strconv.Atoi(m[4])
		if err != nil {
			return [4]int{}, false
		}
		n[3] = v
	}
	return n, true
}

func lessOfficial(a, b [4]int) bool {
	for i := 0; i < 4; i++ {
		if a[i] != b[i] {
			return a[i] < b[i]
		}
	}
	return false
}

func fetchNewestCore(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, coreReleasesURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "infinimesh-platform-operator")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("core releases: %s", resp.Status)
	}
	var releases []ghRelease
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return "", err
	}
	tag, ok := newestOfficial(releases)
	if !ok {
		return "", fmt.Errorf("core releases: no official tag")
	}
	return tag, nil
}
