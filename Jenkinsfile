pipeline {
	agent any;

	environment {
		def go = tool name: "go14", type: "go";
		PATH = "${env.PATH}:${go}/bin";
	}
	stages {
		stage("Run Test") {
			steps {
				def status = sh returnStatus: true, script: 'go test -run xxxxx ./... 2&1> test_compile_out.txt';
				if status != 0 {
					sh "ghr reject < test_compile_out.txt"
				}
				sh "go test ./... -count=1 -race -v -coverprofile coverage.out 2>&1";
			}
		}
	}
}
