
pipeline {
    agent {
        kubernetes {
            cloud 'kubernetes'
            defaultContainer 'go'

            yaml '''
apiVersion: v1
kind: Pod
spec:
  securityContext:
    runAsUser: 1000
    runAsGroup: 1000

  containers:
    - name: go
      image: golang:1.25-bookworm
      command:
        - cat
      tty: true

      env:
        - name: GOPATH
          value: /tmp/gopath
        - name: GOCACHE
          value: /tmp/gocache

      resources:
        requests:
          cpu: 100m
          memory: 256Mi
        limits:
          cpu: 1000m
          memory: 1Gi
'''
        }
    }

    options {
        skipDefaultCheckout(true)
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
                sh 'git log -1 --oneline'
            }
        }

        stage('Test') {
            steps {
                sh 'go version'
                sh 'go test ./...'
            }
        }

        stage('Verify Dockerfile') {
            steps {
                sh 'test -f Dockerfile'
                sh 'echo "Ready for container image building"'
            }
        }
    }
}

