pipeline {
    agent any
    
    environment {
        IMAGE_NAME = 'super-descuentos'   // Nombre de la imagen que se generará
    }
    
    stages {
        // Obtener el código del repositorio
        stage('Checkout') {
            steps {
                // Clona el repo
                git 'https://github.com/caesarnetyet/super-descuentos'
            }
        }

        // Ejecutar los tests
        stage('Run Tests') {
            steps {
                script {
                    // Aquí puedes ejecutar los tests, por ejemplo con npm, pytest o cualquier framework que estés utilizando
                    // Ejemplo con npm (suponiendo que uses Node.js):
                    // sh 'npm install'
                    // sh 'npm test'
                    // Ejemplo con pytest (suponiendo que uses Python):
                    // sh 'pytest'
                    
                    // Si el comando falla, el pipeline se detendrá automáticamente
                    sh 'echo "Ejecutando los tests..."' // Reemplazar por el comando real para ejecutar los tests
                }
            }
        }
        
        // Construir la imagen Docker solo si los tests son exitosos
        stage('Build Docker Image') {
            when {
                branch 'main' // O cualquier otra condición para este paso
            }
            steps {
                script {
                    // Eliminar la imagen Docker si ya existe
                    sh '''
                        if [[ "$(docker images -q $IMAGE_NAME 2> /dev/null)" != "" ]]; then
                            docker rmi -f $IMAGE_NAME
                        fi
                    '''
                    
                    // Construye la imagen Docker a partir del Dockerfile
                    sh 'docker build -t $IMAGE_NAME .'
                }
            }
        }

        stage('Run Docker Container') {
            when {
                branch 'main' // O cualquier otra condición para este paso
            }
            steps {
                script {
                    // Eliminar el contenedor si ya está corriendo
                    sh '''
                        if [[ "$(docker ps -q -f name=$IMAGE_NAME-container)" != "" ]]; then
                            docker stop $IMAGE_NAME-container
                            docker rm $IMAGE_NAME-container
                        fi
                    '''
                    
                    // Corre contenedor en segundo plano con el puerto 8080 mapeado
                    sh 'docker run -d -p 8080:8080 --name $IMAGE_NAME-container $IMAGE_NAME'
                }
            }
        }

    }
    
    post {
        success {
            // Si los tests y los pasos anteriores son exitosos, se muestra este mensaje
            echo 'Pipeline ejecutado correctamente: Los tests pasaron, la imagen y el contenedor fueron creados.'
        }
        
        failure {
            // Si alguno de los pasos falla, se muestra este mensaje
            echo 'Pipeline fallido: Los tests no pasaron, la imagen y el contenedor no fueron creados.'
        }
        
        always {
            // Pasos que siempre deben ejecutarse (limpieza, notificaciones etc...)
            echo 'Pipeline terminado.'
        }
    }
}
