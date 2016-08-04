/*
Segment Kit is a toolkit of packages aimed to solve few problems with
flexibility and good architecture. From the beginning we want to make sure that
every new service we write has some solid grounds rather than repeating the same things every time.
The toolkit aimes to be a foundation and our goal is to keep it as simple as possible.

    func main() {
		// TODO(vinceprignano) example here
    }

Kit make it easy to use out of the box goodies, like config, logs and metrics.

-- Logging

Kit has a logging system already built-in, we have an interface that could be satisfied by
custom logging libraries, if necessary. By default we configure and recommend to use
	https://github.com/cihub/seelog
Seelog make it easier for us to rely on a well maintained library, that has a wide range
of features, like JSON, XML encodings, flexible message formatting (with line numbers, etc),
buffered writers, and so on.

-- Configuration

Kit makes it easy to have a flexible configuration system. The configuration package can be
customized to specific needs. Before running using Run, you need to define a configuration
schema the `config` package.

	func main() {
		// Set the configuration providers
		//
		// Note: this call is optional, by default, the providers are {CommandLine, Environment, File}
		config.SetProviders([]config.ProviderType{config.CommandLine})

		// Run the Service
		kit.Run(schema.Service{
			Handler: func
			Name: "service-name",
			Version: "v1.0.0",
			Config: []config.Value{
				config.Value{
					Key: "desired.key",
					Required: true,
					Type: config.StringType,
					Default: "default value for desired key",
				},
				config.Value{
					Key: "optional.key",
					Type: config.IntType,
				},
			},
	}

Given a key like "first.second", the configuration package queries the providers
in a different way:
	"--firstsecond" for Command Line
	"FIRST_SECOND" for Environment
	"first.second" for File
*/
package kit
