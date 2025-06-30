pipeline {
   agent any

   environment {
       GO_VERSION = '1.22'
       APP_NAME = 'golang-clean-web-api'
   }

   tools {
       go "go${GO_VERSION}"
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
                   if (!tool(name: "go${GO_VERSION}", type: 'Go')) {
                       error "Go version ${GO_VERSION} is not installed."
                   }
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

       stage('Static Analysis') {
           parallel {
               stage('Lint') {
                   steps {
                       dir('src') {
                           sh 'golint ./...'
                       }
                   }
               }
               stage('Test') {
                   steps {
                       dir('src') {
                           sh 'go test -v ./tests/... -coverprofile=coverage.out 2>&1 | go-junit-report -set-exit-code > test-results.xml'
                       }
                   }
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
                   if (sh(script: 'git diff --quiet HEAD~1 src/', returnStatus: true) == 0) {
                       echo 'No changes detected in source code. Skipping Docker build.'
                   } else {
                       dockerImage = docker.build("${APP_NAME}:latest")
                       withCredentials([usernamePassword(credentialsId: 'dockerhub', usernameVariable: 'DOCKER_USER', passwordVariable: 'DOCKER_PASS')]) {
                           sh 'echo $DOCKER_PASS | docker login -u $DOCKER_USER --password-stdin'
                           sh "docker tag ${APP_NAME}:latest $DOCKER_USER/${APP_NAME}:latest"
                           sh "docker push $DOCKER_USER/${APP_NAME}:latest"
                       }
                   }
               }
           }
       }
   }

   post {
       always {
           node('') {
               junit '**/test-results.xml'
               cleanWs()
           }
       }
       failure {
           mail to: 'your@email.com',
                subject: "Failed Pipeline: ${currentBuild.fullDisplayName}",
                body: "Something is wrong with ${env.BUILD_URL}"
       }
   }
}