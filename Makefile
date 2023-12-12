#.PHONY: users station auth run

users:
	@go build -o bin/users cmd/user-service/user_svc_main.go
	@./bin/users

station:
	@go build -o bin/station cmd/charge-station-service/station_srv_main.go
	@./bin/station

navi:
	@go build -o bin/navigation cmd/navigation-service/navigation_svc_main.go
	@./bin/navigation

