package azuredevops

// Task.
type Task struct {
	Author string `json:"author,omitempty"`

	// Where the task appears in Azure DevOps. Use the 'Azure *' categories for Azure DevOps and Azure DevOps Server 2019. Use the other categories for Team Foundation Server 2018 and below.
	Category           string              `json:"category,omitempty"`
	DataSourceBindings []DataSourceBinding `json:"dataSourceBindings,omitempty"`

	// Allows you to define a list of demands that a build agent requires to run this build task.
	Demands []string `json:"demands,omitempty"`

	// Task is deprecated only when the latest version is marked as deprecated. Deprecated tasks appear at the end of searches under a section that is collapsed by default.
	Deprecated bool `json:"deprecated,omitempty"`

	// Detailed description of what your task does
	Description string `json:"description,omitempty"`

	// Execution options for this task
	Execution *Execution `json:"execution,omitempty"`

	// Descriptive name (spaces allowed). Must be <= 40 chars
	FriendlyName string `json:"friendlyName,omitempty"`

	// Describes groups that task properties may be logically grouped by in the UI.
	Groups       []Group `json:"groups,omitempty"`
	HelpMarkDown string  `json:"helpMarkDown,omitempty"`
	HelpURL      string  `json:"helpUrl,omitempty"`

	// A unique guid for this task
	ID     string  `json:"id,omitempty"`
	Inputs []Input `json:"inputs,omitempty"`

	// This is how the task will be displayed within the build step list - you can use variable values by using $(variablename)
	InstanceNameFormat  string    `json:"instanceNameFormat,omitempty"`
	Messages            *Messages `json:"messages,omitempty"`
	MinimumAgentVersion string    `json:"minimumAgentVersion,omitempty"`

	// Name with no spaces
	Name string `json:"name,omitempty"`

	// Describes output variables of task.
	OutputVariables []OutputVariable `json:"outputVariables,omitempty"`

	// Execution options for this task (on Post-Job stage)
	PostJobExecution *Executions `json:"postjobexecution,omitempty"`

	// Execution options for this task (on Pre-Job stage)
	PreJobExecution *Executions `json:"prejobexecution,omitempty"`
	Preview         bool        `json:"preview,omitempty"`
	ReleaseNotes    string      `json:"releaseNotes,omitempty"`

	// Restrictions on tasks
	Restrictions *Restrictions `json:"restrictions,omitempty"`
	RunsOn       []string      `json:"runsOn,omitempty"`
	Schema       string        `json:"$schema,omitempty"`

	// Toggles showing the environment variable editor in the task editor UI. Allows passing environment variables to script based tasks.
	ShowEnvironmentVariables bool               `json:"showEnvironmentVariables,omitempty"`
	SourceDefinitions        []SourceDefinition `json:"sourceDefinitions,omitempty"`

	// Always update this when you release your task, so that the agents utilize the latest code.
	Version    *Version `json:"version,omitempty"`
	Visibility []string `json:"visibility,omitempty"`
}

// Commands Restrictions on available task commands.
type Commands struct {
	Mode string `json:"mode,omitempty"`
}

// DataSourceBinding.
type DataSourceBinding struct {
	CallbackContextTemplate  string      `json:"callbackContextTemplate,omitempty"`
	CallbackRequiredTemplate string      `json:"callbackRequiredTemplate,omitempty"`
	DataSourceName           string      `json:"dataSourceName,omitempty"`
	EndpointID               string      `json:"endpointId,omitempty"`
	EndpointURL              string      `json:"endpointUrl,omitempty"`
	InitialContextTemplate   string      `json:"initialContextTemplate,omitempty"`
	Parameters               *Parameters `json:"parameters,omitempty"`
	RequestContent           string      `json:"requestContent,omitempty"`
	RequestVerb              string      `json:"RequestVerb,omitempty"`
	ResultSelector           string      `json:"resultSelector,omitempty"`
	ResultTemplate           string      `json:"resultTemplate,omitempty"`
	Target                   string      `json:"target,omitempty"`
}

// Executions Execution options for this task.
type Executions struct {
	Node        *Execution `json:"Node,omitempty"`
	Node10      *Execution `json:"Node10,omitempty"`
	Node16      *Execution `json:"Node16,omitempty"`
	PowerShell  *Execution `json:"PowerShell,omitempty"`
	PowerShell3 *Execution `json:"PowerShell3,omitempty"`
}

