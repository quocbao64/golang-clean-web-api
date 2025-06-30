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
                        // Check if sudo is available
                        def sudoAvailable = sh(script: 'which sudo', returnStatus: true) == 0
                        
                        if (sudoAvailable) {
                            // Try to install Go using package manager with sudo
                            def packageManagerInstalled = sh(script: 'sudo apt-get update && sudo apt-get install -y golang-go', returnStatus: true) == 0
                            
                            if (!packageManagerInstalled) {
                                // Fallback: Download and install Go manually with sudo
                                sh "curl -L -o go${GO_VERSION}.linux-amd64.tar.gz https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz"
                                sh "sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz"
                                sh "rm go${GO_VERSION}.linux-amd64.tar.gz"
                            }
                        } else {
                            // No sudo available, try to install in user directory
                            sh "curl -L -o go${GO_VERSION}.linux-amd64.tar.gz https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz"
                            sh "mkdir -p ${env.WORKSPACE}/go-install"
                            sh "tar -C ${env.WORKSPACE}/go-install -xzf go${GO_VERSION}.linux-amd64.tar.gz"
                            sh "rm go${GO_VERSION}.linux-amd64.tar.gz"
                            
                            // Set up environment variables for user installation
                            env.PATH = "${env.WORKSPACE}/go-install/go/bin:${env.PATH}"
                            env.GOROOT = "${env.WORKSPACE}/go-install/go"
                            env.GOPATH = "${env.WORKSPACE}/go"
                            
                            // Verify Go installation
                            sh 'go version'
                            return
                        }
                    }
                    
                    // Set up environment variables for system installation
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
                branch 'master'
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