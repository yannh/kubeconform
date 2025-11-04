package sarif

// Invocation - The runtime environment of the analysis tool run.
type Invocation struct {
	// The environment variables associated with the analysis tool process, expressed as key/value pairs.
	EnvironmentVariables map[string]string `json:"environmentVariables,omitempty"`

	// The account under which the invocation occurred.
	Account *string `json:"account,omitempty"`

	// An array of strings, containing in order the command line arguments passed to the tool from the operating system.
	Arguments []string `json:"arguments,omitempty"`

	// The command line used to invoke the tool.
	CommandLine *string `json:"commandLine,omitempty"`

	// The Coordinated Universal Time (UTC) date and time at which the invocation ended. See "Date/time properties" in the SARIF spec for the required format.
	EndTimeUtc *string `json:"endTimeUtc,omitempty"`

	// An absolute URI specifying the location of the executable that was invoked.
	ExecutableLocation *ArtifactLocation `json:"executableLocation,omitempty"`

	// Specifies whether the tool's execution completed successfully.
	ExecutionSuccessful *bool `json:"executionSuccessful,omitempty"`

	// The process exit code.
	ExitCode *int `json:"exitCode,omitempty"`

	// The reason for the process exit.
	ExitCodeDescription *string `json:"exitCodeDescription,omitempty"`

	// The name of the signal that caused the process to exit.
	ExitSignalName *string `json:"exitSignalName,omitempty"`

	// The numeric value of the signal that caused the process to exit.
	ExitSignalNumber *int `json:"exitSignalNumber,omitempty"`

	// The machine on which the invocation occurred.
	Machine *string `json:"machine,omitempty"`

	// An array of configurationOverride objects that describe notifications related runtime overrides.
	NotificationConfigurationOverrides []*ConfigurationOverride `json:"notificationConfigurationOverrides"`

	// The id of the process in which the invocation occurred.
	ProcessID *int `json:"processId,omitempty"`

	// The reason given by the operating system that the process failed to start.
	ProcessStartFailureMessage *string `json:"processStartFailureMessage,omitempty"`

	// Key/value pairs that provide additional information about the invocation.
	Properties *PropertyBag `json:"properties,omitempty"`

	// The locations of any response files specified on the tool's command line.
	ResponseFiles []*ArtifactLocation `json:"responseFiles,omitempty"`

	// An array of configurationOverride objects that describe rules related runtime overrides.
	RuleConfigurationOverrides []*ConfigurationOverride `json:"ruleConfigurationOverrides"`

	// The Coordinated Universal Time (UTC) date and time at which the invocation started. See "Date/time properties" in the SARIF spec for the required format.
	StartTimeUtc *string `json:"startTimeUtc,omitempty"`

	// A file containing the standard error stream from the process that was invoked.
	Stderr *ArtifactLocation `json:"stderr,omitempty"`

	// A file containing the standard input stream to the process that was invoked.
	Stdin *ArtifactLocation `json:"stdin,omitempty"`

	// A file containing the standard output stream from the process that was invoked.
	Stdout *ArtifactLocation `json:"stdout,omitempty"`

	// A file containing the interleaved standard output and standard error stream from the process that was invoked.
	StdoutStderr *ArtifactLocation `json:"stdoutStderr,omitempty"`

	// A list of conditions detected by the tool that are relevant to the tool's configuration.
	ToolConfigurationNotifications []*Notification `json:"toolConfigurationNotifications"`

	// A list of runtime conditions detected by the tool during the analysis.
	ToolExecutionNotifications []*Notification `json:"toolExecutionNotifications"`

	// The working directory for the invocation.
	WorkingDirectory *ArtifactLocation `json:"workingDirectory,omitempty"`
}

// NewInvocation - creates a new
func NewInvocation() *Invocation {
	return &Invocation{
		Arguments:                          make([]string, 0),
		NotificationConfigurationOverrides: make([]*ConfigurationOverride, 0),
		ResponseFiles:                      make([]*ArtifactLocation, 0),
		RuleConfigurationOverrides:         make([]*ConfigurationOverride, 0),
		ToolConfigurationNotifications:     make([]*Notification, 0),
		ToolExecutionNotifications:         make([]*Notification, 0),
	}
}

