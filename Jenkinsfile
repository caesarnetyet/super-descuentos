pipeline {
    agent any

    environment {
        IMAGE_NAME = 'super-descuentos' // Nombre de la imagen que se generará
    }

    stages {
        // Obtener el código del repositorio
        stage('Checkout') {
            steps {
                git 'https://github.com/caesarnetyet/super-descuentos'
            }
        }

        // Ejecutar los tests
        stage('Run Tests') {
            steps {
                script {
                    // Ejecutar los tests
                    sh 'echo "Ejecutando los tests..."'
                    // Sustituye con tu comando de tests, como:
                    // sh 'npm test' o sh 'pytest'
                }
            }
        }

        // Construir la imagen Docker
        stage('Build Docker Image') {
            steps {
                script {
                    // Eliminar imagen existente si ya existe
                    sh '''
                        if docker images -q $IMAGE_NAME > /dev/null; then
                            docker rmi -f $IMAGE_NAME
                        fi
                    '''
                    // Construir nueva imagen
                    sh 'docker build -t $IMAGE_NAME .'
                }
            }
        }

        // Ejecutar el contenedor Docker
        stage('Run Docker Container') {
            steps {
                script {
                    // Detener y eliminar el contenedor si ya está corriendo
                    sh '''
                        if docker ps -q -f name=$IMAGE_NAME-container > /dev/null; then
                            docker stop $IMAGE_NAME-container
                            docker rm $IMAGE_NAME-container
                        fi
                    '''
                    // Correr el contenedor
                    sh 'docker run -d -p 8080:8080 --name $IMAGE_NAME-container $IMAGE_NAME'
                }
            }
        }
    }

    post {
        success {
            echo 'Pipeline ejecutado correctamente: Los tests pasaron, la imagen y el contenedor fueron creados.'
        }
        failure {
            echo 'Pipeline fallido: Los tests no pasaron, la imagen y el contenedor no fueron creados.'
        }
        always {
            echo 'Pipeline terminado.'
        }
    }
}
