module.exports = (grunt) =>
	cmd = "gorest-example"
	grunt.initConfig
		bgShell:
			_defaults:
				bg: false
			# Start service.
			# Assumes, that the <cmd> binary is "go install"'ed in
			# a path in the $PATH
			"rest-start":
				cmd: cmd
				bg: true
			"rest-stop":
				cmd: "pkill "+cmd
			# Compile and install the binary
			rest:
				cmd: "go install github.com/homburg/"+cmd
			echo:
				cmd: ""
		watch:
			# Watch static files for livereload
			# The server generates the "server.run" file
			# This allows grunt to watch for a fresh server
			server:
				files: ["server.run", "**/*.js", "**/*.css"]
				options:
					livereload: true
			# Build and restart the server on touched *.go-files
			rest:
				files: ["**/*.go"]
				tasks: ['build', 'bgShell:rest-stop', 'bgShell:rest-start']

	grunt.loadNpmTasks 'grunt-bg-shell'
	grunt.loadNpmTasks 'grunt-contrib-watch'

	grunt.registerTask 'default', ['bgShell:echo']
	grunt.registerTask 'build', ['bgShell:rest']
	grunt.registerTask 'start', ['bgShell:rest-start']
