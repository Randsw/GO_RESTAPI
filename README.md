# Simple REST application
## Get all 
curl -i http://localhost:8080/api/v1/records
## Get entry by field
curl -i http://localhost:8080/api/v1/records?name=Sasha
## Get entry by id
curl -i http://localhost:8080/api/v1/records/3
## Add entry
curl -i http://localhost:8080/api/v1/records \
  -d '{"Name":"Tom","Surname":"Brady", "Gender":"Male", "email":"brady@example.com"}'
## Edit entry
curl -i http://localhost:8080/api/v1/records/3 -XPUT \
  -d '{"id":3,"name":"Aaron","surname":"Rodgers","gender":"Male","email":"rodgers_gb@example.com"}'
## Delete entry
curl -i http://localhost:8080/api/v1/records/2 -XDELETE



