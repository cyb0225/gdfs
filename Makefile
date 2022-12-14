###  =================================== Env  ==================================================


###  =================================== Command  ==================================================
# create cobra command doc files
.PYONY: cmd.gendocs
cmd.gendocs: 
	@-echo gendocs...
	@go run cmd/gendocs/main.go
	@-echo success!



# gen xx.pb.go file from xx.proto file
.PYONY: gen.proto
gen.proto:
	@-echo build proto...
	@protoc --go_out=. --go_opt=paths=source_relative \
	 --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	 ./proto/datanode/datanode.proto
	@protoc --go_out=. --go_opt=paths=source_relative \
	 --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	 ./proto/namenode/namenode.proto
	@-echo success!



## start rpc server
.PYONY: build.client
build.client: client.o
	@ go build -o client.o cmd/client/main.go
# @./client.o


port = 50051
.PYONY: build.namenode
build.namenode:
	@ go build -o namenode.o cmd/namenode/main.go


.PYONY: build.datanode
build.datanode: datanode.o
	@ go build -o datanode.o cmd/datanode/main.go


.PYONY: datanode.cluster
datanode.cluster: datanode.o
	@ make clean.logs
	@ go build -o datanode.o cmd/datanode/main.go
	@./datanode.o -p=5000 &
	@./datanode.o -p=5001 &


###  =================================== Build File ==================================================

# create client binary file
## tips: I store binary file in project root directory instead of in ./bin directory
## 		 because, I need to use the config file in the same relative position as usual.
client.o: cmd/client/main.go
	@-echo build client...
	@go build -o client.o cmd/client/main.go
	@-echo success!

# create namenode binary file
namenode.o: cmd/namenode/main.go
	@-echo build namenode...
	@go build -o namenode.o cmd/namenode/main.go
	@-echo success!

# create datandoe binary file
datanode.o: cmd/datanode/main.go
	@-echo build datanode...
	@go build -o datanode.o cmd/datanode/main.go
	@-echo success!


.PYONY: clean
clean:
	@-echo start clean...
	@rm -rf *.o
	@-echo clean over!

.PYONY: clean.logs
clean.logs:
	@rm -r storage/log/*/*.log