install-tools:
	@echo installing tools
	@go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	@go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go install github.com/bufbuild/buf/cmd/buf@latest
	@go install github.com/vektra/mockery/v2@latest
	@go install github.com/go-swagger/go-swagger/cmd/swagger@latest
	@go install github.com/cucumber/godog/cmd/godog@latest
	@echo done

generate:
	@echo running code generation
	@go generate ./...
	@echo done

run-activity:
	docker compose --profile activity up

run-frontend:
	docker compose --env-file ./docker/.env --profile frontend up

run-categories:
	docker compose --profile categories up

run-micro:
	docker compose --env-file ./docker/.env --profile microservices up

run-payments:
	docker compose --profile payments up

run-websocket:
	docker compose --profile websocket up

run-messages:
	docker compose --profile messages up

run-media:
	docker compose --profile media up

run-users:
	docker compose --profile users up

run-wishlists:
	docker compose --profile wishilist up


run-ci:
	docker compose --env-file ./docker/.env --profile ci up  --force-recreate -d


down-micro:
	docker compose --env-file ./docker/.env --profile microservices down

down-ci:
	docker compose --profile ci down

build-ci:
	docker build --no-cache -t middleman-ci --file docker/jenkins/Dockerfile --build-arg service=ci .

build-frontend:
	docker build --no-cache -t middleman-frontend --file docker/Dockerfile.frontend --build-arg service=frontend .

tag-front:
	docker tag middleman-frontend:latest registry.sfx-markt.de/frontend:latest

rebuild-front:
	docker compose --env-file ./docker/.env --profile frontend up --build --force-recreate

down-front:
	docker compose --env-file ./docker/.env --profile frontend down
#down-mic:
#	docker image rm middleman-baskets middleman-cosec middleman-notifications middleman-ordering middleman-payments middleman-search middleman-users

#build-micro: build-reviews build-following build-posts build-categories  build-geocoding build-mailer build-shipping build-products build-wishlists build-offers build-media build-support build-newsletters build-baskets build-activity build-messages build-comments build-notifications build-ordering build-payments build-search build-users
build-micro:  build-geocoding build-mailer build-shipping build-products build-wishlists build-offers build-media build-support build-newsletters build-baskets build-activity build-messages build-comments build-notifications build-ordering build-payments build-search build-users

build-websocket: build-messages build-comments


build-erp:
	docker build -t middleman-erp --file docker/Dockerfile.microservices --build-arg service=erp -t registry.sfx-markt.de/middleman-erp:latest .

build-erp-no-cache:
	docker build --no-cache -t middleman-erp --file docker/Dockerfile.microservices --build-arg service=erp .

build-vectors:
	docker build -t middleman-vectors --file docker/Dockerfile.microservices --build-arg service=vectors -t registry.sfx-markt.de/middleman-vectors:latest .

build-vectors-no-cache:
	docker build --no-cache -t middleman-vectors --file docker/Dockerfile.microservices --build-arg service=vectors .

build-assistants:
	docker build -t middleman-assistants --file docker/Dockerfile.microservices --build-arg service=assistants -t registry.sfx-markt.de/middleman-assistants:latest .

build-assistants-no-cache:
	docker build --no-cache -t middleman-assistants --file docker/Dockerfile.microservices --build-arg service=assistants .
build-managers:
	docker build -t middleman-managers --file docker/Dockerfile.microservices --build-arg service=managers -t registry.sfx-markt.de/middleman-managers:latest .

build-managers-no-cache:
	docker build --no-cache -t middleman-managers --file docker/Dockerfile.microservices --build-arg service=managers .
build-metrics:
	docker build -t middleman-metrics --file docker/Dockerfile.microservices --build-arg service=metrics -t registry.sfx-markt.de/middleman-metrics:latest .

build-metrics-no-cache:
	docker build --no-cache -t middleman-metrics --file docker/Dockerfile.microservices --build-arg service=metrics .

build-services:
	docker build -t middleman-services --file docker/Dockerfile.microservices --build-arg service=services -t registry.sfx-markt.de/middleman-services:latest .

build-services-no-cache:
	docker build --no-cache -t middleman-services --file docker/Dockerfile.microservices --build-arg service=services .

build-reviews:
	docker build -t middleman-reviews --file docker/Dockerfile.microservices --build-arg service=reviews -t registry.sfx-markt.de/middleman-reviews:latest .

