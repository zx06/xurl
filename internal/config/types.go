package config

type File struct {
	Format     string             `yaml:"format"`
	LibraryDir string             `yaml:"library_dir"`
	Profiles   map[string]Profile `yaml:"profiles"`
}

type Profile struct {
	Description    string            `yaml:"description"`
	Format         string            `yaml:"format"`
	BaseURL        string            `yaml:"base_url"`
	Vars           map[string]string `yaml:"vars"`
	Headers        map[string]string `yaml:"headers"`
	AllowPlaintext bool              `yaml:"allow_plaintext"`
	Auth           Auth              `yaml:"auth"`
}

type Auth struct {
	Bearer string     `yaml:"bearer"`
	Basic  *BasicAuth `yaml:"basic"`
	Header *HeaderKey `yaml:"header"`
}

type BasicAuth struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type HeaderKey struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type Resolved struct {
	ConfigPath  string
	ProfileName string
	Format      string
	LibraryDir  string
	HistoryPath string
	Profile     Profile
	AllProfiles map[string]Profile
}

type Options struct {
	ConfigPath string

	CLIProfile    string
	CLIProfileSet bool
	CLIFormat     string
	CLIFormatSet  bool

	EnvProfile string
	EnvFormat  string

	HomeDir string
	WorkDir string
}

type ProfileInfo struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	BaseURL     string `json:"base_url,omitempty"`
}

type ProfileDetail struct {
	Name           string            `json:"name"`
	Description    string            `json:"description,omitempty"`
	BaseURL        string            `json:"base_url,omitempty"`
	Vars           map[string]string `json:"vars,omitempty"`
	Headers        map[string]string `json:"headers,omitempty"`
	AllowPlaintext bool              `json:"allow_plaintext"`
	Auth           AuthMasked        `json:"auth,omitempty"`
}

type AuthMasked struct {
	Bearer string `json:"bearer,omitempty"`
	Basic  *struct {
		Username string `json:"username,omitempty"`
		Password string `json:"password,omitempty"`
	} `json:"basic,omitempty"`
	Header *struct {
		Name  string `json:"name,omitempty"`
		Value string `json:"value,omitempty"`
	} `json:"header,omitempty"`
}
