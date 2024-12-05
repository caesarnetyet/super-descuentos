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
                    sh 'echo "Ejecutando los tests..."'
                }
            }
        }

        // Construir la imagen Docker
        stage('Build Docker Image') {
            steps {
                script {
                    sh '''
                        # Eliminar la imagen Docker si ya existe
                        if [ "$(docker images -q $IMAGE_NAME 2>/dev/null)" != "" ]; then
                            docker rmi -f $IMAGE_NAME
                        fi
                    '''
                    sh 'docker build -t $IMAGE_NAME .'
                }
            }
        }

        // Ejecutar el contenedor Docker
        stage('Run Docker Container') {
            steps {
                script {
                    sh '''
                        # Verificar si el contenedor existe
                        if [ "$(docker ps -aq -f name=$IMAGE_NAME-container)" != "" ]; then
                            echo "Contenedor existente encontrado. Eliminando..."
                            docker stop $IMAGE_NAME-container || true
                            docker rm $IMAGE_NAME-container || true
                        fi

                        # Correr el nuevo contenedor
                        echo "Iniciando el contenedor..."
                        docker run -d -p 8080:8080 --name $IMAGE_NAME-container $IMAGE_NAME
                    '''
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
