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
        
        // Construir la imagen Docker
        stage('Build Docker Image') {
            steps {
                script {
                    // Construye la imagen Docker a partir del Dockerfile
                    sh 'docker build -t $IMAGE_NAME .'
                }
            }
        }
        
        // Ejecutar el contenedor Docker
        stage('Run Docker Container') {
            steps {
                script {
                    // Corre contenedor en segundo plano con el puerto 8080 mapeado
                    sh 'docker run -d -p 8080:8080 --name $IMAGE_NAME-container $IMAGE_NAME'
                }
            }
        }
    }
    
    post {
        always {
            // Pasos que siempre deben ejecutarse (limpieza, notificaciones etc...)
            echo 'Pipeline terminado.'
        }
    }
}