// AddEnvironmentVariable - add a single EnvironmentVariable to the Invocation
func (e *Invocation) AddEnvironmentVariable(key, environmentVariable string) *Invocation {
	e.EnvironmentVariables[key] = environmentVariable
	return e
}

// WithEnvironmentVariables - add a EnvironmentVariables to the Invocation
func (e *Invocation) WithEnvironmentVariables(environmentVariables map[string]string) *Invocation {
	e.EnvironmentVariables = environmentVariables
	return e
}

// WithAccount - add a Account to the Invocation
func (a *Invocation) WithAccount(account string) *Invocation {
	a.Account = &account
	return a
}

// WithArguments - add a Arguments to the Invocation
func (a *Invocation) WithArguments(arguments []string) *Invocation {
	a.Arguments = arguments
	return a
}

// AddArgument - add a single Argument to the Invocation
func (a *Invocation) AddArgument(argument string) *Invocation {
	a.Arguments = append(a.Arguments, argument)
	return a
}

// WithCommandLine - add a CommandLine to the Invocation
func (c *Invocation) WithCommandLine(commandLine string) *Invocation {
	c.CommandLine = &commandLine
	return c
}

// WithEndTimeUtc - add a EndTimeUtc to the Invocation
func (e *Invocation) WithEndTimeUtc(endTimeUtc string) *Invocation {
	e.EndTimeUtc = &endTimeUtc
	return e
}

// WithExecutableLocation - add a ExecutableLocation to the Invocation
func (e *Invocation) WithExecutableLocation(executableLocation *ArtifactLocation) *Invocation {
	e.ExecutableLocation = executableLocation
	return e
}

// WithExecutionSuccessful - add a ExecutionSuccessful to the Invocation
func (e *Invocation) WithExecutionSuccessful(executionSuccessful bool) *Invocation {
	e.ExecutionSuccessful = &executionSuccessful
	return e
}

// WithExitCode - add a ExitCode to the Invocation
func (e *Invocation) WithExitCode(exitCode int) *Invocation {
	e.ExitCode = &exitCode
	return e
}

// WithExitCodeDescription - add a ExitCodeDescription to the Invocation
func (e *Invocation) WithExitCodeDescription(exitCodeDescription string) *Invocation {
	e.ExitCodeDescription = &exitCodeDescription
	return e
}

// WithExitSignalName - add a ExitSignalName to the Invocation
func (e *Invocation) WithExitSignalName(exitSignalName string) *Invocation {
	e.ExitSignalName = &exitSignalName
	return e
}

// WithExitSignalNumber - add a ExitSignalNumber to the Invocation
func (e *Invocation) WithExitSignalNumber(exitSignalNumber int) *Invocation {
	e.ExitSignalNumber = &exitSignalNumber
	return e
}

// WithMachine - add a Machine to the Invocation
func (m *Invocation) WithMachine(machine string) *Invocation {
	m.Machine = &machine
	return m
}

// WithNotificationConfigurationOverrides - add a NotificationConfigurationOverrides to the Invocation
func (n *Invocation) WithNotificationConfigurationOverrides(notificationConfigurationOverrides []*ConfigurationOverride) *Invocation {
	n.NotificationConfigurationOverrides = notificationConfigurationOverrides
	return n
}

// AddNotificationConfigurationOverride - add a single NotificationConfigurationOverride to the Invocation
func (n *Invocation) AddNotificationConfigurationOverride(notificationConfigurationOverride *ConfigurationOverride) *Invocation {
	n.NotificationConfigurationOverrides = append(n.NotificationConfigurationOverrides, notificationConfigurationOverride)
	return n
}

// WithProcessID - add a ProcessID to the Invocation
func (p *Invocation) WithProcessID(processId int) *Invocation {
	p.ProcessID = &processId
	return p
}

