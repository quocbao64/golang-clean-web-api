pipeline {
    agent any

    environment {
        GO_VERSION = '1.22'
        APP_NAME = 'golang-clean-web-api'
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Set up Go') {
            steps {
                script {
                    // Check if Go is already installed
                    def goInstalled = sh(script: 'which go', returnStatus: true) == 0
                    
                    if (!goInstalled) {
                        // Download and install Go
                        sh "wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz"
                        sh "sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz"
                        sh "rm go${GO_VERSION}.linux-amd64.tar.gz"
                    }
                    
                    // Set up environment variables
                    env.PATH = "/usr/local/go/bin:${env.PATH}"
                    env.GOROOT = "/usr/local/go"
                    env.GOPATH = "${env.WORKSPACE}/go"
                    
                    // Verify Go installation
                    sh 'go version'
                }
            }
        }

        stage('Download Dependencies') {
            steps {
                dir('src') {
                    sh 'go mod download'
                }
            }
        }

        stage('Lint') {
            steps {
                dir('src') {
                    sh 'go install golang.org/x/lint/golint@latest'
                    sh 'golint ./...'
                }
            }
        }

        stage('Test') {
            steps {
                dir('src') {
                    sh 'go test -v ./tests/... -coverprofile=coverage.out'
                }
            }
        }

        stage('Build') {
            steps {
                dir('src') {
                    sh 'go build -o ${APP_NAME}'
                }
            }
        }

        stage('Docker Build & Push') {
            when {
                branch 'main'
            }
            steps {
                script {
                    dockerImage = docker.build("${APP_NAME}:latest")
                }
                withCredentials([usernamePassword(credentialsId: 'dockerhub', usernameVariable: 'DOCKER_USER', passwordVariable: 'DOCKER_PASS')]) {
                    sh 'echo $DOCKER_PASS | docker login -u $DOCKER_USER --password-stdin'
                    sh "docker tag ${APP_NAME}:latest $DOCKER_USER/${APP_NAME}:latest"
                    sh "docker push $DOCKER_USER/${APP_NAME}:latest"
                }
            }
        }

//         stage('Deploy') {
//             when {
//                 branch 'main'
//             }
//             steps {
//
//             }
//         }
    }

    post {
        always {
            junit '**/TEST-*.xml'
            cleanWs()
        }
        failure {
            mail to: 'your@email.com',
                 subject: "Failed Pipeline: ${currentBuild.fullDisplayName}",
                 body: "Something is wrong with ${env.BUILD_URL}"
        }
    }
}