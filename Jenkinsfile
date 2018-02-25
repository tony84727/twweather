pipeline {
    agent { docker 'golang'}
    stages {
        stage('Build') {
            steps {
                sh 'go build'
            }
        }
        stage('test') {
            steps {
                sh 'go test'
            }
        }
    }
}