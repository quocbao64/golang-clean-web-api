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
                    sh "go version || wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz && sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz"
                    env.PATH = "/usr/local/go/bin:${env.PATH}"
                }
            }
        }

        stage('Download Dependencies') {
            steps {
                sh 'go mod download'
            }
        }

        stage('Lint') {
            steps {
                sh 'go install golang.org/x/lint/golint@latest'
                sh 'golint ./...'
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
                sh 'go build -o ${APP_NAME}'
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