// Execution.
type Execution struct {
	AdditionalProperties map[string]interface{} `json:"-,omitempty"`
	ArgumentFormat       string                 `json:"argumentFormat,omitempty"`
	Platforms            []interface{}          `json:"platforms,omitempty"`

	// The target file to be executed. You can use variables here in brackets e.g. $(currentDirectory)ilename.ps1
	Target           string `json:"target"`
	WorkingDirectory string `json:"workingDirectory,omitempty"`
}

// Group.
type Group struct {
	DisplayName string `json:"displayName"`
	IsExpanded  bool   `json:"isExpanded,omitempty"`
	Name        string `json:"name"`

	// Allow's you to define a rule which dictates when the group will be visible to a user, for example "variableName1 != \"\" && variableName2 = value || variableName3 NotEndsWith value"
	VisibleRule string `json:"visibleRule,omitempty"`
}

// Input.
type Input struct {
	Aliases []string `json:"aliases,omitempty"`

	// The default value to apply to this input.
	DefaultValue interface{} `json:"defaultValue,omitempty"`

	// Setting this to the name of a group defined in 'groups' will place the input into that group.
	GroupName string `json:"groupName,omitempty"`

	// Help to be displayed when hovering over the help icon for the input. To display URLs use the format [Text To Display](http://URL)
	HelpMarkDown string `json:"helpMarkDown,omitempty"`

	// The text displayed to the user for the input label
	Label string `json:"label"`

	// The variable name to use to store the user-supplied value
	Name       string      `json:"name"`
	Options    *Options    `json:"options,omitempty"`
	Properties *Properties `json:"properties,omitempty"`

	// Whether the input is a required field (default is false).
	Required bool `json:"required,omitempty"`

	// The type that dictates the control rendered to the user.
	Type string `json:"type"`

	// Allow's you to define a rule which dictates when the input will be visible to a user, for example "variableName1 != \"\" && variableName2 = value || variableName3 NotEndsWith value"
	VisibleRule string `json:"visibleRule,omitempty"`
}

// Messages.
type Messages struct{}

// Options.
type Options struct {
	AdditionalProperties map[string]interface{} `json:"-,omitempty"`
}

// OutputVariable.
type OutputVariable struct {
	// Detailed description of the variable
	Description string `json:"description,omitempty"`

	// The variable name
	Name string `json:"name"`
}

// Parameters.
type Parameters struct{}

// Properties.
type Properties struct {
	DisableManageLink             string `json:"DisableManageLink,omitempty"`
	EditableOptions               string `json:"EditableOptions,omitempty"`
	EditorExtension               string `json:"editorExtension,omitempty"`
	EndpointFilterRule            string `json:"EndpointFilterRule,omitempty"`
	IsSearchable                  string `json:"IsSearchable,omitempty"`
	IsVariableOrNonNegativeNumber string `json:"isVariableOrNonNegativeNumber,omitempty"`
	MaxLength                     string `json:"maxLength,omitempty"`
	MultiSelect                   string `json:"MultiSelect,omitempty"`
	MultiSelectFlatList           string `json:"MultiSelectFlatList,omitempty"`
	PopulateDefaultValue          string `json:"PopulateDefaultValue,omitempty"`
	Resizable                     bool   `json:"resizable,omitempty"`
	Rows                          string `json:"rows,omitempty"`
}

// Restrictions Restrictions on tasks.
type Restrictions struct {
	// Restrictions on available task commands
	Commands *Commands `json:"commands,omitempty"`

	// Restrictions on which variables can be set via commands
	SettableVariables *SettableVariables `json:"settableVariables,omitempty"`
}

// SettableVariables Restrictions on which variables can be set via commands.
type SettableVariables struct {
	Allowed []string `json:"allowed,omitempty"`
}

// SourceDefinition.
type SourceDefinition struct {
	AuthKey     string `json:"authKey,omitempty"`
	Endpoint    string `json:"endpoint,omitempty"`
	KeySelector string `json:"keySelector,omitempty"`
	Selector    string `json:"selector,omitempty"`
	Target      string `json:"target,omitempty"`
}

// Version Always update this when you release your task, so that the agents utilize the latest code.
type Version struct {
	Major float64 `json:"Major"`
	Minor float64 `json:"Minor"`
	Patch float64 `json:"Patch"`
}