build-reviews-no-cache:
	docker build --no-cache -t middleman-reviews --file docker/Dockerfile.microservices --build-arg service=reviews .

build-following:
	docker build -t middleman-following --file docker/Dockerfile.microservices --build-arg service=following -t registry.sfx-markt.de/middleman-following:latest .

build-following-no-cache:
	docker build --no-cache -t middleman-following --file docker/Dockerfile.microservices --build-arg service=following .

build-geocoding:
	docker build -t middleman-geocoding --file docker/Dockerfile.microservices --build-arg service=geocoding -t registry.sfx-markt.de/middleman-geocoding:latest .

build-geocoding-no-cache:
	docker build --no-cache -t middleman-geocoding --file docker/Dockerfile.microservices --build-arg service=geocoding .

build-merchant:
	docker build -t middleman-merchant --file docker/Dockerfile.microservices --build-arg service=merchant -t registry.sfx-markt.de/middleman-merchant:latest .

build-merchant-no-cache:
	docker build --no-cache -t middleman-merchant --file docker/Dockerfile.microservices --build-arg service=merchant .

build-offers:
	docker build -t middleman-offers --file docker/Dockerfile.microservices --build-arg service=offers -t registry.sfx-markt.de/middleman-offers:latest .

build-offers-no-cache:
	docker build --no-cache -t middleman-offers --file docker/Dockerfile.microservices --build-arg service=offers .

build-newsletters:
	docker build -t middleman-newsletters --file docker/Dockerfile.microservices --build-arg service=newsletters -t registry.sfx-markt.de/middleman-newsletters:latest .

build-newsletters-no-cache:
	docker build -t middleman-newsletters --file docker/Dockerfile.microservices --build-arg service=newsletters .

build-wishlists:
	docker build -t middleman-wishlists --file docker/Dockerfile.microservices --build-arg service=wishlists -t registry.sfx-markt.de/middleman-wishlists:latest .

build-wishlists-no-cache:
	docker build  --no-cache -t middleman-wishlists --file docker/Dockerfile.microservices --build-arg service=wishlists .

build-cosec:
	docker build -t middleman-cosec --file docker/Dockerfile.microservices --build-arg service=cosec -t registry.sfx-markt.de/middleman-cosec:latest .

build-products:
	docker build -t middleman-products --file docker/Dockerfile.microservices --build-arg service=products -t registry.sfx-markt.de/middleman-products:latest .

build-products-no-cache:
	docker build --no-cache -t middleman-products --file docker/Dockerfile.microservices --build-arg service=products .

build-categories:
	docker build -t middleman-categories --file docker/Dockerfile.microservices --build-arg service=categories -t registry.sfx-markt.de/middleman-categories:latest .

build-categories-no-cache:
	docker build -t middleman-categories --file docker/Dockerfile.microservices --build-arg service=categories .

build-posts:
	docker build -t middleman-posts --file docker/Dockerfile.microservices --build-arg service=posts -t registry.sfx-markt.de/middleman-posts:latest .

build-posts-no-cache:
	docker build --no-cache  -t middleman-posts --file docker/Dockerfile.microservices --build-arg service=posts .

build-activity:
	docker build -t middleman-activity --file docker/Dockerfile.microservices --build-arg service=activity -t registry.sfx-markt.de/middleman-activity:latest .

build-activity-no-cache:
	docker build --no-cache -t middleman-activity --file docker/Dockerfile.microservices --build-arg service=activity -t registry.sfx-markt.de/middleman-activity:latest .

build-notifications:
	docker build -t middleman-notifications --file docker/Dockerfile.microservices --build-arg service=notifications -t registry.sfx-markt.de/middleman-notifications:latest .

build-notifications-no-cache:
	docker build --no-cache -t middleman-notifications --file docker/Dockerfile.microservices --build-arg service=notifications .

build-scheduler:
	docker build -t middleman-scheduler --file docker/Dockerfile.microservices --build-arg service=scheduler -t registry.sfx-markt.de/middleman-scheduler:latest .

build-scheduler-no-cache:
	docker build --no-cache -t middleman-scheduler --file docker/Dockerfile.microservices --build-arg service=scheduler .

build-messages:
	docker build -t middleman-messages --file docker/Dockerfile.microservices --build-arg service=messages -t registry.sfx-markt.de/middleman-messages:latest .

