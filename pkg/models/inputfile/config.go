package inputfile

type Config struct {
	Suman         Suman         `yaml:"suman"`
	SMTP          SMTP          `yaml:"smtp"`
	Dirs          Dirs          `yaml:"dirs"`
	LogLevel      Loglevel      `yaml:"loglevel"`
	Migrate       Migrate       `yaml:"migrate"`
	ErrorHandling ErrorHandling `yaml:"error_handling"`
	Maintenance   Maintenance   `yaml:"maintenance"`
	BootstrapRepo BootstrapRepo `yaml:"bootstrap-repo"`
}

type Suman struct {
	Server              string `yaml:"server"`
	User                string `yaml:"user"`
	Password            string `yaml:"password"`
	Timeout             int    `yaml:"timeout"`
	SslCertificateCheck bool   `yaml:"ssl_certificate_check"`
	RetryCount          int    `yaml:"retry_count"`
}

type SMTP struct {
	Sendmail  bool     `yaml:"sendmail"`
	Receivers []string `yaml:"receivers"`
	Sender    string   `yaml:"sender"`
	Server    string   `yaml:"server"`
}

type Dirs struct {
	LogDir          string `yaml:"log_dir"`
	ScriptsDir      string `yaml:"scripts_dir"`
	UpdateScriptDir string `yaml:"update_script_dir"`
}

type Loglevel struct {
	File   string `yaml:"file"`
	Screen string `yaml:"screen"`
}

type Migrate struct {
	SkipChannels   []string `yaml:"skip_channels"`
	RenameChannels struct {
		PartOrChannelnameToChange string `yaml:"<part or channelname to change>"`
	} `yaml:"rename_channels"`
	ProjectLabels struct {
		TestSp6 string `yaml:"test*sp6"`
		ProdSp6 string `yaml:"prod*sp6"`
	} `yaml:"project_labels"`
}

type ErrorHandling struct {
	Script        string `yaml:"script"`
	Update        string `yaml:"update"`
	Spmig         string `yaml:"spmig"`
	Configupdate  string `yaml:"configupdate"`
	Reboot        string `yaml:"reboot"`
	TimeoutPassed string `yaml:"timeout_passed"`
}

type Maintenance struct {
	WaitBetweenSystems  int                           `yaml:"wait_between_systems"`
	ExcludeForPatch     []string                      `yaml:"exclude_for_patch"`
	SpMigrationProjects map[string]SpMigrationProject `yaml:"sp_migration_project"`
	SpMigrations        map[string]SpMigration        `yaml:"sp_migration"`
	ExceptionSp         map[string]ExceptionSpServer  `yaml:"exception_sp"`
}

type SpMigrationProject struct {
	Project []string
}

type SpMigration struct {
	Project []string
}

type ExceptionSpServer struct {
	Server []string
}

type BootstrapRepo struct {
	Command string `yaml:"command"`
	Repos   map[string]Repo
}

type Repo struct {
	RepoName []string
}
