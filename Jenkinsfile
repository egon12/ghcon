pipeline {
	agent any;

	environment {
		def go = tool name: "go14", type: "go";
		PATH = "${env.PATH}:${go}/bin";
	}
	stages {
		stage("Run Test") {
			steps {
				script {
					def status = sh returnStatus: true, script: 'go test -run xxxxx ./... > test_compile_out.txt 2>&1 ';
					if (status != 0) {
						sh "ghr reject < test_compile_out.txt"
					}
					sh "go test ./... -count=1 -race -v -coverprofile coverage.out 2>&1";
				}
			}
		}
	}
}
