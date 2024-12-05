pipeline {
    agent {
        docker { image 'golang:1.23.2' }
    }

    environment {
        IMAGE_NAME = 'super-descuentos'
    }

    stages {
        // Checkout del código
        stage('Checkout') {
            steps {
                git 'https://github.com/caesarnetyet/super-descuentos'
            }
        }

        // Ejecutar pruebas
        stage('Run Go Tests') {
            steps {
                script {
                    sh 'go test ./...'
                }
            }
        }

        // Construcción de imagen Docker
        stage('Build Docker Image') {
            steps {
                script {
                    sh 'docker build -t $IMAGE_NAME .'
                }
            }
        }

        // Ejecutar contenedor principal
        stage('Run Main Container') {
            steps {
                script {
                    sh '''
                    if [ $(docker ps -aq -f name=$IMAGE_NAME-container) ]; then
                        docker stop $IMAGE_NAME-container || true
                        docker rm $IMAGE_NAME-container || true
                    fi
                    '''
                    sh 'docker run -d -p 8080:8080 --name $IMAGE_NAME-container $IMAGE_NAME'
                }
            }
        }
    }

    post {
        always {
            script {
                sh 'docker-compose down || true'
            }
            echo 'Pipeline terminado.'
        }
    }
}
