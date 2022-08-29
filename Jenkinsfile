pipeline {
    
    triggers {
        pollSCM('') // Enabling being build on Push
    }

    agent any

    options {
        // This is required if you want to clean before build
        skipDefaultCheckout(true)
    }
    
    tools { go '1.18' }
    
    environment { 
        GOBIN = "${JENKINS_HOME}/go/bin"
        PATH = "${PATH}:${GOBIN}"
        GPG_TTY = "${tty}"
    }
    
    stages {
        stage('Pre-build') { //download dependencies and code linting, vetting and formatting
            steps {
                // Clean before build
                cleanWs()
                // We need to explicitly checkout from SCM here
                checkout scm
                echo "Building ${env.JOB_NAME}..."
                sh '''
                    make download
                    make code-check
                '''
            }
        }

        stage('Compile') { //Run build to ensure that there are no compilation error
            steps {
                sh '''
                    make build
                '''
            }
        }

        stage('Test') { //run unit and integration tests
            steps {
                sh '''
                    make test
                '''
            }
        }

        stage ('Release') { //create binaries for different operating systems and release on github
          when {
            buildingTag()
          }

          environment {
            GITHUB_TOKEN = credentials('github-token')
            PASSPHRASE = credentials('gpg-passphrase')
            GPG_KEYGRIP = credentials('gpg-keygrip')
            GPG_FINGERPRINT = credentials('gpg-fingerprint')
          }

          steps {
            sh '''
                /usr/libexec/gpg-preset-passphrase --preset -P $PASSPHRASE $GPG_KEYGRIP
                curl -sfL https://goreleaser.com/static/run | bash 
            '''
          }
        }
    }
}
