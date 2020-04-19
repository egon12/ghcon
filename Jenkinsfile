pipeline {
	agent any;
	stages {
		stage("Test") {
			steps {
				dir (project_dir) {
					sh "go test -run xxxxxx ./...";
					sh "go test ./... -count=1 -race -v -coverprofile coverage.out 2>&1 | go-junit-report --set-exit-code > report.xml"
				}
			}
		}
	}
}
