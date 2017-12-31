sudo -u postgres createdb db3;
psql -h 127.0.0.1 -U neotek db3 < ./database.sql;
go build main.go;
./main;