// WithProcessStartFailureMessage - add a ProcessStartFailureMessage to the Invocation
func (p *Invocation) WithProcessStartFailureMessage(processStartFailureMessage string) *Invocation {
	p.ProcessStartFailureMessage = &processStartFailureMessage
	return p
}

// WithProperties - add a Properties to the Invocation
func (p *Invocation) WithProperties(properties *PropertyBag) *Invocation {
	p.Properties = properties
	return p
}

// WithResponseFiles - add a ResponseFiles to the Invocation
func (r *Invocation) WithResponseFiles(responseFiles []*ArtifactLocation) *Invocation {
	r.ResponseFiles = responseFiles
	return r
}

// AddResponseFile - add a single ResponseFile to the Invocation
func (r *Invocation) AddResponseFile(responseFile *ArtifactLocation) *Invocation {
	r.ResponseFiles = append(r.ResponseFiles, responseFile)
	return r
}

// WithRuleConfigurationOverrides - add a RuleConfigurationOverrides to the Invocation
func (r *Invocation) WithRuleConfigurationOverrides(ruleConfigurationOverrides []*ConfigurationOverride) *Invocation {
	r.RuleConfigurationOverrides = ruleConfigurationOverrides
	return r
}

// AddRuleConfigurationOverride - add a single RuleConfigurationOverride to the Invocation
func (r *Invocation) AddRuleConfigurationOverride(ruleConfigurationOverride *ConfigurationOverride) *Invocation {
	r.RuleConfigurationOverrides = append(r.RuleConfigurationOverrides, ruleConfigurationOverride)
	return r
}

// WithStartTimeUtc - add a StartTimeUtc to the Invocation
func (s *Invocation) WithStartTimeUtc(startTimeUtc string) *Invocation {
	s.StartTimeUtc = &startTimeUtc
	return s
}

// WithStderr - add a Stderr to the Invocation
func (s *Invocation) WithStderr(stderr *ArtifactLocation) *Invocation {
	s.Stderr = stderr
	return s
}

// WithStdin - add a Stdin to the Invocation
func (s *Invocation) WithStdin(stdin *ArtifactLocation) *Invocation {
	s.Stdin = stdin
	return s
}

// WithStdout - add a Stdout to the Invocation
func (s *Invocation) WithStdout(stdout *ArtifactLocation) *Invocation {
	s.Stdout = stdout
	return s
}

// WithStdoutStderr - add a StdoutStderr to the Invocation
func (s *Invocation) WithStdoutStderr(stdoutStderr *ArtifactLocation) *Invocation {
	s.StdoutStderr = stdoutStderr
	return s
}

// WithToolConfigurationNotifications - add a ToolConfigurationNotifications to the Invocation
func (t *Invocation) WithToolConfigurationNotifications(toolConfigurationNotifications []*Notification) *Invocation {
	t.ToolConfigurationNotifications = toolConfigurationNotifications
	return t
}

// AddToolConfigurationNotification - add a single ToolConfigurationNotification to the Invocation
func (t *Invocation) AddToolConfigurationNotification(toolConfigurationNotification *Notification) *Invocation {
	t.ToolConfigurationNotifications = append(t.ToolConfigurationNotifications, toolConfigurationNotification)
	return t
}

// WithToolExecutionNotifications - add a ToolExecutionNotifications to the Invocation
func (t *Invocation) WithToolExecutionNotifications(toolExecutionNotifications []*Notification) *Invocation {
	t.ToolExecutionNotifications = toolExecutionNotifications
	return t
}

// AddToolExecutionNotification - add a single ToolExecutionNotification to the Invocation
func (t *Invocation) AddToolExecutionNotification(toolExecutionNotification *Notification) *Invocation {
	t.ToolExecutionNotifications = append(t.ToolExecutionNotifications, toolExecutionNotification)
	return t
}

// WithWorkingDirectory - add a WorkingDirectory to the Invocation
func (w *Invocation) WithWorkingDirectory(workingDirectory *ArtifactLocation) *Invocation {
	w.WorkingDirectory = workingDirectory
	return w
}
