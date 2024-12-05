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

        // Ejecutar las pruebas Go
        stage('Run Go Tests') {
            steps {
                script {
                    // Ejecuta las pruebas Go, y si alguna falla, termina el pipeline
                    def testResult = sh(script: 'go test ./...', returnStatus: true)
                    if (testResult != 0) {
                        error "Go tests failed. Aborting the pipeline."
                    }
                }
            }
        }

        // Construir la imagen Docker solo si los tests fueron exitosos
        stage('Build Docker Image') {
            steps {
                script {
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
                // Limpiar el entorno después de ejecutar las pruebas
                sh 'docker-compose down || true' // Para asegurar que se detienen los servicios de pruebas si existen
            }

            script {
                // Imprimir un mensaje según el resultado del pipeline
                if (currentBuild.result == 'SUCCESS') {
                    echo 'El proceso se completó correctamente.'
                } else {
                    echo 'El proceso falló. Revisa los errores.'
                }
            }

            echo 'Pipeline terminado.'
        }
    }
}
