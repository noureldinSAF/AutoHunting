package client

import (
	"regexp"

	docker "github.com/fsouza/go-dockerclient"
)

var (
	GeneralPattern = regexp.MustCompile(`(?i)(\"|')?([a-z0-9_-]+)?((key|pass|user|username|pwd|credentials|auth|password|pwd|Ldap|Jenkins|ftp|dotfiles|JDBC|config|connectionstring|ssh|creds|secret|cred|access|Bearer|token|passwd|api|admin|private|bash|aws|s3|cookie)){1,}([a-z0-9 _[:space:]-]+)?(\"|')?(=>|=|:|,|\\+)(( )?(\"|'|return|{))?([a-z0-9 _[:space:]-=\.])+(( )?(\"|'|return|{))`)
)

type Config struct {
	Signatures []*Pattern `yaml:"signatures"`
}

type Pattern struct {
	Name  string         `yaml:"name"`
	Value string         `yaml:"value"`
	regex *regexp.Regexp `yaml:"-"` // Added field with yaml tag to ignore
}



type DockerScan struct {
	client            *docker.Client
	singleVersionScan bool
	imageName         string
	version           string
	workDir           string
	matches           []*SecretMatch}

// SecretMatch represents a secret match with the secret value and file path.
type SecretMatch struct {
	Secret   string `json:"secret"`
	FilePath string `json:"file_path"`
}

type TagsList struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}
