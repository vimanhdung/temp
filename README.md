# Requirement
GoSersion: go1.11 <br>
BeeGo: 1.11

# Clone Repository
git clone https://mor.backlog.jp/git/INDETAIL/MH-INDETAIL-BACKEND.git indetail <br>
# Setting MySQL
With Docker: Build MySql5.7 <br>
cd indetail <br>
git checkout develop <br>
cd ServerSQL <br>
cp env-example .env <br>
docker-compose up -d <br>
---
If you do not have docker: <br>
vi conf/constants.go <br>
-> edit config infomation
# Get Package for beego
go get github.com/go-sql-driver/mysql <br>
go get github.com/boombuler/barcode <br>
go get golang.org/x/crypto/bcrypt
go get github.com/dgrijalva/jwt-go
go get gopkg.in/maddevsio/fcm.v1
# Migrate Database
With Docker: <br>
bee migrate -driver=mysql -conn="root:root@tcp(127.0.0.1:3388)/indetail" <br>
If you do not have docker: <br>
bee migrate -driver=mysql -conn="username:password@tcp(your_ip:db_port)/db_name"
# Run BeeGo
bee run