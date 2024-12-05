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
        
        // Ejecutar el contenedor principal
        stage('Run Main Container') {
            steps {
                script {
                    // Corre el contenedor principal basado en la imagen generada
                    sh 'docker run -d -p 8080:8080 --name $IMAGE_NAME-container $IMAGE_NAME'
                }
            }
        }

        // Ejecutar los servicios definidos en docker-compose.yml
        stage('Run Test Services') {
            steps {
                script {
                    // Levanta los servicios necesarios para pruebas
                    sh 'docker compose up -d'
                }
            }
        }
    }
    
    post {
        always {
            script {
                // Detener los servicios de pruebas
                sh 'docker compose down || true'
                
                // Publica el reporte HTML en la interfaz de Jenkins
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
            echo 'Pipeline terminado.'
        }
    }
}
