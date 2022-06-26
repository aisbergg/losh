package models

import (
	"errors"
	"fmt"
	gourl "net/url"
	"strings"
)

var ErrUnsupportedPlatform = errors.New("unsupported platform")
var ErrInvalidURL = errors.New("invalid URL")

// ProductURL is a URL to a file on a platform.
type ProductURL struct {
	Domain string
	Owner  string
	Repo   string
	Tag    string
	Path   string
}

// NewProductURL creates a new PlatformURL from parts.
func NewProductURL(platform, owner, repo, tag, path string) *ProductURL {
	return &ProductURL{
		Domain: platform,
		Owner:  owner,
		Repo:   repo,
		Tag:    tag,
		Path:   path,
	}
}

// NewProductURLFromURL creates a new PlatformURL from a URL.
func NewProductURLFromURL(url string) (*ProductURL, error) {
	parsedURL, err := gourl.ParseRequestURI(url)
	if err != nil {
		return nil, err
	}

	domain := strings.ToLower(parsedURL.Hostname())
	pathParts := strings.Split(strings.TrimLeft(parsedURL.Path, "/"), "/")
	productURL := &ProductURL{}

	switch domain {
	case "github.com", "raw.githubusercontent.com":
		productURL.Domain = "github.com"
		if domain == "github.com" && len(pathParts) > 4 {
			if pathParts[2] == "blob" || pathParts[2] == "tree" || pathParts[2] == "commit" {
				productURL.Tag = pathParts[3]
				productURL.Path = strings.Join(pathParts[4:], "/")
			} else {
				return nil, ErrInvalidURL
			}
		} else if domain == "raw.githubusercontent.com" {
			productURL.Tag = pathParts[2]
			productURL.Path = strings.Join(pathParts[3:], "/")
		} else {
			return nil, ErrInvalidURL
		}
		productURL.Owner = pathParts[0]
		productURL.Repo = pathParts[1]

	case "gitlab.com":
		productURL.Domain = "gitlab.com"
		// if len(path_parts) >= 5 and path_parts[2] == "-" and path_parts[3] in ["tree", "blob", "raw"]:
		if len(pathParts) > 5 && pathParts[2] == "-" && (pathParts[3] == "tree" || pathParts[3] == "blob" || pathParts[3] == "raw") {
			productURL.Tag = pathParts[4]
			productURL.Path = strings.Join(pathParts[5:], "/")
		} else {
			return nil, ErrInvalidURL
		}
		productURL.Owner = pathParts[0]
		productURL.Repo = pathParts[1]

	case "wikifactory.com":
		productURL.Domain = "wikifactory.com"
		if len(pathParts) >= 2 {
			productURL.Owner = pathParts[0]
			productURL.Repo = pathParts[1]
		}
		if len(pathParts) >= 4 && (pathParts[2] == "file" || pathParts[2] == "files") {
			productURL.Path = strings.Join(pathParts[3:], "/")
		} else if len(pathParts) >= 4 && pathParts[2] == "v" {
			productURL.Tag = pathParts[3]
			if len(pathParts) >= 6 && (pathParts[4] == "file" || pathParts[4] == "files") {
				productURL.Path = strings.Join(pathParts[5:], "/")
			}
		}

	case "certification.oshwa.org", "oshwa.org":
		productURL.Domain = "oshwa.org"
		if len(pathParts) != 1 {
			return nil, ErrInvalidURL
		}
		productURL.Path = pathParts[0]

	default:
		return nil, ErrUnsupportedPlatform

	}

	return productURL, nil
}

// RepositoryURL returns the URL to the repository.
func (pu *ProductURL) RepositoryURL() string {
	switch pu.Domain {

	// format: https://github.com/{owner}/{repo}
	case "github.com":
		return (&gourl.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   fmt.Sprintf("/%s/%s", pu.Owner, pu.Repo),
		}).String()

	// format: https://gitlab.com/{owner}/{repo}
	case "gitlab.com":
		return (&gourl.URL{
			Scheme: "https",
			Host:   "gitlab.com",
			Path:   fmt.Sprintf("/%s/%s", pu.Owner, pu.Repo),
		}).String()

	// format: https://wikifactory.com/{owner}/{repo}
	case "wikifactory.com":
		return (&gourl.URL{
			Scheme: "https",
			Host:   "wikifactory.com",
			Path:   fmt.Sprintf("/%s/%s", pu.Owner, pu.Repo),
		}).String()

	// format: https://certification.oshwa.org/{path}
	case "oshwa.org":
		return (&gourl.URL{
			Scheme: "https",
			Host:   "oshwa.org",
			Path:   fmt.Sprintf("/%s", pu.Path),
		}).String()

	// format: https://{domain}/{owner}/{repo}
	default:
		return (&gourl.URL{
			Scheme: "https",
			Host:   pu.Domain,
			Path:   fmt.Sprintf("/%s/%s", pu.Owner, pu.Repo),
		}).String()
	}

	// TODO: add other platforms
}

// PermaURL returns the download URL for the specific platform.
func (pu *ProductURL) PermaURL() string {
	switch pu.Domain {

	// format: https://raw.githubusercontent.com/{owner}/{repo}/{tag}/{path}
	case "github.com":
		return (&gourl.URL{
			Scheme: "https",
			Host:   "raw.githubusercontent.com",
			Path:   fmt.Sprintf("/%s/%s/%s/%s", pu.Owner, pu.Repo, pu.Tag, pu.Path),
		}).String()

	// format: https://gitlab.com/{owner}/{repo}/-/raw/{tag}/{path}
	case "gitlab.com":
		return (&gourl.URL{
			Scheme: "https",
			Host:   "gitlab.com",
			Path:   fmt.Sprintf("/%s/%s/-/raw/%s/%s", pu.Owner, pu.Repo, pu.Tag, pu.Path),
		}).String()

	// format: https://wikifactory.com/{owner}/{repo}/contributions/{tag}/file/{path}
	case "wikifactory.com":
		url := &gourl.URL{
			Scheme: "https",
			Host:   "projects.fablabs.io",
		}
		if pu.Path == "" {
			url.Path = fmt.Sprintf("/%s/%s/contributions/%s", pu.Owner, pu.Repo, pu.Tag)
		} else {
			url.Path = fmt.Sprintf("/%s/%s/contributions/%s/file/%s", pu.Owner, pu.Repo, pu.Tag, pu.Path)
		}
		return url.String()
	}

	// TODO: add other platforms

	return ""
}
