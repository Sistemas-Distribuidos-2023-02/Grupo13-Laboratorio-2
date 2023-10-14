docker-latinoamerica:
	docker-compose -f docker-compose-latinoamerica.yml build
	docker-compose -f docker-compose-latinoamerica.yml up --remove-orphans

docker-oms:
	docker-compose -f docker-compose-oms.yml build
	docker-compose -f docker-compose-oms.yml up --remove-orphans

docker-dn1:
	docker-compose -f docker-compose-dn1.yml build
	docker-compose -f docker-compose-dn1.yml up --remove-orphans