build-messages-no-cache:
	docker build --no-cache -t middleman-messages --file docker/Dockerfile.microservices --build-arg service=messages .

build-comments:
	docker build --no-cache  -t middleman-comments --file docker/Dockerfile.microservices --build-arg service=comments -t registry.sfx-markt.de/middleman-comments:latest .

build-comments-no-cache:
	docker build -t middleman-comments --file docker/Dockerfile.microservices --build-arg service=comments .

build-users:
	docker build -t middleman-users --file docker/Dockerfile.microservices --build-arg service=users -t registry.sfx-markt.de/middleman-users:latest .

build-users-no-cache:
	docker build --no-cache -t middleman-users --file docker/Dockerfile.microservices --build-arg service=users .

build-baskets:
	docker build -t middleman-baskets --file docker/Dockerfile.microservices --build-arg service=baskets -t registry.sfx-markt.de/middleman-baskets:latest .

build-baskets-no-cache:
	docker build -t middleman-baskets --file docker/Dockerfile.microservices --build-arg service=baskets .

build-ordering:
	docker build -t middleman-ordering --file docker/Dockerfile.microservices --build-arg service=ordering -t registry.sfx-markt.de/middleman-ordering:latest .

build-ordering-no-cache:
	docker build --no-cache -t middleman-ordering --file docker/Dockerfile.microservices --build-arg service=ordering .

build-payments:
	docker build -t middleman-payments --file docker/Dockerfile.microservices --build-arg service=payments -t registry.sfx-markt.de/middleman-payments:latest .

build-payments-no-cache:
	docker build --no-cache -t middleman-payments --file docker/Dockerfile.microservices --build-arg service=payments .

build-search:
	docker build -t middleman-search --file docker/Dockerfile.microservices --build-arg service=search -t registry.sfx-markt.de/middleman-search:latest .

build-search-no-cache:
	docker build --no-cache -t middleman-search --file docker/Dockerfile.microservices --build-arg service=search .


build-support:
	docker build -t middleman-support --file docker/Dockerfile.microservices --build-arg service=support -t registry.sfx-markt.de/middleman-support:latest .

build-support-no-cache:
	docker build --no-cache -t middleman-support --file docker/Dockerfile.microservices --build-arg service=support .

build-media:
	docker build -t middleman-media --file docker/Dockerfile.microservices --build-arg service=media -t registry.sfx-markt.de/middleman-media:latest .


build-media-no-cache:
	docker build --no-cache -t middleman-media --file docker/Dockerfile.microservices --build-arg service=media .

build-mailer:
	docker build -t middleman-mailer --file docker/Dockerfile.mailer --build-arg service=mailer -t registry.sfx-markt.de/middleman-mailer:latest .

build-mailer-no-cache:
	docker build --no-cache -t middleman-mailer --file docker/Dockerfile.mailer --build-arg service=mailer .

build-shipping:
	docker build -t middleman-shipping --file docker/Dockerfile.microservices --build-arg service=shipping -t registry.sfx-markt.de/middleman-shipping:latest .

build-shipping-no-cache:
	docker build --no-cache -t middleman-shipping --file docker/Dockerfile.microservices --build-arg service=shipping .

clean-services:
	docker image rm middleman-baskets middleman-notifications middleman-messages middleman-comments middleman-ordering middleman-payments middleman-search middleman-users
	
prune:
	docker system prune -a -f

rmi:
	docker rmi -f $(docker images -aq)
dmi:
	docker rmi -f $(docker images -aq)

volume:
	docker volume rm -f classified_jsdata classified_pgdata classified_redis_data classified_redisearch_module classified_qdrant_storage

down:
	docker compose down

clear-cache:
	docker buildx prune -f

testing:
	docker compose  --profile testing up

stop-all:
	docker stop $(docker ps -a -q)

remove-all:
	docker rm $(docker ps -a -q)

dangle:
	docker rmi $(sudo docker images -f "dangling=true" -q)

push:
	docker tag middleman-frontend:latest registry.sfx-markt.de/middleman-frontend:latest


build-jenkins:
	docker build   --no-cache   --build-arg HOST_DOCKER_GROUP_ID=988  -t middleman-jenkins:latest   -f docker/jenkins/Dockerfile   .
	docker tag middleman-jenkins:latest registry.sfx-markt.de/middleman-jenkins-custom:latest
	docker push registry.sfx-markt.de/middleman-jenkins-custom:latest