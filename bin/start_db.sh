cd bin/postgres

docker build -t pgdb .

docker run -p 5432:5432 --name testdb -e POSTGRES_PASSWORD=secret_password_here --rm pgdb