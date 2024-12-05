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
                    // Levanta los servicios necesarios para pruebas
                    // Cambia a `docker-compose` si `docker compose` no está disponible
                    sh '''
                    if docker compose version > /dev/null 2>&1; then
                        docker compose up -d
                    else
                        docker-compose up -d
                    fi
                    '''
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
                sh '''
                if docker compose version > /dev/null 2>&1; then
                    docker compose down || true
                else
                    docker-compose down || true
                fi
                '''
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
