
How to run tests
- to run tests normally
  go test
  
- to get test coverage
  to test -cover
  
- to get extended coverage report
  go test -cover -coverprofile="c.out"
  go tool cover -html="c.out" -o "coverage.html"
- then open the generated html file


TODOs
Image-Upload finalisieren
Thumbnails erstellen
Dockerfile erstellen
Image-Download Komponente bauen
EC2-Docker-Instanz konfigurieren
Anwendung einzeln auf EC2 deployen
Authentifizierung einbinden mit Cognito und API Gateway

CodeStar Integration - CodePipeline - CloudFormation


Lesebeschränkung einbinden

Nächsten Schritte:
Übersicht mittels thumbnails
Use Gin engine
Use go mod
Login und User-Management
