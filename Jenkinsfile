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

        // Ejecutar los servicios de pruebas
        stage('Run Test Services') {
            steps {
                script {
                    sh 'docker-compose up -d'
                }
            }
        }

        // Construir la imagen Docker
        stage('Build Docker Image') {
            steps {
                script {
                    sh 'docker build -t $IMAGE_NAME .'
                }
            }
        }
    }

    post {
        always {
            script {
                // Detener y limpiar los servicios de pruebas
                sh 'docker-compose down'
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
