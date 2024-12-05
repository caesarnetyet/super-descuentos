pipeline {
    agent any
    
    environment {
        IMAGE_NAME = 'super-descuentos'
    }
    
    stages {
        stage('Checkout') {
            steps {
                git 'https://github.com/caesarnetyet/super-descuentos'
            }
        }

        stage('Build Docker Image') {
            steps {
                script {
                    sh 'docker build -t $IMAGE_NAME .'
                }
            }
        }

        stage('Run Services') {
            steps {
                script {
                    // Detecta automáticamente si usa `docker compose` o `docker-compose`
                    def composeCommand = sh(script: 'docker compose version > /dev/null 2>&1 && echo "docker compose" || echo "docker-compose"', returnStdout: true).trim()
                    sh "${composeCommand} up -d"
                }
            }
        }
    }

    post {
        always {
            script {
                sh 'docker-compose down || docker compose down || true'
                publishHTML([
                    target: [
                        allowMissing: true,
                        keepAll: true,
                        reportDir: 'e2e/playwright-report',
                        reportFiles: 'index.html',
                        reportName: 'Playwright Test Report'
                    ]
                ])
            }
        }
    }
}
